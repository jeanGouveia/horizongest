package rabbitmq

import (
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection gerencia a conexão com o RabbitMQ
type Connection struct {
	config     Config
	connection *amqp.Connection
	channel    *amqp.Channel
	mu         sync.RWMutex
	closed     bool
}

// NewConnection cria uma nova conexão com o RabbitMQ
func NewConnection(config Config) (*Connection, error) {
	conn := &Connection{
		config: config,
	}

	if err := conn.connect(); err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ connection: %w", err)
	}

	return conn, nil
}

// connect estabelece a conexão com o RabbitMQ
func (c *Connection) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return fmt.Errorf("connection is closed")
	}

	// Criar conexão
	conn, err := amqp.Dial(c.config.URL)
	if err != nil {
		return fmt.Errorf("failed to dial RabbitMQ: %w", err)
	}

	// Criar canal
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	// Habilitar publisher confirms se configurado
	if c.config.EnablePublisherConfirm {
		if err := ch.Confirm(false); err != nil {
			ch.Close()
			conn.Close()
			return fmt.Errorf("failed to enable publisher confirms: %w", err)
		}
	}

	c.connection = conn
	c.channel = ch

	// Configurar listener de fechamento de conexão
	go c.handleConnectionClose(c.connection.NotifyClose(make(chan *amqp.Error, 1)))

	log.Printf("RabbitMQ connected successfully to %s", c.config.URL)
	return nil
}

// handleConnectionClose monitora o fechamento da conexão e tenta reconectar
func (c *Connection) handleConnectionClose(closeChan <-chan *amqp.Error) {
	err := <-closeChan
	if err == nil {
		return
	}

	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return
	}
	c.mu.Unlock()

	log.Printf("RabbitMQ connection closed, attempting to reconnect...")

	for {
		time.Sleep(c.config.ReconnectDelay)

		c.mu.Lock()
		if c.closed {
			c.mu.Unlock()
			return
		}
		c.mu.Unlock()

		reconnectErr := c.connect()
		if reconnectErr == nil {
			log.Printf("RabbitMQ reconnected successfully")
			return
		}

		log.Printf("Failed to reconnect to RabbitMQ: %v, retrying in %v", reconnectErr, c.config.ReconnectDelay)
	}
}

// GetChannel retorna o canal AMQP atual
func (c *Connection) GetChannel() (*amqp.Channel, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.channel == nil || c.connection == nil || c.connection.IsClosed() {
		return nil, fmt.Errorf("connection not established")
	}

	return c.channel, nil
}

// Close fecha a conexão com o RabbitMQ
func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Printf("Error closing RabbitMQ channel: %v", err)
		}
	}

	if c.connection != nil {
		if err := c.connection.Close(); err != nil {
			log.Printf("Error closing RabbitMQ connection: %v", err)
		}
	}

	log.Printf("RabbitMQ connection closed")
	return nil
}

// IsClosed retorna true se a conexão estiver fechada
func (c *Connection) IsClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.closed || (c.connection != nil && c.connection.IsClosed())
}

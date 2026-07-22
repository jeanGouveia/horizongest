-- +goose Up
-- +goose StatementBegin

-- Sprint 4 - Compras: Create tables for suppliers, purchase orders, and receivings

-- Suppliers Table
CREATE TABLE IF NOT EXISTS suppliers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    cnpj VARCHAR(20),
    email VARCHAR(255),
    phone VARCHAR(20),
    address VARCHAR(500),
    city VARCHAR(100),
    state VARCHAR(2),
    zip_code VARCHAR(10),
    notes VARCHAR(1000),
    active BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at INTEGER,
    
    FOREIGN KEY (company_id) REFERENCES companies(id)
);

CREATE INDEX IF NOT EXISTS idx_suppliers_company ON suppliers(company_id);
CREATE INDEX IF NOT EXISTS idx_suppliers_active ON suppliers(active);
CREATE INDEX IF NOT EXISTS idx_suppliers_deleted_at ON suppliers(deleted_at);

-- Purchase Orders Table
CREATE TABLE IF NOT EXISTS purchase_orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    supplier_id INTEGER NOT NULL,
    order_number VARCHAR(50) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    order_date DATETIME NOT NULL,
    expected_date DATETIME,
    received_date DATETIME,
    subtotal REAL NOT NULL DEFAULT 0.0,
    tax REAL NOT NULL DEFAULT 0.0,
    discount REAL NOT NULL DEFAULT 0.0,
    total REAL NOT NULL DEFAULT 0.0,
    notes VARCHAR(1000),
    created_by INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at INTEGER,
    
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (supplier_id) REFERENCES suppliers(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_purchase_orders_company ON purchase_orders(company_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_supplier ON purchase_orders(supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_status ON purchase_orders(status);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_date ON purchase_orders(order_date);
CREATE INDEX IF NOT EXISTS idx_purchase_orders_deleted_at ON purchase_orders(deleted_at);

-- Purchase Order Items Table
CREATE TABLE IF NOT EXISTS purchase_order_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    purchase_order_id INTEGER NOT NULL,
    ingredient_id INTEGER NOT NULL,
    quantity REAL NOT NULL,
    unit VARCHAR(20) NOT NULL,
    unit_price REAL NOT NULL,
    subtotal REAL NOT NULL DEFAULT 0.0,
    received_qty REAL NOT NULL DEFAULT 0.0,
    notes VARCHAR(500),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at INTEGER,
    
    FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id),
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id)
);

CREATE INDEX IF NOT EXISTS idx_purchase_order_items_order ON purchase_order_items(purchase_order_id);
CREATE INDEX IF NOT EXISTS idx_purchase_order_items_ingredient ON purchase_order_items(ingredient_id);
CREATE INDEX IF NOT EXISTS idx_purchase_order_items_deleted_at ON purchase_order_items(deleted_at);

-- Purchase Receivings Table
CREATE TABLE IF NOT EXISTS purchase_receivings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    purchase_order_id INTEGER NOT NULL,
    received_date DATETIME NOT NULL,
    received_by INTEGER NOT NULL,
    notes VARCHAR(1000),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at INTEGER,
    
    FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id),
    FOREIGN KEY (received_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_purchase_receivings_order ON purchase_receivings(purchase_order_id);
CREATE INDEX IF NOT EXISTS idx_purchase_receivings_date ON purchase_receivings(received_date);
CREATE INDEX IF NOT EXISTS idx_purchase_receivings_deleted_at ON purchase_receivings(deleted_at);

-- Purchase Receiving Items Table
CREATE TABLE IF NOT EXISTS purchase_receiving_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    purchase_receiving_id INTEGER NOT NULL,
    purchase_order_item_id INTEGER NOT NULL,
    ingredient_id INTEGER NOT NULL,
    quantity REAL NOT NULL,
    unit VARCHAR(20) NOT NULL,
    unit_price REAL NOT NULL,
    subtotal REAL NOT NULL DEFAULT 0.0,
    notes VARCHAR(500),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at INTEGER,
    
    FOREIGN KEY (purchase_receiving_id) REFERENCES purchase_receivings(id),
    FOREIGN KEY (purchase_order_item_id) REFERENCES purchase_order_items(id),
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id)
);

CREATE INDEX IF NOT EXISTS idx_purchase_receiving_items_receiving ON purchase_receiving_items(purchase_receiving_id);
CREATE INDEX IF NOT EXISTS idx_purchase_receiving_items_order_item ON purchase_receiving_items(purchase_order_item_id);
CREATE INDEX IF NOT EXISTS idx_purchase_receiving_items_ingredient ON purchase_receiving_items(ingredient_id);
CREATE INDEX IF NOT EXISTS idx_purchase_receiving_items_deleted_at ON purchase_receiving_items(deleted_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_purchase_receiving_items_deleted_at;
DROP INDEX IF EXISTS idx_purchase_receiving_items_ingredient;
DROP INDEX IF EXISTS idx_purchase_receiving_items_order_item;
DROP INDEX IF EXISTS idx_purchase_receiving_items_receiving;
DROP TABLE IF EXISTS purchase_receiving_items;
DROP INDEX IF EXISTS idx_purchase_receivings_deleted_at;
DROP INDEX IF EXISTS idx_purchase_receivings_date;
DROP INDEX IF EXISTS idx_purchase_receivings_order;
DROP TABLE IF EXISTS purchase_receivings;
DROP INDEX IF EXISTS idx_purchase_order_items_deleted_at;
DROP INDEX IF EXISTS idx_purchase_order_items_ingredient;
DROP INDEX IF EXISTS idx_purchase_order_items_order;
DROP TABLE IF EXISTS purchase_order_items;
DROP INDEX IF EXISTS idx_purchase_orders_deleted_at;
DROP INDEX IF EXISTS idx_purchase_orders_date;
DROP INDEX IF EXISTS idx_purchase_orders_status;
DROP INDEX IF EXISTS idx_purchase_orders_supplier;
DROP INDEX IF EXISTS idx_purchase_orders_company;
DROP TABLE IF EXISTS purchase_orders;
DROP INDEX IF EXISTS idx_suppliers_deleted_at;
DROP INDEX IF EXISTS idx_suppliers_active;
DROP INDEX IF EXISTS idx_suppliers_company;
DROP TABLE IF EXISTS suppliers;

-- +goose StatementEnd

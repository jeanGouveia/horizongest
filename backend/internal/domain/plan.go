package domain

import "time"

// Plan represents a subscription plan for companies
type Plan struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Slug        string    `json:"slug" gorm:"size:50;uniqueIndex;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       Money     `json:"price" gorm:"type:bigint;default:0"`
	Currency    string    `json:"currency" gorm:"size:3;default:'BRL'"`
	Interval    string    `json:"interval" gorm:"size:20;default:'monthly'"` // monthly, yearly
	MaxUsers    int       `json:"max_users" gorm:"default:1"`
	MaxProducts int       `json:"max_products" gorm:"default:100"`
	Features    string    `json:"features" gorm:"type:text"` // JSON string of features
	Active      bool      `json:"active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name for Plan
func (Plan) TableName() string {
	return "plans"
}

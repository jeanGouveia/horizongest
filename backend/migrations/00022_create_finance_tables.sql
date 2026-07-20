-- Sprint 4 - Financeiro: Create tables for financial transactions and categories

-- Transaction Categories Table
CREATE TABLE IF NOT EXISTS transaction_categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL,
    color VARCHAR(7),
    active BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at INTEGER,
    
    FOREIGN KEY (company_id) REFERENCES companies(id)
);

CREATE INDEX IF NOT EXISTS idx_transaction_categories_company ON transaction_categories(company_id);
CREATE INDEX IF NOT EXISTS idx_transaction_categories_type ON transaction_categories(type);
CREATE INDEX IF NOT EXISTS idx_transaction_categories_active ON transaction_categories(active);
CREATE INDEX IF NOT EXISTS idx_transaction_categories_deleted_at ON transaction_categories(deleted_at);

-- Transactions Table
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    type VARCHAR(20) NOT NULL,
    amount REAL NOT NULL,
    description VARCHAR(1000),
    date DATETIME NOT NULL,
    reference VARCHAR(255),
    created_by INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at INTEGER,
    
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (category_id) REFERENCES transaction_categories(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_transactions_company ON transactions(company_id);
CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category_id);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);
CREATE INDEX IF NOT EXISTS idx_transactions_deleted_at ON transactions(deleted_at);

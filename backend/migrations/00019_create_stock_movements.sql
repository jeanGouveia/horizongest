-- Sprint 4 - Stock Movements: Create tables for stock movements and inventories

-- Stock Movements Table
CREATE TABLE IF NOT EXISTS stock_movements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    ingredient_id INTEGER NOT NULL,
    type VARCHAR(20) NOT NULL,
    quantity REAL NOT NULL,
    previous_stock REAL NOT NULL,
    new_stock REAL NOT NULL,
    reason VARCHAR(500),
    reference_type VARCHAR(50),
    reference_id INTEGER,
    performed_by INTEGER NOT NULL,
    performed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at INTEGER,
    
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id),
    FOREIGN KEY (performed_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_stock_movements_company ON stock_movements(company_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_ingredient ON stock_movements(ingredient_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_reference ON stock_movements(reference_type, reference_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_deleted_at ON stock_movements(deleted_at);

-- Stock Inventories Table
CREATE TABLE IF NOT EXISTS stock_inventories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    inventory_date DATETIME NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    notes VARCHAR(1000),
    performed_by INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    deleted_at INTEGER,
    
    FOREIGN KEY (company_id) REFERENCES companies(id),
    FOREIGN KEY (performed_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_stock_inventories_company ON stock_inventories(company_id);
CREATE INDEX IF NOT EXISTS idx_stock_inventories_date ON stock_inventories(inventory_date);
CREATE INDEX IF NOT EXISTS idx_stock_inventories_status ON stock_inventories(status);
CREATE INDEX IF NOT EXISTS idx_stock_inventories_deleted_at ON stock_inventories(deleted_at);

-- Stock Inventory Items Table
CREATE TABLE IF NOT EXISTS stock_inventory_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    inventory_id INTEGER NOT NULL,
    ingredient_id INTEGER NOT NULL,
    expected_stock REAL NOT NULL,
    actual_stock REAL NOT NULL,
    difference REAL NOT NULL,
    reason VARCHAR(500),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at INTEGER,
    
    FOREIGN KEY (inventory_id) REFERENCES stock_inventories(id),
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id)
);

CREATE INDEX IF NOT EXISTS idx_stock_inventory_items_inventory ON stock_inventory_items(inventory_id);
CREATE INDEX IF NOT EXISTS idx_stock_inventory_items_ingredient ON stock_inventory_items(ingredient_id);
CREATE INDEX IF NOT EXISTS idx_stock_inventory_items_deleted_at ON stock_inventory_items(deleted_at);

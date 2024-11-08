CREATE TABLE IF NOT EXISTS tables (
    id BLOB(16) PRIMARY KEY,
    status TEXT NOT NULL CHECK(status IN ('opened', 'closed'))
);

CREATE TABLE IF NOT EXISTS orders (
    id BLOB(16) PRIMARY KEY,
    table_id BLOB(16) NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('taken', 'done', 'aborted')),
    FOREIGN KEY (table_id) REFERENCES tables(id)
);

CREATE TABLE IF NOT EXISTS preparations (
    id BLOB(16) PRIMARY KEY,
    order_id BLOB(16) NOT NULL,
    menu_item_id BLOB(16) NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('pending', 'in progress', 'ready', 'served', 'aborted')),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (menu_item_id) REFERENCES menu_items(id)
);

CREATE TABLE IF NOT EXISTS menu_items (
    id BLOB(16) PRIMARY KEY,
    name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK(price >= 0)
);

CREATE TABLE IF NOT EXISTS menu_categories (
    id BLOB(16) PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS menu_item_categories (
    item_id BLOB(16) NOT NULL,
    category_id BLOB(16) NOT NULL,
    PRIMARY KEY (item_id, category_id),
    FOREIGN KEY (item_id) REFERENCES menu_items(id),
    FOREIGN KEY (category_id) REFERENCES menu_categories(id)
);

CREATE TABLE IF NOT EXISTS bills (
    id BLOB(16) PRIMARY KEY,
    table_id BLOB(16) NOT NULL,
    total INTEGER NOT NULL CHECK(total >= 0),
    paid INTEGER NOT NULL CHECK(paid >= 0 AND paid <= total),
    status TEXT NOT NULL CHECK(status IN ('pending', 'partially paid', 'paid')),
    FOREIGN KEY (table_id) REFERENCES tables(id)
);

CREATE TABLE IF NOT EXISTS bill_menu_items (
    bill_id BLOB(16) NOT NULL,
    menu_item_id BLOB(16) NOT NULL,
    PRIMARY KEY (bill_id, menu_item_id),
    FOREIGN KEY (bill_id) REFERENCES bills(id),
    FOREIGN KEY (menu_item_id) REFERENCES menu_items(id)
);

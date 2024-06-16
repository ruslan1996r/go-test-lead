CREATE TABLE IF NOT EXISTS clients (
    id INTEGER PRIMARY KEY,
    name TEXT,
    start_date TEXT,
    end_date TEXT,
    priority TEXT CHECK(priority IN ('LOW', 'MEDIUM', 'HIGH')) DEFAULT 'MEDIUM',
    lead_capacity INTEGER
);

CREATE TABLE IF NOT EXISTS leads (
    lead_id TEXT NOT NULL PRIMARY KEY,
    client_id INTEGER NOT NULL,
    start_date TEXT,
    end_date TEXT,
    FOREIGN KEY (client_id) REFERENCES clients(id)
);

CREATE TABLE IF NOT EXISTS migrations (
    timestamp TEXT
)
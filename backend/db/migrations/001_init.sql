-- Table: users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    refresh_token TEXT
);

-- Table: databases
CREATE TABLE IF NOT EXISTS databases (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    host VARCHAR(100) NOT NULL,
    port INT NOT NULL,
    db_username VARCHAR(50) NOT NULL,
    db_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table: backups
CREATE TABLE IF NOT EXISTS backups (
    id SERIAL PRIMARY KEY,
    database_id INT NOT NULL REFERENCES databases(id),
    name VARCHAR(100),
    file_path TEXT NOT NULL,
    file_size FLOAT,
    backup_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    version VARCHAR(20),
    log TEXT
);

-- Table: restores
CREATE TABLE IF NOT EXISTS restores (
    id SERIAL PRIMARY KEY,
    backup_id INT NOT NULL REFERENCES backups(id),
    user_id INT NOT NULL REFERENCES users(id),
    restored_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    log TEXT 
);

-- Table: scheduled_tasks
CREATE TABLE IF NOT EXISTS scheduled_tasks (
    id SERIAL PRIMARY KEY,
    database_id INT NOT NULL REFERENCES databases(id),
    cron_expression VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    last_run_at TIMESTAMP
);

-- Table: alerts
CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    alert_type VARCHAR(50),
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE
);

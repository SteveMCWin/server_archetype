DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    date_created DATE NOT NULL,
    priviledged BOOLEAN DEFAULT FALSE,
    tests_started INTEGER DEFAULT 0,
    tests_completed INTEGER DEFAULT 0,
    all_time_avg_wpm REAL DEFAULT 0.0,
    all_time_avg_acc REAL DEFAULT 0.0
);


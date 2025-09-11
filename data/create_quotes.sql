DROP TABLE IF EXISTS quotes;

CREATE TABLE quotes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source TEXT DEFAULT 'Unknown source',
    quote TEXT NOT NULL,
    len TEXT NOT NULL CHECK(len IN ('s', 'm', 'l', 'xl'))
);


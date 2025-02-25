-- Создание таблицы песен
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group TEXT NOT NULL,
    song TEXT NOT NULL,
    release_date TEXT NOT NULL,
    text TEXT NOT NULL,
    link TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    deletedAt TIMESTAMP
);

-- Создание таблицы куплетов
CREATE TABLE verses (
    id SERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id) ON DELETE CASCADE,
    verse_number INT NOT NULL,
    verse_text TEXT NOT NULL
);
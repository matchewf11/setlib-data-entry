CREATE TABLE IF NOT EXISTS problems (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    section TEXT NOT NULL,
    difficulty TEXT NOT NULL,
    problem TEXT NOT NULL,
	topic TEXT NOT NULL,
	author TEXT NOT NULL,
    question_type TEXT NOT NULL
);

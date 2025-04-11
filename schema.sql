CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	password_hash BYTEA NOT NULL
);

CREATE TABLE sessions (
	user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
	token TEXT UNIQUE NOT NULL
);

CREATE TABLE posts (
	id SERIAL PRIMARY KEY,
	author_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
	created_at TIMESTAMPTZ DEFAULT now(),
	title TEXT NOT NULL,
	body TEXT NOT NULL,
	source TEXT NOT NULL
);

CREATE TABLE replies (
	id SERIAL PRIMARY KEY,
	post_id INTEGER REFERENCES posts(id) ON DELETE CASCADE,
	author_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
	created_at TIMESTAMPTZ DEFAULT now(),
	body TEXT NOT NULL,
	type TEXT CHECK (type IN (
		'opinion',
		'concur',
		'concurj',
		'dissent'
	)) NOT NULL
);

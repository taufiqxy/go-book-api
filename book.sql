CREATE TABLE mst_book (
	id SERIAL PRIMARY KEY,
	title VARCHAR(100),
	author VARCHAR(50),
	release_year VARCHAR(4),
	pages int
);

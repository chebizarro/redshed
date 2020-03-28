CREATE TABLE projects (
	id SERIAL PRIMARY KEY,
	userid VARCHAR(50) NOT NULL,
	title VARCHAR(200) NOT NULL,
	due timestamp without time zone default (now() at time zone 'utc'),
	completed BOOLEAN NOT NULL default(false),
	tags VARCHAR(200) NOT NULL default(''),
	notes VARCHAR(500) NOT NULL default('')
);

CREATE TABLE item (
  id SERIAL PRIMARY KEY,
  name VARCHAR(128) NOT NULL,
  description VARCHAR(512) NOT NULL,
  sku VARCHAR(256) NOT NULL,
  upc VARCHAR(12) NOT NULL
  price FLOAT NOT NULL,
  width FLOAT NOT NULL,
  height FLOAT NOT NULL,
  depth FLOAT NOT NULL,
  weight FLOAT,
  volume FLOAT,
);

CREATE TABLE inventory (
  id SERIAL PRIMARY KEY,
  item_id FOREIGN KEY,

);
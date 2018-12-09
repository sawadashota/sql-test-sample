-- +migrate Up
CREATE TABLE users
(
  id   SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  sex  VARCHAR(10) NOT NULL
);

-- +migrate Down
DROP TABLE users;

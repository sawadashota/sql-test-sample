-- +migrate Up
INSERT INTO
  users (name, sex)
VALUES
  ('Bob', 'male');

-- +migrate Down

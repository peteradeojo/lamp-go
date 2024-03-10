CREATE TABLE logs (
  id BIGSERIAL PRIMARY KEY,
  text TEXT NOT NULL,
  appToken TEXT NOT NULL REFERENCES apps(token),
  level VARCHAR(16) NOT NULL DEFAULT 'error',
  createdAt TIMESTAMP DEFAULT(NOW()),
  updatedAt TIMESTAMP DEFAULT(NOW()),
  context JSON,
  ip TEXT,
  tags JSON
);

CREATE TABLE apps (
  id BIGSERIAL PRIMARY KEY,
  token TEXT,
  userId BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
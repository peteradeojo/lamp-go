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

CREATE TYPE log_level as ENUM ('info', 'error', 'critical', 'debug');

CREATE TABLE system_logs (
  id uuid not null,
  text text not null,
  stack text,
  context json,
  origin text not null,
  level log_level not null,
  from_system bit,
  createdat timestamp not null,
  primary key (id)
);
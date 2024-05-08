CREATE TABLE users (
  id BIGSERIAL primary key,
  email text not null unique,
  githubid text,
  isadmin smallint default 0
);

CREATE TABLE accounts (
  id BIGSERIAL PRIMARY KEY,
  userid bigint not null,
  status int not null,
  tierid int not null,
  foreign key (userid) references users (id),
  foreign key (tierid) references tiers (id)
);

create table tiers (
  id BIGSERIAL PRIMARY KEY,
  name text not null,
  limits text not null,
  generallyavailable smallint default 0
);

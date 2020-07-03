create table schema_version (
  version text not null
);

create extension if not exists hstore;

create table users (
    id serial not null,
    username text not null unique,
    password text,
    language text default 'en_US',
    timezone text default 'UTC',
    last_login_at timestamp with time zone,
    group_id int not null,
    extra hstore,
    primary key (id)
);

create table groups (
    id serial not null,
    name text,
    extra hstore,
    primary key (id)
);

create table plans (
  id serial not null,
  name text not null,
  description text,
  node_id int not null,
  group_id int not null,
  attachments hstore,
  extra hstore,
  create_at timestamp with time zone,
  update_at timestamp with time zone,
  primary key (id)
);

create table plan_logs (
  id bigserial not null,
  log_id bigint not null,
  plan_id bigint not null,
  attachments hstore,
  extra hstore,
  create_at timestamp with time zone,
  update_at timestamp with time zone,

  primary key (plan_id, log_id)
);

create table blobs (
  id bigserial not null,
  filename text,
  content bytea,
  checksum text,
  create_at timestamp with time zone,
  update_at timestamp with time zone,
  primary key (id)
);


create table schema_version (
  version text not null
);

create extension if not exists hstore;

create table plans (
  id serial not null,
  name text not null,
  description text,
  node_id int not null,
  attachments hstore,
  create_at timestamp with time zone,
  update_at timestamp with time zone,
  primary key (id)
);

create table plan_logs (
  id bigserial not null,
  attachments hstore,
  create_at timestamp with time zone,
  update_at timestamp with time zone,

  primary key (id)
);

create table blobs (
  id bigserial not null,
  filename text,
  context text,
  checksum text,
  primary key (id)
);


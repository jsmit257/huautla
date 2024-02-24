\c huautla;

create table vendors (
  uuid varchar(40)  not null primary key,
  name varchar(512) not null unique
);

create table substrates (
  uuid        varchar(40)  not null primary key,
  name        varchar(512) not null,
  type        varchar(25)  not null check (type in ('Grain', 'Bulk')),
  vendor_uuid varchar(40)  not null references vendors(uuid),
  unique(name, vendor_uuid)
);

create table ingredients (
  uuid            varchar(40)  not null primary key,
  name            varchar(512) not null unique
);

create table substrate_ingredients (
  uuid             varchar(40) not null primary key,
  substrate_uuid   varchar(40) not null references substrates(uuid),
  ingredient_uuid  varchar(40) not null references ingredients(uuid),
  unique(substrate_uuid, ingredient_uuid)
);

create table strains (
  uuid        varchar(40)  not null primary key,
  name        varchar(512) not null,
  -- TODO: species     varchar(128) not null,
  vendor_uuid varchar(40)  not null references vendors(uuid),
  unique(name, vendor_uuid)
);

create table strain_attributes (
  uuid         varchar(40)  not null primary key,
  name         varchar(40)  not null,
  value        varchar(512) not null,
  strain_uuid  varchar(40)  not null references strains(uuid),
  unique(name, strain_uuid)
);

create table stages (
  uuid            varchar(40)  not null primary key,
  name            varchar(512) not null unique
);

create table event_types (
  uuid       varchar(40)  not null primary key,
  name       varchar(512) not null,
  severity   varchar(10)  not null check (severity in ('Info', 'Warn', 'Error', 'Fatal', 'RIP')),
  stage_uuid varchar(40)  not null references stages(uuid),
  unique(name, stage_uuid)
);

create table lifecycles (
  uuid                varchar(40)  not null primary key,
  name                varchar(128) not null unique,
  location            varchar(128) not null,
  grain_cost          decimal(8,2) not null,
  bulk_cost           decimal(8,2) not null,
  -- the net weight, fresh or dried; for dried, 1.0-(yield/gross) is how much water they typically contain
  yield               decimal(84,2) not null default 0,
  headcount           decimal(5) not null default 0,
  -- gross the fresh weight, regardless of whether they're sold fresh or dry (see yield)
  gross               decimal(5,2) not null default 0, 
  mtime               timestamp    not null default current_timestamp,
  ctime               timestamp    not null default current_timestamp,
  strain_uuid         varchar(40)  not null references strains(uuid),
  grainsubstrate_uuid varchar(40)  not null references substrates(uuid),
  bulksubstrate_uuid  varchar(40)  not null references substrates(uuid)
);

create table events (
  uuid           varchar(40)  not null primary key,
  temperature    numeric(4,1) not null default 0.0,
  humidity       int          not null default 0,
  mtime          timestamp    not null default current_timestamp,
  ctime          timestamp    not null default current_timestamp,
  lifecycle_uuid varchar(40)  not null references lifecycles(uuid),
  eventtype_uuid varchar(40)  not null references event_types(uuid)
);


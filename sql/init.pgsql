use database huautla;

create table vendors (
  uuid varchar(40)  not null primary key,
  name varchar(512) not null unique
);

create table substrates (
  uuid        varchar(40)  not null primary key,
  name        varchar(512) not null,
  vendor_uuid varchar(40)  not null references vendors(uuid),
  unique (name, vendor_uuid)
);

create table grain_substrates () inherits (substrates);

create table bulk_substrates () inherits (substrates);

create table substrate_ingredients (
  uuid            varchar(40)  not null primary key,
  name            varchar(512) not null unique
);

create table grain_ingredients (
  grain_substrate_uuid varchar(40) not null foreign key references grain_substrates(uuid)
) inherits (substrate_ingredients);

create table bulk_ingredients (
  bulk_substrate_uuid varchar(40) not null foreign key references bulk_substrates(uuid)
) inherits (substrate_ingredients);

create table strains (
  uuid        varchar(40)  not null primary key,
  name        varchar(512) not null,
  vendor_uuid varchar(40)  not null references vendors(uuid),
  unique (name, vendor_uuid)
);

create table strain_attributes (
  uuid         varchar(40)  not null primary key,
  name         varchar(40)  not null,
  value        varchar(512) not null,
  strain_uuid  varchar(40)  not null foreign key references strains(uuid)
);

create table stages (
  uuid            varchar(40)  not null primary key,
  name            varchar(512) not null unique
);

create table event_types (
  uuid       varchar(40)  not null primary key,
  name       varchar(512) not null unique,
  stage_uuid varchar(40)  not null foreign key references stages(uuid)
);

create table lifecycle (
  uuid                 varchar(40)  not null primary key,
  grain_cost           decimal(8,2) not null,
  bulk_cost            decimal(8,2) not null,
  yield                decimal(4,2) not null default 0,
  headcount            decimal(4,2) not null default 0,
  gross                decimal(5,2) not null default 0,
  mtime                datetime     not null default `now`,
  ctime                datetime     not null default `now`,
  strain_uuid          varchar(40)  not null foreign key references strains(uuid),
  grain_substrate_uuid varchar(40)  not null foreign key references grain_substrates(uuid),
  bulk_substrate_uuid  varchar(40)  not null foreign key references bulk_substrates(uuid)
);

create table events (
  uuid            varchar(40) not null primary key,
  temperature     int         not null default 0,
  mtime           datetime    not null default `now`,
  ctime           datetime    not null default `now`,
  lifecycle_uuid  varchar(40) not null foreign key references lifecycle(uuid),
  event_type_uuid varchar(40) not null foreign key references event_types(uuid)
);

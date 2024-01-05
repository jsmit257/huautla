#!/usr/bin/psql -f %

create table event_types (
  id uuid varchar(40) not null primary key
)

create table events (
  id uuid varchar(40) not null primary key
)

create table flush_event (
  id uuid varchar(40) not null primary key
)

create table harvest_event (
  id uuid varchar(40) not null primary key
)

-- this script is slightly inflated so we could test the syntax of the queries 
-- in internal/data/pgsql.yaml; things like ref integrity will be handled in 
-- system tests where it's easier to stage data

\c huautla

-- vendors
insert into vendors(uuid, name)
values('0', 'bass-o-matic'),
      ('1', 'juanita'),
      ('2', 'bogus');

update vendors set name = 'teez-head' where uuid = '0';
select * from vendors;
delete from vendors where uuid = '2';
select * from vendors;

-- substrates
insert into substrates(uuid, name, type, vendor_uuid)
values('0', '5-grain', 'Grain', '0'),
      ('1', 'rye', 'Grain', '1'),
      ('2', 'liquid', 'Bulk', '0'),
      ('3', 'dirt', 'Bulk', '1');

insert
  into substrates(uuid, name, type, vendor_uuid)
select '4', 'other', 'Grain', v.uuid
  from vendors v
 where v.uuid = '0';

update substrates set name = 'bulk' where uuid = '3';
select * from substrates;
delete from substrates where uuid = '3';
select * from substrates;

-- ingredients 
insert into ingredients(uuid, name)
values('0', 'Vermiculite'),
      ('1', 'Maltodextrin'),
      ('2', 'Rye'),
      ('3', 'Millet'),
      ('4', 'Popcorn'),
      ('5', 'Horse Cookies'),
      ('6', 'bogus');

update ingredients set name = 'Manure' where uuid = '5';
select * from ingredients;
delete from ingredients where uuid = '6';
select * from ingredients;

-- substrate_ingredients
insert into substrate_ingredients(uuid, substrate_uuid, ingredient_uuid)
values('0', '0', '2'),
      ('1', '0', '3'),
      ('2', '2', '1'),
      ('3', '1', '2'),
      ('4', '4', '0');

insert into substrate_ingredients (uuid, substrate_uuid, ingredient_uuid)
select '5', s.uuid, i.uuid
  from substrates s,
       ingredients i
 where s.uuid = '4'   -- '4' == other
   and i.uuid = '5';  -- '5' == Manure
select * from substrate_ingredients;

update substrate_ingredients
   set ingredient_uuid = '1'
 where substrate_uuid = '4'
   and ingredient_uuid = '5';
select * from substrate_ingredients;

delete
  from substrate_ingredients
 where substrate_uuid = '2'
   and ingredient_uuid = '1';

select i.uuid,
       i.name
  from ingredients i
  join substrate_ingredients si
    on i.uuid = si.ingredient_uuid
 where si.substrate_uuid = '4';

-- strains
insert into strains(uuid, name, vendor_uuid)
values('0', 'Huautla', '0'),
      ('1', 'Morel', '0'),
      ('2', 'Shitake', '1'),
      ('3', 'Liberty caps', '1');

insert
  into strains(uuid, name, vendor_uuid)
select '4', 'Hakuna Matata', v.uuid
  from vendors v
 where v.uuid = '0';
select * from strains;

update strains set name = 'Hens o'' the Wood' where uuid = '3' ; -- just kidding
select * from strains;

delete from strains where uuid = '2';
select * from strains;

-- strain_attributes
insert into strain_attributes(uuid, name, value, strain_uuid)
values('0', 'contamination resistance', 'high', '0'),
      ('1', 'potency', 'low', '0'),
      ('2', 'color', 'purple', '3');

-- unique names
select distinct name from strain_attributes order by name;

select uuid, name, value
  from strain_attributes sa
 where strain_uuid = '3';

insert
  into strain_attributes (uuid, name, value, strain_uuid)
select '3', 'umami', 'absolutely none', s.uuid
  from strains s
 where s.uuid = '3';
select * from strain_attributes;

update strain_attributes sa
   set value = 'less than 0'
  from strains s
 where sa.name = 'umami'
   and s.uuid = '3';
select * from strain_attributes;

delete from strain_attributes where uuid = '3';
select * from strain_attributes;

-- stages
insert into stages(uuid, name)
values('0', 'Gestation'),
      ('1', 'Colonization'),
      ('2', 'Propagation'),
      ('3', 'Delete me!');
select * from stages;

update stages set name = 'Culturing' where uuid = '2';
select * from stages;

delete from stages where uuid = '3';
select * from stages;

-- event_types 
insert into event_types(uuid, name, severity, stage_uuid)
values('0', 'Humidity', 'Info', '1'),
      ('1', 'Thermal', 'Warn', '1'),
      ('2', 'Crash', 'Error', '1');
select e.uuid,
       e.name,
       s.uuid as stage_uuid,
       s.name as stage_name
  from event_types e
  join stages s
    on e.stage_uuid = s.uuid;

insert
  into event_types(uuid, name, severity, stage_uuid)
select '3',
       'Oh Shit!',
       'Error',
       s.uuid
  from stages s
 where s.uuid = '1';
select * from event_types;

update event_types set name = 'Delete me!' where uuid = '3';
select * from event_types;

delete from event_types where uuid = '3';
select e.name,
       s.uuid as stage_uuid,
       s.name as stage_name
  from event_types e
  join stages s
    on e.stage_uuid = s.uuid
 where e.uuid = '1';
 select * from event_types;

-- lifecycles
insert into lifecycles(uuid, name, location, grain_cost, bulk_cost, yield, headcount, gross, mtime, ctime, strain_uuid, grainsubstrate_uuid, bulksubstrate_uuid)
values('0', '1', '', 0, 0, 0, 0, 0, current_timestamp, current_timestamp, '0', '0', '0'),
      ('1', '2', '', 0, 0, 0, 0, 0, current_timestamp, current_timestamp, '0', '0', '0');
select lc.name,
       lc.location,
       lc.grain_cost,
       lc.bulk_cost,
       lc.yield,
       lc.headcount,
       lc.gross,
       lc.mtime,
       lc.ctime,
       s.uuid as strain_uuid,
       s.name as strain_name,
       sv.uuid as strain_vendor_uuid,
       sv.name as strain_vendor_name,
       gs.uuid as grain_substrate_uuid,
       gs.name as grain_substrate_name,
       gs.type as grain_substrate_type,
       gv.uuid as grain_vendor_uuid,
       gv.name as grain_vendor_name,
       bs.uuid as bulk_substrate_uuid,
       bs.name as bulk_substrate_name,
       bs.type as bulk_substrate_type,
       bv.uuid as bulk_vendor_uuid,
       bv.name as bulk_vendor_name
  from lifecycles lc
  join strains s
    on lc.strain_uuid = s.uuid
  join vendors sv
    on s.vendor_uuid = sv.uuid
  join substrates gs
    on lc.grainsubstrate_uuid = gs.uuid
  join vendors gv
    on gs.vendor_uuid = gv.uuid
  join substrates bs
    on lc.bulksubstrate_uuid = bs.uuid
  join vendors bv
    on bs.vendor_uuid = bv.uuid
 where lc.uuid = '0';

insert
  into lifecycles(
       uuid,
       name,
       location,
       grain_cost,
       bulk_cost,
       yield,
       headcount,
       gross,
       mtime,
       ctime,
       strain_uuid,
       grainsubstrate_uuid,
       bulksubstrate_uuid)
select '2',
       'codename blue',
       'basement',
       100,
       300,
       -1,
       -20,
       0,
       current_timestamp,
       current_timestamp,
       s.uuid,
       gs.uuid,
       bs.uuid
  from strains s,
       substrates gs,
       substrates bs
 where s.uuid = '0'
   and gs.uuid = '1'
   and bs.uuid = '1';
select * from lifecycles;

update lifecycles
   set name = 'robert''s birthday brownie',
       location = 'a galaxy far, far away',
       grain_cost = 1,
       bulk_cost = 2,
       yield = 1000,
       headcount = 10000,
       gross = 12,
       mtime = current_timestamp,
       strain_uuid = s.uuid,
       grainsubstrate_uuid = gs.uuid,
       bulksubstrate_uuid = bs.uuid
  from strains s,
       substrates gs,
       substrates bs
 where s.uuid = '2'
   and gs.uuid = '1'
   and bs.uuid = '1';
select * from lifecycles;

delete from lifecycles where uuid = '2';
select * from lifecycles;

-- events
insert into events( uuid, temperature, humidity, mtime, ctime, lifecycle_uuid, eventtype_uuid)
values('0', 0, 1, current_timestamp, current_timestamp, '0', '1'),
      ('1', 0, 1, current_timestamp, current_timestamp, '1', '0');

-- all-by-lifecycle
select e.uuid,
       e.temperature,
       e.humidity,
       e.mtime,
       e.ctime,
       et.uuid as eventtype_uuid,
       et.name as eventtype_name,
       s.uuid as stage_uuid,
       s.name as stage_name
  from events e
  join event_types et
    on e.eventtype_uuid = et.uuid
  join stages s
    on et.stage_uuid = s.uuid
 where lifecycle_uuid = '1'
 order
    by mtime desc;
select * from events;

--  all-by-eventtype
select e.uuid,
       e.temperature,
       e.humidity,
       e.mtime,
       e.ctime,
       et.uuid as eventtype_uuid,
       et.name as eventtype_name,
       s.uuid as stage_uuid,
       s.name as stage_name
  from events e
  join event_types et
    on e.eventtype_uuid = et.uuid
  join stages s
    on et.stage_uuid = s.uuid
 where et.uuid = '1';
select * from events;

-- select
select e.temperature,
       e.humidity,
       e.mtime,
       e.ctime,
       et.uuid as eventtype_uuid,
       et.name as eventtype_name,
       s.uuid as stage_uuid,
       s.name as stage_name
  from events e
  join event_types et
    on e.eventtype_uuid = et.uuid
  join stages s
    on et.stage_uuid = s.uuid
 where e.uuid = '1';
select * from events;

-- add
insert
  into events(uuid, temperature, humidity, mtime, ctime, lifecycle_uuid, eventtype_uuid)
select '2', 99, 1, current_timestamp, current_timestamp, lc.uuid, et.uuid
  from lifecycles lc,
       event_types et
 where lc.uuid = 1
   and et.uuid = 1;
select * from events;

-- change
update events
   set temperature = 8,
       humidity = 12,
       mtime = current_timestamp,
       eventtype_uuid = et.uuid
  from event_types et
 where events.uuid = 1
   and et.uuid = '2'
select * from events;

\q
-- remove (duh)
delete from events where uuid = ?;
select * from events;

\q


-- run this after `seed.pqsql` to add system-data (e.g. fake data and foreign key relationships)

\c huautla

insert into vendors(uuid, name)
values('-1', 'delete me!');

insert into ingredients(uuid, name)
values('-1', 'delete me!');

insert into substrates(uuid, name, type, vendor_uuid)
values('0', 'Rye', 'Grain', '0'),
      ('1', 'Millet', 'Grain', '0'),
      ('-1', 'delete me!', 'Grain', '0');

insert into substrate_ingredients(uuid, substrate_uuid, ingredient_uuid)
values('0', '0', '2'),
      ('1', '1', '12'),
      ('2', '1', '3');

insert into strains(uuid, name, vendor_uuid)
values('0', 'Morel', '0'),
      ('1', 'Hens o'' the Wood', '0'),
      ('-1', 'delete me!', '0');

insert into strain_attributes(uuid, name, value, strain_uuid)
values('0', 'contamination resistance', 'high', '0'),
      ('1', 'potency', 'low', '0'),
      ('2', 'color', 'purple', '1');

-- event_types 
insert into event_types(uuid, name, severity, stage_uuid)
values('-1', 'delete me!', 'Info', '1');

-- -- lifecycles
-- insert into lifecycles(uuid, name, location, grain_cost, bulk_cost, yield, headcount, gross, mtime, ctime, strain_uuid, grainsubstrate_uuid, bulksubstrate_uuid)
-- values('0', '1', '', 0, 0, 0, 0, 0, current_timestamp, current_timestamp, '0', '0', '0'),
--       ('1', '2', '', 0, 0, 0, 0, 0, current_timestamp, current_timestamp, '0', '0', '0');

-- -- events
-- insert into events( uuid, temperature, humidity, mtime, ctime, lifecycle_uuid, eventtype_uuid)
-- values('0', 0, 1, current_timestamp, current_timestamp, '0', '1'),
--       ('1', 0, 1, current_timestamp, current_timestamp, '1', '0');
-- 

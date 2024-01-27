\c huautla

-- vendors
insert into vendors(uuid, name)
values('0', 'bass-o-matic'),
      ('1', 'juanita');

-- substrates
insert into substrates(uuid, name, type, vendor_uuid)
values('0', '5-grain', 'Grain', '0'),
      ('1', 'rye', 'Grain', '1'),
      ('2', 'liquid', 'Bulk', '0'),
      ('3', 'dirt', 'Bulk', '1');

-- ingredients 
insert into ingredients(uuid, name)
values('0', 'Vermiculite'),
      ('1', 'Maltodextrin'),
      ('2', 'Rye'),
      ('3', 'Millet'),
      ('4', 'Popcorn'),
      ('5', 'Manure');

-- substrate_ingredients
insert into substrate_ingredients(uuid, substrate_uuid, ingredient_uuid)
values('0', '0', '2'),
      ('1', '0', '3'),
      ('2', '2', '1'),
      ('3', '1', '2'),
      ('4', '3', '0'),
      ('5', '3', '5');

-- strains
insert into strains(uuid, name, vendor_uuid)
values('0', 'Huautla', '0'),
      ('1', 'Morel', '0'),
      ('2', 'Shitake', '1'),
      ('3', 'Liberty caps', '1');

-- strain_attributes
insert into strain_attributes(uuid, name, value, strain_uuid)
values('0', 'contamination resistance', 'high', '0'),
      ('1', 'potency', 'low', '0'),
      ('2', 'color', 'purple', '3');

-- stages
insert into stages(uuid, name)
values('0', 'Gestation'),
      ('1', 'Colonization'),
      ('2', 'Culturing');

-- eventtypes 
insert into event_types(uuid, name, severity, stage_uuid)
values('0', 'Humidity', 'Info', '1'),
      ('1', 'Thermal', 'Warn', '1'),
      ('2', 'Crash', 'Error', '1');


-- lifecycles
insert into lifecycles(uuid, name, location, grain_cost, bulk_cost, yield, headcount, gross, mtime, ctime, strain_uuid, grainsubstrate_uuid, bulksubstrate_uuid)
values('0', '1', '', 0, 0, 0, 0, 0, current_timestamp, current_timestamp, '0', '0', '0'),
      ('1', '2', '', 0, 0, 0, 0, 0, current_timestamp, current_timestamp, '0', '0', '0');

-- events
insert into events( uuid, temperature, humidity, mtime, ctime, lifecycle_uuid, eventtype_uuid)
values('0', 0, 1, current_timestamp, current_timestamp, '0', '0'),
      ('1', 0, 1, current_timestamp, current_timestamp, '0', '0');

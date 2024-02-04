-- this script is slightly inflated so we could test the syntax of the queries 
-- in internal/data/pgsql.yaml; things like ref integrity will be handled in 
-- system tests where it's easier to stage data

\c huautla

-- vendors
insert into vendors(uuid, name)
values('0', '127.0.0.1');

-- ingredients 
insert into ingredients(uuid, name)
values('0', 'Vermiculite'),
      ('1', 'Maltodextrin'),
      ('2', 'Rye'),
      ('3', 'White Millet'),
      ('4', 'Popcorn'),
      ('5', 'Manure'),
      ('6', 'Coir'),
      ('7', 'Honey'),
      ('8', 'Agar'),
      ('9', 'Rice Flour'),
      ('10', 'White Milo'),
      ('11', 'Red Milo'),
      ('12', 'Red Millet'),
      ('13', 'Gypsum'),
      ('14', 'Calcium phosphate'),
      ('15', 'Diammonium phosphate');

-- strains
insert into strains(uuid, name, vendor_uuid)
values('0', 'Huautla', '0'),
      ('1', 'Morel', '0'),
      ('2', 'Shitake', '0'),
      ('3', 'Liberty caps', '0');

-- strain_attributes
insert into strain_attributes(uuid, name, value, strain_uuid)
values('0', 'contamination resistance', 'high', '0'),
      ('1', 'potency', 'low', '0'),
      ('2', 'color', 'purple', '3');

-- stages
insert into stages(uuid, name)
values('0', 'Gestation'),
      ('1', 'Colonization'),
      ('2', 'Majority'),
      ('3', 'Vacation');

-- event_types 
insert into event_types(uuid, name, severity, stage_uuid)
values('0', 'Humidity', 'Info', '1'),
      ('1', 'Thermal', 'Warn', '1'),
      ('2', 'Crash', 'Error', '1');

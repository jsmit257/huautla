-- run this after `seed.pqsql` to add system-data (e.g. fake data and foreign key relationships)

\c huautla

insert into vendors(uuid, name)
values('updating substrate', 'updating substrate'),
      ('update me!', 'update me!'),
      ('delete me!', 'delete me!');

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
      ('15', 'Diammonium phosphate'),
      ('update me!', 'update me!'),
      ('delete me!', 'delete me!');

insert into stages(uuid, name)
values('update me!', 'update me!'),
      ('delete me!', 'delete me!');

insert into substrates(uuid, name, type, vendor_uuid)
values('0', 'Rye', 'Grain', '0'),
      ('1', 'Millet', 'Grain', '0'),
      ('2', 'Cedar chips', 'Bulk', '0'),
      ('add ingredient', 'add ingredient', 'Bulk', '0'),
      ('change ingredient', 'change ingredient', 'Grain', '0'),
      ('remove ingredient', 'remove ingredient', 'Bulk', '0'),
      ('update me!', 'update me!', 'Grain', '0'),
      ('delete me!', 'delete me!', 'Grain', '0');

insert into substrate_ingredients(uuid, substrate_uuid, ingredient_uuid)
values('0', '0', '2'),
      ('1', '1', '12'),
      ('2', '1', '3'),
      ('add ingredient', 'add ingredient', '2'),
      ('change ingredient', 'change ingredient', '3'),
      ('change ingredient 2', 'change ingredient', '12'),
      ('remove ingredient', 'remove ingredient', '12'),
      ('remove ingredient 2', 'remove ingredient', '13'),
      ('remove ingredient 3', 'remove ingredient', '14');

insert into strains(uuid, species, name, ctime, vendor_uuid)
values('0', 'M.esculenta', 'Morel', '1970-01-01', '0'),
      ('1', 'G.frondosa', 'Hens o'' the Wood', '1970-01-01', '0'),
      ('add attribute', '', 'add attribute', '1970-01-01', '0'),
      ('change attribute', '', 'change attribute', '1970-01-01', '0'),
      ('remove attribute', '', 'remove attribute', '1970-01-01', '0'),
      ('update me!', '', 'update me!', '1970-01-01', '0'),
      ('delete me!', '', 'delete me!', '1970-01-01', '0');

insert into strain_attributes(uuid, name, value, strain_uuid)
values('0', 'contamination resistance', 'high', '0'),
      ('1', 'headroom (cm)', '25', '0'),
      ('2', 'color', 'purple', '1'),
      ('add attribute', 'existing', 'existing', 'add attribute'),
      ('change attribute', 'color', 'albino', 'change attribute'),
      ('remove attribute 1', 'color', 'red', 'remove attribute'),
      ('remove attribute 2', 'energy', 'pure', 'remove attribute'),
      ('remove attribute 3', 'preferred substrate', 'cats', 'remove attribute');

insert into event_types(uuid, name, severity, stage_uuid)
values('update me!', 'update me!', 'Info', '1'),
      ('delete me!', 'delete me!', 'Info', '1');

insert into lifecycles(uuid, location, strain_cost, grain_cost, bulk_cost, yield, headcount, gross, mtime, ctime, strain_uuid, grainsubstrate_uuid, bulksubstrate_uuid)
values('0', 'reference implementation', 8, 1, 2, 3, 4, 5, '1970-01-01', '1970-01-01', '0', '0', '2'),
      ('1', 'reference implementation 2', 7, 0, 0, 0, 0, 0, '1970-01-01', '1970-01-01', '0', '0', '2'),
      ('add event', 'add event', 6, 0, 0, 0, 0, 0, '1970-01-01', '1970-01-01', '0', '0', '2'),
      ('change event', 'change event', 5, 0, 0, 0, 0, 0, '1970-01-01', '1970-01-01', '0', '0', '0'),
      ('remove event', 'remove event', 4, 0, 0, 0, 0, 0, '1970-01-01', '1970-01-01', '0', '0', '0'),
      ('update me!', 'update me!', 3, 0, 0, 0, 0, 0, '1970-01-01', '1970-01-01', '0', '0', '2'),
      ('delete me!', 'delete me!', 2, 0, 0, 0, 0, 0, '1970-01-01', '1970-01-01', '0', '0', '0');

insert into events(uuid, temperature, humidity, mtime, ctime, lifecycle_uuid, eventtype_uuid)
values('0', 2, 1, '1970-01-01T00:00:00.0Z', '1970-01-01T00:00:00.0Z', '0', '1'),
      ('1', 0, 1, '1970-01-01T00:00:00.0Z', '1970-01-01T00:00:00.0Z', '1', '0'),
      ('2', 0, 8, '1970-01-01T00:00:00.0Z', '1970-01-01T00:00:00.0Z', '0', '0'),
      ('add event', 0, 1, '1970-01-01T00:00:00.0Z', '1970-01-01T00:00:00.0Z', 'add event', '0'),
      ('change event', 0, 8, '1970-01-01T00:00:00.0Z', '1970-01-01T00:00:00.0Z', 'change event', '0'),
      ('remove event 1', 0, 8, '1970-01-01T00:00:00.0Z', '1970-01-01T00:00:00.0Z', 'remove event', '0'),
      ('remove event 2', 0, 8, '1970-01-01T00:00:00.0Z', '1970-01-01T00:00:00.0Z', 'remove event', '0'),
      ('remove event 3', 0, 8, '1970-01-01T00:00:00.0Z', '1970-01-01T00:00:00.0Z', 'remove event', '0'),
      ('update me!', 0, 8, '1970-01-01T00:00:00.0Z', '1970-01-01T00:00:00.0Z', 'update me!', '0');

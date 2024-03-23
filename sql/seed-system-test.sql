-- run this after `seed.pqsql` to add system-data (e.g. fake data and foreign key relationships)

\c huautla

insert into vendors(uuid, name)
values('updating substrate', 'updating substrate'),
      ('update me!', 'update me!'),
      ('delete me!', 'delete me!');

-- ingredients 
insert into ingredients(uuid, name)
values('update me!', 'update me!'),
      ('delete me!', 'delete me!');

insert into stages(uuid, name)
values('update me!', 'update me!'),
      ('delete me!', 'delete me!');

insert into substrates(uuid, name, type, vendor_uuid)
values('0', 'Rye', 'Grain', 'localhost'),
      ('1', 'Cedar chips', 'Bulk', 'localhost'),
      ('2', 'Agar', 'Agar', 'localhost'),
      ('3', 'Liquid', 'Liquid', 'localhost'),
      ('4', 'Millet', 'Grain', 'localhost'),
      ('update generation', 'Update generation', 'Liquid', 'localhost'),
      ('add ingredient', 'add ingredient', 'Bulk', 'localhost'),
      ('change ingredient', 'change ingredient', 'Grain', 'localhost'),
      ('remove ingredient', 'remove ingredient', 'Bulk', 'localhost'),
      ('update me!', 'update me!', 'Grain', 'localhost'),
      ('delete me!', 'delete me!', 'Grain', 'localhost');

insert into substrate_ingredients(uuid, substrate_uuid, ingredient_uuid)
values('0', '4', '2'),
      ('1', '0', '12'),
      ('2', '0', '3'),
      ('add ingredient', 'add ingredient', '2'),
      ('change ingredient', 'change ingredient', '3'),
      ('change ingredient 2', 'change ingredient', '12'),
      ('remove ingredient', 'remove ingredient', '12'),
      ('remove ingredient 2', 'remove ingredient', '13'),
      ('remove ingredient 3', 'remove ingredient', '14');

insert into strains(uuid, species, name, vendor_uuid)
values('0', 'M.esculenta', 'Morel', 'localhost'),
      ('1', 'G.frondosa', 'Hens o'' the Wood', 'localhost'),
      ('spore generation', 'X.test', 'spore generation', 'localhost'),
      ('spore generation 2', 'X.test', 'spore generation 2', 'localhost'),
      ('spore generation 3' , 'X.test', 'spore generation 3', 'localhost'),
      ('clone generation', 'X.test', 'clone generation', 'localhost'),
      ('add strain source', 'X.test', 'add strain source', 'localhost'),
      ('change strain source 1', 'X.test', 'change strain source 1', 'localhost'),
      ('change strain source 0', 'X.test', 'change strain source 0', 'localhost'),
      ('remove strain source', 'X.test', 'remove strain source', 'localhost'),
      ('add attribute', '', 'add attribute', 'localhost'),
      ('change attribute', '', 'change attribute', 'localhost'),
      ('remove attribute', '', 'remove attribute', 'localhost'),
      ('update me!', '', 'update me!', 'localhost'),
      ('delete me!', '', 'delete me!', 'localhost');

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

insert into lifecycles(uuid, location, strain_cost, grain_cost, bulk_cost, yield, headcount, gross, strain_uuid, grainsubstrate_uuid, bulksubstrate_uuid)
values('0', 'reference implementation', 8, 1, 2, 3, 4, 5, '0', '0', '2'),
      ('1', 'reference implementation 2', 7, 0, 0, 0, 0, 0, '0', '0', '2'),
      ('spore', 'spore', 1.1, 0, 0, 0, 0, 0, '0', '0', '2'),
      ('spore 2', 'spore 2', 1.2, 0, 0, 0, 0, 0, '0', '0', '2'),
      ('clone', 'clone', 2.2, 0, 0, 0, 0, 0, '0', '0', '2'),
      ('add event', 'add event', 6, 0, 0, 0, 0, 0, '0', '0', '2'),
      ('change event', 'change event', 5, 0, 0, 0, 0, 0, '0', '0', '0'),
      ('remove event', 'remove event', 4, 0, 0, 0, 0, 0, '0', '0', '0'),
      ('add event source lc', 'add event source', 4, 0, 0, 0, 0, 0, '0', '0', '0'),
      ('update me!', 'update me!', 3, 0, 0, 0, 0, 0, '0', '0', '1'),
      ('delete me!', 'delete me!', 2, 0, 0, 0, 0, 0, '0', '0', '0');

insert into generations(uuid, platingsubstrate_uuid, liquidsubstrate_uuid)
values('0', '2', '3'),
      ('1', '2', '3'),
      ('2', '2', '3'),
      ('3', '2', '3'),
      ('4', '2', '3'),
      ('change_source_fail_type', '2', '3'),
      ('add event source', '2', '3'),
      ('add source', '2', '3'),
      ('fail source check', '2', '3'),
      ('change source', '2', '3'),
      ('remove source', '2', '3'),
      ('add gen event', '2', '3'),
      ('has events', '2', '3'),
      ('change event', '2', '3'),
      ('remove gen event', '2', '3'),
      ('get gen event', '2', '3'),
      ('update me!', '2', '3'),
      ('delete me!', '2', '3')
;

insert into events(uuid, temperature, humidity, observable_uuid, eventtype_uuid)
values('0', 2, 1, '0', '1'),
      ('1', 0, 1, '1', '0'),
      ('2', 0, 8, '0', '0'),
      ('spore print', 20, 21, 'spore', '26'),
      ('spore print 2', 10, 11, 'spore', '26'),
      ('spore print 3', 10, 11, 'spore', '26'),
      ('clone', 10, 11, '0', '27'),
      ('add event', 0, 1, 'add event', '0'),
      ('change event', 0, 8, 'change event', '0'),
      ('remove event 1', 0, 8, 'remove event', '0'),
      ('remove event 2', 0, 8, 'remove event', '0'),
      ('remove event 3', 0, 8, 'remove event', '0'),
      ('get gen event 0', 0, 8, 'get gen event', '0'),
      ('get gen event 1', 0, 8, 'get gen event', '0'),
      ('get gen event 2', 0, 8, 'get gen event', '0'),
      ('add gen event 0', 0, 8, 'add gen event', '0'),
      ('remove gen event 0', 0, 8, 'remove gen event', '0'),
      ('remove gen event 1', 0, 8, 'remove gen event', '0'),
      ('remove gen event 2', 0, 8, 'remove gen event', '0'),
      ('add spore event source 0', 0, 8, 'add event source lc', '26'),
      ('add spore event source 1', 0, 8, 'add event source lc', '26'),
      ('add spore event source 2', 0, 8, 'add event source lc', '26'),
      ('add clone event source 0', 0, 8, 'add event source lc', '27'),
      ('add clone event source 1', 0, 8, 'add event source lc', '27'),
      ('update me!', 0, 8, 'update me!', '0');

insert into sources(uuid, type, progenitor_uuid, generation_uuid)
values('0', 'Spore', 'spore print', '0'),
      ('1', 'Spore', 'spore print 2', '0'),
      ('2', 'Spore', 'spore generation', '1'),
      ('3', 'Clone', 'clone generation', '2'),
      ('4', 'Clone', 'clone', '3'),
      ('5', 'Spore', 'spore generation 2', '4'),
      ('6', 'Spore', 'spore generation 3', '4'),
      ('add source', 'Spore', 'add strain source', 'add source'),
      ('change source 0', 'Spore', 'change strain source 0', 'change source'),
      ('change source 1', 'Spore', 'change strain source 1', 'change source'),
      ('change_source_fail_type 0', 'Spore', 'change strain source 1', 'change_source_fail_type'),
      ('delete me!', 'Spore', 'remove strain source', 'remove source')
;

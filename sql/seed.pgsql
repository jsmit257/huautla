-- seed the db with a few standard values that probably anyone will need; 
-- if not, you can just delete them later; used by system-test so be 
-- careful what you change

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

-- stages
insert into stages(uuid, name)
values('0', 'Gestation'),
      ('1', 'Colonization'),
      ('2', 'Majority'),
      ('3', 'Vacation');

-- event_types 
insert into event_types(uuid, name, severity, stage_uuid)
values('0', 'Condensation', 'Warn', '3'),
      ('1', 'Fruiting', 'Info', '1'),
      ('2', 'Crashed', 'Error', '1'),
      ('3', 'Sunset', 'RIP', '2'),
      ('4', 'Spore printing', 'Info', '0'),
      ('5', 'Innoculating agar substrate', 'Info', '0'),
      ('6', 'Innoculating liquid substrate', 'Info', '0');

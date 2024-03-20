-- seed the db with a few standard values that probably anyone will need; 
-- if not, you can just delete them later; used by system-test so be 
-- careful what you change

\c huautla

-- vendors
insert into vendors(uuid, name, website)
values('0', '127.0.0.1', 'https://localhost:8080/');

-- stages
insert into stages(uuid, name)
values('0', 'Gestation'),
      ('1', 'Colonization'),
      ('2', 'Majority'),
      ('3', 'Vacation');

-- event_types 
insert into event_types(uuid, name, severity, stage_uuid)
values('0', 'Agar sampling', 'Info', '0'),
      ('1', '1/2 colonization', 'Info', '0'),
      ('2', '100% colonization', 'Info', '0'),
      ('3', 'Agar mold', 'Error', '0'),
      ('4', 'Agar bacteria', 'Error', '0'),
      ('5', 'Liquid innoculation', 'Info', '0'),
      ('6', 'Liquid colonization', 'Info', '0'),
      ('7', 'Liquid mold', 'Error', '0'),
      ('8', 'Liquid bacteria', 'Error', '0'),
      ('9', 'Innoculation', 'Info', '1'),
      ('10', '1/3 colonized', 'Info', '1'),
      ('11', '100% colonized', 'Info', '1'),
      ('12', 'Redistribute substrate', 'Info', '1'),
      ('13', 'Binning', 'Info', '2'),
      ('14', '50% colonization', 'Info', '2'),
      ('15', 'Pinning', 'Info', '2'),
      ('16', 'Fruiting', 'Info', '2'),
      ('17', 'Resting', 'Info', '2'),
      ('18', 'Mold', 'Fatal', '2'),
      ('19', 'Sunset', 'RIP', '2'),
      ('20', 'Freeze', 'Info', '3'),
      ('21', 'Bacteria', 'Info', '3'),
      ('22', 'Mold', 'Info', '3'),
      ('23', 'Thaw', 'Info', '3');

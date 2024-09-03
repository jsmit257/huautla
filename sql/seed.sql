-- seed the db with a few standard values that probably anyone will need; 
-- if not, you can just delete them later; used by system-test so be 
-- careful what you change

\c huautla

insert into vendors(uuid, name, website)
values('localhost', '127.0.0.1', 'https://localhost:8080/');

insert into stages(uuid, name)
values('0', 'Gestation'),
      ('1', 'Colonization'),
      ('2', 'Majority'),
      ('3', 'Vacation'),
      ('4', 'Any');

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

insert into substrates(uuid, name, type, vendor_uuid)
values('no-op', 'N/A', 'plating', 'localhost');

insert into event_types(uuid, name, severity, stage_uuid)
values('0', 'Agar sampling', 'Begin', '0'),
      ('10', '33% colonization', 'Info', '4'),
      ('1', '50% colonization', 'Info', '4'),
      ('2', '100% colonization', 'Info', '4'),
      ('3', 'Mold', 'Error', '4'),
      ('yeast', 'Yeast', 'Error', '0'),
      ('4', 'Agar bacteria', 'Error', '0'), -- this may be recoverable
      ('5', 'Liquid innoculation', 'Begin', '0'),
      ('9', 'Innoculation', 'Begin', '1'),
      ('12', 'Redistribute substrate', 'Info', '1'),
      ('13', 'Binning', 'Begin', '2'),
      ('15', 'Pinning', 'Info', '2'),
      ('16', 'Fruiting', 'Info', '2'),
      ('17', 'Harvesting', 'Info', '2'),
      ('18', 'Resting', 'Info', '2'),
      ('20', 'Sunset', 'RIP', '2'),
      ('21', 'Chill', 'Begin', '3'),
      ('22', 'Freeze', 'Error', '3'),
      ('23', 'Bacteria', 'Fatal', '4'),
      ('24', 'Mold', 'Fatal', '3'),
      ('25', 'Thaw', 'Info', '3'),
      ('26', 'Spore print', 'Generation', '2'),
      ('27', 'Clone', 'Generation', '4'),
      ('28', 'Photo', 'Info', '4');

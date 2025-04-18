package data

type sqlMap map[string]map[string]string

var psqls = sqlMap{

	"event": {
		"all-by-observable": `
      select e.uuid,
             e.temperature,
             e.humidity,
             e.mtime at time zone 'utc',
             e.ctime at time zone 'utc',
             et.uuid as eventtype_uuid,
             et.name as eventtype_name,
             et.severity as eventtype_severity,
             s.uuid as stage_uuid,
             s.name as stage_name
       from  events e
       join  event_types et
         on  e.eventtype_uuid = et.uuid
       join  stages s
         on  et.stage_uuid = s.uuid
      where  e.observable_uuid = $1
      order
         by  e.mtime desc`,
		"all-by-eventtype": `
      select  e.uuid,
              e.temperature,
              e.humidity,
              e.mtime,
              e.ctime,
              et.uuid as eventtype_uuid,
              et.name as eventtype_name,
              et.severity as eventtype_severity,
              s.uuid as stage_uuid,
              s.name as stage_name
        from  events e
        join  event_types et
          on  e.eventtype_uuid = et.uuid
        join  stages s
          on  et.stage_uuid = s.uuid
       where  et.uuid = $1`,
		"notes-and-photos": `
      select  e.uuid as event_uuid,
              n.uuid as note_uuid,
              n.note as note_note,
              n.mtime as note_mtime,
              n.ctime as note_ctime,
              p.uuid as photo_uuid,
              p.filename as photo_filename,
              p.mtime as photo_mtime,
              p.ctime as photo_ctime,
              pn.uuid as photonote_uuid,
              pn.note as photonote_note,
              pn.mtime as photonote_mtime,
              pn.ctime as photonote_ctime
        from  events e
        left
        join  notes n 
          on  e.uuid = n.notable_uuid
        left
        join  photos p
          on  e.uuid = p.photoable_uuid
        left
        join  notes pn
          on  p.uuid = pn.notable_uuid
       where  e.observable_uuid = $1
         and  coalesce(n.uuid, p.uuid) is not null
       order
          by  e.uuid, n.mtime, p.mtime, pn.mtime`,
		"select": `
    select e.uuid,
           e.temperature,
           e.humidity,
           e.mtime,
           e.ctime,
           et.uuid as eventtype_uuid,
           et.name as eventtype_name,
           et.severity as eventtype_severity,
           s.uuid as stage_uuid,
           s.name as stage_name
      from events e
      join event_types et
        on e.eventtype_uuid = et.uuid
      join stages s
        on et.stage_uuid = s.uuid
     where e.uuid = $1`,
		"add": `
      insert
        into  events(uuid, temperature, humidity, mtime, ctime, observable_uuid, eventtype_uuid)
      select  $1, $2, $3, $4, $5, $6, et.uuid
        from  event_types et
      where  et.uuid = $7`,
		"change": `
      update  events e
         set  temperature = $1,
              humidity = $2,
              mtime = $3,
              eventtype_uuid = et.uuid
        from  event_types et
       where  e.uuid = $4
         and  et.uuid = $5`,
		"remove": `delete from events where uuid = $1`,
	},

	"eventtype": {
		"select-all": `
      select  e.uuid,
              e.name,
              e.severity,
              s.uuid as stage_uuid,
              s.name as stage_name
        from  event_types e
        join  stages s
          on  e.stage_uuid = s.uuid
      order
          by  s.name, e.name`,
		"select": `
    select e.uuid,
           e.name,
           e.severity,
           s.uuid as stage_uuid,
           s.name as stage_name
      from event_types e
      join stages s
        on e.stage_uuid = s.uuid
     where e.uuid = $1`,
		"insert": `
    insert
      into event_types(uuid, name, severity, stage_uuid)
    select $1,
           $2,
           $3,           
           s.uuid
      from stages s
     where s.uuid = $4`,
		"update": `
      update  event_types 
         set  name = $1,
              severity = $2,
              stage_uuid = $3,
              mtime = current_timestamp
       where  uuid = $4`,
		"delete": `delete from event_types where uuid = $1`,
	},

	"generation": {
		"ndx": `
      select  g.uuid,
              ps.uuid as plating_id,
              ps.name as plating_name,
              ps.type as plating_type,
              psv.uuid as plating_vendor_uuid,
              psv.name as plating_vendor_name,
              psv.website as plating_vendor_website,
              ls.uuid as liquid_uuid,
              ls.name as liquid_name,
              ls.type as liquid_type,
              lsv.uuid as liquid_vendor_uuid,
              lsv.name as liquid_vendor_name,
              lsv.website as liquid_vendor_website,
              s.uuid as source_uuid,
              s.type,
              lc.uuid as observable_id,
              coalesce(lc.strain_uuid, s.progenitor_uuid) as strain_uuid,
              st.name as strain_name,
              st.species as strain_species,
              st.ctime as strain_ctime,
              stv.uuid as strain_vendor_uuid,
              stv.name as strain_vendor_name,
              stv.website as strain_vendor_website,
              g.mtime as generation_mtime,
              g.ctime as generation_ctime,
              g.dtime as generation_dtime
        from  generations g
        join  substrates ps
          on  g.platingsubstrate_uuid = ps.uuid
        join  vendors psv
          on  ps.vendor_uuid = psv.uuid
        join  substrates ls
          on  g.liquidsubstrate_uuid = ls.uuid
        join  vendors lsv
          on  ls.vendor_uuid = lsv.uuid
        left
        join  sources s
          on  g.uuid = s.generation_uuid
        left
        join  events e
          on  s.progenitor_uuid = e.uuid
        left
        join  lifecycles lc
          on  e.observable_uuid = lc.uuid
        left
        join  strains st
          on  st.uuid = coalesce(lc.strain_uuid, s.progenitor_uuid)
        left
        join  vendors stv
          on  st.vendor_uuid = stv.uuid
       order
          by  g.mtime`,
		// just goes to show you can solve every problem with a union
		"select": `
        with  strain_sources as (
      select  so.generation_uuid,
              st.uuid as strain_uuid
        from  sources so
        join  strains st
          on  so.progenitor_uuid = st.uuid
        union
      select  so.generation_uuid,
              st.uuid
        from  sources so
        join  events ev
          on  so.progenitor_uuid = ev.uuid
        join  lifecycles lc
          on  ev.observable_uuid = lc.uuid
        join  strains st
          on  lc.strain_uuid = st.uuid
      ),      events_filter as (
      select  g.uuid as generation_uuid
        from  generations g
        join  events e
          on  g.uuid = e.observable_uuid
       where  e.eventtype_uuid = coalesce($5, e.eventtype_uuid)
       union
      select  g.uuid 
        from  generations g
        left
        join  events e
          on  g.uuid = e.observable_uuid
       where  e.uuid is null
         and  $5 is null
      )
      select  distinct
              g.uuid,
              ps.uuid as plating_id,
              ps.name as plating_name,
              ps.type as plating_type,
              psv.uuid as plating_vendor_uuid,
              psv.name as plating_vendor_name,
              psv.website as plating_vendor_website,
              ls.uuid as liquid_uuid,
              ls.name as liquid_name,
              ls.type as liquid_type,
              lsv.uuid as liquid_vendor_uuid,
              lsv.name as liquid_vendor_name,
              lsv.website as liquid_vendor_website,
              g.mtime,
              g.ctime,
              g.dtime
        from  generations g
        join  events_filter ef
          on  g.uuid = ef.generation_uuid
        join  substrates ps
          on  g.platingsubstrate_uuid = ps.uuid
        join  vendors psv
          on  ps.vendor_uuid = psv.uuid
        join  substrates ls
          on  g.liquidsubstrate_uuid = ls.uuid
        join  vendors lsv
          on  ls.vendor_uuid = lsv.uuid
        join  strain_sources ss
          on  g.uuid = ss.generation_uuid
       where  g.uuid = coalesce($1, g.uuid)
         and  ss.strain_uuid = coalesce($2, ss.strain_uuid)
         and  ps.uuid = coalesce($3, ps.uuid)
         and  ls.uuid = coalesce($4, ls.uuid)
       union
      select  distinct
              g.uuid,
              ps.uuid as plating_id,
              ps.name as plating_name,
              ps.type as plating_type,
              psv.uuid as plating_vendor_uuid,
              psv.name as plating_vendor_name,
              psv.website as plating_vendor_website,
              ls.uuid as liquid_uuid,
              ls.name as liquid_name,
              ls.type as liquid_type,
              lsv.uuid as liquid_vendor_uuid,
              lsv.name as liquid_vendor_name,
              lsv.website as liquid_vendor_website,
              g.mtime,
              g.ctime,
              g.dtime
        from  generations g
        join  events_filter ef
          on  g.uuid = ef.generation_uuid
        join  substrates ps
          on  g.platingsubstrate_uuid = ps.uuid
        join  vendors psv
          on  ps.vendor_uuid = psv.uuid
        join  substrates ls
          on  g.liquidsubstrate_uuid = ls.uuid
        join  vendors lsv
          on  ls.vendor_uuid = lsv.uuid
        left
        join  strain_sources ss
          on  g.uuid = ss.generation_uuid
       where  ss.generation_uuid is null
         and  g.uuid = coalesce($1, g.uuid)
         and  $2 is null
         and  ps.uuid = coalesce($3, ps.uuid)
         and  ls.uuid = coalesce($4, ls.uuid)`,
		"insert": `
      insert  into generations(uuid, platingsubstrate_uuid, liquidsubstrate_uuid, mtime, ctime)
      select  $1,
              ps.uuid,
              ls.uuid,
              $4,
              $4
        from  substrates ps,
              substrates ls
       where  ps.type = 'plating'
         and  ls.type = 'liquid'
         and  ps.uuid = $2
         and  ls.uuid = $3`,
		"update": `
      update  generations g
         set  platingsubstrate_uuid = ps.uuid,
              liquidsubstrate_uuid = ls.uuid,
              mtime = $4
        from  substrates ps,
              substrates ls
       where  ps.type = 'plating'
         and  ls.type = 'liquid'
         and  ps.uuid = $1
         and  ls.uuid = $2
         and  g.uuid = $3`,
		"delete": "update generations set dtime = current_timestamp where uuid = $1",
	},

	"ingredient": {
		"select-all": `select uuid, name from ingredients order by name`,
		"select":     `select name from ingredients where uuid = $1`,
		"insert":     `insert into ingredients(uuid, name) values($1, $2)`,
		"update":     `update ingredients set name = $1 where uuid = $2`,
		"delete":     `delete from ingredients where uuid = $1`,
	},

	"lifecycle": {
		// it's an ugly, bad precedent, except that it saves a lot of hits to the db
		"index": `
      select distinct 
             l.uuid,
             l.location,
             l.mtime,
             l.ctime,
             s.uuid as strain_uuid,
             s.species as strain_species,
             s.name as strain_name,
             s.ctime as strain_ctime,
             v.uuid as strain_vendor_uuid,
             v.name as strain_vendor_name,
             v.website as strain_vendor_website,
             e.uuid as event_uuid,
             e.temperature,
             e.humidity,
             e.mtime at time zone 'utc',
             e.ctime at time zone 'utc',
             et.uuid as eventtype_uuid,
             et.name as eventtype_name,
             et.severity as eventtype_severity,
             st.uuid as stage_uuid,
             st.name as stage_name
       from  lifecycles l
       join  strains s
         on  l.strain_uuid = s.uuid
       join  vendors v 
         on  s.vendor_uuid = v.uuid
       left
       join  events e 
         on  e.observable_uuid = l.uuid
        and  e.eventtype_uuid in ('sunset', 'sporeprint', 'clone')
       left
       join  event_types et
         on  e.eventtype_uuid = et.uuid
       left
       join  stages st
         on  et.stage_uuid = st.uuid
      order
         by  l.mtime desc, l.uuid`,
		"select": `
        with  event_filter as (
      select  lc.uuid as lifecycle_uuid
        from  lifecycles lc
        join  events e
          on  lc.uuid = e.observable_uuid
       where  e.eventtype_uuid = coalesce($5, e.eventtype_uuid)
       union  
      select  lc.uuid
        from  lifecycles lc
        join  events e
          on  lc.uuid = e.observable_uuid
       where  e.uuid is null
         and  $5 is null
      )
      select  lc.uuid,
              lc.location,
              lc.strain_cost,
              lc.grain_cost,
              lc.bulk_cost,
              lc.yield,
              lc.headcount,
              lc.gross,
              lc.mtime at time zone 'utc',
              lc.ctime at time zone 'utc',
              s.uuid as strain_uuid,
              s.species as strain_species,
              s.name as strain_name,
              s.generation_uuid,
              s.ctime as strain_ctime,
              s.dtime as strain_dtime,
              sv.uuid as strain_vendor_uuid,
              sv.name as strain_vendor_name,
              sv.website as strain_vendor_website,
              gs.uuid as grain_substrate_uuid,
              gs.name as grain_substrate_name,
              gs.type as grain_substrate_type,
              gv.uuid as grain_vendor_uuid,
              gv.name as grain_vendor_name,
              gv.website as grain_vendor_website,
              bs.uuid as bulk_substrate_uuid,
              bs.name as bulk_substrate_name,
              bs.type as bulk_substrate_type,
              bv.uuid as bulk_vendor_uuid,
              bv.name as bulk_vendor_name,
              bv.website as bulk_vendor_website
        from  lifecycles lc
        left
        join  event_filter ef
          on  lc.uuid = ef.lifecycle_uuid
        join  strains s
          on  lc.strain_uuid = s.uuid
        join  vendors sv
          on  s.vendor_uuid = sv.uuid
        join  substrates gs
          on  lc.grainsubstrate_uuid = gs.uuid
        join  vendors gv
          on  gs.vendor_uuid = gv.uuid
        join  substrates bs
          on  lc.bulksubstrate_uuid = bs.uuid 
        join  vendors bv
          on  bs.vendor_uuid = bv.uuid
       where  lc.uuid = coalesce($1, lc.uuid)
         and  s.uuid = coalesce($2, s.uuid)
         and  gs.uuid = coalesce($3, gs.uuid)
         and  bs.uuid = coalesce($4, bs.uuid)`,
		"insert": `
      insert
        into lifecycles(
             uuid,
             location,
             strain_cost,
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
      select $1,
             $2,
             $3,
             $4,
             $5,
             $6,
             $7,
             $8,
             $9,
             $10,
             s.uuid,
             gs.uuid,
             bs.uuid
       from  strains s,
             substrates gs,
             substrates bs
      where  s.uuid = $11
        and  gs.uuid = $12
        and  gs.type = 'grain'
        and  bs.uuid = $13
        and  bs.type = 'bulk'`,
		"update": `
      update lifecycles
        set location = $1,
            strain_cost = $2,
            grain_cost = $3,
            bulk_cost = $4,
            yield = $5,
            headcount = $6,
            gross = $7,
            mtime = $8,
            strain_uuid = s.uuid,
            grainsubstrate_uuid = gs.uuid,
            bulksubstrate_uuid = bs.uuid
        from strains s,
            substrates gs,
            substrates bs
      where s.uuid = $9
        and gs.uuid = $10
        and gs.type = 'grain'
        and bs.uuid = $11
        and bs.type = 'bulk'
        and lifecycles.uuid = $12`,
		"delete": `delete from lifecycles where uuid = $1`,
	},

	"note": {
		"get": `
      select  uuid,
              note,
              mtime,
              ctime
        from  notes
       where  notable_uuid = $1
       order
          by  mtime desc`,
		"add": `
      insert into notes(uuid, note, notable_uuid, mtime, ctime)
      values ($1, $2, $3, $4, $4)`,
		"change": `
      update  notes
         set  note = $1,
              mtime = $2
       where  uuid = $3`,
		"remove": `delete from notes where uuid = $1`,
	},

	"photo": {
		"get": `
      select  p.uuid,
              p.filename,
              p.mtime,
              p.ctime,
              n.uuid as note_uuid,
              n.note,
              n.mtime as note_mtime,
              n.ctime as note_ctime
        from  photos p
        left
        join  notes n
          on  n.notable_uuid = p.uuid
       where  p.photoable_uuid = $1
       order
          by  p.mtime desc, p.uuid, n.mtime desc`,
		"add": `
      insert into photos(uuid, filename, photoable_uuid, mtime, ctime)
      values ($1, $2, $3, $4, $5)`,
		"change": `
      update  photos
         set  filename = $1,
              mtime = $2
       where  uuid = $3`,
		"remove": `delete from photos where uuid = $1`,
	},

	"source": {
		"get": `
      select  s.uuid,
              s.type, 
              s.progenitor_uuid,
              lc.uuid as lifecycle_uuid,
              st.uuid as strain_uuid,
              st.name as strain_name,
              st.species,
              st.ctime as strain_ctime,
              st.dtime as strain_dtime,
              v.uuid as strain_vendor_uuid,
              v.name as strain_vendor_name,
              v.website as strain_vendor_website
        from  sources s
        left
        join  events e
          on  s.progenitor_uuid = e.uuid
        left
        join  lifecycles lc
          on  e.observable_uuid = lc.uuid
        join  strains st
          on  st.uuid = coalesce(lc.strain_uuid, s.progenitor_uuid)
        join  vendors v
          on  st.vendor_uuid = v.uuid
       where  s.generation_uuid = $1`,
		"add": `
      insert
        into  sources(uuid, type, progenitor_uuid, generation_uuid)
      values  ($1, $2, $3, $4)`,
		"change": `
      update  sources s
         set  type = $1,
              mtime = current_timestamp
       where  s.uuid = $2`,
		"delete": `delete from sources where uuid = $1`,
		"strain-from-event": `
      select  lc.strain_uuid
        from  lifecycles lc
        join  events e
          on  lc.uuid = e.observable_uuid
       where  e.uuid = $1`,
	},

	"stage": {
		"select-all": `select uuid, name from stages order by name`,
		"select":     `select name from stages where uuid = $1`,
		"insert":     `insert into stages(uuid, name) values($1, $2)`,
		"update":     `update stages set name = $1 where uuid = $2`,
		"delete":     `delete from stages where uuid = $1`,
	},

	"strain": {
		"select-all": `
      select  s.uuid,
              s.species,
              s.name,
              s.ctime,
              s.dtime,
              v.uuid as vendor_uuid,
              v.name as vendor_name,
              v.website as vendor_website,
              s.generation_uuid
        from  strains s
        join  vendors v
          on  s.vendor_uuid = v.uuid
       order
          by  s.name`,
		"select": `
      select  s.uuid,
              s.species,
              s.name,
              s.ctime,
              s.dtime,
              v.uuid as vendor_uuid,
              v.name as vendor_name,
              v.website as vendor_website,
              s.generation_uuid
        from  strains s
        join  vendors v
          on  s.vendor_uuid = v.uuid
       where  s.uuid = coalesce($1, s.uuid)
         and  v.uuid = coalesce($2, v.uuid)`,
		"insert": `
      insert
        into  strains(uuid, species, name, ctime, vendor_uuid)
      select  $1, $2, $3, $4, v.uuid
        from  vendors v
       where  v.uuid = $5`,
		"update": `
      update  strains 
         set  species = $1, 
              name = $2,
              vendor_uuid = $3
       where  uuid = $4`,
		"delete": `update strains set mtime = current_timestamp, dtime = current_timestamp where uuid = $1`,
		// XXX: combine this with the general select?
		"generated-strain": `
      select  s.uuid,
              s.species,
              s.name,
              s.ctime,
              s.dtime,
              v.uuid as vendor_uuid,
              v.name as vendor_name,
              v.website as vendor_website
        from  strains s 
        join  vendors v
          on  s.vendor_uuid = v.uuid
       where  s.generation_uuid = $1
       order
          by  s.name, s.ctime`,
		"update-gen-strain": `
      update  strains
         set  generation_uuid = $1
       where  uuid = $2`,
	},

	"strainattribute": {
		"get-unique-names": `select distinct name from strain_attributes order by name`,
		"all": `
      select uuid, name, value
        from strain_attributes sa
       where strain_uuid = $1
       order
          by name`,
		"add": `
    insert
      into strain_attributes (uuid, name, value, strain_uuid)
    select $1, $2, $3, s.uuid
      from strains s
     where s.uuid = $4`,
		"change": `
    update strain_attributes
       set value = $1,
           name = $2
     where uuid = $3`,
		"remove": `delete from strain_attributes where uuid = $1`,
	},

	"substrate-ingredient": {
		"all": `
    select i.uuid,
           i.name
      from ingredients i
      join substrate_ingredients si
        on i.uuid = si.ingredient_uuid
     where si.substrate_uuid = $1`,
		"add": `
    insert
      into substrate_ingredients (uuid, substrate_uuid, ingredient_uuid)
    select $1, s.uuid, i.uuid
      from substrates s
      join ingredients i
        on s.uuid = $2
       and i.uuid = $3`,
		"change": `
    update substrate_ingredients
       set ingredient_uuid = $1
     where substrate_uuid = $2
       and ingredient_uuid = $3`,
		"remove": `
    delete
      from substrate_ingredients
     where substrate_uuid = $1
       and ingredient_uuid = $2`,
	},

	"substrate": {
		"select-all": `
      select  s.uuid,
              s.name,
              s.type,
              v.uuid as vendor_uuid,
              v.name as vendor_name,
              v.website as vendor_website
        from  substrates s
        join  vendors v
          on  s.vendor_uuid = v.uuid
       order
          by  s.name`,
		"select": `
      select  s.uuid,
              s.name,
              s.type,
              v.uuid as vendor_uuid,
              v.name as vendor_name,
              v.website as vendor_website
        from  substrates s
        join  vendors v
          on  s.vendor_uuid = v.uuid
       where  s.uuid = coalesce($1, s.uuid)
         and  v.uuid = coalesce($2, v.uuid)`,
		"insert": `
      insert
        into substrates(uuid, name, type, vendor_uuid)
      select $1, $2, $3, v.uuid
        from vendors v
      where v.uuid = $4`,
		"update": `
      update  substrates s
         set  name = $1,
              type = $2,
              vendor_uuid = v.uuid
        from  vendors v 
       where  v.uuid = $3
         and  s.uuid = $4`,
		"delete": `delete from substrates where uuid = $1`,
	},

	"timestamp": {
		"touch":    `update %s set mtime = $1 where uuid = $2`,
		"update":   `update %s set %s where uuid = $1`,
		"undelete": `update %s set dtime = null where uuid = $1`,
	},

	"vendor": {
		"select-all": `select uuid, name, website from vendors order by name`,
		"select":     `select uuid, name, website from vendors where uuid = $1`,
		"insert":     `insert into vendors(uuid, name, website) values($1, $2, $3)`,
		"update":     `update vendors set name = $1, website = $2 where uuid = $3`,
		"delete":     `delete from vendors where uuid = $1`,
	},
}

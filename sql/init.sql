\c huautla;

create table uuids (
  uuid  varchar(40) not null primary key,
  mtime timestamp   not null default current_timestamp,
  ctime timestamp   not null default current_timestamp
);

begin; /** base table creation */
  -- anything that donates DNA to a generation
  create table progenitors () inherits(uuids);

  -- sort of like a pub/sub where generations and lifecycles create 
  -- events, and the events table is the observer that captures them
  create table observables () inherits(uuids);

  -- anything you want to keep notes on - which could be most things
  create table notables() inherits(uuids);

  begin; /** base table constraints */
    create function noinsert()
    returns trigger
    as
    $$
    begin
      raise exception 'base tables cannot be changed directly';
    end
    $$
    language plpgsql;

    create trigger UUIDInserter
      before  insert
          on  uuids
        for  each statement
    execute  function noinsert();

    create trigger  ProgenitorInserter
      before  insert
          on  progenitors
        for  each statement
    execute  function noinsert();

    create trigger  ObservableInserter
      before  insert
          on  observables
        for  each statement
    execute  function noinsert();

    create trigger NotableInserter
    before insert
        on notables
      for each statement
  execute function noinsert();
  end;
end;

create table vendors (
  uuid    varchar(40)  not null primary key,
  name    varchar(512) not null unique,
  website varchar(512) not null default ''
) inherits(progenitors);

create table substrates (
  uuid        varchar(40)  not null primary key,
  name        varchar(512) not null,
  type        varchar(25)  not null check (type in ('Agar', 'Liquid', 'Grain', 'Bulk')),
  vendor_uuid varchar(40)  not null references vendors(uuid),
  unique(name, vendor_uuid)
) inherits(uuids);

create table ingredients (
  uuid            varchar(40)  not null primary key,
  name            varchar(512) not null unique
) inherits(uuids);

create table substrate_ingredients (
  uuid             varchar(40) not null primary key,
  substrate_uuid   varchar(40) not null references substrates(uuid),
  ingredient_uuid  varchar(40) not null references ingredients(uuid),
  unique(substrate_uuid, ingredient_uuid)
) inherits(uuids);

create table strains (
  uuid        varchar(40)  not null primary key,
  species     varchar(128) not null default '',
  name        varchar(512) not null,
  vendor_uuid varchar(40)  not null references vendors(uuid),
  unique(name, vendor_uuid, ctime)
) inherits(progenitors);

create table strain_attributes (
  uuid         varchar(40)  not null primary key,
  name         varchar(40)  not null,
  value        varchar(512) not null,
  strain_uuid  varchar(40)  not null references strains(uuid),
  unique(name, strain_uuid)
) inherits(uuids);

create table stages (
  uuid            varchar(40)  not null primary key,
  name            varchar(512) not null unique
) inherits(uuids);

create table event_types (
  uuid       varchar(40)  not null primary key,
  name       varchar(512) not null,
  severity   varchar(10)  not null check (severity in ('Begin', 'Info', 'Warn', 'Error', 'Fatal', 'RIP', 'Generation')),
  stage_uuid varchar(40)  not null references stages(uuid),
  unique(name, stage_uuid)
) inherits(uuids);

create table lifecycles (
  uuid                varchar(40)  not null primary key,
  location            varchar(128) not null,
  strain_cost         decimal(8,2) not null,
  grain_cost          decimal(8,2) not null,
  bulk_cost           decimal(8,2) not null,
  -- the net weight, fresh or dried; for dried, 1.0-(yield/gross) is how much water they typically contain
  yield               decimal(84,2) not null default 0,
  headcount           decimal(5) not null default 0,
  -- gross the fresh weight, regardless of whether they're sold fresh or dry (see yield)
  gross               decimal(5,2) not null default 0, 
  strain_uuid         varchar(40)  not null references strains(uuid),
  grainsubstrate_uuid varchar(40)  not null references substrates(uuid),
  bulksubstrate_uuid  varchar(40)  not null references substrates(uuid),
  unique(location, ctime)
) inherits(observables, notables);

create table events (
  uuid            varchar(40)  not null primary key,
  temperature     numeric(4,1) not null default 0.0,
  humidity        int          not null default 0,
  observable_uuid varchar(40)  not null /*references observables(uuid)*/,
  eventtype_uuid  varchar(40)  not null references event_types(uuid)
) inherits(progenitors, notables);

create table event_photos (
  uuid       varchar(40) not null primary key,
  filename   varchar(44) not null unique,
  event_uuid varchar(40) not null references events(uuid)
) inherits(uuids, notables);

create table generations (
  uuid                  varchar(40) not null primary key,
  platingsubstrate_uuid varchar(40) not null references substrates(uuid),
  liquidsubstrate_uuid  varchar(40) not null references substrates(uuid)
) inherits(observables, notables);

alter table strains add generation_uuid varchar(40) null references generations(uuid);

create table sources (
  uuid            varchar(40) not null primary key,
  type            varchar(8)  not null check (type in ('Clone', 'Spore')),
  progenitor_uuid varchar(40) not null /*references progenitors(uuid)*/,
  generation_uuid varchar(40) not null references generations(uuid),
  unique(progenitor_uuid, generation_uuid)
) inherits(uuids);

create table notes (
  uuid varchar(40) not null primary key,
  note text not null,
  notable_uuid varchar(40) not null
) inherits(uuids);

begin; /** progenitor constraints */
  create function progenitordelete()
  returns trigger
  as
  $$
  begin
    if exists (select 1 from sources s where s.progenitor_uuid = old.uuid) then
      raise exception 'foreign key violation';
    end if;
    return old;
  end
  $$
  language plpgsql;

  create trigger StrainProgenitorDelete
    before  delete
        on  strains
      for  each row
  execute  function progenitordelete();

  create  trigger EventProgenitorDelete
    before  delete
        on  events
      for  each row
  execute  function progenitordelete();
end;

begin; /** observable constraints */
  create function observabledelete()
  returns trigger
  as
  $$
  begin
    if exists (select 1 from events e where e.observable_uuid = old.uuid) then
      raise exception 'foreign key violation';
    end if;
    return old;
  end
  $$
  language plpgsql;

  create trigger LifecycleObservableDelete
    before  delete
        on  lifecycles
      for  each row
  execute  function observabledelete();

  create trigger GenerationObservableDelete
    before  delete
        on  generations
      for  each row
  execute  function observabledelete();
end;

begin; /** notable constraints */
  create function notabledelete()
  returns trigger
  language plpgsql
  as
  $$
  begin
    if exists (select 1 from notes n where n.notable_uuid = old.uuid) then
      raise exception 'foreign key violation';
    end if;
    return old;
  end
  $$;

  create trigger LifecycleNotableDelete
    before delete
        on lifecycles
      for each row
  execute function notabledelete();

  create trigger EventNotableDelete
    before delete
        on events
      for each row
  execute function notabledelete();

  create trigger EventPhotoNotableDelete
    before delete
        on event_photos
      for each row
  execute function notabledelete();

  create trigger GenerationNotableDelete
    before delete
        on generations
      for each row
  execute function notabledelete();
end;

begin; /** source constraints */
  create function sourcechange() 
  returns trigger
      as
  $$
  begin
    if not exists (select 1 from progenitors p where p.uuid = new.progenitor_uuid) then
      raise exception 'no existing progenitor';
    elsif exists (
      select  1
        from  sources s
      where  s.generation_uuid = new.generation_uuid
        and  s.type != new.type
        and  s.uuid != new.uuid
    ) then
      raise exception 'source types can''t be mixed';
    elsif exists (
      select  1
        from  sources s
        join  (
                select  2 as cap,
                        'Spore' as type
                union
                select  1,
                        'Clone'
              ) limits
          on  s.type = limits.type
      where  s.generation_uuid = new.generation_uuid
        and  s.uuid != new.uuid
      group
          by  limits.cap
      having  limits.cap = count(s.uuid)
    ) then
      raise exception 'too many sources for this generation';
    elsif exists (
      select  1
        from  events e
        join  event_types t
          on  e.eventtype_uuid = t.uuid
      where  e.uuid = new.progenitor_uuid
        and  t.severity != 'Generation'
    ) then
      raise exception 'event is not a generation type';
    end if;
    return new;
  end
  $$
  language plpgsql
  ;

  create trigger CheckSource
    before  insert or update
        on  sources
      for  each row
  execute  function sourcechange();
end;

begin; /** event constraints */
  create function eventchange()
  returns  trigger
      as
  $$
  begin
    if exists (select 1 from observables o where o.uuid = new.observable_uuid) then
      return new;
    end if;
    raise exception 'foreign key violation';
  end
  $$
  language plpgsql;

  create trigger CheckObservable 
  before  insert or update of observable_uuid
      on  events
     for  each row
 execute function eventchange();
end;

begin; /** note constraints */
  create function notechange()
  returns  trigger
      as
  $$
  begin
    if exists (select 1 from notables n where n.uuid = new.notable_uuid) then
      return new;
    end if;
    return null;
    -- raise exception 'foreign key violation';
  end
  $$
  language plpgsql;

  create trigger CheckNotable
    before insert or update of notable_uuid
        on notes
       for each row
   execute function notechange();
end;
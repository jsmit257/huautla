## Huautla
Named after the town in Jiménez that was the bridge between native traditions and what qualifies as a modern understanding of psychedelics. We neither recommend nor encourag the use of psychedelics, and that's not the point of this project. Our name is just a nod to a very influential force in human evolution.

### Overview
The goal of this project is to track analytics observing the lifecycle of fungi. We wanted it to be more multi-purpose - like for plants and animals and such - but that seems impractical. There are 3 major lifecycle events that we track with their own sets of details.

* Innoculation: this is the period between the actual innoculation and 100% colonization in the grain substrate. The grain doesn't have to me mixed with bulk substrate at the end of this event, it could be iced for a while, which is why this is a separate lifecycle event
* Fruiting: basically covers everything that happens after the grain and bulk substrates are mixed and set in a bin to fully colonize, up to the time the colony is spent. Spores and/or samples for cloning could be gathered during this phase, but they don't necessarily need to be cultured immediately, which is why this event is separate from culturing
* Culturing: could be using spore prints or actual mycelia to create (preferably) liquid cultures for future innoculations, and then we're back at step one. Cultures can also be stored for quite some time too, so that's why this is its own event

But more on all that later (or elsewhere)

### Basics
These are just the high-level topics you can read in more detail in more appropriate locations

* Templates
* Strains
* Events
```sql
create table event_types (
  id varchar(40) not null primary key,
  label varchar(255) not null
)
create table events (
  id varchar(40) not null primary key,
  ctime timestamp not null,
  temperature decimal(3,0) not null,
  type varchar(40) not null references event_types(id),
  ...
)
```

### Downloading

### Installing

### Configuring

### Running

### Testing

### Contributing

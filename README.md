## Huautla
Named after the town in Jiménez that was the bridge between native traditions and what qualifies as a modern understanding of psychedelics. We neither recommend nor encourage the use of psychedelics, and that's not the point of this project. Our name is just a nod to an interesting story, and a very influential force in human evolution.

### Overview
The goal of this project is to track analytics observing the lifecycle of fungi. We wanted it to be more multi-purpose - like for plants and animals and such - but that seems impractical. There are 3 major lifecycle events that we track with their own sets of details.

* Innoculation: this is the period between the actual innoculation and 100% colonization in the grain substrate. The grain doesn't have to me mixed with bulk substrate at the end of this event, it could be iced for a while, which is why this is a separate lifecycle event
* Fruiting: basically covers everything that happens after the grain and bulk substrates are mixed and set in a bin to fully colonize, up to the time the colony is spent. Spores and/or samples for cloning could be gathered during this phase, but they don't necessarily need to be cultured immediately, which is why this event is separate from culturing
* Culturing: could be using spore prints or actual mycelia to create (preferably) liquid cultures for future innoculations, and then we're back at step one. Cultures can also be stored for quite some time too, so that's why this is its own event

But more on all that later (or elsewhere)

### Basics
???

### Object Model
These are the database tables described as a golang object tree; the sql perspective is [here](./sql/init.pgsql); no, this isn't really yaml

```yaml
vendor: &vendor
  uuid: surrogate key
  name: must be unique

strainattribute: &strainattribute
  uuid: surrogate key
  name: the key in the key/value pair, can't be changed after it's created
  value: arbitrary value - use it for anything

strain: &strain
  uuid: surrogate key
  name: just a name, without a vendor it doesn't mean much
  vendor: *vendor ### vendors give reputation to names, so (name+vendor) is unique in this table
  attributes: list of *strainattribute

ingredient: &ingredient
  uuid: surrogate key
  name: unique in this table, but can be shared by many substrates

substrate: &substrate
  uuid: surrogate key
  name: the name the vendor gave it
  type: constrained to ('Bulk', 'Grain', TBD)
  vendor: *vendor ### substrates are unique by (name+vendor) for the same reasons as strains
  ingredients: list of 0 or more *ingredient

grainsubstrate: &grainsubstrate *substrate
  ### these are just substrates that pass an API level check (not a check constraint on the database, yet) that only allows types of 'Grain'

bulksubstrate: &bulksubstrate
  ### like grainsubstrate, except the type check is for 'Bulk'

stage: &stage
  uuid: surrogate key
  name: unique to this table

eventtype: &eventtype
  uuid: surrogate key
  name: something short and descriptive
  severity: label for quick filtering of log noise
  stage: *stage ### unique with name, b/c some types wouldn't apply to some stages, and there may be dublicate names

event: &event
  uuid: surrogate key
  temperature:
  humidity:
  mtime: date/time when this record was last modified
  ctime: date/time when this record was created
  eventtype: *eventtype

lifecycle:
  uuid: surrogate key
  name: unique identifier ### probably going away or being otherwise refactored
  location: where the bin/bag/jar/agar is being stored
  grain_cost:
  bulk_cost:
  yield: in grams, how much dried or fresh product was shipped (see gross)
  count: how many caps were harvested?
  market_price: TBD - this should be the market value when the lifecycle begins; it may go up, it may go down, the only thing you know for sure is this is where it's at now; once entered, you cacn't change this, but that's not an invitation to speculate
  gross: product fresh weight, mostly to anticipate weight loss from dehydration
  mtime: date/time when this record was last modified
  ctime: date/time when this record was created
  strain: *strain
  grainsubstrate: *grainsubstrate
  bulksubstrate: *bulksubstrate
  events: list of *event
```
The root of anything interesting is the Lifecycle. It's generally a path from inception to completion, but it's a 'cycle' and not a 'span' because it doesn't always start at the beginning, and it will eventually support spawning one lifecycle to create new lifecycles, which makes them more tree-like than discrete instances (though ironically not cyclical). But that isn't really finished yet

### Configuring
Basically, fill out this form
```gol
	Config struct {
		PGHost string
		PGUser string
		PGPass string
		PGPort uint
		PGSSL  string
	}
```
and submit it to huautla.New, along with a logger. There's a workable reference implementation in the system test [init](./tests/system/main_test.go) function.

### Using
Comments are inlined in the [api](./types/api.go) source, rather than have to maintain that and a README separately

### Testing
TBD: tag docker postgres+huautla images with code releases

Until then, the [docker-compose](./docker-compose.yml) should be pretty portable to client environments, but you'll need to change the location of the [sql](./sql/) and [bin](./bin) script sources. You only need enough to run the `install` service

### Contributing

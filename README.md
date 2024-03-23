## [Huautla](https://github.com/jsmit257/huautla)
Named after the town in Jiménez that was the bridge between native traditions and what qualifies as a modern understanding of psychedelics. There's an interesting story of people and politics behind their noteriety, so like any memorial, this is a nod to what we consider to be a remarkable moment in history.

### Overview
The goal of this project is to track analytics observing the lifecycle of fungi. There's a lot more to consider before we could call this list of attributes comprehensive for many fungi, but for now we've captured the few events that cover most of our needs. Even this highly abstract overview of a colony over time is missing many species-specific attributes. We'll continue to evolve the analytics as we encounter more distinct charecteristics.

More details of a lifecycle are described in the [object model](#object-model), but at a higher level, lifecycles evolve through some of the following 4 stages:

* Gestation: These are the range of events between culturing spores, and having a live colony in stasis, so to speak, that's ready to be introduced to an environment that can support growth to the point of procreation. Simply and roughly put, this is the time between a spore-print, and life outside of sugar-water
* Colonization: Maybe 'adolescence' is a better name. Probably sealed in a bag of grain(s), woodchips or ???, this is when cosmopolitan organisms establish a colony, share resources, build strong contamination resistance - all in the safety of closed bubble.
* Majority: This is when the colony is established and strong and can grow to the excess needed to fruit new generations. An example is mixing colonized adolescents from the previous stage with porous, moist, nutritive bulk substrate in a bed. They can consume much more than the colony needs, so the extra will be used for procreation.
* Vacation: A sort of suspended animation when a fungus is ahead of schedule and needs to stop working. Fortunately, they chill well with minimal damage, and simply wake up when they get warm again. So time spent on `Vacation` within a lifecycle isn't fully counted into a normal temporal model for a particular strain, but it still needs to be predictable in terms of how long they can chill, and whether they end with fruits, or the colony dies.

### Requirements
To develop and test locally, you'll need at least:
* Git (you only need an account if you plan to change anything)
* Golang environment
* docker/docker-compose (sometimes they're packaged separately, get both)
* postgres-client (optional) if you'd like to query the running container from localhost, for debugging, etc
* postgres-server (optional) if you want to work/test on a local database like one that you use for other things, with scheduled backups, etc; see [Docker](#docker) for a turnkey standalone database container
* `make` probably a GNU-compatible one
* bash-compatible shell (bourne doesn't handle arrays nicely); changing this would require other users to use your shell, s0 we prefer `bash`; `zsh` is ubiquitous, but `tcsh` is a bit much

### Docker
The only build artifacts from this project are the docker images at [dockerhub](https://hub.docker.com/repository/docker/jsmit257/huautla/tags?page=1&ordering=last_updated). Proper semantic/commit-sha versioning isn't currently supported, but there will always be at least `...:lkg` (last-known good) and `...:lts` (long-term-support) versions. Last-known good is tagged when `make install-system-test` successfully seeds the test data. LKG is pushed to the remote as `jsmit257/huautla:lkg` when `make system-test` succeeds.

Only the minimal seed data is captured in the image, the test data is lost when the test container exits. 

Additional entrypoints are also packaged in the image for the purpose of persistence management. References implementations for all the following features are documented in [cffc standalone](https://github.com/jsmit257/centerforfunguscontrol/standalone/docker-compose.yml). The scripts themselves have descriptive errors, where possible.

- [migration](./bin/migration-entrypoint.sh) restores data from a running source to a destination. In the reference implementation, the source is the minimally pre-seeded huautla database, and the destination is a short-lived empty database with a volume mounted from the host - hence, persistent. This should only be run once - the service issues appropriate errors for troubleshooting.
  #### parameters:
    - `SOURCE_HOST`: (required) data source hostname; typically, the hostname of a venilla instance of `jsmit257:lkg` or similar
    - `SOURCE_PORT`: (optional, default:5432) postgres port on the source host
    - `SOURCE_USER`: (optional, default:postgres) user with admin privileges on the source server
    - `DEST_HOST`: (optional, default:localhost) instance to seed from the source; the service fails with a descriptive error if a database already exists
    - `SOURCE_PORT`: (optional, default:5432) postgres port on the destination host
    - `SOURE_USER`: (optional, default:postgres) user with admin privileges on the destination server
- [backup](./bin/backup-entrypoint.sh) backup a running instance to the `/pgbackups` directory. This is mostly only useful if that directory is mapped to a persistent volume.
  #### parameters:
    - `SOURCE_HOST`: (required) data source hostname
    - `SOURCE_PORT`: (optional, default:5432) postgres port on the source host
    - `SOURCE_USER`: (optional, default:postgres) user with admin privileges on the source server
    - `POSTGRES_PASSWORD`: (required) password for the source user. This is typically supplied by a `docker-compose.yml` and/or `.env` file. These passwords are typically stock `root` when used with a vanilla `postgres:<whatever>` tag, but attempt to avoid saving this anywhere including command history if it's at-all sensitive
- [restore](./bin/restore-entrypoint.sh) restores an archive to a running instance of the huautla database This is mostly only useful if that directory is mapped to a persistent volume.
  #### parameters:
    - `DEST_HOST`: (required) data destination hostname
    - `DEST_PORT`: (optional, default:5432) postgres port on the destination host
    - `DEST_USER`: (optional, default:postgres) user with admin privileges on the destination server
    - `POSTGRES_PASSWORD`: (required) password for the destination user. The same recommendations and caveats apply as with `backup`
    - `RESTORE_POINT`: (required) see the reference implementation for a description of this parameter

TODO: webhook with github/et al.

### Local Database

### Using
Public bindings are consolidated in the [api](./types/api.go) and [data types](./types/data.go).

### Object Model
These are the database tables described as a golang object tree; the sql perspective is [here](./sql/init.pgsql), and there are [cliff-notes](./sql/cliff-notes); no, this isn't really yaml

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
Basically, fill out this form and submit it to `huautla.New()`, along with a logger.
```go
Config struct {
  PGHost string
  PGUser string
  PGPass string
  PGPort uint
  PGSSL  string
}
```
There's a workable reference implementation in the system test [init()](./tests/system/main_test.go) function.

### Testing
- `make unit` obviously handles the unit-testing - i.e. how the persistence-bindings respond to cretain events from the database server
- `make system-test` loads additional data, partly to make sure any referential- or other integrity-constraints aren't violated, then runs the [system tests](./tests/system) to veryfy basic CRUD opeartions, including all possible errors thrown from the database.

### Contributing
Please do!!! The [license](./LICENCE.md) is open and free with attribution.
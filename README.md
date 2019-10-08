# gQB
[![coverage](https://godoc.org/github.com/jobstoit/gqb?status.svg)](https://godoc.org/github.com/jobstoit/gqb)
[![pipeline status](https://gitlab.com/jobstoit/gsql/badges/master/pipeline.svg)](https://gitlab.com/jobstoit/gsql/commits/master)
[![coverage report](https://gitlab.com/jobstoit/gsql/badges/master/coverage.svg)](https://gitlab.com/jobstoit/gsql/commits/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/jobstoit/gqb)](https://goreportcard.com/report/github.com/jobstoit/gqb)

The generated Query Builder creates a [query builder](https://github.com/ultraware/qb) and if needed a migration using a yaml configuration.

The gqb takes takes the following flags as arguments:
```
-migrate		Writes the configuration to the specified sql file
			takes the output file as argument

-db			Directly inserts the configured structure as migration
			into the database using the DB_DRIVER and DB_URL enviroment
			variables as flags for this mode

-model			Writes the configruration to git.ultraware.nl/NiseVoid/qb
			models and takes the output as argument

-pkg			Specifies the package name for the qbmodels (default model)

-dvr			Specifies the driver for sql generation (default postgres)
```

## Configuration
Create a configuration file based on the following yaml structure (this is an example configuration):
```yaml
pkg: models                             # optional
driver: mysql                           # optional (default postgres)
tables:
  users:
    id: int, primary
    email: varchar, unique
    password: varchar
    first_name: varchar(50)
    last_name: varchar(100)
    role: role                          # foreign key (enum)

  posts:
    id: int, primary
    created_by: users.id                # foreign key
    created_at: datetime, default(NOW)
    updated_at: datetime, default(NOW)
    title: varchar
    subtitle: varchar, nullable
    context: text

  images:
    id: int, primary
    title: varchar
    created_at: datetime, default(NOW)
    updated_at: datetime, default(NOW)
    path: varchar

  post_images:
    id: int, primary
    image_id: images.id                 # foreign key
    post_id: posts(id)                  # another foreign key
    description: varchar

enums:
  role:
  - ADMIN
  - MODERATOR
  - GENERAL_USER
```

# Installation
Install this using go get:
```bash
$ go get -u github.com/jobstoit/gqb
```

Then start using it, try:
```bash
$ gqb
```

# Todo
- [x] Create configuration reader
- [x] Create SQL generator
- [x] Create NiseVoid/qb generator
- [ ] Create a database differential inspector (config against the current state of the database)
- [ ] Create a diffential SQL generator & executor
- [ ] Create a seeder & seeder configuration

# gQB
[![coverage](https://godoc.org/github.com/jobstoit/gqb?status.svg)](https://godoc.org/github.com/jobstoit/gqb)
[![pipeline status](https://gitlab.com/jobstoit/gsql/badges/master/pipeline.svg)](https://gitlab.com/jobstoit/gsql/commits/master)
[![coverage report](https://gitlab.com/jobstoit/gsql/badges/master/coverage.svg)](https://gitlab.com/jobstoit/gsql/commits/master)
[![go report](https://goreportcard.com/badge/github.com/jobstoit/gqb)](https://goreportcard.com/report/github.com/jobstoit/gqb)

The generated Query Builder creates a query builder and if needed a migration using a yaml configuration.

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
pkg: models                         # optional
driver: mysql                       # optional (default postgres)
tables:
  user:
    id: int, primary
    email: varchar, unique
    password: varchar
    first_name: varchar(50)
    last_name: varchar(100)
    role: role
  post:
    id: int, primary
    created_by: user.id             # foreign key
    created_at: datetime
    updated_at: datetime
    title: varchar
    subtitle: varchar, nullable
    context: text
  images:
    id: int, primary
    title: varchar
    created_at: datetime,
    updated_at: datetime
    path: varchar
  post_images:
    id: int, primary
    image_id: image.id              # foreign key
    post_id: post.id                # foreign key
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
$ go get -u gitlab.com/jobstoit/gqb
```

Then start using it, try:
```bash
$ gqb
```

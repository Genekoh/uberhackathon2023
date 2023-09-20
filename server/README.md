# HOW TO SETUP

#### This project uses GO version 1.23, python and sqlite3. Please make sure to have them installed

First, while in the server directory (this directory) run:

```
go get
```

to fetch all dependencies.

Then run

```
python migrateup.py
```

to set up the database.

And finally, start the server by running:

```
go run cmd/*
```

By default, this should be running on localhost:8080

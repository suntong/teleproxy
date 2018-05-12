
teleproxy
=========

[![GoCard][1]][2]
[![GitHub license][3]][4]

[1]: https://goreportcard.com/badge/LeKovr/teleproxy
[2]: https://goreportcard.com/report/github.com/LeKovr/teleproxy
[3]: https://img.shields.io/badge/license-MIT-blue.svg
[4]: LICENSE

[teleproxy](https://github.com/LeKovr/teleproxy) - Telegram proxy bot.

**WARNING:** Current version of this project is not intended for production use. This is an MVP (minimum viable product).
Refactoring, tests, docs and more than 1 committer is required for getting this project production-ready.

This service

* runs as telegram bot
* gets user messages
* forwards them to telegram group
* forwards replies to user

Features
--------

* [x] database storage for users and messages
* [x] autoregister all senders with short numerical id
* [x] group members can enable and disable users
* [x] message templates based on text/template

### ToDo

* [ ] tests
* [ ] correct reply on 'joined the group via invite link'
* [ ] file transfer

Install
-------

```
go get github.com/LeKovr/teleproxy
```

### Download

See [Latest release](https://github.com/LeKovr/teleproxy/latest)

Usage
-----

```
# create default config
make .env
```
Edit .env to suit your needs

### Without docker
```
# create postgresql database (see man createdb)
# ...

# run standalone
make run
```

### With docker
Required postgresql available via docker network (DCAPE_NET in .env).

```
# create database
make db-create DCAPE_DB=running_postgresql_container_name

# build docker image and run docker container
make up
```

## See also

```
# show all Makefile targets
make help
```

We use [dcape](https://github.com/dopos/dcape) in our systems.

License
-------

The MIT License (MIT), see [LICENSE](LICENSE).

Copyright (c) 2016 Alexey Kovrizhkin lekovr+teleproxy@gmail.com

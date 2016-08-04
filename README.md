
dbrpc
=====

[![GoCard][1]][2]
[![GitHub license][3]][4]

[1]: https://goreportcard.com/badge/LeKovr/teleproxy
[2]: https://goreportcard.com/report/github.com/LeKovr/teleproxy
[3]: https://img.shields.io/badge/license-MIT-blue.svg
[4]: LICENSE

[teleproxy](https://github.com/LeKovr/teleproxy) - Telegram proxy bot.

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

* tests

Install
-------

```
go get github.com/LeKovr/teleproxy
```

### Download

See [Latest release](https://github.com/LeKovr/teleproxy/latest)

License
-------

The MIT License (MIT), see [LICENSE](LICENSE).

Copyright (c) 2016 Alexey Kovrizhkin ak@elfire.ru

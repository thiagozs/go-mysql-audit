#### MySQL Proxy Audit

# go-mysql-audit - A very simple mysql proxy.

* Deep mysql connection alive
* Reuse authorized connection

## Helpers and commands

* Build a binary file

```sh
$: make build
```

* Helper commands

```sh
$:./paudit --help
Run proxy server for mysql

Usage:
  paudit [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  runserver   Run proxy for mysql

Flags:
  -h, --help   help for proxy

Use "paudit [command] --help" for more information about a command.
```

* Running proxy server

```sh
$:./paudit runserver --mysql=3306 --proxy=33060 --debug=false
2021/10/26 14:06:43 [SERVER] Proxy MySQL Audit
2021/10/26 14:06:43 [SERVER] version: beta
2021/10/26 14:06:43 [SERVER] build  : 38dc57f
2021/10/26 14:06:43 [SERVER] proxy server on host=33060, mysql server listening host=3306...
```

## Versioning and license

Our version numbers follow the [semantic versioning specification](http://semver.org/). You can see the available versions by checking the [tags on this repository](https://github.com/thiagozs/go-mysql-audit/tags). For more details about our license model, please take a look at the [LICENSE](LICENSE) file.

**2021**, thiagozs

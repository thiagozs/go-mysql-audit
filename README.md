# Simple mysql audit

The main goal is to learn the MySQL Protocol by implementing it.

The plan:
- [x] Implement TCP Proxy as a starting point
- [ ] Implement state machine
- [x] Implement query/query data buffering
- [ ] Implement plugins

Packets decode/encode todo:
- [x] Handshake Packet
- [ ] Authorization Packet


go version go1.7

To try it, just clone, and run:

```
go run main.go
```

## Versioning and license

Our version numbers follow the [semantic versioning specification](http://semver.org/). You can see the available versions by checking the [tags on this repository](https://github.com/thiagozs/go-mysql-audit/tags). For more details about our license model, please take a look at the [LICENSE](LICENSE) file.

2021, thiagozs

# kairosdb-streamer
[![Build Status](https://travis-ci.org/swco/kairosdb-streamer.svg)](https://travis-ci.org/swco/kairosdb-streamer)

Takes messages from stdin and uploads them to kairosdb. This is designed to work with the fluentd exec plugin.

# Installation
```bash
export GOPATH=~/go
go get github.com/swco/kairosdb-streamer
go test github.com/swco/kairosdb-streamer
go install github.com/swco/kairosdb-streamer
```

This will install kairosdb-streamer to `~/go/bin/kairosdb-streamer`

# Example fluentd config

```xml
<match kairosdb>
  type exec
  command /usr/local/bin/kairosdb-streamer -host "localhost:4242"
  format json
  buffer_path /var/log/fluent/kairosdb_streamer.*.buffer
</match>
```

See the [fluentd out_exec documentation](http://docs.fluentd.org/articles/out_exec) for more options.

# Development
To simulate a kairosdb instance in development you can use socat:

```bash
socat TCP-LISTEN:4242,fork -
```

metrics are expected to be read from a file or stdin in the following format:

```json
{"timestamp":1427291847309,"name":"memcache.status.connections","value":427.0,"tags":{"function":"cache","datacenter":"DC1","host":"host1.com","serverid":"HOST1"}}
{"timestamp":1427291847309,"name":"memcache.status.connections","value":200.0,"tags":{"function":"cache","datacenter":"DC2","host":"host2.com","serverid":"HOST2"}}
```

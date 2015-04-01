# kairosdb-streamer
Takes messages from stdin and uploads them to kairosdb. This is designed to work with the fluentd exec plugin.

# Installation
```bash
go get github.com/swco/kairosdb-streamer
go install github.com/swco/kairosdb-streamer
```

# Development
To simulate a kairosdb instance in development you can use socat:

    socat TCP-LISTEN:4242,fork -

metrics are expected to be read from a file in the following format:

    {"timestamp":1427291847309,"name":"memcache.status.connections","value":427.0,"tags":{"function":"cache","datacenter":"DC1","host":"host1.com","serverid":"HOST1"}}
    {"timestamp":1427291847309,"name":"memcache.status.connections","value":200.0,"tags":{"function":"cache","datacenter":"DC2","host":"host2.com","serverid":"HOST2"}}
    

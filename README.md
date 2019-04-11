# DNS Proxy
A simple DNS proxy written in go based on [github.com/miekg/dns](https://github.com/miekg/dns)

## How to use it

## Arguments

```
    -config		path to config filename
```

## Config file format

```json
{
    "listen_host" : "0.0.0.0",
    "listen_port" : "53",
    "default_dns": "8.8.8.8:53",
    "domains": {
        ".*example1.com": "192.168.30.11:53",
        ".*example2.com": "192.168.30.12:53",
    }
}
```
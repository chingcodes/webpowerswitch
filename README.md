# webpowerswitch

A simple client for controlling an ethernet [WebPowerSwitch](webpowerswitch.com).

__Note: This only works for pre-API ethernet switches, and not API-enabled WiFi ones.__

## Getting Started

### Building

```
go build webpowerswitch.go
```

### Usage

```
Usage: ./webpowerswitch [FLAGS] CMD OUTLET|all

   -addr string
        Address of Web Switch (default "192.168.0.100")
  -password string
        Password used to login (default "1234")
  -port int
        HTTP port of Web Switch (default 80)
  -user string
        Username to login (default "admin")
CMD
Command for switch [on|off]

OUTLET
Either outlet number or 'all'
```
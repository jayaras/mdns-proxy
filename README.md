# mdns-proxy
My home network and lab make heavy use of mDNS either directly (avahi & MacOS) or indirectly via https://github.com/blake/external-mdns.  This is all well and good except.  1.  I have several services that do not support mDNS (IoT Devices, AndroidTV doesn't seem to look cross bunet boundary mDNS reflection. and when I am VPN'ed in) So I found mdns-proxy out there modified it up to meet my needs and here we are.

This work is a fork of the original mdns-proxy code with the following changes.
* It is a restructure of the original code to test out some a couple patterns I was interested in.
* updates the mDNS package to adopt an issue with reflected mDNS traffic across subnet boundary/ethernet broadcast domains.  
* It also adds the server to have a recursive mode. 
* moved to cobra/viper for the flag and configuration.

## Running (Docker)
`docker run ghcr.io/jayaras/mdns-proxy/mdns-proxy`
## Building code
`go build`
## Building a container
`ko build`
## Configuration & Flags
```
Usage:
  mdns-proxy [flags]

Flags:
  -h, --help               help for mdns-proxy
  -i, --ip string          ip address to listen on (default "0.0.0.0")
  -p, --port uint16        dns server udp port (default 5345)
  -r, --recursive          enable recursive resolver (default true)
  -t, --timeout duration   timeout for mdns response (default 4s)
  -u, --upstream string    upstream DNS Server (default "192.168.1.1:53")
  -z, --zone string        authoritive dns zone (default "mdns.")
```
## Configuration via Environment Variables.  
Environment variables are also an option to change the settings.  The following variables are available to and allign to the coresponding flag settings after the MDNS_PROXY_ prefix.

```MDNS_PROXY_IP```

```MDNS_PROXY_PORT```

```MDNS_PROXY_RECURSIVE```

```MDNS_PROXY_TIMEOUT```

```MDNS_PROXY_UPSTREAM```

```MDNS_PROXY_ZONE```
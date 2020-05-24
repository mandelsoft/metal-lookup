# An HTTP Metadata Mapper for *kipxe* Discovery API

The project [*kpixe*](https://github.com/mandelsoft/kipxe) provides
a simple kubernetes based resource matching engine to serve *iPXE* 
requests based on request parameters.

To feed the matching and substitution engine request metadata is required.
It is basically taken from the request parameters, but *kipxe* also allows
to enrich this metadata based on a 
[Discovery API](https://github.com/mandelsoft/kipxe/blob/master/README.md#the-discovery-api).

One possibility here is the usage of REST calls for this purpose.
This project provides such a metadata enricher based of the
[metal api](https://github.com/metal-stack/metal-go).

It tries to match the given machine UUID and/or mac addresses to find
the machine in the metal database. If found, it completes the UUID and/or
mac addresses and adds some additional properties from the database.

These properties can then be used by the matching and substitution engine of
*kipxe* to provide and/or generate appropriate content for the iPXE requests.

# Command Line Reference

```
lookup machine objects

Usage:
  <options> [flags]

Flags:
      --bind-address-http string   HTTP server bind address
      --config string              config file
      --driver string              Driver URL
      --grace-period duration      Grace period for shutdown (default 2m0s)
  -h, --help                       help for <options>
      --hmac string                HMAC
      --metalconfig string         config file for metal-api
      --server-port-http int       HTTP server port (serving /healthz, /metrics, ...)
      --token string               Token
  -v, --version                    version for <options>
```


# paymail-server

Pretty-much the simplest possible implementation of a [paymail](https://bsvalias.org/) server.

Perfect for the self-hoster.

Why should you rely on third-parties to serve your decentralized address in a decentralized network?

If you consider this useful please send a satoshi or two to [jonas@cowboyprogrammer.org](payto:jonas@cowboyprogrammer.org)

## Installation

Build the binary with

```
go build
```

Then set it up with the provided system service and nginx snippet.

Also remember to set up a DNS SRV record as specified by https://bsvalias.org/02-01-host-discovery.html

# paymail-server

Pretty-much the simplest possible implementation of a [paymail](https://bsvalias.org/) server.

Perfect for the self-hoster.

Why should you rely on third-parties to serve your decentralized address in a decentralized network?

If you consider this useful please send some satoshis to [jonas@cowboyprogrammer.org](payto:jonas@cowboyprogrammer.org)

## Installation

Run

```
go get gitlab.com/spacecowboy/paymail-server
```

and you'll have the binary in `~/go/bin/paymail-server`

## Configuration

Remember to set up a DNS SRV record as specified by https://bsvalias.org/02-01-host-discovery.html

See `config-example.toml` for the configuration options.

## Running it

Execute the server like so

```
paymail-server --config=/path/to/config.toml
```

or preferrably - use the provided system service and nginx config.

The service assumes that you copy the binary as follows

```
mkdir -p /var/lib/paymail
cp paymail-server /var/lib/paymail/
```

And the config in /etc

```
mkdir -p /etc/paymail
cp config.toml /etc/paymail/config.toml
```

Place the system service where it belongs and start it

```
cp paymail.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable --now paymail.service
```

You should now be able to curl the server to check that it works:

```
curl localhost:26245/.well-known/bsvalias
```

## Make it publically available

You need to configure Nginx (or whatever you like) with a valid SSL-certificate and proxy requests to the local paymail-server.
See `paymail_nginx.conf` for an example.

In case you already have a site running on the host you want to expose paymail on - just include the `paymail_nginx.snippet` in your existing config.
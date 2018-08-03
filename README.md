# Attach custom domain to ngrok tunnel
Ngrok is an awesome tool which let you expose port on your local machine to the
internet. Ngrok gives you a random hostname for each tunnel.

This tool enables you to attach subdomains against those tunnels.

Currently supported DNS includes GoDaddy only.

## Build and install from source
This requires Golang to be installed

```bash
$ go get -v github.com/alyyousuf7/ngrok-custom-domain/cmd/ngrokcd
```

## Usage
`ngrokcd` requires `ngrok` to be already running.

`ngrokcd` can be used in three different ways mentioned below.

### Attach a single subdomain to the first tunnel
This is useful if you expose only one tunnel to the internet using
`ngrok http 8080` and wish to have a subdomain against it.

```bash
$ ngrokcd attach-first [subdomain] [domain] [service=godaddy] [key] [secret]
```

### Attach a single subdomain to a specific tunnel
This is useful if you expose only one tunnel to the internet using ngrok
configuration file and wish to have a subdomain against it.

In ngrok configuration, you mention a tunnel name. In the following
configuration, it is `mytunnel`.

```yaml
authtoken: [your ngrok authtoken]
tunnels:
  mytunnel:
    proto: http
    addr: 8080
```

```bash
$ ngrokcd attach [ngrok tunnel name] [subdomain] [domain] [service=godaddy] [key] [secret]
```

### Attach multiple subdomains to a number of tunnels
This is useful if you expose multiple tunnels to the internet and wish to attach
subdomains to couple of them.

You must have already used ngrok configuration file to expose them. For example
you have the following tunnels:

```yaml
authtoken: [your ngrok authtoken]
tunnels:
  api:
    proto: http
    addr: 4000
  web:
    proto: http
    addr: 8080
```

Now to attach one subdomain to `api` and two subdomains to `web`, you have to
modify the configuration file (`~/.ngrok2/ngrok.yml`) to something like this:

```yaml
authtoken: [your ngrok authtoken]
tunnels:
  api:
    proto: http
    addr: 4000
    records:
    - myapp-api
  web:
    proto: http
    addr: 8080
    records:
    - myapp-app
    - myapp-www
records:
  myapp-api:
    dns: personal-domain
    cname: api
  myapp-app:
    dns: personal-domain
    cname: app
  myapp-www:
    dns: personal-domain
    cname: www
dns:
  personal-domain:
    domain: aliyousuf.com
    service: godaddy
    key: [godaddy api key]
    secret: [godaddy api secret]
```

```bash
$ ngrokcd
2018/08/02 18:30:20 waiting for ngrok
2018/08/02 18:30:20 ngrok is up
2018/08/02 18:30:22 updating movies: 0c7d38d0.ngrok.io
2018/08/02 18:30:24 updating torrent: facae89f.ngrok.io
2018/08/02 18:30:25 updating media: 56cb7796.ngrok.io
```

## Flags / Configuration
```
  -c, --config string     ngrok YAML config file (default "/home/[user]/.ngrok2/ngrok.yml")
  -h, --help              help for ngrokcd
  -H, --host string       ngrok hostname (default "localhost")
  -p, --port int          ngrok port (default 4040)
  -r, --retries int       retry count (default 5)
  -d, --delay int         retry delay (ms) (default 1000)
```

## License
MIT
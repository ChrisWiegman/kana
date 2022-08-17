// nolint
package appSetup 

const (
DYNAMIC_TOML = `[tls.options]

[tls.options.default]
minVersion = "VersionTLS12"
sniStrict = true

[[tls.certificates]]
certFile = "/var/certs/kana.site.pem"
keyFile = "/var/certs/kana.site.key"
`
TRAEFIK_TOML = `[log]
level = "INFO"

[providers]
[providers.docker]
endpoint = "unix:///var/run/docker.sock"
exposedByDefault = false
network = "kana"
[providers.file]
filename = "/etc/traefik/dynamic.toml"

[api]
dashboard = true
debug = true
insecure = true

[entryPoints]
[entryPoints.web]
address = ":80"
[entryPoints.web.http]
[entryPoints.web.http.redirections]
[entryPoints.web.http.redirections.entryPoint]
scheme = "https"
to = "websecure"

[entryPoints.websecure]
address = ":443"
`
)

// nolint
package setup 

const (
DYNAMIC_TOML = `[tls.options]

[tls.options.default]
minVersion = "VersionTLS12"
sniStrict = true

[[tls.certificates]]
certFile = "/var/certs/kana.ca.pem"
keyFile = "/var/certs/kana.ca.key"
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
[entryPoints.web-secure]
address = ":443"
`
)

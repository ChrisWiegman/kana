// nolint
package setup 

const (
DYNAMIC_TOML = `[tls.options]

[tls.options.default]
minVersion = "VersionTLS12"
sniStrict = true

[[tls.certificates]]
certFile = "/var/certs/cert.pem"
keyFile = "/var/certs/key.pem"
`
TRAEFIK_TOML = `[log]
level = "INFO"

[providers]
[providers.docker]
endpoint = "unix:///var/run/docker.sock"
exposedByDefault = false
network = "kana"

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

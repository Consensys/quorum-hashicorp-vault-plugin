backend "file" {
  path = "/vault/file"
}

listener "tcp" {
  address = "vault:8200"
  tls_disable = true
}

default_lease_ttl = "12h"
max_lease_ttl = "24h"
api_addr = "http://vault:8200"
plugin_directory = "/vault/plugins"
log_level = "Debug"
ui = true

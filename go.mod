module github.com/consensys/quorum-hashicorp-vault-plugin

go 1.15

require (
	github.com/consensys/gnark-crypto v0.5.0
	github.com/consensys/quorum v2.7.0+incompatible
	github.com/ethereum/go-ethereum v1.10.13
	github.com/golang/mock v1.4.3
	github.com/hashicorp/go-hclog v0.9.2
	github.com/hashicorp/vault/api v1.0.5-0.20200117231345-460d63e36490
	github.com/hashicorp/vault/sdk v0.1.14-0.20200305172021-03a3749f220d
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
)

replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.4.1+incompatible

// Containous forks
replace (
	github.com/abbot/go-http-auth => github.com/containous/go-http-auth v0.4.1-0.20200324110947-a37a7636d23e
	github.com/go-check/check => github.com/containous/check v0.0.0-20170915194414-ca0bf163426a
)

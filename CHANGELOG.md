# CHANGELOG

All notable changes to this project will be documented in this file.

##  v1.1.4 (2022-01-03)
- Support for optional SSL mode within init script
- Support for optional kv-v2 engine within init script
- Include vault-agent init script

##  v1.1.3 (2021-12-22)
### 🛠 Bug fixes
- Allow empty namespace in migration script

##  v1.1.2 (2021-10-6)
### 🆕 Features
- Added migration namespace script

##  v1.1.1 (2021-10-8)
### 🆕 Features
- Publishing of docker hub images

##  v1.1.0 (2021-10-6)

### ⚠ BREAKING CHANGES
- Renamed `bn254` by `babyjubjub`
- Migrated `/ethereum` namespace to `/keys`

##  v1.0.0 (2021-09-22)
### 🆕 Features
- Support for key operations on Hashicorp
- Supported curves: `bn254` and `secp256k1`
- Supported signing algorithms: `ecdsa` and `eddsa`

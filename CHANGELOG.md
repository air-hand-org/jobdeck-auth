# Changelog

## [1.2.3](https://github.com/air-hand-org/jobdeck-auth/compare/v1.2.2...v1.2.3) (2026-10-24)


### Bug Fixes

* **build:** isolate mise environment and remove unused actor lookup ([#32](https://github.com/air-hand-org/jobdeck-auth/issues/32)) ([b2b5a97](https://github.com/air-hand-org/jobdeck-auth/commit/b2b5a97571409edf0f91cd8ca777ddccc107bbd8))
* **build:** use mise base image for auth builds ([#30](https://github.com/air-hand-org/jobdeck-auth/issues/30)) ([8b9b824](https://github.com/air-hand-org/jobdeck-auth/commit/8b9b824379700af1d0491d9dd9674bf3e286d047))
* **deps:** update dependency go to v1.27.3 ([#39](https://github.com/air-hand-org/jobdeck-auth/issues/39)) ([9f823ab](https://github.com/air-hand-org/jobdeck-auth/commit/9f823ab9dac1c3ae5368580181e6d0035a1c91e1))
* **deps:** update module connectrpc.com/connect to v1.21.1 ([#47](https://github.com/air-hand-org/jobdeck-auth/issues/47)) ([9dfcca6](https://github.com/air-hand-org/jobdeck-auth/commit/9dfcca6361c9f773d2d6575ba0c9663b983d7296))
* **deps:** update module github.com/air-hand-org/api-contracts/packages/jobdeck/auth to v1.3.3 ([#34](https://github.com/air-hand-org/jobdeck-auth/issues/34)) ([c8ceaf1](https://github.com/air-hand-org/jobdeck-auth/commit/c8ceaf14d3365f13f91d0d4957f1c1291aec053c))
* **deps:** update module github.com/exaring/otelpgx to v0.13.0 ([#49](https://github.com/air-hand-org/jobdeck-auth/issues/49)) ([6dc3bf5](https://github.com/air-hand-org/jobdeck-auth/commit/6dc3bf5348955bba30038084c5f67a18a3916384))
* **deps:** update module github.com/jackc/pgx/v5 to v5.11.0 ([#50](https://github.com/air-hand-org/jobdeck-auth/issues/50)) ([34efdf8](https://github.com/air-hand-org/jobdeck-auth/commit/34efdf8019f231e2c4f9ecc64a7e67c93c0f4039))
* **deps:** update module github.com/stretchr/testify to v1.12.1 ([#51](https://github.com/air-hand-org/jobdeck-auth/issues/51)) ([7d327f1](https://github.com/air-hand-org/jobdeck-auth/commit/7d327f18c78c8938dfb095c50dc77dbae56bbd37))
* **deps:** update module golang.org/x/crypto to v0.58.0 ([#53](https://github.com/air-hand-org/jobdeck-auth/issues/53)) ([6fd5774](https://github.com/air-hand-org/jobdeck-auth/commit/6fd5774d42aec63892174260e9bf80a00ec5098b))
* **deps:** update module google.golang.org/grpc to v1.84.0 ([#54](https://github.com/air-hand-org/jobdeck-auth/issues/54)) ([87407fd](https://github.com/air-hand-org/jobdeck-auth/commit/87407fded16133f97c227e13a3b4148e70602f7b))
* **deps:** update testcontainers-go monorepo to v0.44.0 ([#56](https://github.com/air-hand-org/jobdeck-auth/issues/56)) ([2df1603](https://github.com/air-hand-org/jobdeck-auth/commit/2df16038efe467a79dfa522e5df8ad926e3156a1))
* **renovate:** remove api-contracts release age delay ([#33](https://github.com/air-hand-org/jobdeck-auth/issues/33)) ([cb8091d](https://github.com/air-hand-org/jobdeck-auth/commit/cb8091d026a21f75f7a3a3abd8c2ef31c145793f))

## [1.2.2](https://github.com/air-hand-org/jobdeck-auth/compare/v1.2.1...v1.2.2) (2026-10-21)


### Bug Fixes

* package application and dependency licenses ([#19](https://github.com/air-hand-org/jobdeck-auth/issues/19)) ([13f4d86](https://github.com/air-hand-org/jobdeck-auth/commit/13f4d862415e06d1c1ae92cc3c5475efa3ab281e))

## [1.2.1](https://github.com/air-hand-org/jobdeck-auth/compare/v1.2.0...v1.2.1) (2026-09-30)


### Bug Fixes

* replace auth Envoy proxy with Connect-Go ([#17](https://github.com/air-hand-org/jobdeck-auth/issues/17)) ([eb51a72](https://github.com/air-hand-org/jobdeck-auth/commit/eb51a726f4483d02e182aa872266aa28ff4504c7))

## [1.2.0](https://github.com/air-hand-org/jobdeck-auth/compare/v1.1.1...v1.2.0) (2026-09-18)


### Features

* add auth Envoy proxy ([#11](https://github.com/air-hand-org/jobdeck-auth/issues/11)) ([f3eef27](https://github.com/air-hand-org/jobdeck-auth/commit/f3eef270e97f50ee7d2a860fe16e40d1f0ddcb61))

## [1.1.1](https://github.com/air-hand-org/jobdeck-auth/compare/v1.1.0...v1.1.1) (2026-09-15)


### Bug Fixes

* update auth package ([#9](https://github.com/air-hand-org/jobdeck-auth/issues/9)) ([956241e](https://github.com/air-hand-org/jobdeck-auth/commit/956241eab95a0ba2077d8c96142715a2316d4806))

## [1.1.0](https://github.com/air-hand-org/jobdeck-auth/compare/v1.0.5...v1.1.0) (2026-09-11)


### Features

* **auth:** expose JWKS endpoint ([#7](https://github.com/air-hand-org/jobdeck-auth/issues/7)) ([61b2e7f](https://github.com/air-hand-org/jobdeck-auth/commit/61b2e7fe8275aed5bd5e0983ec79ba637fa10354))

## [1.0.5](https://github.com/air-hand-org/jobdeck-auth/compare/v1.0.4...v1.0.5) (2026-09-08)


### Bug Fixes

* change event type ([#5](https://github.com/air-hand-org/jobdeck-auth/issues/5)) ([7791da8](https://github.com/air-hand-org/jobdeck-auth/commit/7791da87fc5aa9f233186a6759b1fe1f2dca9054))
* set up standalone auth CI and image publishing ([#2](https://github.com/air-hand-org/jobdeck-auth/issues/2)) ([4cbb333](https://github.com/air-hand-org/jobdeck-auth/commit/4cbb33327ca9a037bbb8e875f7c20be1fcb87e74))

## Changelog

For changes before this service was split into its own repository, see the
[jobdeck changelog](https://github.com/air-hand-org/jobdeck/blob/main/CHANGELOG.md).

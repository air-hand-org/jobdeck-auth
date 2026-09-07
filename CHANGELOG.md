# Changelog

## [1.1.0](https://github.com/air-hand-org/jobdeck-auth/compare/v1.0.4...v1.1.0) (2026-09-07)


### Features

* add GCP bootstrap infrastructure ([#36](https://github.com/air-hand-org/jobdeck-auth/issues/36)) ([ff1bcf4](https://github.com/air-hand-org/jobdeck-auth/commit/ff1bcf45dfc5d45f57616636963fada37b1c67d6))
* add job worker execution loop ([#5](https://github.com/air-hand-org/jobdeck-auth/issues/5)) ([4a6aa36](https://github.com/air-hand-org/jobdeck-auth/commit/4a6aa36702c151cdfb67ab10407adb46a3ed33e5))
* add PostgreSQL runtime entrypoints ([#44](https://github.com/air-hand-org/jobdeck-auth/issues/44)) ([005d0e6](https://github.com/air-hand-org/jobdeck-auth/commit/005d0e680aa8c01985bfa206a0ebabc1d461932c))
* add release container images ([#40](https://github.com/air-hand-org/jobdeck-auth/issues/40)) ([67534af](https://github.com/air-hand-org/jobdeck-auth/commit/67534af342bcf8453766d7c6948e313fd8546a1c))
* apply bootstrap independently ([#78](https://github.com/air-hand-org/jobdeck-auth/issues/78)) ([6a7177d](https://github.com/air-hand-org/jobdeck-auth/commit/6a7177d0688946b142bf16ba20cc9b78f422877c))
* avoid ssrf by forward proxy ([#34](https://github.com/air-hand-org/jobdeck-auth/issues/34)) ([ac4b4e3](https://github.com/air-hand-org/jobdeck-auth/commit/ac4b4e3f212a1c63c0a6c9b70b6a82e459cd0b16))
* e2e test ([#32](https://github.com/air-hand-org/jobdeck-auth/issues/32)) ([bd94376](https://github.com/air-hand-org/jobdeck-auth/commit/bd94376d45a582e1e25ee2175837842919060f85))
* **iac:** deploy application services to Cloud Run ([#56](https://github.com/air-hand-org/jobdeck-auth/issues/56)) ([eee7735](https://github.com/air-hand-org/jobdeck-auth/commit/eee7735010bd90d6a37c1c310a96cb7c4fe3cef1))
* implement auth ([#22](https://github.com/air-hand-org/jobdeck-auth/issues/22)) ([9d6ccfe](https://github.com/air-hand-org/jobdeck-auth/commit/9d6ccfec2f5f3d9c8350b0846ae0852a62c9a010))
* implement execution ([#3](https://github.com/air-hand-org/jobdeck-auth/issues/3)) ([072feae](https://github.com/air-hand-org/jobdeck-auth/commit/072feae2e2f73b54491294d625c7f817e7dffb4d))
* retry job ([#8](https://github.com/air-hand-org/jobdeck-auth/issues/8)) ([616f719](https://github.com/air-hand-org/jobdeck-auth/commit/616f71952fa0b7db75770628cb3202e6d1fa38c4))
* UI/BFF ([#14](https://github.com/air-hand-org/jobdeck-auth/issues/14)) ([4a67f31](https://github.com/air-hand-org/jobdeck-auth/commit/4a67f3160d1bcc76a31029facfcc327b309066e3))
* use HTTPS for local service communication ([#51](https://github.com/air-hand-org/jobdeck-auth/issues/51)) ([ebaa38e](https://github.com/air-hand-org/jobdeck-auth/commit/ebaa38ef67b24891439b44badea8c73263b34093))


### Bug Fixes

* link protobuf other repo ([#96](https://github.com/air-hand-org/jobdeck-auth/issues/96)) ([6b7c1c1](https://github.com/air-hand-org/jobdeck-auth/commit/6b7c1c14e79ea66e9c00042777b9088fcb725cd2))
* optimize CI image builds ([#47](https://github.com/air-hand-org/jobdeck-auth/issues/47)) ([136e5d6](https://github.com/air-hand-org/jobdeck-auth/commit/136e5d6947a112e9c85fd9ffa8b675a00d836021))
* parallelize slow CI tests and publish releases ([#46](https://github.com/air-hand-org/jobdeck-auth/issues/46)) ([f37743e](https://github.com/air-hand-org/jobdeck-auth/commit/f37743ecdc6d0a34f6e4e65d8c7f44627d04f005))
* **release:** publish stable releases ([#54](https://github.com/air-hand-org/jobdeck-auth/issues/54)) ([93cc9e1](https://github.com/air-hand-org/jobdeck-auth/commit/93cc9e11014efeef69b858b17d595144ce49325c))
* remove wait bootstrap job ([#80](https://github.com/air-hand-org/jobdeck-auth/issues/80)) ([a332cbf](https://github.com/air-hand-org/jobdeck-auth/commit/a332cbfbb2f3fdd0da93108851c406035f348dda))
* restore CI dependency caches concurrently ([#49](https://github.com/air-hand-org/jobdeck-auth/issues/49)) ([e8c12a6](https://github.com/air-hand-org/jobdeck-auth/commit/e8c12a6f8abbdfa1dd2b242114e632a4c61687bc))
* run Go linters per module ([#98](https://github.com/air-hand-org/jobdeck-auth/issues/98)) ([00096dd](https://github.com/air-hand-org/jobdeck-auth/commit/00096dd2677cf9d4a53c3659a280eb493302c771))
* split update job ([#2](https://github.com/air-hand-org/jobdeck-auth/issues/2)) ([18dec6f](https://github.com/air-hand-org/jobdeck-auth/commit/18dec6fd4aa216e40e6ba8390d545faf6ae21ae1))

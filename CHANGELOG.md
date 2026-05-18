# [2.1.0](https://github.com/shivanshkc/squelette/compare/v2.0.0...v2.1.0) (2026-05-18)


### Features

* config validation ([cdb11ec](https://github.com/shivanshkc/squelette/commit/cdb11ece3930b559c30c787c1f173797e90f00ec))

# [2.0.0](https://github.com/shivanshkc/squelette/compare/v1.3.0...v2.0.0) (2026-03-17)


### Bug Fixes

* **ci:** correct Makefile container command ([96fdc89](https://github.com/shivanshkc/squelette/commit/96fdc89f25e8c059db29564678c4b777cc6dde12))
* **ci:** golangci lint version upgrade ([ecac786](https://github.com/shivanshkc/squelette/commit/ecac78651facf2c8de82980819da2e359981b5f2))
* **docs:** add more instructions to readme, update server accordingly ([41c3618](https://github.com/shivanshkc/squelette/commit/41c3618bd518c5d688577539bb23535341819739))
* **http:** add Hijacker implementation for websockets, http.Server refactoring ([7afe452](https://github.com/shivanshkc/squelette/commit/7afe452690a3dc1a03226efb16f0eb08b8371345))
* **http:** add http verb to the /api route ([c95ee1e](https://github.com/shivanshkc/squelette/commit/c95ee1e4829b16f10b84997afe7a36856ad030df))
* **http:** ensure shutdown in all cases, fix tests ([3bc961f](https://github.com/shivanshkc/squelette/commit/3bc961f0d36e58905328ffe2ab995621b20d34b7))


### Features

* **config:** simplify config management, no viper ([6ad0e28](https://github.com/shivanshkc/squelette/commit/6ad0e281c799af44f29b2b822e1f10a96cd345af))
* new framework introduced ([5f0778d](https://github.com/shivanshkc/squelette/commit/5f0778d05974a69adea33e88e19daa6b02aadb04))


### BREAKING CHANGES

* new framework
* **config:** yaml config is no longer supported

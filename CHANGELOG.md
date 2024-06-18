## [1.0.6](https://github.com/kilianpaquier/gitlab-storage-cleaner/compare/v1.0.5...v1.0.6) (2024-06-18)


### Chores

* **deps:** bump github.com/panjf2000/ants/v2 in the minor-patch group ([a4faaaf](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/a4faaafe190ff8d9e9330058d916d8d1b3b88fd7))

## [1.0.5](https://github.com/kilianpaquier/gitlab-storage-cleaner/compare/v1.0.4...v1.0.5) (2024-06-16)


### Chores

* **deps:** upgrade cobra and validator dependencies ([76907e3](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/76907e3862f0d8d127ae2b4bde344117d6427c23))
* **deps:** upgrade Dockerfile go version ([034f863](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/034f863b16c2c046df244dd67e0da292f3e3ba99))

## [1.0.4](https://github.com/kilianpaquier/gitlab-storage-cleaner/compare/v1.0.3...v1.0.4) (2024-06-01)


### Chores

* **deps:** upgrade validator dependency ([4ac8cf5](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/4ac8cf5c5fb76f6366dedd4f56958bb7a687f9bc))

## [1.0.3](https://github.com/kilianpaquier/gitlab-storage-cleaner/compare/v1.0.2...v1.0.3) (2024-05-13)


### Chores

* **deps:** bump github.com/xanzy/go-gitlab in the minor-patch group ([494939a](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/494939a2e3f150481eb27056c67878a7ce459db4))

## [1.0.2](https://github.com/kilianpaquier/gitlab-storage-cleaner/compare/v1.0.1...v1.0.2) (2024-05-08)


### Chores

* **deps:** bump github.com/xanzy/go-gitlab in the minor-patch group ([01ca109](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/01ca10973145a7c403da7ba24350a4bdefb9b866))
* **deps:** bump github.com/xanzy/go-gitlab in the minor-patch group ([b11e858](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/b11e8588382ed9f4ad49021158fb614dea27adbe))
* **deps:** upgrade toolchain to go1.22.3 ([10c37b6](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/10c37b6910538b42a27b1011ddc91512829d392f))

## [1.0.1](https://github.com/kilianpaquier/gitlab-storage-cleaner/compare/v1.0.0...v1.0.1) (2024-05-01)


### Bug Fixes

* **gitlab-ci:** bad template stage and run paths ([37085da](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/37085da922e87526e5ef4b9e8a32f91dec246d56))


### Chores

* **deps:** bump github.com/go-playground/validator/v10 ([1fa29d5](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/1fa29d5be2e366fd2eef3b9b94dc25d8997779e5))
* **deps:** bump github.com/xanzy/go-gitlab in the all group ([e951b12](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/e951b125ebbd0985ce35ef107f7e29fb9ec58f9d))
* **deps:** bump github.com/xanzy/go-gitlab in the all group ([9c7cd40](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/9c7cd405fb709d6b132424f8be867d3500a7c3de))
* **deps:** bump golangci/golangci-lint-action from 4 to 5 ([2e42b83](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/2e42b836a47a4bde997d88dd369360487c7221d2))
* **deps:** upgrade dependencies ([eac3d85](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/eac3d85140259131f6ecb86edcb031f366e9d9c8))
* **golangci:** remove govet deleted option ([3836ef3](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/3836ef342d60dc22124e9a0ff9226ca7fcbdc67d))
* **go:** update go to 1.22.2 ([4818496](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/4818496c8b4c78fd918ea583b424852e7c6d6d87))

## 1.0.0 (2024-03-30)


### Features

* **ci:** add codecov workflow ([6d47d8c](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/6d47d8c960a2bb17285952b10eea4c0a4332657d))
* **ci:** add release branches handling ([0e065d9](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/0e065d9b9a70d4913ce9fd8721e1a97dcf08a2aa))
* import project from gitlab ([8ff431e](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/8ff431e8faf8de6f217792f9edfe55b72a0597d8))


### Bug Fixes

* **ci:** add checkout and setup go for docker build job ([9779b2d](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/9779b2dc1f9dc53dc88f8e7a6a0e4380da4172ac))
* **ci:** bad image name for gitlab-ci template ([deb288e](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/deb288e9360177583831131810c2a17f10780aff))
* **ci:** codecov config in subdir really doesn't work ([707f11c](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/707f11c8d236d182a209f5697f7df9cb7d52ea73))
* **ci:** gitlab-ci template image name ([ebf38ac](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/ebf38ac3ab604d77e69644ee0079b01daf678712))
* **ci:** handle correctly dependabot codecov ignore ([530330b](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/530330bf9a033fdfb449591b7058a162a9e4d6a8))
* **golang:** invalid order in Dockerfile instructions ([ba5b93a](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/ba5b93aa6257c499e16c3b51aa923e88f9bd850a))
* **release:** add v prefix on github workflows version output ([a18128d](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/a18128da45cc01f6715c1ae61874f7d89b0ab069))
* **release:** missing operand to compute release boolean ([cbd4836](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/cbd4836f8984a1a252db241ff6c31edf028e5ec1))
* **release:** only add major and major.version docker images on release branches versions ([e4e992f](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/e4e992fa514b1c17bce56fc41f67df7d2f94892a))
* **release:** prefix present on expected exported boolean ([f8bab79](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/f8bab79dec2f6f2ba4679c34c237ed70fbca9839))
* **release:** release branches output not matching between boolean and string ([8ff62f1](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/8ff62f1273f6993c36278e3163f8cb958b7808f6))
* **release:** release branches output not matching between boolean and string ([221ec05](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/221ec05e0162f64821a81aa1955861defb209ea2))
* **release:** try avoiding same tag two times in docker hub ([99f943f](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/99f943f84551b9839db9c4c7e43bfb893abca137))
* **release:** v prefix present on release computing while version doesn't have it ([4b607fd](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/4b607fd3120ffc677288238da43be40e45507605))


### Reverts

* "chore: migrate docker image to ghcr.io" because of private ([1cee781](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/1cee78145eab7eab7d6bf2ac6a5aa4da0a4e0715))
* **ci:** "revert: "chore: migrate docker image to ghcr.io" because of private" ([025f9a6](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/025f9a6f4885f532c17a8098ba411bf9bdd477f6))


### Documentation

* **readme:** add commands ([95f5041](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/95f50419902e783c9bbd00e6234a979e66a83289))
* remove specific install command ([14bf3e7](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/14bf3e7d55e6e3974fc9202b3e5f3dfb31315786))


### Chores

* **ci:** add docker-hadolint and docker-trivy analysis ([3dbd584](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/3dbd584a6f1c4d3d8565f9c7fae4ad6369d1cde0))
* **ci:** add strategy execution for tests ([ba612e6](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/ba612e6557ffd99b1f8485b98c1a2c7bcbbc9990))
* **ci:** regenerate layout ([de43b27](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/de43b278baa57b204099f2b61838210adc7dc112))
* **ci:** remove build/ci directory ([e6cefac](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/e6cefac9b1e7d32a0c94b05510bf60107f87a0b9))
* **ci:** update golangci rules ([49539d0](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/49539d0d3356dc7887ce3f5203571b0460731af6))
* **deps:** add dependabot ([cb15db6](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/cb15db696be9edecf0f8e7368f1e259255e5b620))
* **deps:** update dependencies ([fbb7cac](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/fbb7cac5b950eed08e311781fdf20a193530bf71))
* migrate docker image to ghcr.io ([e1b0114](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/e1b01148819ea5947b48b41628a4a71a7d1837c4))
* **release:** v1.0.0-alpha.1 [skip ci] ([f0b2af6](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/f0b2af6a2cd728afb391b84c0447a6063a2b57f3))
* **release:** v1.0.0-alpha.2 [skip ci] ([d53ac45](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/d53ac4575b456cac4c02d1cd7998bbc72c75d016))
* **release:** v1.0.0-alpha.3 [skip ci] ([a9bcecc](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/a9bcecc0c9ff39618de15b84b1b79dbfc7b633f4))

## [1.0.0-alpha.3](https://github.com/kilianpaquier/gitlab-storage-cleaner/compare/v1.0.0-alpha.2...v1.0.0-alpha.3) (2024-03-03)


### Bug Fixes

* **ci:** bad image name for gitlab-ci template ([deb288e](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/deb288e9360177583831131810c2a17f10780aff))
* **ci:** gitlab-ci template image name ([ebf38ac](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/ebf38ac3ab604d77e69644ee0079b01daf678712))


### Reverts

* "chore: migrate docker image to ghcr.io" because of private ([1cee781](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/1cee78145eab7eab7d6bf2ac6a5aa4da0a4e0715))
* **ci:** "revert: "chore: migrate docker image to ghcr.io" because of private" ([025f9a6](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/025f9a6f4885f532c17a8098ba411bf9bdd477f6))

## [1.0.0-alpha.2](https://github.com/kilianpaquier/gitlab-storage-cleaner/compare/v1.0.0-alpha.1...v1.0.0-alpha.2) (2024-03-03)


### Chores

* **ci:** add strategy execution for tests ([ba612e6](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/ba612e6557ffd99b1f8485b98c1a2c7bcbbc9990))
* migrate docker image to ghcr.io ([e1b0114](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/e1b01148819ea5947b48b41628a4a71a7d1837c4))

## 1.0.0-alpha.1 (2024-03-03)


### Features

* **ci:** add codecov workflow ([6d47d8c](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/6d47d8c960a2bb17285952b10eea4c0a4332657d))
* **ci:** add release branches handling ([0e065d9](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/0e065d9b9a70d4913ce9fd8721e1a97dcf08a2aa))
* import project from gitlab ([8ff431e](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/8ff431e8faf8de6f217792f9edfe55b72a0597d8))


### Bug Fixes

* **ci:** add checkout and setup go for docker build job ([9779b2d](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/9779b2dc1f9dc53dc88f8e7a6a0e4380da4172ac))
* **golang:** invalid order in Dockerfile instructions ([ba5b93a](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/ba5b93aa6257c499e16c3b51aa923e88f9bd850a))
* **release:** add v prefix on github workflows version output ([a18128d](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/a18128da45cc01f6715c1ae61874f7d89b0ab069))
* **release:** missing operand to compute release boolean ([cbd4836](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/cbd4836f8984a1a252db241ff6c31edf028e5ec1))
* **release:** only add major and major.version docker images on release branches versions ([e4e992f](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/e4e992fa514b1c17bce56fc41f67df7d2f94892a))
* **release:** prefix present on expected exported boolean ([f8bab79](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/f8bab79dec2f6f2ba4679c34c237ed70fbca9839))
* **release:** release branches output not matching between boolean and string ([8ff62f1](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/8ff62f1273f6993c36278e3163f8cb958b7808f6))
* **release:** release branches output not matching between boolean and string ([221ec05](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/221ec05e0162f64821a81aa1955861defb209ea2))
* **release:** try avoiding same tag two times in docker hub ([99f943f](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/99f943f84551b9839db9c4c7e43bfb893abca137))
* **release:** v prefix present on release computing while version doesn't have it ([4b607fd](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/4b607fd3120ffc677288238da43be40e45507605))


### Documentation

* remove specific install command ([14bf3e7](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/14bf3e7d55e6e3974fc9202b3e5f3dfb31315786))


### Chores

* **deps:** add dependabot ([cb15db6](https://github.com/kilianpaquier/gitlab-storage-cleaner/commit/cb15db696be9edecf0f8e7368f1e259255e5b620))

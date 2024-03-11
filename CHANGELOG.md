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

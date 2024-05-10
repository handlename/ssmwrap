# Changelog

## [v2.0.1](https://github.com/handlename/ssmwrap/compare/v2.0.0...v2.0.1) - 2024-05-10
- go import path for v2 by @handlename in https://github.com/handlename/ssmwrap/pull/76
- Update license URL to not depend on branch name by @handlename in https://github.com/handlename/ssmwrap/pull/78
- chore(deps): bump github.com/aws/aws-sdk-go-v2/config from 1.27.9 to 1.27.12 by @dependabot in https://github.com/handlename/ssmwrap/pull/75

## [v2.0.0](https://github.com/handlename/ssmwrap/compare/v1.2.2...v2.0.0) - 2024-05-02
- migrate to aws-sdk-go-v2 by @handlename in https://github.com/handlename/ssmwrap/pull/61
- test parsing flags by @handlename in https://github.com/handlename/ssmwrap/pull/63
- io/ioutil is deprecated by @handlename in https://github.com/handlename/ssmwrap/pull/67
- Handle SIGINT by @handlename in https://github.com/handlename/ssmwrap/pull/68
- Group dependabot PRs by @handlename in https://github.com/handlename/ssmwrap/pull/69
- Reorganize flags by @handlename in https://github.com/handlename/ssmwrap/pull/70
- Sort parameters on test by @handlename in https://github.com/handlename/ssmwrap/pull/73

## [v1.2.1](https://github.com/handlename/ssmwrap/compare/v1.2.1...v1.2.1) - 2024-03-22

## 1.2.0 (2022-03-03)

- build with go v1.17.7
- add `-env-entire-path` option #51 #53 #55
- add release binaries for arm64
- add lisence file #45
- update package aws/aws-sdk-go

## 1.1.1 (2020-05-11)

- fix panic without command #30
- fix useless export for versioned name #29
- update dependencies #28

## 1.1.0 (2020-05-11)

- release from GitHub Actions. there are no changes for ssmwrap itself #31

## 1.0.3 (2020-01-14)

- ssmwrap reads sharde config file (~/.aws/config) #24
- update dependencies #23

## 1.0.2 (2019-12-26)

- now, -file option enabled without -path/-names #21 #22
- update dependencies #20

## 1.0.1 (2019-10-25)

- update dependencies #15 #17
- build with go 1.13

## 1.0.0 (2019-03-04)

- add `-names` option #14
- remove public function `FetchParameters` and `FetchParametersByNames`

## 0.7.0 (2019-02-12)

- returns exit code 1 when error occurred #13

## 0.6.0 (2019-01-18)

- add `-file` option #9
- add `-recursive`/`-no-recursive` options #12

## 0.5.0 (2018-09-06)

- add `-env`/`-no-env`/`-env-prefix` options #5
- add library interface `Export` #6

## 0.4.0 (2018-07-13)

- configurations via environment variables #4

## 0.3.1 (2018-07-11)

- build without cgo. 0.2.1 is not worked...

## 0.3.0 (2018-07-04)

- added -retries flag #3

## 0.2.1 (2018-06-28)

- build without cgo

## 0.2.0 (2018-06-27)

- ssm parameters takes precedence over the current environment variables

## 0.1.0 (2018-06-25)

- First release

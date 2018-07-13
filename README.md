# ssmwrap

[![CircleCI](https://circleci.com/gh/handlename/ssmwrap.svg?style=svg)](https://circleci.com/gh/handlename/ssmwrap)

ssmwrap execute commands with environment variables loaded from AWS SSM Parameter Store.

## Usage

```console
$ ssmwrap -paths /production -- your_command
```

## Install

```console
$ go get github.com/handlename/ssmwrap/cmd/ssmwrap
```

## Options

```console
$ ssmwrap -help
Usage of ./cmd/ssmwrap/ssmwrap:
  -paths string
    	comma separated parameter paths (default "/")
  -prefix string
    	prefix for environment variables
  -retries int
    	number of times of retry
  -version
    	display version
```

### Environment Variables

All of command line options can be set via environment variables.

```console
$ SSMWRAP_PATHS=/production ssmwrap ...
```

means,

```console
$ ssmwrap -paths /production ...
```

If there are command line options too, these takes priority.

```console
$ SSMWRAP_PREFIX=FOO_ ssmwrap -prefix BAR_ env
BAR_SOMETHING=...
```

## Motivation

There are some tools to use values stored in AWS System Manager Parameter Store,
but I couldn't find that manipulate values including newline characters correctly.

ssmwrap runs your command through syscall.Exec, not via shell,
so newline characters are treated as part of a environment value.

## Licence

MIT

## Special Thanks

@fujiwara has gave me an idea of ssmwrap.

## Author

@handlename (https://github.com/handlename)

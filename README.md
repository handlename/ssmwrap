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

## Motivation

There are some tools to use values stored in AWS System Manager Parameter Store,
but I couldn't find that manipulate values including newline characters correctly.

ssmwrap runs your command through syscall.Exec, not via shell,
so newline characters are treated as part of a environment value.

## Licence

MIT

## Author

@handlename (https://github.com/handlename)

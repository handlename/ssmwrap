# ssmwrap

[![CircleCI](https://circleci.com/gh/handlename/ssmwrap.svg?style=svg)](https://circleci.com/gh/handlename/ssmwrap) [![Documentation](https://godoc.org/github.com/handlename/ssmwrap?status.svg)](https://godoc.org/github.com/handlename/ssmwrap)

ssmwrap execute commands with output values loaded from AWS SSM Parameter Store to somewhere.

Supported output targets:

- environment variables

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
  -env
    	export values as environment variables (default true)
  -env-prefix string
    	prefix for environment variables
  -file value
    	write values as file
    	format:  Name=VALUE_NAME,Path=FILE_PATH,Mode=FILE_MODE,Gid=FILE_GROUP_ID,Uid=FILE_USER_ID
    	example: Name=/foo/bar,Path=/etc/bar,Mode=600,Gid=123,Uid=456
  -names string
        comma separated parameter names
  -no-env
    	disable export to environment variables
  -no-recursive
    	retrieve values just under -paths only
  -paths string
        comma separated parameter paths
  -prefix string
    	alias for -env-prefix
  -recursive
    	retrieve values recursively (default true)
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

## Usage as a library

`ssmwrap.Export()` fetches parameters from SSM and export those to envrionment variables.

```go
err := ssmwrap.Export(ssmwrap.ExportOptions{
	Paths: []string{"/path/"},
	Retries: 3,
	Prefix: "SSM_",
})
if err != nil {
	// ...
}
foo := os.Getenv("SSM_FOO")  // a value of /path/foo in SSM
```

## Licence

MIT

## Special Thanks

@fujiwara has gave me an idea of ssmwrap.

## Author

@handlename (https://github.com/handlename)

# ssmwrap

[![Documentation](https://godoc.org/github.com/handlename/ssmwrap?status.svg)](https://godoc.org/github.com/handlename/ssmwrap)

ssmwrap execute commands with output values loaded from AWS SSM Parameter Store to somewhere.

Supported output targets:

- environment variables
- files

## Usage

```console
$ ssmwrap \
	-env 'path=/production/*' \
	-file 'path=/production/ssl_cert,to=/etc/ssl/cert.pem,mode=0600' \
	-- app
```

## Install

Download binary from [releases](https://github.com/handlename/ssmwrap/releases)

or

```console
$ brew install handlename/tap/ssmwrap
```

or

```console
$ go install github.com/handlename/ssmwrap/cmd/ssmwrap/v2@latest
```

## Options

```console
$ ssmwrap -help
Usage of ssmwrap:
  -env rule
    	Alias of rule flag with `type=env`.
  -file rule
    	Alias of rule flag with `type=file`.
  -retries int
    	Number of times of retry. Default is 0
  -rule path
    	Set rule for exporting values. multiple flags are allowed.
    	format: path=...,type={env,file}[,to=...][,entirepath={true,false}][,prefix=...][,mode=...][,gid=...][,uid=...]
    	parameters:
    	        path: [required]
    	              Path of parameter store.
    	              If path ends with no-slash character, only the value of the path will be exported.
    	              If `path` ends with `/**/*`, all values under the path will be exported.
    	              If `path` ends with `/*`, only top level values under the path will be exported.
    	        type: [required]
    	              Destination type. `env` or `file`.
    	          to: [required for `type=file`]
    	              Destination path.
    	              If `type=env`, `to` is name of exported environment variable.
    	              If `type=env`, but `to` is not set, `path` will be used as name of exported environment variable.
    	              If `type=file`, `to` is path of file to write.
    	  entirepath: [optional, only for `type=env`]
    	              Export entire path as environment variables name.
    	              If `entirepath=true`, all values under the path will be exported. (/path/to/param -> PATH_TO_PARAM)
    	              If `entirepath=false`, only top level values under the path will be exported. (/path/to/param -> PARAM)
    	      prefix: [optional, only for `type=env`]
    	              Prefix for exported environment variable.
    	        mode: [optional, only for `type=file`]
    	              File mode. Default is 0644.
    	         gid: [optional, only for `type=file`]
    	              Group ID of file. Default is current user's Gid.
    	         uid: [optional, only for `type=file`]
    	              User ID of file. Default is current user's Uid.
  -version
    	Display version and exit
```

### Environment Variables

All of command line options can be set via environment variables.

```console
$ SSMWRAP_ENV='path=/production/*' ssmwrap ...
```

is same as,

```console
$ ssmwrap -env 'path=/production/*' ...
```

You can set multiple options by add suffix like '_1', '_2', '_3'...

```console
$ SSMWRAP_ENV_1='path=/production/app/*' SSMWRAP_ENV_2='path=/production/db/*' ssmwrap ...
```

## Migration from v1.x to v2.x

On v2, options flags are reformed.

### Output to environment variables

Flags for output to environment variables are consolidated to `-env` flag.

```conosle
# v1
$ ssmwrap \
	-paths '/foo,/bar' \
	-env-entire-path \
	-- ...

# v2
$ ssmwrap \
	-env 'path=/foo/*,entirepath=true' \
	-env 'path=/bar/*,entirepath=true' \
	-- ...
```

### Output to files

Flags for output to files are remaining as `-file` flag, but format is changed.

```conosle
# v1
$ ssmwrap -file 'Path=/foo/value,Name=/path/to/file,Mode=0600' -- ...

# v2
$ ssmwrap -file 'path=/foo/value,to=/path/to/file,mode=0600' -- ...
```

### General output rules

Added new flag `-rule` that can be used for all type of output.
Flag `-env` and `-file` are alias of `-rule` flag.

```conosle
# by -env and -file flag
$ ssmwrap \
	-env 'path=/foo/*' \
	-file 'path=/bar/value,to=/path/to/file' \
	-- ...

# by -rule flag (same as above)
$ ssmwrap \
	-rule 'type=env,path=/foo/*' \
	-rule 'type=file,path=/bar/value,to=/path/to/file' \
	-- ...
```

## Motivation

There are some tools to use values stored in AWS System Manager Parameter Store,
but I couldn't find that manipulate values including newline characters correctly.

ssmwrap runs your command through syscall.Exec, not via shell,
so newline characters are treated as part of a environment value.

## Usage as a library

`ssmwrap.Export()` fetches parameters from SSM and export those to envrionment variables.
Please check [example](./examples/lib/main.go).

## License

see [LICENSE](https://github.com/handlename/ssmwrap?tab=MIT-1-ov-file#readme) file.

## Special Thanks

@fujiwara has gave me an idea of ssmwrap.

## Author

@handlename (https://github.com/handlename)

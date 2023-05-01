# regexpp-go

A regular expression parser for latest ECMAScript written in Go.

Work in progress...

## Run parser tests

Compare snapshots and actually results:

```sh
go test ./internal/parser
```

Updates snapshots:

```sh
UPDATE=true go test ./internal/parser
```

Run for only one fixture:

```sh
TARGET=alternative1 go test ./internal/parser
```

## Prior art

- https://github.com/mysticatea/regexpp

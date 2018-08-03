# jg - A command-line JSON processor in Go

## Install

_`go` 1.6 or later is needed for install from source code._

```bash
go get -u github.com/azusa0127/jg
go install github.com/azusa0127/jg
```

## Example

```bash
# formatting
$ echo '{"abc":123,"cde":"foo"}' | jg
{
  "abc": 123,
  "cde": "foo"
}

# jsonpath
$ echo '{"abc":123,"cde":"foo"}' | jg $.abc
123

# non pipeing
$ jg '{\"abc\":123,\"cde\":\"foo\"}' $.abc
123
```

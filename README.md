# jg - A commandline json processor in Go

## Install

```bash
go get -u github.com/azusa0127/jg
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

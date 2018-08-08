# gsp - a general stream processor for command-line

A tiny tool to focus line-delimited stream processing on CLI.

## Install

_`go` 1.6 or later is needed for install from source code._

```bash
go get -u github.com/azusa0127/gsp
go install github.com/azusa0127/gsp
```

## Usage

### Prettify JSON

```bash
echo '{"n":123,"t":"foo"}\n{"n":789,"t":"bar"}' | gsp
```

### Query JSON with [JSMEPath](http://jmespath.org/)

```bash
echo '{"n":123,"t":"foo"}\n{"n":789,"t":"bar"}' | gsp n
```

```bash
gsp -q=n '{"n":123,"t":"foo"}'
```

```bash
gsp '{"n":123,"t":"foo"}' n
```

### Query JSON with [JSONPath](http://goessner.net/articles/JsonPath/index.html)

```bash
echo '{"n":123,"t":"foo"}\n{"n":789,"t":"bar"}' | gsp -qe jsonpath $.n
```

### Encode with Base64

```bash
echo '{"n":123,"t":"foo"}\n{"n":789,"t":"bar"}' | gsp -base64e
```

### Decode with Base64

```bash
echo '{"n":123,"t":"foo"}\n{"n":789,"t":"bar"}' | gsp -base64e | gsp -base64d
```

### Convert JSON to YAML

```bash
cat sample.json | gsp -o yaml
```

```bash
cat sample.json | gsp -toyaml
```

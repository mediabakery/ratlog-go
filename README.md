# :rat: ratlog-go

Application Logging for Rats, Humans and Machines

[![Go Report Card](https://goreportcard.com/badge/github.com/mediabakery/ratlog-go?style=flat-square)](https://goreportcard.com/report/github.com/mediabakery/ratlog-go)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/mediabakery/ratlog-go)
[![Release](https://img.shields.io/github/release/mediabakery/ratlog-go.svg?style=flat-square)](https://github.com/mediabakery/ratlog-go/releases/latest)


golang implementation of the [ratlog spec](https://github.com/ratlog/ratlog-spec)

There are two way you can use ratlog.

```golang
var b bytes.Buffer
rlog := ratlog.New(&b)
rlog.Log(Props{message: test.Data.Message, tags: test.Data.Tags, fields: fields})
```

```golang
var b bytes.Buffer
rlog := ratlog.New(&b)
rlog.Message("Hello World").Fields("a": "1", "b": "2").Tag("debug", "error").Log()
```

You can provide a log instance with an additional tag

```golang
var b bytes.Buffer
rlog := ratlog.New(&b)
errorLog := rlog.Tag("error")
errorLog.Message("Hello World").Fields("a": "2", "b": "3").Tag("bla", "baz").Log()
// [error|debug|bla|baz] Hello World | a: 2 | b: 3
```

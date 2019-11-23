# :rat: ratlog-go

Application Logging for Rats, Humans and Machines

golang implementation of the [ratlog spec](https://github.com/ratlog/ratlog-spec)

There are two way you can use ratlog.

```golang
var b bytes.Buffer
ratlog := Ratlog{writer: &b}
ratlog.Log(Props{message: test.Data.Message, tags: test.Data.Tags, fields: fields})
```

```golang
var b bytes.Buffer
ratlog := Ratlog{writer: &b}
ratlog.Message("Hello World").Fields("a": "1", "b": "2").Tag("debug", "error").Log()
```

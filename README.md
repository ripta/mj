# MJ

`Make JSON` from the command line.

The program takes a list of key=value parameter on the command line:

```shell
$ mj foo=bar
{"foo":"bar"}

$ mj foo=bar baz=quux
{"baz":"quux","foo":"bar"}
```

The only requirement being that the keys be unique:

```shell
$ mj a=b a=c
Error: Key path was already assigned
```

A key can be a simple string, or a .-separated path of arbitrary depth:

```shell
$ mj foo.bar=baz
{"foo":{"bar":"baz"}}
```

There are command-line options to select a different separator for key-value pairs and for key paths:

```shell
$ mj foo:bar=baz
{"foo:bar":"baz"}

$ mj -s=: foo:bar=baz
{"foo":"bar=baz"}

$ mj -p=: foo:bar=baz
{"foo":{"bar":"baz"}}

$ mj -p='->' 'foo->bar=baz'
{"foo":{"bar":"baz"}}
```

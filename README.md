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

There are command-line options to select a different separator for key-value
pairs and for key paths:

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

If a key starts with `-`, it will be interpreted as a command line flag. To
prevent that, use `--`:

```shell
$ mj -really=why
Error: flag provided but not defined: -really

$ mj -- -really=why
{"-really":"why"}
```

There is also support for slices on the last level:

```shell
$ mj foo[]=abc foo[]=def
{"foo":["abc","def"]}
```

But the operation is not supported on deeply-nested objects yet:

```shell
$ mj foo[].bar=abc foo[].bar=def
mj: encountered error while processing argument #0: "foo[].bar=abc"
	underlying error: while processing key path [foo[] bar]: cannot set key "bar" to "abc" on []interface {}: not supported
```

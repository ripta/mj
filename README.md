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

And be careful that numeric keys are _not_ interpreted as array indices:

```shell
$ mj foo.0=bar foo.1=baz
{"foo":{"0":"bar","1":"baz"}}

$ mj foo.0.bar=baz
{"foo":{"0":{"bar":"baz"}}}
```

By default, values are interpreted as-is:

```shell
$ mj foo=@bar.txt
{"foo":"@bar.txt"}
```

However, you can designate certain prefixes to read values from a file,
using the `-r` flag. The default prefix is blank, which disables file reading,
as shown above.

```shell
$ echo hello-world > bar.txt
$ mj -r=@ foo=@bar.txt
{"foo":"hello-world\n"}
```

By default, all values are treated as strings. Use type suffixes to create other JSON primitives:

- `:string` - String value (default, can be omitted)
- `:int` - Integer number
- `:float` - Floating-point number
- `:bool` - Boolean (`true`/`false`, `1`/`0`)
- `:null` - Null value (value part ignored)

```shell
# Numbers
$ mj age:int=25
{"age":25}

$ mj price:float=19.99
{"price":19.99}

$ mj temp:int=-10
{"temp":-10}

# Booleans
$ mj active:bool=true
{"active":true}

$ mj deleted:bool=false
{"deleted":false}

# Null
$ mj value:null=
{"value":null}

# Mixed types
$ mj name=Alice age:int=30 premium:bool=true
{"age":30,"name":"Alice","premium":true}

# With nested objects
$ mj user.name=Bob user.age:int=25 user.verified:bool=true
{"user":{"age":25,"name":"Bob","verified":true}}

# Arrays with types
$ mj scores:int[]=95 scores:int[]=87 scores:int[]=92
{"scores":[95,87,92]}
```

Use the `-t` flag to change the type separator (default `:`):

```shell
$ mj -t=/ age/int=25
{"age":25}
```

This is useful when your keys contain colons:

```shell
$ mj -t=/ "url:port"=8080 "url:host"=localhost
{"url:host":"localhost","url:port":"8080"}

$ mj -t=/ url:port/int=8080 url:host=localhost
{"url:host":"localhost","url:port":8080}
```

Invalid values for typed fields will produce errors:

```shell
$ mj age:int=not-a-number
Error: cannot parse "not-a-number" as int

$ mj active:bool=maybe
Error: cannot parse "maybe" as bool
```


# `yo` 

[![Go Report Card](https://goreportcard.com/badge/github.com/lucasepe/yo?style=flat-square)](https://goreportcard.com/report/github.com/lucasepe/yo) &nbsp;&nbsp;&nbsp; [![Release](https://img.shields.io/github/release/lucasepe/yo.svg?style=flat-square)](https://github.com/lucasepe/yo/releases/latest)

> An alternative syntax to generate YAML (or JSON) from commandline.

The ultimate commanline YAML (or JSON) generator! ... I'm kidding of course! but I'd like to know what you think.

https://youtu.be/QL6DsCLFQ30

## Usage

```sh
$ yo 'apiVersion="example.lucasepe.it/v1alpha1" kind=Project metadata={ namespace=default name=example-project } spec.replicas=1'
apiVersion: example.lucasepe.it/v1alpha1
kind: Project
metadata:
  name: example-project
  namespace: default
spec:
  replicas: 1
```

...or you can use piping (using `echo` here just to show the syntax, but you can use `cat` too):

```sh
$ echo 'apiVersion="example.lucasepe.it/v1alpha1" kind=Project metadata={ namespace=default name=example-project } spec.replicas=1' | yo
apiVersion: example.lucasepe.it/v1alpha1
kind: Project
metadata:
  name: example-project
  namespace: default
spec:
  replicas: 1
```

...or you can try the interactive mode:

```sh
$ yo 
[hit CTRL+d to finish]
>> apiVersion = "example.lucasepe.it/v1alpha1"
>> kind = Project
>> metadata.namespace = default
>> metadata.name = example-project
>> spec.replicas = 1

apiVersion: example.lucasepe.it/v1alpha1
kind: Project
metadata:
  name: example-project
  namespace: default
spec:
  replicas: 1
```

👉 interactive mode supports history! (try hitting the up arrow ⬆️).

# Syntax Overview

- a field is a key/value pair
- curly braces hold objects
- square brackets hold arrays

## fields

> A field is defined by: `IDENTIFIER = VALUE` .

- field key/value pairs have a equal `=` between them as in `key = value` 
- each field is eventually separated by space (zero, one or more does not matter)

```sh
$ yo firstName = Scarlett lastName = Johansson
```

generates...

```yaml
firstName: Scarlett
lastName: Johansson
```

- booleans, integeres, floating numbers are automatically resolved
- put the text beween quotes `"` to enter spaces and others unicode chars
  - es. `proverb = "interface{} says nothing"`

```sh
$ yo fullName="Scarlett Johansson" age=36 hot=true
```

generates...

```yaml
age: 36
fullName: Scarlett
hot: true
```

## objects

> An object is defined by: `IDENTIFIER = { fields... }`.

- begin a new object using the left curly brace `{`
- close the object with a right curly brace `}`

```sh
$ yo user = { name=foo age=30 active=true address = { zip="123" country=IT } }
```

generates...

```yaml
user:
  active: true
  address:
    country: IT
    zip: 123
  age: 30
  name: foo
```

- you can also use dotted notation (and/or eventually mix things!)

```sh
$ yo user = { name=foo age=30 active=true address.zip="123" address.country=IT }
```

```sh
$ yo user.name=foo user.age=30 user.active=true user.address = {zip="123" country=IT}
```

```sh
$ yo user.name=foo user.age=30 user.active=true user.address.zip="123" user.address.country=IT
```

All the previous examples produce the same result...it's up to you to find your way.


## arrays

> An array is defined by: `IDENTIFIER = [ fields...]`.

- begin a new array using the left square brace `[`
- end the array with a right quare brace `]`

```sh
$ yo tags = [ foo bar qix ]
```

```yaml
tags:
- foo
- bar
- qix
```

You can create an array of object too:

```sh
$ yo pets = [ { name=Dash kind=cat age=3 } {name=Harley kind=dog age=4} ]
```

```yaml
pets:
- age: 3
  kind: cat
  name: Dash
- age: 4
  kind: dog
  name: Harley
```


# How to install?

In order to use the `yo` command, compile it using the following command:

```bash
go get -u github.com/lucasepe/yo
```

This will create the executable under your `$GOPATH/bin` directory.

## Ready-To-Use Releases 

If you don't want to compile the sourcecode yourself, [Here you can find the tool already compiled](https://github.com/lucasepe/yo/releases/latest) for:

- MacOS
- Linux
- Windows


# `yo` 

[![Go Report Card](https://goreportcard.com/badge/github.com/lucasepe/yo?style=flat-square)](https://goreportcard.com/report/github.com/lucasepe/yo) &nbsp;&nbsp;&nbsp; [![Release](https://img.shields.io/github/release/lucasepe/yo.svg?style=flat-square)](https://github.com/lucasepe/yo/releases/latest)

> An alternative syntax to generate YAML (or JSON) from commandline.

The ultimate commanline YAML (or JSON) generator! ... I'm kidding of course! but I'd like to know what you think.

## Usage

### Interactive mode

```sh
$ yo eval
[hit CTRL+d to finish]
>> apiVersion=v1
>> kind=Secret
>> metadata.name=mysecret
>> type=Opaque
>> data.username=(b64enc "USER")
>> data.password=(b64enc "PASS")
```

```yaml
apiVersion: example.lucasepe.it/v1alpha1
kind: Project
metadata:
  name: example-project
  namespace: default
spec:
  replicas: 1
```

👉 interactive mode supports history! (try hitting the up arrow ⬆️).

## With piping (example using `cat`):

```sh
$ cat testdata/sample1.yo | ./yo eval
apiVersion: v1
data:
  password: UEFTUw==
  username: VVNFUg==
kind: Secret
metadata:
  name: mysecret
type: Opaque
```


# Syntax Overview

- a field is a key/value pair
- curly braces hold objects
- square brackets hold arrays

## fields

> A field is defined by: `IDENTIFIER = VALUE` .

- field key/value pairs have a equal `=` between them as in `key = value` 
- each field is eventually separated by space (zero, one or more does not matter)

```sh
$ yo eval 'firstName=Scarlett lastName=Johansson'
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
$ yo eval 'fullName="Scarlett Johansson" hot=true'
```

generates...

```yaml
fullName: Scarlett Johansson
hot: true
```

## objects

> An object is defined by: `IDENTIFIER = { fields... }`.

- begin a new object using the left curly brace `{`
- close the object with a right curly brace `}`

```sh
$ yo eval 'user = { name=foo age=30 active=true address = { zip="123" country=IT } }'
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
$ yo eval 'user = { name=foo age=30 active=true address.zip="123" address.country=IT }'
```

```sh
$ yo eval 'user.name=foo user.age=30 user.active=true user.address = {zip="123" country=IT}'
```

```sh
$ yo eval 'user.name=foo user.age=30 user.active=true user.address.zip="123" user.address.country=IT'
```

All the previous examples produce the same result...it's up to you to find your way.


## arrays

> An array is defined by: `IDENTIFIER = [ fields...]`.

- begin a new array using the left square brace `[`
- end the array with a right quare brace `]`

```sh
$ yo eval 'tags = [ foo bar qix ]'
```

```yaml
tags:
- foo
- bar
- qix
```

You can create an array of object too:

```sh
$ yo eval 'pets = [ { name=Dash kind=cat age=3 } {name=Harley kind=dog age=4} ]'
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

# Built-in functions

`yo` has also built-in handy functions

```sh
$ yo funcs
+--------------+-------------------------------------------------------------------------+
| FUNCTION     | SUMMARY                                                                 |
+--------------+-------------------------------------------------------------------------+
| b64dec       | Decode with Base64                                                      |
|              |                                                                         |
|              | >> yo eval 'key = ( b64dec "Hello world!" )'                            |
+--------------+-------------------------------------------------------------------------+
| b64enc       | Encode or decode with Base64                                            |
|              |                                                                         |
|              | >> yo eval 'key = ( b64enc "Hello world!" )'                            |
+--------------+-------------------------------------------------------------------------+
| coalesce     | Takes a list of values and returns the first non-empty one.             |
|              |                                                                         |
|              | >> yo eval 'key = ( coalesce .name .parent.name "Scarlett" )'           |
+--------------+-------------------------------------------------------------------------+
| contains     | Test to see if one string is contained inside of another.               |
|              |                                                                         |
|              | >> yo eval 'key = ( contains "cat" "catch" )'                           |
+--------------+-------------------------------------------------------------------------+
| date         | Formats a date.                                                         |
|              |                                                                         |
|              | >> yo eval 'key = ( now | date "2006-01-02" )'                          |
+--------------+-------------------------------------------------------------------------+
| dateInZone   | Formats a date with a timezone.                                         |
|              |                                                                         |
|              | >> yo eval 'key = ( dateInZone "2006-01-02" (now) "UTC" )'              |
+--------------+-------------------------------------------------------------------------+
| dateModify   | Modify a date.                                                          |
|              |                                                                         |
|              | >> yo eval 'key = ( now | dateModify "-1.5h" )'                         |
+--------------+-------------------------------------------------------------------------+
| decr         | Decrement an integer value by one.                                      |
|              |                                                                         |
|              | >> yo eval 'key = ( decr 9 )'                                           |
+--------------+-------------------------------------------------------------------------+
| default      | Set a simple default value.                                             |
|              |                                                                         |
|              | >> yo eval 'key = ( default "foo" .Bar )'                               |
+--------------+-------------------------------------------------------------------------+
| duration     | Formats seconds as a time.Duration.                                     |
|              |                                                                         |
|              | >> yo eval 'key = ( duration 95 )'                                      |
+--------------+-------------------------------------------------------------------------+
| env          | Reads an environment variable.                                          |
|              |                                                                         |
|              | >> yo eval 'key = ( env "HOME" )'                                       |
+--------------+-------------------------------------------------------------------------+
| expandenv    | Substitute environment variables in a string.                           |
|              |                                                                         |
|              | >> yo eval 'key = ( expandenv "Your path is set to $PATH" )'            |
+--------------+-------------------------------------------------------------------------+
| has          | Test to see if a list has a particular element.                         |
|              |                                                                         |
|              | >> yo eval 'key = ( has 4 $myList )'                                    |
+--------------+-------------------------------------------------------------------------+
| hasPrefix    | Test whether a string has a given prefix.                               |
|              |                                                                         |
|              | >> yo eval 'key = ( hasPrefix "cat" "catch" )'                          |
+--------------+-------------------------------------------------------------------------+
| hasSuffix    | Test whether a string has a given suffix.                               |
|              |                                                                         |
|              | >> yo eval 'key = ( hasSuffix "tch" "catch" )'                          |
+--------------+-------------------------------------------------------------------------+
| incr         | Increment an integer value by one.                                      |
|              |                                                                         |
|              | >> yo eval 'key = ( incr 7 )'                                           |
+--------------+-------------------------------------------------------------------------+
| join         | Join a list of strings into a single string, with the given separator.  |
|              |                                                                         |
|              | >> yo eval 'key = ( join "-" .Names )'                                  |
+--------------+-------------------------------------------------------------------------+
| lower        | Convert the entire string to lowercase.                                 |
|              |                                                                         |
|              | >> yo eval 'key = ( lower "HELLO" )'                                    |
+--------------+-------------------------------------------------------------------------+
| now          | The current date/time.                                                  |
|              |                                                                         |
|              | >> yo eval 'key = ( now | date "2006-01-02" )'                          |
+--------------+-------------------------------------------------------------------------+
| pathListSep  | Returns OS-specific path list separator.                                |
|              |                                                                         |
|              | >> yo eval 'key = ( pathSep )'                                          |
+--------------+-------------------------------------------------------------------------+
| pathSep      | Returns OS-specific path separator.                                     |
|              |                                                                         |
|              | >> yo eval 'key = ( pathSep )'                                          |
+--------------+-------------------------------------------------------------------------+
| rand         | Returns a random integer value from min (inclusive) to max (exclusive). |
|              |                                                                         |
|              | >> yo eval 'key = ( rand 8 16 )'                                        |
+--------------+-------------------------------------------------------------------------+
| regexFind    | Return the first match of the regular expression.                       |
|              |                                                                         |
|              | >> yo eval 'key = ( regexFind "[a-zA-Z][1-9]" "abcd1234" )'             |
+--------------+-------------------------------------------------------------------------+
| regexFindAll | Returns a slice of all matches of the regular expression.               |
|              |                                                                         |
|              | >> yo eval 'key = ( regexFindAll "[2,4,6,8]" "123456789" -1 )'          |
+--------------+-------------------------------------------------------------------------+
| regexMatch   | Test if the input string matches a regular expression.                  |
|              |                                                                         |
|              | >> yo eval 'key = ( regexMatch "dog$" "bulldog" )'                      |
+--------------+-------------------------------------------------------------------------+
| repeat       | Repeat a string multiple times.                                         |
|              |                                                                         |
|              | >> yo eval 'key = ( repeat 3 "hello" )'                                 |
+--------------+-------------------------------------------------------------------------+
| replace      | Perform simple string replacement.                                      |
|              |                                                                         |
|              | >> yo eval 'key = ( "I Am Henry VIII" | replace " " "-" )'              |
+--------------+-------------------------------------------------------------------------+
| sha1sum      | Computes the SHA1 digest of a specified string                          |
|              |                                                                         |
|              | >> yo eval 'key = ( sha1sum "Hello world!" )'                           |
+--------------+-------------------------------------------------------------------------+
| sha256sum    | Computes the SHA256 digest of a specified string                        |
|              |                                                                         |
|              | >> yo eval 'key = ( sha256sum "Hello world!" )'                         |
+--------------+-------------------------------------------------------------------------+
| split        | Split a string into a list of strings.                                  |
|              |                                                                         |
|              | >> yo eval 'key = ( split "$" "foo$bar$baz$bar" )'                      |
+--------------+-------------------------------------------------------------------------+
| substr       | Get a substring from a string.                                          |
|              |                                                                         |
|              | >> yo eval 'key = ( substr 0 5 "hello world" )'                         |
+--------------+-------------------------------------------------------------------------+
| tempDir      | Returns the default directory to use for temporary files.               |
|              |                                                                         |
|              | >> yo eval 'key = ( tempDir )'                                          |
+--------------+-------------------------------------------------------------------------+
| title        | Convert to title case.                                                  |
|              |                                                                         |
|              | >> yo eval 'key = ( title "hello world" )'                              |
+--------------+-------------------------------------------------------------------------+
| toDate       | Converts a string to a date.                                            |
|              |                                                                         |
|              | >> yo eval 'key = ( toDate "2006-01-02" "2017-12-31" )'                 |
+--------------+-------------------------------------------------------------------------+
| trim         | Removes space from either side of a string.                             |
|              |                                                                         |
|              | >> yo eval 'key = ( trim "   hello    " )'                              |
+--------------+-------------------------------------------------------------------------+
| trimPrefix   | Trim just the prefix from a string.                                     |
|              |                                                                         |
|              | >> yo eval 'key = ( trimPrefix "-" "-hello" )'                          |
+--------------+-------------------------------------------------------------------------+
| trimSuffix   | Trim just the suffix from a string.                                     |
|              |                                                                         |
|              | >> yo eval 'key = ( trimSuffix "-" "hello-" )'                          |
+--------------+-------------------------------------------------------------------------+
| uniq         | Generate a list with all of the duplicates removed.                     |
|              |                                                                         |
|              | >> yo eval 'key = ( split "$" "foo$bar$baz$bar" | uniq )'               |
+--------------+-------------------------------------------------------------------------+
| unixEpoch    | Returns the seconds since the unix epoch.                               |
|              |                                                                         |
|              | >> yo eval 'key = ( now | unixEpoch )'                                  |
+--------------+-------------------------------------------------------------------------+
| upper        | Convert the entire string to uppercase.                                 |
|              |                                                                         |
|              | >> yo eval 'key = ( upper "hello" )'                                    |
+--------------+-------------------------------------------------------------------------+
| uuid         | Generate UUID v4 universally unique IDs.                                |
|              |                                                                         |
|              | >> yo eval 'key = ( uuid )'                                             |
+--------------+-------------------------------------------------------------------------+
```

example:

```sh
$ yo eval 'key = ( upper "hello" )'
key: HELLO
```

# Yes but i want JSON!

Ok, use the `-j / --json` flag:

```sh
$ yo eval 'pets = [ { name=Dash kind=cat age=3 } {name=Harley kind=dog age=4} ]' -j
```

```json
{
   "pets": [
      {
         "age": 3,
         "kind": "cat",
         "name": "Dash"
      },
      {
         "age": 4,
         "kind": "dog",
         "name": "Harley"
      }
   ]
}
```

# How to install?

In order to use the `yo` command, compile it using the following command:

```sh
$ go install github.com/lucasepe/yo@latest
```

This will create the executable under your `$GOPATH/bin` directory.

## Ready-To-Use Releases 

If you don't want to compile the sourcecode yourself, [Here you can find the tool already compiled](https://github.com/lucasepe/yo/releases/latest) for:

- MacOS
- Linux
- Windows


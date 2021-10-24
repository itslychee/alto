# alto, a music organizer

alto is a program built for audio management. It's purpose is to provide the user the means to write a 
path construct, or format, to copy/rename audio to a custom path.

For Example:

```bash
$ alto -format "{%filename%}" -source path/to/source -destination path/to/destination -operation rename
```

Would move all files under `path/to/source` to `path/to/destination` with their original filename, name collisions
are automatically handled, but I intend to make this behavior customizable, as people may have different interpretations
of this behavior.

# Installing

Simply grab a binary from [Releases](/releases), if it is available. Otherwise go to [building](#Building) to learn
how to build alto from source.

# Building

You will need [Go](https://golang.org), [Git](https://git-scm.com), and a working internet connection so dependencies
in [`go.mod`](go.mod) can be installed.

```bash
$ git clone https://github.com/ItsLychee/alto
$ cd alto
$ go build ./cmd/alto
$ ./alto -help
# ...
```

# Path constructs

You may be a bit confused about what the `-format` argument does, which is understandable. This 
section is dedicated to teach you on how you can use alto to its fullest to achieve the organization meant 
for you.

## What is a path construct?

A path construct, as it implies, is a result of alto's processing of the `-format` string with the
current file's metadata being used as reference. So when this manual brings up stuff like "omitting from the path
construct" it simply means that *X* value won't be in the final result.

## Variables

As you may have guessed, `%filename%` is a variable. Proper variables must have ASCII-only identifiers wrapped around with
`%`, so while `%name%` is valid, `%こんにちは%` is not, although this requirement may change in the future.

### List of default variables

You are provided with variables representing metadata, metadatic variables are just pipelines to the methods
in [dhowden's tag metadata interface](https://pkg.go.dev/github.com/dhowden/tag#Metadata), and alto will be kept
updated to the latest and stable release.


* `title`
* `artist`
* `album`
* `albumartist`
* `genre`
* `composer`
* `year`
* `tracknumber`
* `tracktotal`
* `discnumber`
* `disctotal`
* `comment`
* `format`
* `filetype`
* `filename` _*this variable is not handled by `tag`, but it's just a variable of the name of the current file*_

The variable names should be self-explantory, you can refer the link above to get a
grasp of what each variable does based on the name of it.

## Groups

Groups are a collection of fields which are separated by `|`, and enwrapped by `{` and `}`. A group's job is to start on
the first field and see if it has a value, if it does not it will keep iterating over the list of fields until it finds a non-nil
field. If it does not find a viable field, it will simply just be omitted from the path construct

### Field

A field is a collection of string literals and variables. Unlike the outside of a group where only string
literals and groups are parsed, variables are also parsed. So `%variable% foobar` in a field would 
be represented `Variable StringLiteral` in a field, but outside one it would be `StringLiteral`.

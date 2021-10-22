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

The variable names should be self-explantory, you can refer to [this](https://pkg.go.dev/github.com/dhowden/tag#Metadata) to get a
grasp of what each variable does based on the name of it.

## Groups

*soon:tm:*
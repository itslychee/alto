# alto, a music organizer

alto is a program built for audio management. It's purpose is to provide the user the means to create
a path construct to move individual audio files to a select path while being provided with the metadata of the file
through [variables](#variables).

![APNG showcasing alto](.github/assets/showcase.png)

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

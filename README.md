#

# Installation

```
go install github.com/mkg0/bouldering@latest
```

# Usage

check cli guideline with `bouldering -h` to see all the commands and descriptions

```
bouldering profile add
bouldering profile list
bouldering profile remove
bouldering book
bouldering auto-book
bouldering enable-remote-booking
bouldering disable-remote-booking
```

[![asciicast](https://asciinema.org/a/HuQxtB0RHvBM5GFjk2UKXuuzd.svg)](https://asciinema.org/a/HuQxtB0RHvBM5GFjk2UKXuuzd)

# Development

Clone the repo in Github first and mind the change username in below commands

```sh
git clone git@github.com:mkg0/bouldering.git $GOPATH/src/github.com/mkg0/bouldering
cd $GOPATH/src/github.com/mkg0/bouldering
echo > temp_bouldering.tmp
```

Then you can run commands with `go run . {command}`
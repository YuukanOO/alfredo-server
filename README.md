# Alfredo: Your private smart home assistant

## Installation

Only from go for now, clone this repository as `alfredo` in `$GOPATH/src/github.com/YuukanOO`, on linux:

```
mkdir -p $GOPATH/src/github.com/YuukanOO
cd $GOPATH/src/github.com/YuukanOO
git clone https://github.com/YuukanOO/alfredo-server.git alfredo
```

And then install the package with `go install github.com/YuukanOO/alfredo`.

## Usage

First, starts by configuring the config file provided in this repository. Adjust it to suit your environment by following guidelines in the comments.

When it's done, run it with:

```
alfredo -c TOML_CONFIG_FILE run
```

*TODO: documentation*
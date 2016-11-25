# Alfredo: Your private smart home assistant

## Installation

Only from go for now, clone this repository as `alfredo` in `$GOPATH/src/github.com/YuukanOO`, on linux:

```
mkdir -p $GOPATH/src/github.com/YuukanOO
cd $GOPATH/src/github.com/YuukanOO
git clone https://github.com/YuukanOO/alfredo-server.git alfredo
cd $GOPATH/src/github.com/YuukanOO/alfredo
go install ./... && npm i
```

And then install the package with `go install github.com/YuukanOO/alfredo`.

## Usage

First, starts by configuring the config file provided in this repository. Adjust it to suit your environment by following guidelines in the comments.

When it's done, run it with:

```
alfredo -c TOML_CONFIG_FILE run
```

## How does it works?

Alfredo is a web server which will parse an `adapters.json` file to provides adapters to a client application. It provides an easy to use REST interface to manage rooms, devices and to execute commands as per defined in the adapters file.

It has been developed to be easy to extend by customizing the adapters file.

## How to define new adapters?

*TODO: more documentation*
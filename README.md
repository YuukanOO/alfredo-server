# Alfredo: Your private smart home assistant

## Installation

Only from go for now, clone this repository as `alfredo` in `$GOPATH/src/github.com/YuukanOO`, on linux:

```
mkdir -p $GOPATH/src/github.com/YuukanOO
cd $GOPATH/src/github.com/YuukanOO
git clone https://github.com/YuukanOO/alfredo-server.git alfredo
cd $GOPATH/src/github.com/YuukanOO/alfredo
go install ./...
```

## Usage

First, starts by configuring the config file provided in this repository. Adjust it to suit your environment by following guidelines in the comments.

When it's done, run it with:

```
alfredo -c TOML_CONFIG_FILE run
```

## Development

```
mkdir -p $GOPATH/src/github.com/YuukanOO
cd $GOPATH/src/github.com/YuukanOO
git clone https://github.com/YuukanOO/alfredo-server.git alfredo
cd $GOPATH/src/github.com/YuukanOO/alfredo
go get ./...
```

### Code organization

Alfredo is build around Domain Driven Design development.

- `webapp/`: Contains all code related to the web app,
- `identity/`: Identity and access context related to security,
- `registry/`: Registry context represents a house registry of all connected devices,

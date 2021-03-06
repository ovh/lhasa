# appcatalog

## Technical overview

This project is currently at very early stage and under active development. It is mainly written in golang and in angular5

## Contact

### Authors

* Rayene Ben Rayana <rayene.ben-rayana@corp.ovh.com>
* Fabien Meurillon <fabien.meurillon@corp.ovh.com>
* Yannick Roffin <yannick.roffin@corp.ovh.com>

## How to build

make

### Build requirements

* Git
* [Go installation](https://golang.org/doc/install) and [workspace](https://golang.org/doc/code.html#Workspaces) (`GOROOT` and `GOPATH` correctly set)
* GNU Make
* [dep](https://github.com/golang/dep) - Go dependency management tool
* [Angular 5 CLI](https://angular.io/guide/quickstart) - Management CLI for angular 5

### Steps

```
cd $GOPATH/src/github.com/ovh
git clone git@github.com/ovh/lhasa.git
cd lhasa && make
```

### Run

#### Locally

```
cd $GOPATH/src/github.com/ovh/lhasa
make run
```

#### Locally in full stack mode

```
cd $GOPATH/src/github.com/ovh/lhasa/webui
make live
```

```
cd $GOPATH/src/github.com/ovh/lhasa/api
make live
```

```
cd $GOPATH/src/github.com/ovh/lhasa-companion/api
make live
```

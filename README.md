# Docker Usage CLI Plugin
This plugin will list all docker images currently deployed in a foundation that
the calling user has access to.
### Complilation

```bash
go get github.com/ecsteam/docker-usage
cd $GOPATH/src/github.com/ecsteam/docker-usage

GOOS=darwin go build -o docker-usage-plugin-macosx
GOOS=linux go build -o docker-usage-plugin-linux
GOOS=windows go build -o docker-usage-plugin-windows.exe
```
### Installation
```bash
cf install-plugin ./docker-usage-plugin-${YOUR_OS}
```

### Usage
```
$ cf docker-usage 
```

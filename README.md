# List All Apps CF CLI Plugin
This plugin will list all running apps in the foundation (if you run it as admin), and display the memory and disk quota. It is pretty easy to modify the code to add other fields as required.

### Complilation

```bash
go get github.com/jtgammon/list-all-apps
cd $GOPATH/src/github.com/jtgammon/list-all-apps

GOOS=darwin go build -o docker-usage-plugin-macosx
GOOS=linux go build -o docker-usage-plugin-linux
GOOS=windows go build -o docker-usage-plugin-windows.exe
```
### Installation
```bash
cf install-plugin ./list-all-apps-${YOUR_OS}
```

### Usage
```
$ cf list-all-apps
```

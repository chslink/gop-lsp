CGO_ENABLED=0
GOARCH=amd64
buildStamp=$(date '+%Y-%m-%d %I:%M:%S')
# echo $buildStamp
gitHash=$(git describe --long --dirty --abbrev=14 --always)
# echo $gitHash
goVersion=$(go version)
# echo $goVersion
flags="-s -w -X 'main.buildStamp=$buildStamp' -X 'main.gitHash=$gitHash' -X 'main.goVersion=$goVersion'"
echo $flags
go build -ldflags "$flags" -o bin/linuxgoplsp main.go
upx -9 bin/linuxgoplsp

GOOS=darwin
go build -ldflags "$flags" -o bin/macgoplsp main.go
upx -9 bin/macgoplsp

GOOS=windows
go build -ldflags "$flags" -o bin/goplsp.exe main.go
upx -9 bin/goplsp.exe

cp bin/* ../gop-lsp-vscode/server/
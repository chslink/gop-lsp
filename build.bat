SET CGO_ENABLED=0
SET GOARCH=amd64
set buildStamp=%date:~0,4%-%date:~5,2%-%date:~8,2% %time:~0,2%:%time:~3,2%:%time:~6,2%
echo %buildStamp%
for /f "tokens=*" %%a in ('git describe --long --dirty --abbrev^=14 --always') do set gitHash=%%a
echo %gitHash%
for /f "tokens=*" %%b in ('go version') do set goVersion=%%b
echo %goVersion%
set flags="-s -w -X 'main.buildStamp=%buildStamp%' -X 'main.gitHash=%gitHash%' -X 'main.goVersion=%goVersion%'"
go build -ldflags %flags% -o bin\goplsp.exe main.go
upx -9 bin\goplsp.exe

SET GOOS=darwin
go build -ldflags %flags% -o bin\macgoplsp main.go
upx -9 bin\macgoplsp

SET GOOS=linux
go build -ldflags %flags% -o bin\linuxgoplsp main.go
upx -9 bin\linuxgoplsp

copy /y bin\* ..\gop-lsp-vscode\server\
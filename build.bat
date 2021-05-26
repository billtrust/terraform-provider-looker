@ECHO OFF

SET version=%1

IF "%version%" == "" (
    ECHO Version must be passed in
    EXIT /B 1
)

SET GOARCH=amd64
SET CGO_ENABLED=0

ECHO Using verion: %version%...

:: Build Linux
ECHO Building for linux...
SET GOOS=linux

go build -o bin/linux/terraform-provider-looker_v%version%

:: Build OSX (Darwin)
ECHO Building for OSX/darwin...
SET GOOS=darwin

go build -o bin/darwin/terraform-provider-looker_v%version%

:: Build Windows
ECHO Building for windows...
SET GOOS=windows

go build -o bin/win/terraform-provider-looker_v%version%.exe

ECHO Building complete!
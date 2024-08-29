@ECHO OFF

@SET GOOS=windows
@SET GOARCH=amd64
@SET FILENAME=vrtc3_win64.zip
go build -ldflags "-s -w" -trimpath && 7z a -mx9 -sdel %FILENAME% vrtc3.exe

@SET GOOS=windows
@SET GOARCH=386
@SET FILENAME=vrtc3_win32.zip
go build -ldflags "-s -w" -trimpath && 7z a -mx9 -sdel %FILENAME% vrtc3.exe

@SET GOOS=windows
@SET GOARCH=arm64
@SET FILENAME=vrtc3_win_arm64.zip
go build -ldflags "-s -w" -trimpath && 7z a -mx9 -sdel %FILENAME% vrtc3.exe

@SET GOOS=linux
@SET GOARCH=amd64
@SET FILENAME=vrtc3_linux_amd64
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%

@SET GOOS=linux
@SET GOARCH=386
@SET FILENAME=vrtc3_linux_i386
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%

@SET GOOS=linux
@SET GOARCH=arm64
@SET FILENAME=vrtc3_linux_arm64
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%

@SET GOOS=linux
@SET GOARCH=arm
@SET GOARM=7
@SET FILENAME=vrtc3_linux_arm
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%

@SET GOOS=linux
@SET GOARCH=arm
@SET GOARM=6
@SET FILENAME=vrtc3_linux_armv6
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%

@SET GOOS=linux
@SET GOARCH=mipsle
@SET FILENAME=vrtc3_linux_mipsel
go build -ldflags "-s -w" -trimpath -o %FILENAME% && upx %FILENAME%

@SET GOOS=darwin
@SET GOARCH=amd64
@SET FILENAME=vrtc3_mac_amd64.zip
go build -ldflags "-s -w" -trimpath && 7z a -mx9 -sdel %FILENAME% vrtc3

@SET GOOS=darwin
@SET GOARCH=arm64
@SET FILENAME=vrtc3_mac_arm64.zip
go build -ldflags "-s -w" -trimpath && 7z a -mx9 -sdel %FILENAME% vrtc3

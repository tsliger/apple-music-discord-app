@echo off

setlocal
set URL=https://dotnet.microsoft.com/en-us/download/dotnet/thank-you/runtime-8.0.17-windows-x64-installer
set DEST_DIR=.\src-tauri\dotnet
set FILENAME=dotnet-runtime-8.0.17-win-x64.exe
set DEST=%DEST_DIR%\%FILENAME%

:: Create the folder if it doesn't exist
if not exist "%DEST_DIR%" (
    mkdir "%DEST_DIR%"
)


:: Use PowerShell to download the file
powershell -Command "Invoke-WebRequest -Uri '%URL%' -OutFile '%DEST%'"

echo Download complete.
endlocal

cd windows-apple-music-info
if errorlevel 1 (
    echo Directory ./windows-apple-music-info does not exist. Exiting.
    exit /b 1
)

echo Building Windows dependencies
dotnet publish -c Release -o ..\src-tauri\
copy "..\src-tauri\windows-apple-music-info.exe" "..\go-am-discord-rpc"

cd ..

cd go-am-discord-rpc
if errorlevel 1 (
    echo Directory ./go-am-discord-rpc does not exist. Exiting.
    exit /b 1
)

echo Building Go dependencies
make %*

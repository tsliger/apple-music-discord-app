#!/bin/bash

URL="https://dotnet.microsoft.com/en-us/download/dotnet/thank-you/runtime-8.0.17-windows-x64-installer"
DEST_DIR="./src-tauri/dotnet"
FILENAME="dotnet-runtime-8.0.17-win-x64.exe"
DEST="$DEST_DIR/$FILENAME"

# Create the folder if it doesn't exist
mkdir -p "$DEST_DIR"

# Download the file
curl -L "$URL" -o "$DEST"

echo "Download complete."

cd ./windows-apple-music-info

echo "Running: dotnet publish -c Release -r win-x64 ..."
dotnet publish -c Release -r win-x64 /p:EnableWindowsTargeting=true -o ../src-tauri/
echo "dotnet publish completed with exit code $?"

cd ..

cd ./go-am-discord-rpc || {
  echo "Directory ./go-am-discord-rpc does not exist. Exiting."
  exit 1
}

echo "Building Go dependancies"
make "$@"

#! /bin/bash

GOARCH=386 go build -o macos32/EFFWords
GOOS=windows GOARCH=386 go build -o win32/EFFWords.exe
GOOS=linux GOARCH=386 go build -o nix32/EFFWords

# zipum
zip -r macos32{.zip,}
zip -r win32{.zip,}
zip -r nix32{.zip,}
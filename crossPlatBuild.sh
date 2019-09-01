#! /bin/bash

GOOS=darwin GOARCH=386 go build -o bin/macos32/EFFWords
GOOS=windows GOARCH=386 go build -o bin/win32/EFFWords.exe
GOOS=linux GOARCH=386 go build -o bin/nix32/EFFWords

GOOS=darwin GOARCH=amd64 go build -o bin/macos64/EFFWords
GOOS=windows GOARCH=amd64 go build -o bin/win64/EFFWords.exe
GOOS=linux GOARCH=amd64 go build -o bin/nix64/EFFWords

# zipum
zip -r bin/macos32{.zip,}
zip -r bin/win32{.zip,}
zip -r bin/nix32{.zip,}
zip -r bin/macos64{.zip,}
zip -r bin/win64{.zip,}
zip -r bin/nix64{.zip,}
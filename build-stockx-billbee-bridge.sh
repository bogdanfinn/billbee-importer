#!/bin/sh

echo 'Clear Build Archiv'
mkdir ./stockx-billbee-bridge || true
rm -rf ./stockx-billbee-bridge-release/*


echo 'Build OSX'
go build -o ./stockx-billbee-bridge-release/stockx-billbee-bridge-macos ./cmd/stockx-billbee-bridge/main.go

echo 'Build Linux'
GOOS=linux GOARCH=amd64 go build -o ./stockx-billbee-bridge-release/stockx-billbee-bridge-linux ./cmd/stockx-billbee-bridge/main.go

echo 'Build Windows'
GOOS=windows GOARCH=amd64 go build -o ./stockx-billbee-bridge-release/stockx-billbee-bridge-windows.exe ./cmd/stockx-billbee-bridge/main.go

echo 'Copy config files'
cp ./cmd/stockx-billbee-bridge/config.dist.yml ./stockx-billbee-bridge-release/config.dist.yml
cp ./cmd/stockx-billbee-bridge/example.csv ./stockx-billbee-bridge-release/example.csv
cp ./cmd/stockx-billbee-bridge/README.md ./stockx-billbee-bridge-release/README.md

echo 'Zip Build'

zip -r stockx-billbee-bridge-$1.zip ./stockx-billbee-bridge-release

rm -rf ./stockx-billbee-bridge-release
rm -rf ./stockx-billbee-bridge
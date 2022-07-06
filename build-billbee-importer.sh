#!/bin/sh

echo 'Clear Build Archiv'
mkdir ./billbee-importer || true
rm -rf ./billbee-importer-release/*


echo 'Build OSX'
go build -o ./billbee-importer-release/billbee-importer-macos ./cmd/billbee-importer/main.go

echo 'Build Linux'
GOOS=linux GOARCH=amd64 go build -o ./billbee-importer-release/billbee-importer-linux ./cmd/billbee-importer/main.go

echo 'Build Windows'
GOOS=windows GOARCH=amd64 go build -o ./billbee-importer-release/billbee-importer-windows.exe ./cmd/billbee-importer/main.go

echo 'Copy config files'
cp ./cmd/billbee-importer/config.dist.yml ./billbee-importer-release/config.dist.yml
cp ./cmd/billbee-importer/stockx-example.csv ./billbee-importer-release/stockx-example.csv
cp ./cmd/billbee-importer/alias-example.csv ./billbee-importer-release/alias-example.csv
cp ./cmd/billbee-importer/README.md ./billbee-importer-release/README.md

echo 'Zip Build'

zip -r billbee-importer-$1.zip ./billbee-importer-release

rm -rf ./billbee-importer-release
rm -rf ./billbee-importer
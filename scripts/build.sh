#!/bin/sh
srcPath="cmd"
pkgFile="main.go"
outputPath="build"
entrypoint="redshed"
outputApp="redshed"
output="$outputPath/$outputApp"
src="$srcPath/$entrypoint/$pkgFile"

printf "\nBuilding: $outputApp\n"
time go build -o $output $src
printf "\nBuilt: $outputApp size:"
ls -lah $output | awk '{print $5}'
printf "\nDone building: $outputApp\n\n"

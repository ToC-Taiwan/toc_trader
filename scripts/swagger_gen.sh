#!/bin/bash

echo 'package main' > ./tradebot.go
swag init -g pkg/routers/swagger.go -o ./docs
rm -rf ./tradebot.go

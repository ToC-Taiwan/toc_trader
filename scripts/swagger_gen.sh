#!/bin/bash

echo 'package main' > ./tradebot.go
rm -rf ./docs
swag init -g ./pkg/routers/swagger.go ./docs
rm -rf ./tradebot.go

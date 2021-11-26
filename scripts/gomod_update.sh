#!/bin/bash

rm go.mod
rm go.sum

go mod init gitlab.tocraw.com/root/toc_trader
go mod tidy

# python ./scripts/master_mod.py

# rm go.mod
# rm go.sum

# mv temp_go.mod go.mod
# go mod tidy

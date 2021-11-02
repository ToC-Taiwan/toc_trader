#!/bin/bash

git clone git@gitlab.tocraw.com:root/trade_bot_protobuf.git

/Users/timhsu/dev_projects/tools/protoc/bin/protoc -I=. --go_out=. ./trade_bot_protobuf/src/*.proto

rm -rf trade_bot_protobuf

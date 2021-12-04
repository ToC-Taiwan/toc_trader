#!/bin/bash

go-callvis \
    -focus gitlab.tocraw.com/root/toc_trader/pkg/modules/tradebot \
    -skipbrowser \
    -file=./assets/callvis \
    ./cmd || exit 1

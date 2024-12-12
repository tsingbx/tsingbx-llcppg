#!/bin/bash
set -e

llgo install ./_xtool/llcppsymg
llgo install ./_xtool/llcppsigfetch
go install ./cmd/gogensig

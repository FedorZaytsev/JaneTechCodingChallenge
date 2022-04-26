#!/usr/bin/env bash

cd cmd/server
go build && ./server -c ../../configs/config.yaml
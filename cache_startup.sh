#!/bin/bash
source ~/.profile
cd web_cache
go mod tidy
go run main.go
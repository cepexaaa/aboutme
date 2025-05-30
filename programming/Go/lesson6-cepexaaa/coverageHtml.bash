#!/bin/bash
 go test -v -coverprofile cover.out ./... -race
 go tool cover -html cover.out -o cover.html
 open cover.html
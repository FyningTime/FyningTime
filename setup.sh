#!/usr/bin/env bash

go get fyne.io/fyne/v2@latest
go install fyne.io/tools/cmd/fyne@latest

go mod tidy

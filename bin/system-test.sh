#!/bin/sh

# go test -v -run Test_Vendorer -run Test_Ingredienter -run Test_Stager ./tests/main.go
go test \
  ./tests/system/main_test.go \
  ./tests/system/vendorer_test.go

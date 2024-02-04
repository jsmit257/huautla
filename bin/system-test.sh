#!/bin/sh

go test -v -run Test_Vendorer -run Test_Ingredienter -run Test_Stager ./tests/main.go

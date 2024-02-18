#!/bin/bash

files=(
  ./tests/system/main_test.go
  ./tests/system/vendorer_test.go
  ./tests/system/substrater_test.go
  ./tests/system/ingredienter_test.go
  ./tests/system/substrateingredienter_test.go
  ./tests/system/stager_test.go
  ./tests/system/eventtyper_test.go
  ./tests/system/eventer_test.go
  ./tests/system/strainer_test.go
  ./tests/system/strainattributer_test.go
  ./tests/system/lifecycler_test.go
)

# go test \
#   -run Test_GetLifecycleEvents \
#   -run Test_SelectByEventType \
#   -run Test_RemoveEvent \
#   "${files[@]}"
go test "${files[@]}"
  

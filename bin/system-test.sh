#!/bin/bash

files=(
  ./tests/system/main_test.go
  ./tests/system/vendorer_test.go
  ./tests/system/substrater_test.go
  ./tests/system/ingredienter_test.go
  ./tests/system/substrateingredienter_test.go
  ./tests/system/stager_test.go
  ./tests/system/eventtyper_test.go
  ./tests/system/observer_test.go
  ./tests/system/strainer_test.go
  ./tests/system/strainattributer_test.go
  ./tests/system/lifecycler_test.go
  ./tests/system/lifecycleeventer_test.go
  ./tests/system/generationer_test.go
  ./tests/system/generationeventer_test.go
  ./tests/system/sourcer_test.go
  ./tests/system/noter_test.go
  ./tests/system/photoer_test.go
)

go test "${files[@]}"

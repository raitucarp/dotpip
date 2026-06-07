#!/bin/bash
sed -i 's/func TestVectorOptionsCoverage(t \*testing\.T) {/func TestVectorOptionsCoverage(_ \*testing.T) {/' options_test.go

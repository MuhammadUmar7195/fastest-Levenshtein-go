#!/bin/bash
# Differential Fuzzer Harness
# This script runs the differential fuzzer located in the root test suite.
go test -run TestDifferentialAgainstOriginal -v ../...

#!/bin/bash

OUTPUT="compiler"
trap "rm -f $OUTPUT" EXIT

clang++ $(ls *.cpp) -o $OUTPUT
./$OUTPUT

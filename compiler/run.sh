#!/bin/bash

OUTPUT="compiler"
trap "rm -f $OUTPUT" EXIT

clang++ -g $(ls *.cpp) `llvm-config --cxxflags --ldflags --system-libs --libs core` -o $OUTPUT
./$OUTPUT

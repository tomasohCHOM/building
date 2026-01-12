#!/bin/bash

OUTPUT="compiler"
trap "rm -f $OUTPUT" EXIT

clang++ $(ls *.cpp) `llvm-config-17 --cxxflags --ldflags --system-libs --libs core orcjit native` -o $OUTPUT
./$OUTPUT

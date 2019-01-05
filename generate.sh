#!/bin/bash

set -ex

# Setup greet pb file
echo running...
protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
protoc calculator/calcpb/calc.proto --go_out=plugins=grpc:.

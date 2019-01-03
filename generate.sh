#!/bin/bash

# Setup greet pb file
protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

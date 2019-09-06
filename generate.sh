#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=plugins:.
protoc calculator/calculatorpb/calculator.proto --go_out=plugins:.
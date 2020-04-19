#!/bin/bash

./SimElevatorServer &
cd src/
go run main.go
cd ../

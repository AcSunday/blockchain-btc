#!/bin/bash
rm ./*.db
go build -o blockchain
./blockchain

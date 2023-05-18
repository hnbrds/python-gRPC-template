#!/bin/bash
go run gateway_server.go &
python grpc_server.py

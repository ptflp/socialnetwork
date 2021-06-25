#!/bin/sh
openssl genrsa -f4 -out private 4096
openssl rsa -in private -outform PEM -pubout -out public
#!/bin/sh
openssl genrsa -out private 2048
openssl rsa -in private -pubout -out public
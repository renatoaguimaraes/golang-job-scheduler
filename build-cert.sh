#!/bin/bash
# file: build-cert.sh

openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout cert/server-ca-key.pem -out cert/server-ca-cert.pem -subj "/C=BR/OU=Server CA/CN=localhost/emailAddress=renatoaguimaraes@gmail.com"
openssl x509 -in cert/server-ca-cert.pem -noout -text
openssl req -newkey rsa:4096 -nodes -keyout cert/server-key.pem -out cert/server-req.pem -subj "/C=BR/OU=Server/CN=localhost/emailAddress=renatoaguimaraes@gmail.com"
openssl x509 -req -in cert/server-req.pem -days 60 -CA cert/server-ca-cert.pem -CAkey cert/server-ca-key.pem -CAcreateserial -out cert/server-cert.pem -extfile cert/server-ext.conf
openssl x509 -in cert/server-cert.pem -noout -text

openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout cert/client-ca-key.pem -out cert/client-ca-cert.pem -subj "/C=BR/OU=Client CA/CN=localhost/emailAddress=renatoaguimaraes@gmail.com"
openssl x509 -in cert/client-ca-cert.pem -noout -text
openssl req -newkey rsa:4096 -nodes -keyout cert/client-key.pem -out cert/client-req.pem -subj "/C=BR/OU=Client/CN=localhost/emailAddress=renatoaguimaraes@gmail.com"
openssl x509 -req -in cert/client-req.pem -days 60 -CA cert/client-ca-cert.pem -CAkey cert/client-ca-key.pem -CAcreateserial -out cert/client-cert.pem -extensions v3_req -extfile cert/client-ext.conf
openssl x509 -in cert/client-cert.pem -noout -text
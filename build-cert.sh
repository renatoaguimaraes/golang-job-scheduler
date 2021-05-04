#!/bin/bash
# file: build-cert.sh

# server ca
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout cert/server-ca-key.pem -out cert/server-ca-cert.pem -subj "/C=BR/OU=Server CA/CN=localhost/emailAddress=renatoaguimaraes@gmail.com"
openssl req -newkey rsa:4096 -nodes -keyout cert/server-key.pem -out cert/server-req.pem -subj "/C=BR/OU=Server/CN=localhost/emailAddress=renatoaguimaraes@gmail.com"
openssl x509 -req -in cert/server-req.pem -days 60 -CA cert/server-ca-cert.pem -CAkey cert/server-ca-key.pem -CAcreateserial -out cert/server-cert.pem -extfile cert/server-ext.conf

# authorized client ca
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout cert/client-ca-key.pem -out cert/client-ca-cert.pem -subj "/C=BR/OU=Client CA/CN=localhost/emailAddress=renatoaguimaraes@gmail.com"

# authorized user with admin role
openssl req -newkey rsa:4096 -nodes -keyout cert/client-key.pem -out cert/client-req.pem -subj "/C=BR/OU=Client/CN=localhost/emailAddress=renatoaguimaraes@gmail.com"
openssl x509 -req -in cert/client-req.pem -days 60 -CA cert/client-ca-cert.pem -CAkey cert/client-ca-key.pem -CAcreateserial -out cert/client-cert.pem -extensions v3_req -extfile cert/client-ext.conf

# authorized user with user role
openssl req -newkey rsa:4096 -nodes -keyout cert/client-user-key.pem -out cert/client-user-req.pem -subj "/C=BR/OU=Client/CN=localhost/emailAddress=renatoaguimaraes@gmail.com"
openssl x509 -req -in cert/client-user-req.pem -days 60 -CA cert/client-ca-cert.pem -CAkey cert/client-ca-key.pem -CAcreateserial -out cert/client-user-cert.pem -extensions v3_req -extfile cert/client-user-ext.conf
openssl x509 -in cert/client-user-cert.pem -noout -text

# unauthorized ca
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout cert/client-unauth-ca-key.pem -out cert/client-unauth-ca-cert.pem -subj "/C=BR/OU=Unauthorized CA/CN=localhost/emailAddress=renatoaguimaraes@gmail.com"

# unauthorized user
openssl req -newkey rsa:4096 -nodes -keyout cert/client-unauth-user-key.pem -out cert/client-unauth-user-req.pem -subj "/C=BR/OU=Unauthorized Client/CN=localhost/emailAddress=renatoaguimaraes@gmail.com"
openssl x509 -req -in cert/client-unauth-user-req.pem -days 60 -CA cert/client-unauth-ca-cert.pem -CAkey cert/client-unauth-ca-key.pem -CAcreateserial -out cert/client-unauth-user-cert.pem -extensions v3_req -extfile cert/client-user-ext.conf

# Monetary Data

This sample will attempt to insert documents using the [BSON Type Decimal](https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.9.0/bson/primitive#Decimal128) from the official MongoDB driver. 

For arithmetic operations with decimal data type numbers in Golang I'll be using ussing https://github.com/shopspring/decimal .

A working docker mongodb container is required in order to run this program, it can be initiallized with `docker run -d -p 27017:27017 --name m1 mongo`.

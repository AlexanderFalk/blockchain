FROM golang:jessie

# Installing needed dependency:
RUN go get github.com/boltdb/bolt/...

# Setting the working directory in the container
WORKDIR /blockchain

# Copying the blockchain code into the container
COPY block.go /blockchain/block.go
COPY blockchain.go /blockchain/blockchain.go
COPY pow.go /blockchain/pow.go
COPY transaction.go /blockchain/transaction.go
COPY CLI.go /blockchain/CLI.go
COPY utils.go /blockchain/utils.go
COPY run.go /blockchain/run.go

# Making port 8080 available to the host
EXPOSE 8080

# Building and running the blockchain when the container is started
RUN go build 

# /blockchain/block.go /blockchain/blockchain.go /blockchain/pow.go /blockchain/
# transaction.go /blockchain/CLI.go /blockchain/run.go

ENTRYPOINT ./blockchain createblockchain --address Jules
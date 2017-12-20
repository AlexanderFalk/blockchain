# Blockchain - Group14

This repository contains Group14's solution to the 2017 System Integration [blockchain assignment](https://github.com/datsoftlyngby/soft2017fall-system-integration-teaching-material/blob/e7729438dd0a3fa1c4cbc2a7b1d3651e8fc4600f/lecture_notes/12-Blockchain_Intro.ipynb).



## Project requirements

### General

* The blockchain has to be able to run on a network of distributed nodes (containers, virtual machines or cloud machines).
* It should be possible to perform simple blockchain *actions*. Examples include:
  * Adding transactions.
  * Block mining.
    * At least two versions for mining - used for Proof of Work.


### Formal

* The network has to be P2P network (peer-to-peer).

* The blockchain network should contain at least four (4) nodes.

* A reproducible setup is required. In this case, a bash script that instantiates the blockchain with example nodes and executes a test scenario.

* A screencast or screendumps that demonstrate the blockchain functionality.

* References for **all** information sources used for development/implementation.


### Optional

* User interface.​


### Own specifications

* We are writing the blockchain in [Go](https://golang.org/#).

* We are using Docker containers for our network of nodes.

* We are launching the four peers/Docker containers using `docker-compose`

  ​



## Get started

To setup and test our Blockchain simulation, follow these steps.

### Prerequisites

* You need to have [Docker](https://docs.docker.com/engine/installation/) installed. 

* You need to have [Git](https://git-scm.com/downloads) installed.

* You need to have a Bash shell installed (Terminal, iTerm2, Git Bash etc.).


### Installation

1. Open your terminal. 

2. Make sure Docker is booted on your machine.

3. Go ahead an clone this repository in a folder of your choosing:

   ```sh
   $ git clone https://github.com/AlexanderFalk/blockchain.git
   ```

4. Go into your new local directory:

   ```sh
   $ cd blockchain/
   ```

5. Give our setup script permission to be executed:

   ```sh
   $ chmod +x ./setup.sh
   ```

6. Execute our setup script:

   ```sh
   $ ./setup.sh
   ```



The script will now automatically build a proper Docker image that handles building the Go files, while simultaneously launching four separate containers, acting as peers.

# Overview of concepts
Our MVP Blockchain is built in Golang. Golang is very modern, when it comes to network programming and that is one of the main reason why we chose this language. 
To explain how the blockchain works, we will go through what each file is during and what concepts that are covered by this program. 
## Blocks
In blockchain a block is what store valuable information. In our example it stores transactions, which is the essence of any cryptocurrency. It also contains a timestamp, the hash of the block, the hash of the previous block, and a nonce. A nonce is an arbitrary number that is used in cryptography. In our case we will use it as an incremental value with our total hash value of all our block information. This is also called HashCash, which is a Proof-of-Work algorithm. It is a key concept within blockchain and it is used to make mining “hard”. The algorithm is used to fulfill a requirement of the hash value. In case of Bitcoin, the requirement is that “the first 20 bits of a hash must be zeros”. When you hash the block information at first, you might not hit that requirement, and therefore you need to use this algorithm. The idea is quite simple; you take the block hash value, increment it with one, and if it matches the requirement, then the block is validated. If it doesn’t, then you add one more and calculate the hash value. This incremental is happening until the requirement is met. 
In our program, the target bits for the requirement is set to 16. This is done due to the time it takes to validate the block. 
## Blockchain
A blockchain is a simple database with a certain structure. It’s ordered and has a back-linked list. This means that the blocks are stored in the insertion order and each block is linked to the previous one. 
## Persistence
Since we don’t want to store the blockchain in memory, we’ve used a database called: BoltDB. BoltDB is a minimalistic golang database that doesn’t require a running server. Exactly what we need, since the blockchain must be decentralized, meaning, every active peer needs a copy of it. 
BoltDB is a key/value storage, so there are no tables like in relational databases. This data is stored in what BoltDB is calling “buckets”. The buckets are collections of key/value pairs within the database. Another nice feature about BoltDB is that it is data-type independent. Everything is stored in byte arrays. In our case, we serialize our data first and then we add it to the bucket. To get the data out, we deserialize the data. This is done to get the block information back, when we need it from the blockchain. 
The bucket is just a file on the disk, with the extension: “.db”.
Before we can add a blockchain to the database file, we need an existing block. There must be at least one block to generate a blockchain and that block is called “genesis block”. New blocks has to refer to a previous block and therefore we need to “start” the chain, so other can use it. 

## Transactions
This is the hard part. This is where it gets intense. A transaction is a combination of inputs and outputs (and amount). 
An input is a record of which an address was used to send assets to an unknown output. It is a new transaction, which is referencing to an output of previous transactions. 
An amount is even telling. What is being transferred. 
An output is where coins are stored. This is the address of the retriever. 
To send assets you need two things: an asset address and a private key. The asset address is kept secret.
To make it understandable, this metaphor gives a quite good picture: Think of your asset address as a safe deposit box with a glass front. Everyone can see and knows what is in it, but only the private key can unlock it and take things out or put things in. This is also known as a wallet, but it is not implemented in your MVP. 
The funny thing about assets and transactions are that they don’t exist as a single entity. They exist as records of transactions, increasing and decreasing in values. You can end with many different transactions tied to your asset address. Let’s create a short example: 
Jane sent Alice two assets, John sent her four, and Mia sent her a single. These are separate transactions. They are not combined in Alice’s wallet. They are simply different transactions records. 
So, when Alice wants to sent Mia assets, her wallet will use the transaction records with different amount and add up to the number of assets, that she wants to send. 
If Alice wants to send (the input) Mia 1.5 asset, she will have a problem. None of her previous transactions match up that amount. Even not when they are combined. She can’t split them. You can only spend the whole output of a transaction. 
What she will have to do in this case is that she will have to send one of the incoming transactions, in this case Jane’s transaction, and then send it with a new input of 0.5, that will point to her with the “leftovers” as an output. 
This is the core concept of the blockchain and it is very hard to wrap your head around. 

## Putting it all together
For our blockchain you will be able to:  
*	Create a new blockchain
*	Write to a database file
*	Mine a block and validate the Proof-of-Work
*	Use the blockchain as a CLI command
Below you will see the usage of the blockchain, which will describe how you use the blockchain. 


### Usage

To run the blockchain locally, type the following in the terminal:

``` 
$ go build
```
You can create a blockchain with a new DB file, containing a genesis block, by typing:
```
$ ./blockchain createblockchain --address Jules
```
You can print out the chain. This will open the database, look at the tip, print it out, take the next "tip" and print it out, and continue until our Next() function is empty:
```
$ ./blockchain printchain
```
Get the balance of a user. Here we find all of the unspent output transactions that an address has. Unspent means that the outputs weren't referenced in any inputs. We will go through the database file and if the address (used as our private key for now) can unlock an output, then we will append it to our balance struct:
```
$ ./blockchain getbalance --address Ivan
```
To send from one user to another. When we send assets, we need to create a new transaction, put it in a block, and mine the block. When we executed the _createblockchain_ command, we created a genesis block. So this is all we have for now. When we use the _send_ command, we find all of the spendable outputs for the _From_ address and check if they have enough assets based on the calculation. We create two outputs when we are sending assets
* One that's locked with the receiver address. The transfer of assets
* One that's locked with the senders address. This is our change, which is explained above. This is only created if the unspent output  amount exceeds the required amount. Then at last we will put it into the blockchain:
```
$ ./blockchain send -from Ivan -to Pedro -amount 6
```
You'll then be able to retrieve the balance again for both Ivan and Pedro. Here you will see that Ivan will 4 assets, while Pedro has 6
```
$ ./blockchain getbalance --address Ivan
```

## Docker usage
If you followed the instructions of the previous section, then four peers should have been launched on your machine as Docker containers. You should now be able to see these containers as processes if you type:

```sh
$ docker ps
```

Now, to try out the blockchain you need to run a blockchain command on one of the peers. For example, you could try to create a transaction to "Jules" with peer number one:

```sh
$ docker exec -d peer1 /blockchain/block createblockchain --address Jules
```

This command will execute the `createblockchain` command on the container named "peer1".

...

## Issues
We are aware of the fact that the blockchain can't be decentralized to multiple peers. We are able to communicate between the peers, but we haven't been able to setup a functional Peer-to-peer network, where each has a copy of the blockchain. One of the main issues are, that the nodes haven't been split into categories. We should've had a:
* Central node - A node all nodes connect to, and it's the node that sends data between other nodes
* Full node - Validate blocks that miners have mined. This one has a full copy of the blockchain. This is also the node that helps the other nodes to discover each other.  
* Two mining nodes - Storing new transactions in mememory pool and when there're enough of transactions, it'll mine a new block. 

As you see above, we've an idea of using Docker, but without actually using it as a P2P network. We simply couldn't get it up and running probably. 
We also tried with vagrant and virtual machines. Here we instantiated 4 VM's, which got the blockchain CLI copied onto them, so they were able to execute it. The biggest issue here was that the program was not configured to let the nodes run with separated roles. This would've required a lot of changes in the blockchain CLI. 
Below you can see the vagrantfile:
```
# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|

  	config.vm.box = "ubuntu/trusty64"
  	config.vm.network "private_network", type: "dhcp"
  	config.vm.provision "file", source: "./", destination: "$HOME/blockchain"
  	config.vm.define "FullNode", primary: true do |server|
      server.vm.network "private_network", ip: "192.168.20.2"
      server.vm.network "forwarded_port", guest: 8080, host: 9001
      server.vm.provider "virtualbox" do |vb|
        vb.memory = "1024"
        vb.cpus = "1"
      end
    end
    config.vm.define "CentralNode" do |client|
      client.vm.network "private_network", ip: "192.168.20.3"
      client.vm.network "forwarded_port", guest: 8080, host: 9002
      client.vm.provider "virtualbox" do |vb|
        vb.memory = "1024"
        vb.cpus = "1"
      end
	 end
    config.vm.define "MineNode" do |client|
      client.vm.network "private_network", ip: "192.168.20.4"
      client.vm.network "forwarded_port", guest: 8080, host: 9003
      client.vm.provider "virtualbox" do |vb|
        vb.memory = "1024"
        vb.cpus = "1"
      end
	  end
    config.vm.define "MineNode" do |client|
      client.vm.network "private_network", ip: "192.168.20.5"
      client.vm.network "forwarded_port", guest: 8080, host: 9004
      client.vm.provider "virtualbox" do |vb|
        vb.memory = "1024"
        vb.cpus = "1"
      end
    end
end
```



## References

- https://hackernoon.com/learn-blockchains-by-building-one-117428612f46
- https://docs.docker.com/compose/compose-file/
- https://jeiwan.cc/posts/building-blockchain-in-go-part-1/
- https://github.com/datsoftlyngby/soft2017fall-lsd-teaching-material/blob/master/lecture_notes/03-Containers%20and%20VMs.ipynb

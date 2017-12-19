# Blockchain - Group14

This repository contains Group14's solution to the 2017 System Integration [blockchain assignment](https://github.com/datsoftlyngby/soft2017fall-system-integration-teaching-material/blob/e7729438dd0a3fa1c4cbc2a7b1d3651e8fc4600f/lecture_notes/12-Blockchain_Intro.ipynb).



## Project requirements

### General

* The blockchain has to be able to run on a network of distributed nodes (containers, virtual machines or cloud machines).
* It should be possible to perform simple blockchain *actions*. Examples include:
  * Adding transactions.
  * Block mining.
    * At least two versions for mining - used for Proof of Work.

      ​

### Formal

* The network has to be P2P network (peer-to-peer).

* The blockchain network should contain at least four (4) nodes.

* A reproducible setup is required. In this case, a bash script that instantiates the blockchain with example nodes and executes a test scenario.

* A screencast or screendumps that demonstrate the blockchain functionality.

* References for **all** information sources used for development/implementation.

  ​

### Optional

* User interface.

  ​

##Own specifications

* We are writing the blockchain in [Go](https://golang.org/#).
* We are using Docker containers for our network of nodes.
* We are launching the four peers/Docker containers using `docker-compose`



## Get started

To setup and test our Blockchain simulation, follow these steps.



### Prerequisites

* You need to have [Docker](https://docs.docker.com/engine/installation/) installed. 

* You need to have [Git](https://git-scm.com/downloads) installed.

* You need to have a Bash shell installed (Terminal, iTerm2, Git Bash etc.).

  ​

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



### Usage

If you followed the instructions of the previous section, then four peers should have been launched on your machine as Docker containers. You should now be able to see these containers as processes if you type:

```sh
$ docker ps
```



Now, to try out the blockchain you need to run a blockchain command on one of the peers. For example, you could try to create a transaction to "Jules" with peer number one:

```sh
$ docker exec -d peer1 /blockchain/block createblockchain --address Jules
```

...



## References

- https://hackernoon.com/learn-blockchains-by-building-one-117428612f46
- https://docs.docker.com/compose/compose-file/
- https://jeiwan.cc/posts/building-blockchain-in-go-part-1/
- https://github.com/datsoftlyngby/soft2017fall-lsd-teaching-material/blob/master/lecture_notes/03-Containers%20and%20VMs.ipynb
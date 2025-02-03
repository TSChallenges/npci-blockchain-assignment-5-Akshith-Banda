==================================================
Hyperledger Fabric Test Network and Node.js Gateway
==================================================

This document provides step-by-step instructions to:

1. Set up and start the Hyperledger Fabric test network using the `fabric-samples` repository.
2. Use the Node.js Gateway client to interact with the network.


Step 1: Clone the `fabric-samples` Repository
============================================
Execute `curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh`

.. code-block:: bash

       ./install-fabric.sh

this will install the fabric samples, binaries and the docker images


Step 2: Set Up the Test Network
===============================
The `test-network` directory contains scripts to deploy a basic Fabric network.

1. Navigate to the `test-network` directory:

   .. code-block:: bash

       cd test-network

2. Setup the network with certificate authority:

   .. code-block:: bash

       ./network.sh up createChannel -c mychannel -ca

Step 3: Write smart contract and deploy chaincode
===============================

1. Develop smart contract according to the given problem statement:
    a. CreateAsset and ReadAsset methods have already been provided.
    b. UpdateAsset and DeleteAsset methods have been added along with an helper mehtod AssetExists.

The smart contract can be found in `/workspaces/npci-blockchain-assignment-5-Akshith-Banda/fabric-samples/chaincode` dir

   .. code-block:: bash
        cd test-network
       ./network.sh deployCC -ccn asset-management -ccp ../chaincode/asset-management/go -ccl go

   This will deploy the `asset-management` chaincode written in JavaScript.

Step 4: Set Up the Node.js Gateway Client
=========================================
The Node.js Gateway client allows you to interact with the Fabric network programmatically.

1. Follow the readme instructins and create client directory in test-network. Create 2 files
    a. enrollAppUser.js - this will create the appuser identity in the ./client/wallet/ , which is used by gateway client to connect to gateway server.
    b. app.js - this will conect to fabric network and execute CreateAsset, ReadAsset, updateAsset and DeleteAsset transactions.

2. Run enrollAppUser.js:

   .. code-block:: bash

      node enrollAppUser.js

3. Run app.js:

   .. code-block:: bash

      node app.js

   This command outputs the following result:
   `Create transaction has been submitted
    Asset details: {"id":"asset1","owner":"Alice","value":100}
    update transaction has been submitted
    Asset details: {"id":"asset1","owner":"Alice","value":150}
    Delete transaction has been submitted
    Failed to submit transaction: Error: asset not found: asset1`
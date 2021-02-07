# Golang-BullyAlgorithm

### Introduction
This project focuses on skeleton implemention of the Bully Algorithm in Leader Election methodology in Distributed Systems using Remote Procedure Calls(RPC) in Go Language.

### Implementation
Major implementation contraints or strategies involve-
* Each node is aware of the other's node network information, i.e. IP address and port
* There is one file which needs to be built and the program needs to run on several terminal windows, each representing a node in the network. IP addresses is 127.0.0.1 and port numbers range from 3000-3004.
* For ease of implementation, 5 nodes are allowed by default. More nodes can be added with minimal changes in the code.
* By default, 5th node is the coordinator.
* To invoke elections, close the program of coordinator using *ctrl + c* and try communicating to the coordinator from other nodes.
* Upon the recovery of the old coordinator, elections can be invoked to elect the new coordinator again.

### Output interface
Output interface contains 3 major components-
1. Selecting id of the node
1. Give an option to state if the node just recovered from a crash
1. Communicate with the coordinator

![Output interface] (/screenshots/output_interface.png)

### RPC methods:
* **Election**
Handles election received from the nodes with lower id.by sending them OK message along with the node invoking election. Checks like `no_election_invoked: bool` are present to prevent multiple elections by a single host.

* **NewCoordinator**
This function is used by the coodinator by calling this function as broadcast in other nodes to update the coordinator id as the last stage of the Bully Algorithm.

* **HandleCommunication**
This function is called by the particpants to communicate to the coordinator and get response from the function. Fail in calling or getting response this function from coordinator triggers the election.

### Screenshots
_Fig: When coordinator node 5 is active and the other nodes are able to communicate with the coordinator_
![Coordinator active] (/screenshots/coordinator_active.png)

_Fig: When coordinator node is closed, Elections! 4 is the new coordinator after the elections._
![Election] (/screenshots/elections.png)

_Fig: Post election communication. 4 is the new coordinator and 3 is able to communicate with the new coordinator._
![Post elections] (/screenshots/post_election.png)

_Fig: Recovering from crash. 5 is coordinator again and 1 is able to communicate with new coordinator._
![Recovery elections] (/screenshots/recovery_elections.png)

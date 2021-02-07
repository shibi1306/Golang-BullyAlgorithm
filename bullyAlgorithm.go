
package main

import (
	"fmt"
	"net"
	"net/rpc"
	"log"
)

// Reply for getting response like OK messages from the RPC calls. 
type Reply struct{
	Data string
}

// Core Bully Algorithm structure. It contains functions registered for RPC.
// it also contains information regarding other sites
type BullyAlgorithm struct{
	my_id int
	coordinator_id int
	ids_ip map[int]string
}

// if a site has already invoked election, it doesnt need to start elections again
var no_election_invoked = true

// This is the election function which is invoked when a smaller host id requests for an election to this host
func (bully *BullyAlgorithm) Election(invoker_id int, reply *Reply) error{
	fmt.Println("Log: Receiving election from", invoker_id)
	if invoker_id < bully.my_id{
		fmt.Println("Log: Sending OK to", invoker_id)
		reply.Data = "OK"				// sends OK message to the small site
		if no_election_invoked{
			no_election_invoked = false
			go invokeElection()			// invokes election to its higher hosts
		}
	}
	return nil
}

var superiorNodeAvailable = false				// Toggled when any superior host sends OK message

// This function invokes the election to its higher hosts. It sends its Id as the parameter while calling the RPC
func invokeElection(){
	for id,ip := range bully.ids_ip{
		reply := Reply{""}
		if id > bully.my_id{
			fmt.Println("Log: Sending election to", id)
			client, error := rpc.Dial("tcp",ip)
			if error != nil{
				fmt.Println("Log:", id, "is not available.")
				continue
			}
			err := client.Call("BullyAlgorithm.Election", bully.my_id, &reply)
			if err != nil{
				fmt.Println(err)
				fmt.Println("Log: Error calling function", id, "election")
				continue
			}
			if reply.Data == "OK"{				// Means superior host exists
				fmt.Println("Log: Received OK from", id)
				superiorNodeAvailable = true
			}
		}
	}
	if !superiorNodeAvailable{					// if no superior site is active, then the host can make itself the coordinator
		makeYourselfCoordinator()
	}
	superiorNodeAvailable = false
	no_election_invoked = true					// reset the election invoked
}

// This function is called by the new Coordinator to update the coordinator information of the other hosts
func (bully *BullyAlgorithm) NewCoordinator(new_id int, reply *Reply) error{
	bully.coordinator_id = new_id 
	fmt.Println("Log:", bully.coordinator_id, "is now the new coordinator")
	return nil
}

func (bully *BullyAlgorithm) HandleCommunication(req_id int, reply *Reply) error{
	fmt.Println("Log: Receiving communication from", req_id)
	reply.Data = "OK"
	return nil
}

func communicateToCoordinator(){
	coord_id := bully.coordinator_id
	coord_ip := bully.ids_ip[coord_id]
	fmt.Println("Log: Communicating coordinator", coord_id)
	my_id := bully.my_id
	reply := Reply{""}
	client, err := rpc.Dial("tcp", coord_ip)
	if err != nil{
		fmt.Println("Log: Coordinator",coord_id, "communication failed!")
		fmt.Println("Log: Invoking elections")
		invokeElection()
		return
	}
	err = client.Call("BullyAlgorithm.HandleCommunication", my_id, &reply)
	if err != nil || reply.Data != "OK"{
		fmt.Println("Log: Communicating coordinator", coord_id, "failed!")
		fmt.Println("Log: Invoking elections")
		invokeElection()
		return
	}
	fmt.Println("Log: Communication received from coordinator", coord_id)
}

// This function is called when the host decides that it is the coordinator.
// it broadcasts the message to all other hosts and updates the leader info, including its own host.
func makeYourselfCoordinator(){
	reply := Reply{""}
	for id, ip := range bully.ids_ip{
		client, error := rpc.Dial("tcp", ip)
		if error != nil{
			fmt.Println("Log:", id, "communication error")
			continue
		}
		client.Call("BullyAlgorithm.NewCoordinator", bully.my_id, &reply)
	}
}

// Core object of bully algorithm initialized with all ip addresses of all other sites in the network
var bully = BullyAlgorithm{
	my_id: 		1,
	coordinator_id: 5,
	ids_ip: 	map[int]string{	1:"127.0.0.1:3000", 2:"127.0.0.1:3001", 3:"127.0.0.1:3002", 4:"127.0.0.1:3003", 5:"127.0.0.1:3004"}}


func main(){
	my_id := 0
	fmt.Printf("Enter the site id[1-5]: ")			// initialize the host id at the run time
	fmt.Scanf("%d", &my_id)
	bully.my_id = my_id
	my_ip := bully.ids_ip[bully.my_id]
	address, err := net.ResolveTCPAddr("tcp", my_ip) 
	if err != nil{
		log.Fatal(err)
	}
	inbound, err := net.ListenTCP("tcp", address)
	if err != nil{
		log.Fatal(err)
	}
	rpc.Register(&bully)
	fmt.Println("server is running with IP address and port number:", address)
	go rpc.Accept(inbound) // Accepting connections from other hosts.

	reply := ""
	fmt.Printf("Is this node recovering from a crash?(y/n): ")	// Recovery from crash.
	fmt.Scanf("%s", &reply)
	if reply == "y"{
		fmt.Println("Log: Invoking Elections")
		invokeElection()
	}

	random := ""
	for{
		fmt.Printf("Press enter for %d to communicate with coordinator.\n", bully.my_id)
		fmt.Scanf("%s", &random)
		communicateToCoordinator()
		fmt.Println("")
	}
	fmt.Scanf("%s", &random)
}

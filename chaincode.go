package main


//not forked
//besides the normal packages, we need to import shim and peer
//which can be found within fabric (git)
import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// we implement chaincode under 'ourchain' structure
type ourchain struct {
}

// we represent metadata with 'ourdata' structure
type entry struct {
	ObjectType  string `json:"docType"`
	ID          string `json:"id"`          //ID for the entry
	Hash        string `json:"hash"`        //primary ID for each entry
	Application string `json:"application"` //participating app
	NodeIP      string `json:"nodeIP"`      //IP of device device that created entry
	Owner       string `json:"owner"`       //username
	Updated     int    `json:"updated"`     //0 is outdated, 1 for updated
}

//main innitiate smartcontract ourchain
func main() {
	err := shim.Start(new(ourchain))
	if err != nil {
		fmt.Print("Error starting the chaincode, reason: %s", err)
	}
}

//Initialize the ledger
func (v *ourchain) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

//invoke a given function
func (v *ourchain) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, arguments := stub.GetFunctionAndParameters()
	fmt.Println("function is running....")

	if function == "initentry" {
		return v.initentry(stub, arguments)
		//} else if function == "setentry" {
		//return v.setentry(stub, arguments)
	} else if function == "readentry" {
		return v.readentry(stub, arguments)
	}

	fmt.Println("invoke did not find func: " + function)
	return shim.Error("function not found")
}

//innitialize chaincode entry
func (v *ourchain) initentry(stub shim.ChaincodeStubInterface, arguments []string) peer.Response {

	//detect errors in entry
	var err error
	if len(arguments) != 6 {
		return shim.Error("Incorrect # of arguments, expecting 6")
	}
	if len(arguments[0]) <= 0 {
		return shim.Error("ID must be 1 string")
	}
	if len(arguments[1]) <= 0 {
		return shim.Error("hash must be a string")
	}
	if len(arguments[2]) <= 0 {
		return shim.Error("Application must be a string")
	}
	if len(arguments[3]) <= 0 {
		return shim.Error("Node must be a string")
	}
	if len(arguments[4]) <= 0 {
		return shim.Error("Owner must be a string")
	}
	if len(arguments[5]) <= 0 {
		return shim.Error("Updated must be a string")
	}
	id := arguments[0]
	hash := arguments[1]
	application := arguments[2]
	node := arguments[3]
	owner := arguments[4]
	updated, err := strconv.Atoi(arguments[5])

	//check if entry already exists
	entrybytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Fail to get entry: " + err.Error())
	} else if entrybytes != nil {
		return shim.Error("This entry already exists.")
	}

	//create entry and marshal to JSON
	objecttype := "entry"
	entry := &entry{objecttype, id, hash, application, node, owner, updated}
	entryJSONbytes, err := json.Marshal(entry)
	if err != nil {
		return shim.Error(err.Error())
	}

	//save entry to state/blockchain
	err = stub.PutState(id, entryJSONbytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//success message
	fmt.Println("end initiation")
	return shim.Success(nil)
}

//read an entry into the ledger
func (v *ourchain) readentry(stub shim.ChaincodeStubInterface, arguments []string) peer.Response {
	var id, jsonResp string
	var err error

	if len(arguments) != 1 {
		return shim.Error("Incorrect number of arguments. expecting ID")
	}
	id = arguments[0]
	entrybytes, err := stub.GetState(id)
	if err != nil {
		jsonResp = "failed to get entry"
		return shim.Error(jsonResp)
	} else if entrybytes == nil {
		jsonResp = "No entry for " + id + " was found"
		return shim.Error(jsonResp)
	}
	return shim.Success(entrybytes)
}

/*
package main

import (
	ft "fmt"
)

func main() {
	ab := 1
	Ab := 2
	ft.Println(ab)
	ft.Println(Ab)
}
*/

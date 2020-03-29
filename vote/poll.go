package main

import (
	"fmt"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type poll struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	Status		string 	`json:"status"`
	NumVotes	int 	`json:"numVotes"`
}

type pollPrivateDetails struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	Salt 		string 	`json:"salt"`
	PollHash 	string 	`json:"pollHash"`
}

// ===========================================================
// initPoll - create a new poll and store into chaincode state
// ===========================================================

func (vc *VoteChaincode) initPoll(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	type pollTransientInput struct {
		PollID			string 	`json:"pollID"`
		Salt 			string 	`json:"salt"`
		PollHash 		string 	`json:"pollHash"`
	}

	fmt.Println("- start init vote")

	if len(args) != 0 {
		return shim.Error("Private data should be passed in transient map.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	pollJsonBytes, success := transMap["poll"]
	if !success {
		return shim.Error("poll must be a key in the transient map")
	}

	if len(pollJsonBytes) == 0 {
		return shim.Error("poll value in transient map cannot be empty JSON string")
	}

	var pollInput pollTransientInput
	err = json.Unmarshal(pollJsonBytes, &pollInput)
	if err != nil {
		return shim.Error("failed to decode JSON of: " + string(pollJsonBytes))
	}

	if len(pollInput.PollID) == 0 {
		return shim.Error("poll ID field must be a non-empty string")
	}

	if len(pollInput.Salt) == 0 {
		return shim.Error("salt field must be a non-empty string")
	} 

	if len(pollInput.PollHash) == 0 {
		return shim.Error("poll hash field must be a non-empty string")
	}

	existingPollAsBytes, err := stub.GetPrivateData("collectionPoll", pollInput.PollID)
	if err != nil {
		return shim.Error("Failed to get vote: " + err.Error())
	} else if existingPollAsBytes != nil {
		fmt.Println("This poll already exists: " + pollInput.PollID)
		return shim.Error("This poll already exists: " + pollInput.PollID)
	}

	poll := &poll{
		ObjectType: "poll",
		PollID: pollInput.PollID,
		Status: "ongoing",
		NumVotes: 0,
	}

	pollJSONasBytes, err := json.Marshal(poll)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutPrivateData("collectionPoll", pollInput.PollID, pollJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// create a composite key for poll private details collections using the poll ID and salt
	pollPrivateDetailsCompositeKey, err := stub.CreateCompositeKey("poll", []string{pollInput.PollID, pollInput.Salt})
	if err != nil {
		return shim.Error("Failed to create composite key for poll private details: " + err.Error())
	}

	pollPrivateDetails := &pollPrivateDetails {
		ObjectType: "pollPrivateDetails",
		PollID: pollInput.PollID,
		Salt: pollInput.Salt,
		PollHash: pollInput.PollHash,
	}

	pollPrivateDetailsBytes, err := json.Marshal(pollPrivateDetails)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutPrivateData(
		"collectionPollPrivateDetails", 
		pollPrivateDetailsCompositeKey, 
		pollPrivateDetailsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//register event
	err = stub.SetEvent("initEvent", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init poll (success)")
	return shim.Success(nil)

}

// ===================================================================
// getPoll - retrieve poll metadata from chaincode state
// ===================================================================

func (vc *VoteChaincode) getPoll(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting poll ID to query")
	}

	pollID := args[0]

	// ==== retrieve the vote ====
	pollAsBytes, err := stub.GetPrivateData("collectionPoll", pollID)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get state for " + pollID + "\"}")
	} else if pollAsBytes == nil {
		return shim.Error("{\"Error\":\"Vote does not exist: " + pollID + "\"}")
	}

	return shim.Success(pollAsBytes)
}

// ========================================================================
// getPollPrivateDetails - retrieve poll data IPFS CID from chaincode state
// ========================================================================

func (vc *VoteChaincode) getPollPrivateDetails(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting cc key to query")
	}

	iterator, err := stub.GetPrivateDataByPartialCompositeKey("collectionPollPrivateDetails", "poll", args)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get poll private details by partial composite key\"}")
	} else if iterator == nil {
		return shim.Error("{\"Error\":\"Poll private details with partial composite key do not exist\"}")
	}
	defer iterator.Close()

	kv, err := iterator.Next()
	if err != nil {
		return shim.Error("Failed to iterate over iterator: " + err.Error())
	}
	privateDetailsKey := kv.GetKey()

	pollAsBytes, err := stub.GetPrivateData("collectionPollPrivateDetails", privateDetailsKey)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get state for " + privateDetailsKey + "\"}")
	} else if pollAsBytes == nil {
		return shim.Error("{\"Error\":\"Poll does not exist: " + privateDetailsKey + "\"}")
	}

	return shim.Success(pollAsBytes)
}

// =============================================================
// updatePollStatus - change "Status" attribute of a poll object
// =============================================================

func (vc *VoteChaincode) updatePollStatus(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	fmt.Println("- begin update poll")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting poll ID and new poll status")
	}

	pollID := args[0]
	status := args[1]

	var p poll

	existingPollAsBytes, err := stub.GetPrivateData("collectionPoll", pollID)
	if err != nil {
		return shim.Error("Failed to get associated poll: " + err.Error())
	} else if existingPollAsBytes == nil {
		return shim.Error("Poll does not exist: " + pollID)
	}

	err = json.Unmarshal(existingPollAsBytes, &p)
	if err != nil {
		return shim.Error(err.Error())
	}

	p.Status = status
	pollJSONasBytes, err := json.Marshal(p)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutPrivateData("collectionPoll", pollID, pollJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("updateEvent", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end update poll (success)")
	return shim.Success(nil)
}
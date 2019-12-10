package main

import (
	"bytes"
	"fmt"
	"strings"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)


type VoteChaincode struct {
}


type vote struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	VoterID		string 	`json:"voterID"`
	VoterSex 	string 	`json:"voterID"`
	VoterAge	int 	`json:"voterID"`
}

type votePrivateDetails struct {
	ObjectType 	string 	`json:"docType"`
	PollID		string 	`json:"pollID"`
	VoterID		string 	`json:"voterID"`
	VoteHash 	string 	`json:"voteHash"`	// hash(ipfsHash + salt)

}

func main() {
	err := shim.Start(new(VoteChaincode))
	if err != nil {
		fmt.Printf("Error starting Vote chaincode: %s", err)
	}
}

// ============================
// Init - initializes chaincode
// ============================
func (vc *VoteChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// =================================================
// Invoke - starting point for chaincode invocations
// =================================================
func (vc *VoteChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + fn)

	switch fn {
	case "initVote":
		return vc.initVote(stub, args)
	case "getVote":
		return vc.getVote(stub, args)
	case "changeVote":
		return vc.changeVote(stub, args)
	case "queryVotesByPoll":							// parametrized rich query w/ poll ID
		return vc.queryVotesByPoll(stub, args)
	case "queryVotesByVoter":							// parametrized rich query w/ voter ID
		return vc.queryVotesByVoter(stub, args)			
	case "queryVotes":									// ad hoc rich query
		return vc.queryVotes(stub, args)
	}

	fmt.Println("invoke did not find fn: " + fn)
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initVote - create a new vote and store into chaincode state
// ============================================================
func (vc *VoteChaincode) initVote(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// ==== Input sanitation ====
	fmt.Println("- start init vote")
	for _, arg := range args {
		if len(arg) == 0 {
			return shim.Error("Arguments must be non-empty strings")
		}
	}

	pollID := args[0]
	voterID := args[1]
	voteHash := args[2]
	// unique key of the vote data consists of the poll ID + the voter ID
	voteKey := pollID + voterID

	// ==== Check if the vote already exists ====
	voteAsBytes, err := stub.GetState(voteKey)
	if err != nil {
		return shim.Error("Failed to get vote: " + err.Error())
	} else if voteAsBytes != nil {
		return shim.Error("This vote already exists: " + voteKey)
	}

	// ==== Create vote object and marshal to JSON ====
	objectType := "vote"
	vote := &vote{objectType, pollID, voterID, voteHash}
	voteJSONasBytes, err := json.Marshal(vote)
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Put vote to state ====
	err = stub.PutState(voteKey, voteJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end init vote (success)")
	return shim.Success(nil)
}

// =================================================
// getVote - retrieve vote hash from chaincode state
// =================================================

func (vc *VoteChaincode) getVote(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting vote key to query")
	}

	voteKey := args[0]

	// ==== retrieve the vote ====
	voteAsBytes, err := stub.GetState(voteKey)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get state for " + voteKey + "\"}")
	} else if voteAsBytes == nil {
		return shim.Error("{\"Error\":\"Vote does not exist: " + voteKey + "\"}")
	}

	return shim.Success(voteAsBytes)
}

// =================================================
// changeVote - replace vote hash with new vote hash
// =================================================

func (vc *VoteChaincode) changeVote(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting vote key and new vote hash")		
	}

	fmt.Println("- begin change vote")

	voteKey := args[0]
	newVoteHash := args[1]

	// ==== retrieve the old vote ====
	voteAsBytes, err := stub.GetState(voteKey)
	if err != nil {
		return shim.Error("{\"Error\":\"Failed to get state for " + voteKey + "\"}")
	} else if voteAsBytes == nil {
		return shim.Error("{\"Error\":\"Vote does not exist: " + voteKey + "\"}")
	}

	// Unmarshal old vote into a new vote object
	changedVote := vote{}
	err = json.Unmarshal(voteAsBytes, &changedVote)
	if err != nil {
		return shim.Error(err.Error())
	}

	// update new vote object's hash to new hash
	changedVote.VoteHash = newVoteHash

	changedVoteJSONasBytes, err := json.Marshal(changedVote)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(voteKey, changedVoteJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end change vote (success)")
	return shim.Success(changedVoteJSONasBytes)
}

// ===========================================================================================
// Taken from fabric-samples/marbles_chaincode.go.
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// =========================================================================================
// Taken from fabric-samples/marbles_chaincode.go.
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// ===== Parametrized rich queries =========================================================

// =========================================================================================
// queryVotesByPoll takes the poll ID as a parameter, builds a query string using
// the passed poll ID, executes the query, and returns the result set.
// =========================================================================================
func (vc *VoteChaincode) queryVotesByPoll(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: Poll ID")
	}

	pollID := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"vote\",\"pollID\":\"%s\"}}", pollID)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

// =========================================================================================
// queryVotesByPoll takes the voter ID as a parameter, builds a query string using
// the passed voter ID, executes the query, and returns the result set.
// =========================================================================================	
func (vc *VoteChaincode) queryVotesByVoter(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1: Voter ID")
	}

	voterID := strings.ToLower(args[0])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"vote\",\"voterID\":\"%s\"}}", voterID)
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

// ===== Ad hoc rich queries ===============================================================

// =========================================================================================
// Taken from fabric-samples/marbles_chaincode.go.
// queryVotes uses a query string to perform a query for votes.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// =========================================================================================
func (vc *VoteChaincode) queryVotes(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]
	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(queryResults)
}

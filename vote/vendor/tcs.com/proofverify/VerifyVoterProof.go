package proofverify

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

//VerifyVoterProof function receives did as argument and verifies that the identity has required attributes.
func VerifyVoterProof(Did string) (bool, error) {

	//Check validity of arguments
	DidSize := len(Did)
	if DidSize == 0 {
		return false, errors.New("Empty DID received")
	}
	if (DidSize) != 22 {
		return false, errors.New("Size of DID is not 22")
	}

	//Initialize
	type Attributes struct {
		AppName string `json:"app_name"`
		AppID   string `json:"app_id"`
	}
	type IndyResponse struct {
		Status     string     `json:"status"`
		Attributes Attributes `json:"attributes"`
	}
	ProofAttributes := "app_name,app_id"
	RequiredAppName := "voter"
	RequiredAppID := "101"

	//Prepare Payload for Indy
	IndyURL := "http://10.0.1.14:7997/verify_proof"
	Payload := []byte("{\"proof_attr\" : \"" + ProofAttributes + "\",\"their_did\" : \"" + Did + "\"}")
	Request, _ := http.NewRequest("POST", IndyURL, bytes.NewBuffer(Payload))
	Request.Header.Add("content-type", "text/plain")
	Response, err := http.DefaultClient.Do(Request)
	if err != nil || Response == nil || Response.StatusCode != 200 {
		return false, errors.New("Error connecting to Indy Server i!!!!!!!!!!!!!")
	}
	defer Response.Body.Close()

	//Validate Response from Indy
	Body, _ := ioutil.ReadAll(Response.Body)
	ResponseJSON := IndyResponse{}
	err = json.Unmarshal(Body, &ResponseJSON)
	if err != nil {
		return false, errors.New("Error unmarshaling Indy response")
	}
	if ResponseJSON.Status != "true" {
		return false, errors.New("Proof verification failed: attributes missing")
	}
	if !(ResponseJSON.Attributes.AppName == RequiredAppName && ResponseJSON.Attributes.AppID == RequiredAppID) {
		return false, errors.New("Attribute values didn't match")
	}
	return true, nil
}

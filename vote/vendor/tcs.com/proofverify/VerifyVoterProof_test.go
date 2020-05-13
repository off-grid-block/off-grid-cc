package proofverify

import "testing"

func TestEmptyInput(t *testing.T) {

	emptyResult, err := VerifyVoterProof("")
	if emptyResult != false {
		t.Errorf("failed, expected false, but got %v and error is %v", emptyResult, err)
	} else {
		t.Logf("success, expected false and got %v and error is %v", emptyResult, err)
	}
}

func TestInvalidDid(t *testing.T) {

	result, err := VerifyVoterProof("123456789")
	if result != false {
		t.Errorf("failed, expected false, but got %v and error is %v", result, err)
	} else {
		t.Logf("success, expected false and got %v and error is %v", result, err)
	}
}

func TestFailure(t *testing.T) {
	result, err := VerifyVoterProof("YcCtf91BT4MzGKiMtsmGJu")
	if result != false {
		t.Errorf("failed, expected false, but got %v and error is %v", result, err)
	} else {
		t.Logf("success, expected false and got %v and error is %v", result, err)
	}
}

func TestSuccess(t *testing.T) {
	result, err := VerifyVoterProof("5YHg62YsxmfXxFysGsocaE")
	if result != true {
		t.Errorf("failed, expected true, but got %v and error is %v", result, err)
	} else {
		t.Logf("success, expected true and got %true and error is %v", result, err)
	}
}

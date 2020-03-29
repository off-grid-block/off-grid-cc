# Guide

### initVote ###
**Description**:  
Initialize vote object on Fabric network.  

**Usage**: see cc/vote.go  


### GetVote ###
**Description**:  
Reads vote metadata (poll ID, voter ID, voter sex, voter age) stored in collectionVote.  

**Usage**: see cc/vote.go  


### getVotePrivateDetails ###
**Description**:  
Reads a vote's IPFS CID and randomly generated salt, which are stored in collectionVotePrivateDetails.  

**Usage**: see cc/vote.go  


### getVotePrivateDetailsHash ###
**Description**:  
Returns the hash of the collectionVotePrivateDetails entry keyed by the passed poll ID and voter ID args. The hash is retrieved from the channel ledger.  

**Usage**: see cc/vote.go  


### queryVotesByPoll ###
**Description**:  
Returns the metadata of all the votes cast for poll with ID pollID.  

**Usage**: see cc/vote.go  


### queryVotePrivateDetailsByPoll ###
**Description**:  
Returns the IPFS CIDs of all the votes cast for poll with ID pollID.  

**Usage**: see cc/vote.go  


### queryVotesByVoter ###
**Description**:  
Returns the metadata of all the votes cast by voter with ID voterID.  

**Usage**: see cc/vote.go  


### initPoll ###
**Description**:  
Creates a new poll object, splitting the object's metadata into two groups: the poll ID, poll status, and number of votes cast are committed to the private database of collectionPoll (which Org1 and Org2 can both access); the poll ID, randomly generated salt, and IPFS CID of the poll's contents are committed to the private database of collectionPollPrivateDetails (which only Org1 can access).  

**Usage**: see blockchain/functions.go  


### getPoll ###
**Description**:  
Reads poll metadata (poll ID, poll status, number of votes cast) stored in collectionPoll.  

**Usage**: see blockchain/functions.go  


### getPollPrivateDetails ###
**Description**:  
Reads a poll's IPFS CID and randomly generated salt, which are stored in collectionPollPrivateDetails.  

**Usage**: see blockchain/functions.go  


### updatePollStatus ###
**Description**:  
Updates the status of a poll. Possible statuses are: ongoing, paused, closed.  

**Usage**: see blockchain/functions.go  


### queryVotes ###
**Description**:  
Ad hoc rich query.  
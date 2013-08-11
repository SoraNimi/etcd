package main

import (
	"encoding/json"
	"github.com/coreos/go-raft"
	"net/http"
)

//-------------------------------------------------------------
// Handlers to handle raft related request via raft server port
//-------------------------------------------------------------

// Get all the current logs
func GetLogHttpHandler(w http.ResponseWriter, req *http.Request) {
	u, _ := nameToRaftURL(raftServer.Name())
	debugf("[recv] GET %s/log", u)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(raftServer.LogEntries())
}

// Response to vote request
func VoteHttpHandler(w http.ResponseWriter, req *http.Request) {
	rvreq := &raft.RequestVoteRequest{}
	err := decodeJsonRequest(req, rvreq)
	if err == nil {
		u, _ := nameToRaftURL(raftServer.Name())
		debugf("[recv] POST %s/vote [%s]", u, rvreq.CandidateName)
		if resp := raftServer.RequestVote(rvreq); resp != nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
	warnf("[vote] ERROR: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
}

// Response to append entries request
func AppendEntriesHttpHandler(w http.ResponseWriter, req *http.Request) {
	aereq := &raft.AppendEntriesRequest{}
	err := decodeJsonRequest(req, aereq)

	if err == nil {
		u, _ := nameToRaftURL(raftServer.Name())
		debugf("[recv] POST %s/log/append [%d]", u, len(aereq.Entries))
		if resp := raftServer.AppendEntries(aereq); resp != nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
			if !resp.Success {
				debugf("[Append Entry] Step back")
			}
			return
		}
	}
	warnf("[Append Entry] ERROR: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
}

// Response to recover from snapshot request
func SnapshotHttpHandler(w http.ResponseWriter, req *http.Request) {
	aereq := &raft.SnapshotRequest{}
	err := decodeJsonRequest(req, aereq)
	if err == nil {
		u, _ := nameToRaftURL(raftServer.Name())
		debugf("[recv] POST %s/snapshot/ ", u)
		if resp := raftServer.RequestSnapshot(aereq); resp != nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
	warnf("[Snapshot] ERROR: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
}

// Response to recover from snapshot request
func SnapshotRecoveryHttpHandler(w http.ResponseWriter, req *http.Request) {
	aereq := &raft.SnapshotRecoveryRequest{}
	err := decodeJsonRequest(req, aereq)
	if err == nil {
		u, _ := nameToRaftURL(raftServer.Name())
		debugf("[recv] POST %s/snapshotRecovery/ ", u)
		if resp := raftServer.SnapshotRecoveryRequest(aereq); resp != nil {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
	warnf("[Snapshot] ERROR: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
}

// Get the port that listening for etcd connecting of the server
func EtcdURLHttpHandler(w http.ResponseWriter, req *http.Request) {
	u, _ := nameToRaftURL(raftServer.Name())
	debugf("[recv] Get %s/etcdURL/ ", u)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(argInfo.EtcdURL))
}

// Response to the join request
func JoinHttpHandler(w http.ResponseWriter, req *http.Request) {

	command := &JoinCommand{}

	if err := decodeJsonRequest(req, command); err == nil {
		debugf("Receive Join Request from %s", command.Name)
		dispatch(command, &w, req, false)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Response to the name request
func NameHttpHandler(w http.ResponseWriter, req *http.Request) {
	u, _ := nameToRaftURL(raftServer.Name())

	debugf("[recv] Get %s/name/ ", u)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(raftServer.Name()))
}

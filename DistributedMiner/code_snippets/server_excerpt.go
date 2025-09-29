// server_excerpt.go
// Excerpted scheduler/failover logic from a fault-tolerant distributed miner.
// Framework- and course-specific APIs redacted; minimal types stubbed for illustration.

package miner

import (
	"container/list"
	"encoding/json"
	"errors"
	"log"
	"math"
)

type MsgType int

const (
	NodeJoin MsgType = iota
	WorkRequest
	WorkResult
)

type TaskMessage struct {
	Type         MsgType
	Data         []byte // opaque workload payload
	Lower, Upper uint64 // work range [Lower, Upper)
	Hash, Nonce  uint64 // result fields for WorkResult
}

// ReliableServer abstracts the underlying reliable UDP layer (course API redacted).
type ReliableServer interface {
	Read() (connID int, payload []byte, err error)
	Write(connID int, payload []byte) error
	Close() error
}

type server struct {
	net            ReliableServer
	maxChunkSize   uint64
	clientMessages map[int]*TaskMessage // clientID -> remaining request window
	miners         map[int]*Job         // minerID  -> currently assigned job (or nil if idle)
	clients        *list.List           // round-robin queue of clientIDs awaiting chunks
	results        map[int]*TaskMessage // clientID -> best result so far
	idleMinerJobs  *list.List           // jobs to reassign if a miner drops
	log            *log.Logger
}

type Job struct {
	clientConnID int
	upper        uint64
	message      *TaskMessage
}

func newServer(net ReliableServer, logger *log.Logger) (*server, error) {
	if net == nil {
		return nil, errors.New("nil network")
	}
	return &server{
		net:            net,
		maxChunkSize:   10000, // tunable batch size
		clientMessages: make(map[int]*TaskMessage),
		miners:         make(map[int]*Job),
		clients:        list.New().Init(),
		results:        make(map[int]*TaskMessage),
		idleMinerJobs:  list.New().Init(),
		log:            logger,
	}, nil
}

// distributeJob assigns the next available job to miner connID using round-robin chunking.
func distributeJob(srv *server, minerID int) {
	// If there's a previously interrupted job, assign it first.
	if back := srv.idleMinerJobs.Back(); back != nil {
		job := back.Value.(*Job)
		srv.miners[minerID] = job
		if b, err := json.Marshal(job.message); err == nil {
			_ = srv.net.Write(minerID, b)
		}
		srv.idleMinerJobs.Remove(back)
		return
	}

	// Otherwise, pull next client from RR queue and assign a chunk.
	last := srv.clients.Back()
	if last == nil {
		return // nothing to do
	}
	clientID := last.Value.(int)
	req := srv.clientMessages[clientID]
	lower, upper := req.Lower, req.Upper
	newUpper := lower + srv.maxChunkSize
	if newUpper > upper {
		newUpper = upper
	}

	chunk := &TaskMessage{Type: WorkRequest, Data: req.Data, Lower: lower, Upper: newUpper}
	srv.miners[minerID] = &Job{clientConnID: clientID, upper: newUpper, message: chunk}
	req.Lower = newUpper // advance client window

	if b, err := json.Marshal(chunk); err == nil {
		_ = srv.net.Write(minerID, b)
	}

	// If more work remains for this client, rotate them to the front; else remove.
	if req.Lower < req.Upper {
		srv.clients.MoveToFront(last)
	} else {
		srv.clients.Remove(last)
	}
}

func deleteClientFromQueue(srv *server, clientID int) {
	for e := srv.clients.Front(); e != nil; e = e.Next() {
		if e.Value.(int) == clientID {
			srv.clients.Remove(e)
			return
		}
	}
}

func deleteClientIdleJobs(srv *server, clientID int) {
	for e := srv.idleMinerJobs.Front(); e != nil; {
		next := e.Next()
		if e.Value.(*Job).clientConnID == clientID {
			srv.idleMinerJobs.Remove(e)
		}
		e = next
	}
}

// handleDisconnect cleans up on client/miner drop and redistributes any in-flight work.
func handleDisconnect(srv *server, connID int, lastMsg *TaskMessage) {
	switch lastMsg.Type {
	case NodeJoin: // miner dropped
		if j := srv.miners[connID]; j != nil {
			srv.idleMinerJobs.PushFront(j) // requeue interrupted job
		}
		delete(srv.miners, connID)

	case WorkRequest: // client dropped
		delete(srv.clientMessages, connID)
		deleteClientFromQueue(srv, connID)
		deleteClientIdleJobs(srv, connID)

	case WorkResult:
		// treat like miner drop in this handler
	}
}

// mainLoop illustrates message handling
func mainLoop(srv *server) {
	defer srv.net.Close()

	for {
		connID, payload, err := srv.net.Read()
		var msg *TaskMessage
		_ = json.Unmarshal(payload, &msg)

		if err != nil {
			// On read error, infer role from last known message (if any) and clean up.
			if last, ok := srv.clientMessages[connID]; ok {
				handleDisconnect(srv, connID, last)
			}
			continue
		}

		switch msg.Type {
		case NodeJoin: // miner joins
			srv.clientMessages[connID] = msg
			srv.miners[connID] = nil
			distributeJob(srv, connID)

		case WorkRequest: // client submits work
			clientID := connID
			srv.clients.PushFront(clientID)
			srv.clientMessages[clientID] = msg
			// initialize best result
			srv.results[clientID] = &TaskMessage{Type: WorkResult, Hash: math.MaxUint64, Nonce: math.MaxUint64}
			// assign any idle miners
			for minerID := range srv.miners {
				if srv.miners[minerID] == nil {
					distributeJob(srv, minerID)
				}
			}

		case WorkResult: // miner result
			minerID := connID
			job := srv.miners[minerID]
			srv.miners[minerID] = nil
			distributeJob(srv, minerID) // immediately keep miners busy

			// aggregate min hash across chunks
			clientID := job.clientConnID
			if msg.Hash < srv.results[clientID].Hash {
				srv.results[clientID] = msg
			}
			// if this chunk reached the client's original upper bound, send final result
			if job.upper == srv.clientMessages[clientID].Upper {
				final := srv.results[clientID]
				if b, err := json.Marshal(final); err == nil {
					_ = srv.net.Write(clientID, b)
				}
				delete(srv.clientMessages, clientID)
				delete(srv.results, clientID)
			}
		}
	}
}

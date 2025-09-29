// miner_excerpt.go
// Excerpted worker logic from a distributed mining system.
// Course-specific libraries redacted; minimal interfaces/types included for clarity.

package miner

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

// --- Minimal types (replace framework imports) ---

type MsgType int

const (
	NodeJoin MsgType = iota
	WorkRequest
	WorkResult
)

type TaskMessage struct {
	Type         MsgType
	Data         string
	Lower, Upper uint64 // work range: [Lower, Upper) (upper-exclusive)
	Hash, Nonce  uint64 // filled for WorkResult
}

type ReliableClient interface {
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
}

// dial is a placeholder for a reliable transport (e.g., UDP+reliability layer).
func dial(hostport string) (ReliableClient, error) {
	return nil, errors.New("transport redacted in public excerpt")
}

// --- Core Logic ---

// hash is a stand-in for the assignment's PoW hash function.
func hash(data string, nonce uint64) uint64 {
	h := sha256.New()
	_, _ = h.Write([]byte(data))
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], nonce)
	_, _ = h.Write(b[:])
	sum := h.Sum(nil)
	return binary.LittleEndian.Uint64(sum[:8]) // take 64 bits
}

func mine(data string, lower, upper uint64) (minHash, bestNonce uint64) {
	minHash = ^uint64(0)             // max uint64
	for n := lower; n < upper; n++ { // upper-exclusive
		if h := hash(data, n); h < minHash {
			minHash, bestNonce = h, n
		}
	}
	return
}

func runMiner(hostport string) error {
	client, err := dial(hostport)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer client.Close()

	// Example: send a join message (protocol details redacted).
	_ = client.Write([]byte(`{"Type":0}`)) // NodeJoin

	for {
		b, err := client.Read()
		if err != nil {
			// read error: server closed or network issue (retry/backoff omitted for brevity)
			return err
		}

		var msg TaskMessage
		if err := json.Unmarshal(b, &msg); err != nil {
			continue
		}
		if msg.Type != WorkRequest {
			continue
		}

		h, nonce := mine(msg.Data, msg.Lower, msg.Upper)
		result := TaskMessage{Type: WorkResult, Hash: h, Nonce: nonce}
		out, _ := json.Marshal(result)
		if err := client.Write(out); err != nil {
			return err
		}
	}
}

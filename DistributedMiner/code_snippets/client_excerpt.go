// client_excerpt.go
// Excerpted request/response client for a distributed miner.
// Course-specific libraries redacted; minimal interfaces/types included.

package miner

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// --- Minimal protocol & transport stubs ---

type MsgType int

const (
	WorkRequest MsgType = iota
	WorkResult
)

type TaskMessage struct {
	Type         MsgType
	Data         string
	Lower, Upper uint64 // work range [Lower, Upper) (upper-exclusive)
	Hash, Nonce  uint64 // set on WorkResult
}

type ReliableClient interface {
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
}

// dial is intentionally redacted for public posting.
func dial(hostport string) (ReliableClient, error) {
	return nil, errors.New("transport redacted in public excerpt")
}

// --- Client Orchestration ---

func runClient(hostport, message string, maxNonce uint64) error {
	c, err := dial(hostport)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer c.Close()

	req := TaskMessage{
		Type:  WorkRequest,
		Data:  message,
		Lower: 0,
		Upper: maxNonce, // upper-exclusive; pass maxNonce+1 for inclusive semantics
	}
	payload, _ := json.Marshal(req)
	if err := c.Write(payload); err != nil {
		return fmt.Errorf("write request: %w", err)
	}

	b, err := c.Read()
	if err != nil {
		return fmt.Errorf("read result: %w", err)
	}
	var res TaskMessage
	if err := json.Unmarshal(b, &res); err != nil || res.Type != WorkResult {
		return fmt.Errorf("invalid result")
	}

	// Human-friendly output
	fmt.Printf("MinHash=%d Nonce=%d\n", res.Hash, res.Nonce)
	return nil
}

func main() {
	// Usage: client <host:port> <message> <maxNonce>
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <hostport> <message> <maxNonce>\n", os.Args[0])
	}
	flag.Parse()
	if flag.NArg() != 3 {
		flag.Usage()
		os.Exit(2)
	}
	hostport := flag.Arg(0)
	message := flag.Arg(1)
	maxN, err := strconv.ParseUint(flag.Arg(2), 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "maxNonce must be uint64: %v\n", err)
		os.Exit(2)
	}
	if err := runClient(hostport, message, maxN); err != nil {
		fmt.Fprintf(os.Stderr, "client error: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"log"
	"time"

	socketio "github.com/SavvasMohito/go-socket.io-client"
)

// IxClient implements an Ix client similar to your Python version.
type IxClient struct {
	url         string
	username    string
	password    string
	client      *socketio.Client
	isConnected bool
	relationID  string
	dataLogging bool
}

// IxOption defines a function type for setting options
type IxOption func(*IxClient)

// WithDataLogging enables data logging
func WithDataLogging(enable bool) IxOption {
	return func(ix *IxClient) {
		ix.dataLogging = enable
	}
}

// WithRelationID sets the relation ID
func WithRelationID(relationID string) IxOption {
	return func(ix *IxClient) {
		if relationID != "" {
			ix.relationID = relationID
		}
	}
}

// NewIxClient returns a new IxClient instance.
func NewIxClient(url, username, password string, opts ...IxOption) *IxClient {
	ix := &IxClient{
		url:         url,
		username:    username,
		password:    password,
		relationID:  "",
		dataLogging: false,
	}

	// Apply options
	for _, opt := range opts {
		opt(ix)
	}

	auth := map[string]string{"ix_username": username, "ix_password": password}
	if ix.relationID != "" {
		auth = map[string]string{"ix_username": username, "ix_password": password, "relation_id": ix.relationID}
	}

	socketOpts := &socketio.ClientOptions{
		Path: "/internal/ws",
		Auth: auth,
	}

	client, err := socketio.NewClient(url, socketOpts)
	if err != nil {
		log.Fatalf("ERROR: Failed to create Ix client: %v", err)
	}

	ix.client = client

	return ix
}

// Connect initiates a connection to the Ix server.
func (ix *IxClient) Connect() error {
	err := ix.client.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	// Wait briefly to ensure handlers are registered
	time.Sleep(100 * time.Millisecond)
	log.Println("INFO: Connected to Ix!")
	ix.isConnected = true
	return nil
}

// Disconnect terminates the Ix Client's connection.
func (ix *IxClient) Disconnect() {
	if ix.client != nil && ix.isConnected {
		log.Println("INFO: Disconnected from Ix.")
		ix.client.Close()
	}
}

// Send transmits a map of data to the Ix Interface.
// (If using a relation connection, the data is sent using the configured emit path.)
func (ix *IxClient) Send(data map[string]interface{}) {
	if !ix.isConnected {
		log.Println("ERROR: Please make a client connection before sending any data.")
		return
	}
	if data == nil {
		log.Println("ERROR: Data cannot be empty. Please send a valid map.")
		return
	}

	emitPath := "xapp_local_emit"
	if ix.relationID != "" {
		emitPath = "xapp_relation_emit"
	}

	err := ix.client.Emit(emitPath, data)
	if err != nil {
		log.Printf("ERROR: Failed to emit data: %v", err)
		return
	}

	if ix.dataLogging {
		suffix := ""
		if ix.relationID != "" {
			suffix = fmt.Sprintf("to relation %s", ix.relationID)
		}
		log.Printf("INFO: Successfully emitted %v %s\n", data, suffix)
	}
}

package main

import (
	"log"
	"time"
)

func main() {
	client := NewIxClient(
		"http://10.68.123.135:80",
		"go_test",
		"041Rz-o1tzP--WrTnywGvw",
		WithDataLogging(true),                 // This will log all data being sent
		WithRelationID("rel28edd4a35aa553e5"), // This will forward data to a remote relation
	)

	// Connect and wait for connection to establish
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Wait a bit longer to ensure connection is ready
	time.Sleep(1 * time.Second)

	testData := map[string]interface{}{
		"attenuation": 69,
		"health":      0.420,
	}

	// Send data multiple times to ensure it's getting through
	for i := 0; i < 3; i++ {
		client.Send(testData)
		time.Sleep(500 * time.Millisecond)
	}

	// Wait longer before disconnecting
	time.Sleep(2 * time.Second)
	client.Disconnect()
}

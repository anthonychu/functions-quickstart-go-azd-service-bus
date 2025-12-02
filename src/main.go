package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// InvokeRequest represents the request payload from the Azure Functions host
type InvokeRequest struct {
	Data     map[string]json.RawMessage `json:"Data"`
	Metadata map[string]interface{}     `json:"Metadata"`
}

// InvokeResponse represents the response payload to the Azure Functions host
type InvokeResponse struct {
	Outputs     map[string]interface{} `json:"Outputs"`
	Logs        []string               `json:"Logs"`
	ReturnValue interface{}            `json:"ReturnValue"`
}

func main() {
	// Get the port from environment variable
	port := os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if port == "" {
		port = "8080"
	}

	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create a Gin router
	router := gin.Default()

	// Register the Service Bus queue trigger handler
	router.POST("/serviceBusQueueTrigger", serviceBusQueueTriggerHandler)
	fmt.Printf("Go server listening on port: %s\n", port)
	if err := router.Run(":" + port); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}
}

// serviceBusQueueTriggerHandler handles Service Bus queue trigger invocations
func serviceBusQueueTriggerHandler(c *gin.Context) {
	var invokeRequest InvokeRequest

	// Parse the incoming request
	if err := c.ShouldBindJSON(&invokeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Extract the message from the request data
	var message interface{}
	if msgData, ok := invokeRequest.Data["message"]; ok {
		if err := json.Unmarshal(msgData, &message); err != nil {
			// If it's not JSON, treat it as a string
			message = string(msgData)
		}
	}

	logs := []string{}
	logs = append(logs, "Go ServiceBus Queue trigger start processing a message")
	fmt.Printf("Go ServiceBus Queue trigger start processing a message: %v\n", message)

	// Log message metadata if available
	if messageId, ok := invokeRequest.Metadata["MessageId"]; ok {
		fmt.Printf("MessageId: %v\n", messageId)
		logs = append(logs, "MessageId: "+messageId.(string))
	}
	if enqueuedTime, ok := invokeRequest.Metadata["EnqueuedTimeUtc"]; ok {
		fmt.Printf("EnqueuedTimeUtc: %v\n", enqueuedTime)
	}
	if deliveryCount, ok := invokeRequest.Metadata["DeliveryCount"]; ok {
		fmt.Printf("DeliveryCount: %v\n", deliveryCount)
	}

	// Simulate 5-second processing time
	time.Sleep(5 * time.Second)

	logs = append(logs, "Go ServiceBus Queue trigger end processing a message")
	fmt.Printf("Go ServiceBus Queue trigger end processing a message\n")

	// Create the response
	invokeResponse := InvokeResponse{
		Outputs:     make(map[string]interface{}),
		Logs:        logs,
		ReturnValue: nil,
	}

	c.JSON(http.StatusOK, invokeResponse)
}

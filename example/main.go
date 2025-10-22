package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	graphiti "github.com/pentagi/graphiti-go-client"
)

// This example demonstrates how to use the Graphiti Go client.
//
// Important: The /messages endpoint processes data asynchronously. This example
// polls for episodes to verify data was successfully created before searching.
//
// Troubleshooting: If you see "No episodes were created" errors:
// 1. Check server logs for "Error executing Neo4j query: Driver closed"
// 2. Ensure Neo4j is running and properly configured
// 3. Verify the Graphiti server has a persistent database connection
// 4. Check that the async worker is processing jobs successfully

func main() {
	// Create a client with extended timeout for long-running operations
	client := graphiti.NewClient("http://localhost:8000", graphiti.WithTimeout(60*time.Second))

	// Health check
	fmt.Println("=== Health Check ===")
	health, err := client.HealthCheck()
	if err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	fmt.Printf("Status: %s\n\n", health.Status)

	// Create a unique group ID for this example
	groupID := uuid.New().String()
	fmt.Printf("Using group ID: %s\n\n", groupID)

	// Add messages
	fmt.Println("=== Adding Messages ===")
	messages := []graphiti.Message{
		{
			Content:   "I love hiking in the mountains on weekends.",
			Author:    "Alice",
			Timestamp: time.Now().Add(-2 * time.Hour),
		},
		{
			Content:   "That sounds great! Do you have a favorite trail?",
			Author:    "Assistant",
			Timestamp: time.Now().Add(-90 * time.Minute),
		},
		{
			Content:   "Yes, I particularly enjoy the Pacific Crest Trail. I try to go there every summer.",
			Author:    "Alice",
			Timestamp: time.Now().Add(-60 * time.Minute),
		},
	}

	addResult, err := client.AddMessages(graphiti.AddMessagesRequest{
		GroupID:  groupID,
		Messages: messages,
	})
	if err != nil {
		log.Fatalf("Failed to add messages: %v", err)
	}
	fmt.Printf("%s: %v\n\n", addResult.Message, addResult.Success)

	// Wait for processing and verify data exists (poll for episodes)
	fmt.Println("Waiting for messages to be processed...")
	maxAttempts := 10
	pollInterval := 5 * time.Second
	var episodes []graphiti.Episode

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("  Polling for episodes (attempt %d/%d)...\n", attempt, maxAttempts)
		episodes, err = client.GetEpisodes(groupID, 10)
		if err != nil {
			log.Printf("  Warning: Failed to get episodes: %v", err)
		} else if len(episodes) > 0 {
			fmt.Printf("  ✓ Found %d episodes, processing complete!\n\n", len(episodes))
			break
		}

		if attempt < maxAttempts {
			time.Sleep(pollInterval)
		}
	}

	if len(episodes) == 0 {
		log.Fatalf("Timeout: No episodes were created after %v. The async job may have failed.", time.Duration(maxAttempts)*pollInterval)
	}

	// Search for facts
	fmt.Println("=== Searching for Facts ===")
	searchResult, err := client.Search(graphiti.SearchQuery{
		Query:    "What does the user like to do?",
		MaxFacts: 5,
		GroupIDs: &[]string{groupID},
	})
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	fmt.Printf("Found %d facts:\n", len(searchResult.Facts))
	for i, fact := range searchResult.Facts {
		fmt.Printf("%d. %s\n   (from: %s, created: %s)\n",
			i+1, fact.Fact, fact.Name, fact.CreatedAt.Format(time.RFC3339))
	}
	fmt.Println()

	// Get memory from messages
	fmt.Println("=== Getting Memory ===")
	memoryMessages := []graphiti.Message{
		{
			Content:   "What hobbies does the user have?",
			Author:    "User",
			Timestamp: time.Now(),
		},
	}
	memoryResponse, err := client.GetMemory(graphiti.GetMemoryRequest{
		GroupID:  groupID,
		MaxFacts: 10,
		Messages: memoryMessages,
	})
	if err != nil {
		log.Fatalf("Failed to get memory: %v", err)
	}
	fmt.Printf("Retrieved %d facts from memory:\n", len(memoryResponse.Facts))
	for i, fact := range memoryResponse.Facts {
		fmt.Printf("%d. %s\n", i+1, fact.Fact)
	}
	fmt.Println()

	// Add an entity node
	fmt.Println("=== Adding Entity Node ===")
	entityUUID := uuid.New().String()
	node, err := client.AddEntityNode(graphiti.AddEntityNodeRequest{
		UUID:    entityUUID,
		GroupID: groupID,
		Name:    "User Interests",
		Summary: "The user's hobbies and interests",
	})
	if err != nil {
		log.Fatalf("Failed to add entity node: %v", err)
	}
	fmt.Printf("Created entity node: %s (UUID: %s)\n\n", node.Name, node.UUID)

	// Display episodes (already fetched during polling)
	fmt.Println("=== Episodes Summary ===")
	fmt.Printf("Total episodes: %d\n", len(episodes))
	for i, episode := range episodes {
		fmt.Printf("%d. %s: %s\n", i+1, episode.Name, episode.Content)
	}
	fmt.Println()

	// Cleanup: delete the group
	fmt.Println("=== Cleanup ===")
	deleteResult, err := client.DeleteGroup(groupID)
	if err != nil {
		log.Printf("Warning: Failed to delete group: %v", err)
	} else {
		fmt.Printf("%s: %v\n", deleteResult.Message, deleteResult.Success)
	}
}

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	graphiti "github.com/vxcontrol/graphiti-go-client"
)

// This example demonstrates the advanced search capabilities of the Graphiti Go client.
// It shows how to use specialized search methods for temporal queries, entity relationships,
// diverse results, episodes, successful tools, recent context, and entity label filtering.

func main() {
	// Create a client
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

	// Add some security test messages to demonstrate advanced search
	fmt.Println("=== Adding Security Test Messages ===")
	messages := []graphiti.Message{
		{
			Content:   "Performed nmap scan on target 192.168.1.100. Found open ports: 22 (SSH), 80 (HTTP), 443 (HTTPS).",
			Author:    "pentester",
			Timestamp: time.Now().Add(-2 * time.Hour),
		},
		{
			Content:   "Detected Apache 2.4.41 running on port 80 with potential CVE-2021-41773 vulnerability.",
			Author:    "scanner",
			Timestamp: time.Now().Add(-90 * time.Minute),
		},
		{
			Content:   "Successfully exploited CVE-2021-41773 using metasploit. Gained shell access.",
			Author:    "exploit_tool",
			Timestamp: time.Now().Add(-60 * time.Minute),
		},
		{
			Content:   "Performed privilege escalation using CVE-2021-3156 (sudo vulnerability).",
			Author:    "exploit_tool",
			Timestamp: time.Now().Add(-45 * time.Minute),
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

	// Wait for processing
	fmt.Println("Waiting for messages to be processed...")
	maxAttempts := 12
	pollInterval := 5 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("  Polling (attempt %d/%d)...\n", attempt, maxAttempts)
		episodes, err := client.GetEpisodes(groupID, 10)
		if err != nil {
			log.Printf("  Warning: Failed to get episodes: %v", err)
		} else if len(episodes) > 0 {
			fmt.Printf("  ✓ Found %d episodes, ready for advanced search!\n\n", len(episodes))
			break
		}

		if attempt == maxAttempts {
			log.Printf("Warning: No episodes found after waiting. Continuing anyway...\n\n")
		}
		time.Sleep(pollInterval)
	}

	// Example 1: Temporal Window Search
	fmt.Println("=== Example 1: Temporal Window Search ===")
	fmt.Println("Searching for activities in the last 3 hours...")
	temporalResult, err := client.TemporalWindowSearch(graphiti.TemporalSearchRequest{
		Query:      "What security tests were performed?",
		GroupID:    &groupID,
		TimeStart:  time.Now().Add(-3 * time.Hour),
		TimeEnd:    time.Now(),
		MaxResults: 10,
	})
	if err != nil {
		log.Printf("Temporal search failed: %v\n", err)
	} else {
		fmt.Printf("Found %d edges, %d nodes, %d episodes\n",
			len(temporalResult.Edges), len(temporalResult.Nodes), len(temporalResult.Episodes))
		fmt.Printf("Time window: %s to %s\n",
			temporalResult.TimeWindow.Start.Format(time.RFC3339),
			temporalResult.TimeWindow.End.Format(time.RFC3339))
		for i, episode := range temporalResult.Episodes {
			score := 0.0
			if i < len(temporalResult.EpisodeScores) {
				score = temporalResult.EpisodeScores[i]
			}
			fmt.Printf("  - Episode (score: %.2f): %s\n", score, episode.Content[:min(80, len(episode.Content))])
		}
	}
	fmt.Println()

	// Example 2: Diverse Results Search
	fmt.Println("=== Example 2: Diverse Results Search ===")
	fmt.Println("Getting diverse attack vectors with high diversity...")
	diverseResult, err := client.DiverseResultsSearch(graphiti.DiverseSearchRequest{
		Query:          "Show different attack methods and vulnerabilities",
		GroupID:        &groupID,
		DiversityLevel: "high",
		MaxResults:     10,
	})
	if err != nil {
		log.Printf("Diverse search failed: %v\n", err)
	} else {
		fmt.Printf("Found %d diverse edges, %d nodes, %d episodes\n",
			len(diverseResult.Edges), len(diverseResult.Nodes), len(diverseResult.Episodes))
		for i, edge := range diverseResult.Edges {
			score := 0.0
			if i < len(diverseResult.EdgeMMRScores) {
				score = diverseResult.EdgeMMRScores[i]
			}
			fmt.Printf("  - Edge (MMR: %.2f): %s\n", score, edge.Fact[:min(80, len(edge.Fact))])
		}
	}
	fmt.Println()

	// Example 3: Episode Context Search
	fmt.Println("=== Example 3: Episode Context Search ===")
	fmt.Println("Searching for tool execution records...")
	episodeResult, err := client.EpisodeContextSearch(graphiti.EpisodeContextSearchRequest{
		Query:            "Show tool execution results",
		GroupID:          &groupID,
		IncludeToolCalls: true,
		MaxResults:       10,
	})
	if err != nil {
		log.Printf("Episode search failed: %v\n", err)
	} else {
		fmt.Printf("Found %d episodes with %d mentioned nodes\n",
			len(episodeResult.Episodes), len(episodeResult.MentionedNodes))
		for i, episode := range episodeResult.Episodes {
			score := 0.0
			if i < len(episodeResult.RerankerScores) {
				score = episodeResult.RerankerScores[i]
			}
			fmt.Printf("  - Episode (score: %.2f) from %s: %s\n",
				score, episode.Source, episode.Content[:min(80, len(episode.Content))])
		}
	}
	fmt.Println()

	// Example 4: Successful Tools Search
	fmt.Println("=== Example 4: Successful Tools Search ===")
	fmt.Println("Finding successful exploits with min 1 mention...")
	toolsResult, err := client.SuccessfulToolsSearch(graphiti.SuccessfulToolsSearchRequest{
		Query:       "Find successful exploits",
		GroupID:     &groupID,
		MinMentions: 1,
		MaxResults:  10,
	})
	if err != nil {
		log.Printf("Successful tools search failed: %v\n", err)
	} else {
		fmt.Printf("Found %d edges with sufficient mentions\n", len(toolsResult.Edges))
		for i, edge := range toolsResult.Edges {
			mentions := 0.0
			if i < len(toolsResult.EdgeMentionCounts) {
				mentions = toolsResult.EdgeMentionCounts[i]
			}
			fmt.Printf("  - Edge (mentions: %.0f): %s\n", mentions, edge.Fact[:min(80, len(edge.Fact))])
		}
	}
	fmt.Println()

	// Example 5: Recent Context Search
	fmt.Println("=== Example 5: Recent Context Search ===")
	fmt.Println("Getting recent discoveries from the last 6 hours...")
	recentResult, err := client.RecentContextSearch(graphiti.RecentContextSearchRequest{
		Query:         "What was discovered recently?",
		GroupID:       &groupID,
		RecencyWindow: "6h",
		MaxResults:    10,
	})
	if err != nil {
		log.Printf("Recent context search failed: %v\n", err)
	} else {
		fmt.Printf("Found %d recent edges, %d nodes, %d episodes\n",
			len(recentResult.Edges), len(recentResult.Nodes), len(recentResult.Episodes))
		fmt.Printf("Recency window: %s to %s\n",
			recentResult.TimeWindow.Start.Format(time.RFC3339),
			recentResult.TimeWindow.End.Format(time.RFC3339))
		for i, edge := range recentResult.Edges {
			score := 0.0
			if i < len(recentResult.EdgeScores) {
				score = recentResult.EdgeScores[i]
			}
			fmt.Printf("  - Recent edge (score: %.2f): %s\n", score, edge.Fact[:min(80, len(edge.Fact))])
		}
	}
	fmt.Println()

	// Example 6: Entity By Label Search
	fmt.Println("=== Example 6: Entity By Label Search ===")
	fmt.Println("Searching for SERVICE and VULNERABILITY entities...")
	labelResult, err := client.EntityByLabelSearch(graphiti.EntityByLabelSearchRequest{
		Query:      "Find vulnerable services",
		GroupID:    &groupID,
		NodeLabels: []string{"SERVICE", "VULNERABILITY", "IP_ADDRESS"},
		MaxResults: 20,
	})
	if err != nil {
		log.Printf("Entity label search failed: %v\n", err)
	} else {
		fmt.Printf("Found %d nodes, %d edges\n", len(labelResult.Nodes), len(labelResult.Edges))
		for i, node := range labelResult.Nodes {
			score := 0.0
			if i < len(labelResult.NodeScores) {
				score = labelResult.NodeScores[i]
			}
			fmt.Printf("  - Node (score: %.2f): %s [labels: %v]\n",
				score, node.Name, node.Labels)
		}
	}
	fmt.Println()

	// Optional: Entity Relationships Search (requires a known entity UUID)
	// This example shows the pattern, but won't execute without a real entity UUID
	fmt.Println("=== Example 7: Entity Relationships Search (Pattern) ===")
	fmt.Println("To use EntityRelationshipsSearch, you need a center node UUID.")
	fmt.Println("Example usage:")
	fmt.Println(`
  result, err := client.EntityRelationshipsSearch(graphiti.EntityRelationshipSearchRequest{
      Query:          "What is connected to this service?",
      GroupID:        &groupID,
      CenterNodeUUID: "your-entity-uuid-here",
      MaxDepth:       2,
      NodeLabels:     &[]string{"VULNERABILITY", "EXPLOIT"},
      MaxResults:     20,
  })
  if err != nil {
      log.Printf("Entity relationships search failed: %v\n", err)
  } else {
      fmt.Printf("Found %d related nodes at various distances\n", len(result.Nodes))
  }
	`)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

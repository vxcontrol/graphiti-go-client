# Graphiti Go Client

A Go client library for the Graphiti HTTP API.

## Installation

```bash
go get github.com/vxcontrol/graphiti-go-client
```

## Usage

### Creating a Client

```go
import (
    "time"
    graphiti "github.com/vxcontrol/graphiti-go-client"
)

// Create a client with default settings
client := graphiti.NewClient("http://localhost:8000")

// Create a client with custom timeout
client := graphiti.NewClient("http://localhost:8000",
    graphiti.WithTimeout(60 * time.Second))

// Create a client with a custom HTTP client
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}
client := graphiti.NewClient("http://localhost:8000",
    graphiti.WithHTTPClient(httpClient))
```

### Health Check

```go
health, err := client.HealthCheck()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Status: %s\n", health.Status)
```

### Search for Facts

```go
result, err := client.Search(graphiti.SearchQuery{
    Query:    "Tell me about user preferences",
    MaxFacts: 10,
})
if err != nil {
    log.Fatal(err)
}

for _, fact := range result.Facts {
    fmt.Printf("Fact: %s\n", fact.Fact)
}
```

### Search with Group Filtering

```go
groupIDs := []string{"group-123", "group-456"}
result, err := client.Search(graphiti.SearchQuery{
    GroupIDs: &groupIDs,
    Query:    "user settings",
    MaxFacts: 5,
})
```

### Add Messages

**⚠️ Important:** The `/messages` endpoint is asynchronous. Messages are queued and processed by a background worker. Data may not be immediately available after this call returns.

```go
messages := []graphiti.Message{
    {
        Content:   "Hello, how are you?",
        Author:    "User",
        Timestamp: time.Now(),
    },
    {
        Content:   "I'm doing great, thank you!",
        Author:    "Assistant",
        Timestamp: time.Now(),
    },
}

result, err := client.AddMessages(graphiti.AddMessagesRequest{
    GroupID:  "my-group-id",
    Messages: messages,
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%s: %v\n", result.Message, result.Success)

// Wait for processing by polling for episodes
maxAttempts := 10
for attempt := 1; attempt <= maxAttempts; attempt++ {
    episodes, err := client.GetEpisodes("my-group-id", 10)
    if err == nil && len(episodes) > 0 {
        fmt.Println("Messages processed successfully!")
        break
    }
    time.Sleep(5 * time.Second)
}
```

### Add an Entity Node

```go
uuid := "entity-uuid-123"
node, err := client.AddEntityNode(graphiti.AddEntityNodeRequest{
    UUID:    uuid,
    GroupID: "my-group-id",
    Name:    "User Preferences",
    Summary: "Contains user's preferred settings",
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created node: %s\n", node.UUID)
```

### Get Memory from Messages

```go
messages := []graphiti.Message{
    {
        Content:   "What were my settings?",
        Author:    "User",
        Timestamp: time.Now(),
    },
}

response, err := client.GetMemory(graphiti.GetMemoryRequest{
    GroupID:  "my-group-id",
    MaxFacts: 10,
    Messages: messages,
})
if err != nil {
    log.Fatal(err)
}

for _, fact := range response.Facts {
    fmt.Printf("Fact: %s (from %s)\n", fact.Fact, fact.Name)
}
```

### Get Episodes

```go
episodes, err := client.GetEpisodes("my-group-id", 5)
if err != nil {
    log.Fatal(err)
}

for _, episode := range episodes {
    fmt.Printf("Episode: %s - %s\n", episode.Name, episode.Content)
}
```

### Get a Specific Entity Edge

```go
fact, err := client.GetEntityEdge("edge-uuid-123")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Fact: %s\n", fact.Fact)
```

### Delete Operations

```go
// Delete an entity edge
result, err := client.DeleteEntityEdge("edge-uuid-123")

// Delete an episode
result, err := client.DeleteEpisode("episode-uuid-123")

// Delete a group
result, err := client.DeleteGroup("group-id-123")

// Clear all data (use with caution!)
result, err := client.Clear()
```

## Advanced Search Methods

The client provides specialized search methods for advanced querying and analysis.

### Temporal Window Search

Search for all relevant context (facts, entities, and agent responses) from a specific time window:

```go
result, err := client.TemporalWindowSearch(graphiti.TemporalSearchRequest{
    Query:      "What attacks were performed?",
    GroupID:    &groupID,
    TimeStart:  time.Now().Add(-24 * time.Hour),
    TimeEnd:    time.Now(),
    MaxResults: 15,
})
if err != nil {
    log.Fatal(err)
}

for _, edge := range result.Edges {
    fmt.Printf("Fact: %s (score: %.2f)\n", edge.Fact, result.EdgeScores[i])
}
```

### Entity Relationships Search

Find all relationships and related entities starting from a specific discovered entity using graph traversal:

```go
result, err := client.EntityRelationshipsSearch(graphiti.EntityRelationshipSearchRequest{
    Query:          "What vulnerabilities are related to this service?",
    GroupID:        &groupID,
    CenterNodeUUID: "service-uuid-123",
    MaxDepth:       2,
    NodeLabels:     &[]string{"VULNERABILITY", "EXPLOIT"},
    MaxResults:     20,
})
if err != nil {
    log.Fatal(err)
}

if result.CenterNode != nil {
    fmt.Printf("Center: %s\n", result.CenterNode.Name)
}
for i, node := range result.Nodes {
    fmt.Printf("Related: %s (distance: %.2f)\n", node.Name, result.NodeDistances[i])
}
```

### Diverse Results Search

Get diverse, non-redundant results using Maximal Marginal Relevance (MMR) to prevent receiving repetitive information:

```go
result, err := client.DiverseResultsSearch(graphiti.DiverseSearchRequest{
    Query:          "Find different attack vectors",
    GroupID:        &groupID,
    DiversityLevel: "high", // "low", "medium", or "high"
    MaxResults:     10,
})
if err != nil {
    log.Fatal(err)
}

for i, edge := range result.Edges {
    fmt.Printf("Fact: %s (MMR score: %.2f)\n", edge.Fact, result.EdgeMMRScores[i])
}
```

### Episode Context Search

Search through complete agent responses, reasoning, and tool execution records:

```go
result, err := client.EpisodeContextSearch(graphiti.EpisodeContextSearchRequest{
    Query:            "Show me nmap scan results",
    GroupID:          &groupID,
    AgentTypes:       &[]string{"pentester"},
    IncludeToolCalls: true,
    MaxResults:       10,
})
if err != nil {
    log.Fatal(err)
}

for i, episode := range result.Episodes {
    fmt.Printf("Episode: %s (score: %.2f)\n", episode.Content, result.RerankerScores[i])
}
```

### Successful Tools Search

Find successful tool executions and attack patterns, prioritizing facts that led to successful exploitation:

```go
result, err := client.SuccessfulToolsSearch(graphiti.SuccessfulToolsSearchRequest{
    Query:       "Find successful exploits",
    GroupID:     &groupID,
    ToolNames:   &[]string{"metasploit", "sqlmap"},
    MinMentions: 2,
    MaxResults:  15,
})
if err != nil {
    log.Fatal(err)
}

for i, edge := range result.Edges {
    fmt.Printf("Fact: %s (mentions: %.0f)\n", edge.Fact, result.EdgeMentionCounts[i])
}
```

### Recent Context Search

Retrieve the most recent relevant context, biased toward recent actions and discoveries:

```go
result, err := client.RecentContextSearch(graphiti.RecentContextSearchRequest{
    Query:         "What was discovered recently?",
    GroupID:       &groupID,
    RecencyWindow: "24h", // "1h", "6h", "24h", or "7d"
    MaxResults:    10,
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Searching from %s to %s\n", result.TimeWindow.Start, result.TimeWindow.End)
for i, edge := range result.Edges {
    fmt.Printf("Recent fact: %s (score: %.2f)\n", edge.Fact, result.EdgeScores[i])
}
```

### Entity By Label Search

Search for specific entity types (IPs, services, vulnerabilities, tools, etc.) with label-based filtering:

```go
result, err := client.EntityByLabelSearch(graphiti.EntityByLabelSearchRequest{
    Query:      "Find all vulnerable services",
    GroupID:    &groupID,
    NodeLabels: []string{"SERVICE", "VULNERABILITY"},
    EdgeTypes:  &[]string{"HAS_VULNERABILITY", "EXPLOITS"},
    MaxResults: 25,
})
if err != nil {
    log.Fatal(err)
}

for i, node := range result.Nodes {
    fmt.Printf("Entity: %s [%s] (score: %.2f)\n", 
        node.Name, strings.Join(node.Labels, ", "), result.NodeScores[i])
}
```

## Types

### Message

```go
type Message struct {
    Content           string    // The message content
    UUID              *string   // Optional UUID
    Name              string    // Optional name for episodic node
    Author            string    // The author/entity that created this message
    Timestamp         time.Time // Message timestamp
    SourceDescription string    // Optional source description
}
```

### SearchQuery

```go
type SearchQuery struct {
    GroupIDs *[]string // Optional group IDs to filter
    Query    string    // Search query text
    MaxFacts int       // Maximum number of facts to return (default: 10)
}
```

### FactResult

```go
type FactResult struct {
    UUID      string     // Unique identifier
    Name      string     // Fact name
    Fact      string     // The actual fact text
    ValidAt   *time.Time // When fact became valid
    InvalidAt *time.Time // When fact became invalid
    CreatedAt time.Time  // Creation timestamp
    ExpiredAt *time.Time // Expiration timestamp
}
```

### Advanced Search Types

#### NodeResult

```go
type NodeResult struct {
    UUID      string    // Node UUID
    Name      string    // Entity name
    Labels    []string  // Entity type labels (e.g., ["SERVICE", "WEB"])
    Summary   string    // Node summary/description
    CreatedAt time.Time // Creation timestamp
}
```

#### EdgeResult

```go
type EdgeResult struct {
    UUID           string     // Edge UUID
    Name           string     // Relationship name
    Fact           string     // The fact/relationship description
    SourceNodeUUID string     // Source entity UUID
    TargetNodeUUID string     // Target entity UUID
    ValidAt        *time.Time // When relationship became valid
    InvalidAt      *time.Time // When relationship became invalid
    CreatedAt      time.Time  // Creation timestamp
    ExpiredAt      *time.Time // Expiration timestamp
}
```

#### EpisodeResult

```go
type EpisodeResult struct {
    UUID              string    // Episode UUID
    Content           string    // Episode content (agent response, tool output, etc.)
    Source            string    // Source type (e.g., "tool", "agent")
    SourceDescription string    // Detailed source description
    CreatedAt         time.Time // Creation timestamp
    ValidAt           time.Time // When episode occurred
}
```

#### CommunityResult

```go
type CommunityResult struct {
    UUID      string    // Community UUID
    Name      string    // Community name
    Summary   string    // Community summary
    CreatedAt time.Time // Creation timestamp
}
```

## Error Handling

All client methods return an error as the last return value. Always check for errors:

```go
result, err := client.Search(query)
if err != nil {
    // Handle error
    log.Printf("Search failed: %v", err)
    return
}
// Use result
```

## Examples

- **[Basic Example](./example/main.go)**: Complete working demonstration of the client, including proper handling of asynchronous operations and data verification.
- **[Advanced Search Example](./advanced_search_example/advanced_search_example.go)**: Comprehensive demonstration of all advanced search methods including temporal queries, entity relationships, diverse results, episode context, successful tools, recent context, and entity label filtering.

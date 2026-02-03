# Graphiti Go Client

A Go client library for the Graphiti HTTP API.

## Features

- Full coverage of Graphiti HTTP API endpoints
- Optional Langfuse observation tracking for monitoring and debugging
- Configurable HTTP client and timeouts
- Type-safe request and response structures

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

### Langfuse Integration

The client supports optional Langfuse observation tracking for monitoring and debugging. You can attach an `Observation` object to any of the following operations:

- `Search` - Track search queries
- `GetMemory` - Track memory retrieval operations
- `AddMessages` - Track message ingestion
- `AddEntityNode` - Track entity node creation

**⚠️ Important:** The `Observation.ID` and `Observation.TraceID` must be valid UUIDs that correspond to actual observation and trace objects in your Langfuse instance. These IDs are used to link Graphiti operations to Langfuse traces for monitoring and debugging.

Example:

```go
import "github.com/google/uuid"

// Create an observation with UUIDs that exist in Langfuse
// These IDs should come from your Langfuse SDK after creating
// an observation/trace in your Langfuse instance
observation := &graphiti.Observation{
    ID:      "existing-observation-uuid-from-langfuse",
    TraceID: "existing-trace-uuid-from-langfuse",
    Time:    time.Now(),
}

// Use it in any supported operation
result, err := client.Search(graphiti.SearchQuery{
    Query:       "my search query",
    MaxFacts:    10,
    Observation: observation,
})
```

The observation tracking is completely optional - all operations work without it.

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

### Search with Group Filtering and Observation Tracking

```go
groupIDs := []string{"group-123", "group-456"}

// Optional: Link to existing Langfuse observation
// IDs must correspond to actual observation/trace in Langfuse
observation := &graphiti.Observation{
    ID:      "existing-observation-uuid",
    TraceID: "existing-trace-uuid",
    Time:    time.Now(),
}

result, err := client.Search(graphiti.SearchQuery{
    GroupIDs:    &groupIDs,
    Query:       "user settings",
    MaxFacts:    5,
    Observation: observation, // Optional
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

// Optional: Link to existing Langfuse observation
// IDs must correspond to actual observation/trace in Langfuse
observation := &graphiti.Observation{
    ID:      "existing-observation-uuid",
    TraceID: "existing-trace-uuid",
    Time:    time.Now(),
}

result, err := client.AddMessages(graphiti.AddMessagesRequest{
    GroupID:     "my-group-id",
    Messages:    messages,
    Observation: observation, // Optional
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

// Optional: Link to existing Langfuse observation
// IDs must correspond to actual observation/trace in Langfuse
observation := &graphiti.Observation{
    ID:      "existing-observation-uuid",
    TraceID: "existing-trace-uuid",
    Time:    time.Now(),
}

node, err := client.AddEntityNode(graphiti.AddEntityNodeRequest{
    UUID:        uuid,
    GroupID:     "my-group-id",
    Name:        "User Preferences",
    Summary:     "Contains user's preferred settings",
    Observation: observation, // Optional
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

// Optional: Link to existing Langfuse observation
// IDs must correspond to actual observation/trace in Langfuse
observation := &graphiti.Observation{
    ID:      "existing-observation-uuid",
    TraceID: "existing-trace-uuid",
    Time:    time.Now(),
}

response, err := client.GetMemory(graphiti.GetMemoryRequest{
    GroupID:     "my-group-id",
    MaxFacts:    10,
    Messages:    messages,
    Observation: observation, // Optional
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

## Types

### Observation

```go
type Observation struct {
    ID      string    // Observation UUID from Langfuse
    TraceID string    // Trace UUID from Langfuse
    Time    time.Time // Observation timestamp
}
```

**Note:** The `Observation` type is used for integrating with Langfuse for tracking and observability. The `ID` and `TraceID` must be valid UUIDs corresponding to actual observation and trace objects in your Langfuse instance. This field is optional in all requests.

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
    GroupIDs    *[]string    // Optional group IDs to filter
    Query       string       // Search query text
    MaxFacts    int          // Maximum number of facts to return (default: 10)
    Observation *Observation // Optional Langfuse observation for tracking
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

### GetMemoryRequest

```go
type GetMemoryRequest struct {
    GroupID        string       // Group ID
    MaxFacts       int          // Maximum number of facts to return
    CenterNodeUUID *string      // Optional center node UUID
    Messages       []Message    // Messages for context
    Observation    *Observation // Optional Langfuse observation for tracking
}
```

### AddMessagesRequest

```go
type AddMessagesRequest struct {
    GroupID     string       // Group ID
    Messages    []Message    // Messages to add
    Observation *Observation // Optional Langfuse observation for tracking
}
```

### AddEntityNodeRequest

```go
type AddEntityNodeRequest struct {
    UUID        string       // Entity UUID
    GroupID     string       // Group ID
    Name        string       // Entity name
    Summary     string       // Optional entity summary
    Observation *Observation // Optional Langfuse observation for tracking
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

## Example

See the [example](./example/main.go) for a complete working demonstration of the client, including proper handling of asynchronous operations and data verification.

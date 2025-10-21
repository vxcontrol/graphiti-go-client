# Graphiti Go Client

A Go client library for the Graphiti HTTP API.

## Installation

```bash
go get github.com/pentagi/graphiti-go-client
```

## Usage

### Creating a Client

```go
import (
    "time"
    graphiti "github.com/pentagi/graphiti-go-client"
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
        Content:  "Hello, how are you?",
        RoleType: graphiti.RoleTypeUser,
        Timestamp: time.Now(),
    },
    {
        Content:  "I'm doing great, thank you!",
        RoleType: graphiti.RoleTypeAssistant,
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
        Content:  "What were my settings?",
        RoleType: graphiti.RoleTypeUser,
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

## Types

### Message

```go
type Message struct {
    Content           string    // The message content
    UUID              *string   // Optional UUID
    Name              string    // Optional name for episodic node
    RoleType          RoleType  // user, assistant, or system
    Role              *string   // Optional custom role (user name, bot name, etc.)
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

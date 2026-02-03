package graphiti

import "time"

// Observation represents Langfuse observation object to link
type Observation struct {
	ID      string    `json:"id"`
	TraceID string    `json:"trace_id"`
	Time    time.Time `json:"time"`
}

// Message represents a message in the system
type Message struct {
	Content           string    `json:"content"`
	UUID              *string   `json:"uuid,omitempty"`
	Name              string    `json:"name,omitempty"`
	Author            string    `json:"author"`
	Timestamp         time.Time `json:"timestamp"`
	SourceDescription string    `json:"source_description,omitempty"`
}

// Result represents a generic result response
type Result struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// SearchQuery represents a search query request
type SearchQuery struct {
	GroupIDs    *[]string    `json:"group_ids,omitempty"`
	Query       string       `json:"query"`
	MaxFacts    int          `json:"max_facts,omitempty"`
	Observation *Observation `json:"observation,omitempty"`
}

// FactResult represents a fact result from the graph
type FactResult struct {
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	Fact      string     `json:"fact"`
	ValidAt   *time.Time `json:"valid_at,omitempty"`
	InvalidAt *time.Time `json:"invalid_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`
}

// SearchResults represents the results of a search query
type SearchResults struct {
	Facts []FactResult `json:"facts"`
}

// GetMemoryRequest represents a request to get memory
type GetMemoryRequest struct {
	GroupID        string       `json:"group_id"`
	MaxFacts       int          `json:"max_facts,omitempty"`
	CenterNodeUUID *string      `json:"center_node_uuid"`
	Messages       []Message    `json:"messages"`
	Observation    *Observation `json:"observation,omitempty"`
}

// GetMemoryResponse represents the response from getting memory
type GetMemoryResponse struct {
	Facts []FactResult `json:"facts"`
}

// AddMessagesRequest represents a request to add messages
type AddMessagesRequest struct {
	GroupID     string       `json:"group_id"`
	Messages    []Message    `json:"messages"`
	Observation *Observation `json:"observation,omitempty"`
}

// AddEntityNodeRequest represents a request to add an entity node
type AddEntityNodeRequest struct {
	UUID        string       `json:"uuid"`
	GroupID     string       `json:"group_id"`
	Name        string       `json:"name"`
	Summary     string       `json:"summary,omitempty"`
	Observation *Observation `json:"observation,omitempty"`
}

// EntityNode represents an entity node in the graph
type EntityNode struct {
	UUID      string                 `json:"uuid"`
	GroupID   string                 `json:"group_id"`
	Name      string                 `json:"name"`
	Summary   string                 `json:"summary,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	Labels    []string               `json:"labels,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// Episode represents an episode in the graph
type Episode struct {
	UUID              string                 `json:"uuid"`
	GroupID           string                 `json:"group_id"`
	Name              string                 `json:"name"`
	Content           string                 `json:"content"`
	Source            string                 `json:"source"`
	SourceDescription string                 `json:"source_description,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	ValidAt           time.Time              `json:"valid_at"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

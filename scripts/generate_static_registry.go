package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

// Simple registry format (GitHub instructions)
// {
//   "servers": [ { "id": "github", "name": "GitHub MCP Server", ... } ],
//   "total_count": 2,
//   "updated_at": "2025-09-09T12:00:00Z"
// }

type SimpleServer struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    ManifestURL string    `json:"manifest_url"`
    Categories  []string  `json:"categories,omitempty"`
    Version     string    `json:"version,omitempty"`
    ReleaseDate time.Time `json:"release_date,omitempty"`
    Latest      bool      `json:"latest,omitempty"`
}

type Curated struct {
    Servers    []SimpleServer `json:"servers"`
    TotalCount int            `json:"total_count"`
    UpdatedAt  time.Time      `json:"updated_at"`
}

func main() {
    sourcePath := "data/curated_servers.json"
    outputPath := "docs/v0/servers/index.json"

    raw, err := os.ReadFile(sourcePath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to read curated_servers.json: %v\n", err)
        os.Exit(1)
    }

    var curated Curated
    if err := json.Unmarshal(raw, &curated); err != nil {
        fmt.Fprintf(os.Stderr, "failed to parse curated_servers.json (expecting simple registry format): %v\n", err)
        os.Exit(1)
    }

    // Populate counts & timestamps if missing / zero.
    if curated.TotalCount == 0 {
        curated.TotalCount = len(curated.Servers)
    }
    // Always refresh updated_at to reflect regeneration time.
    curated.UpdatedAt = time.Now().UTC().Truncate(time.Second)

    if err := os.MkdirAll("docs/v0/servers", 0o755); err != nil {
        fmt.Fprintf(os.Stderr, "failed to ensure output dir: %v\n", err)
        os.Exit(1)
    }

    outBytes, err := json.MarshalIndent(curated, "", "  ")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to marshal index.json: %v\n", err)
        os.Exit(1)
    }
    if err := os.WriteFile(outputPath, outBytes, 0o640); err != nil { // 0640 for read by others if served statically
        fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", outputPath, err)
        os.Exit(1)
    }
    fmt.Fprintln(os.Stdout, "Generated simple registry docs/v0/servers/index.json")
}

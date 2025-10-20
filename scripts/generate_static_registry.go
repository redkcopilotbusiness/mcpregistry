package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

// Minimal structures mirroring needed fields
type ServerJSON struct {
    Schema      string            `json:"$schema"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Title       string            `json:"title,omitempty"`
    Repository  json.RawMessage   `json:"repository,omitempty"`
    Version     string            `json:"version"`
    WebsiteURL  string            `json:"websiteUrl,omitempty"`
    Packages    []json.RawMessage `json:"packages,omitempty"`
    Icons       []json.RawMessage `json:"icons,omitempty"`
    Meta        json.RawMessage   `json:"_meta,omitempty"`
    Remotes     []json.RawMessage `json:"remotes,omitempty"`
}

type OfficialMeta struct {
    Status          string    `json:"status"`
    PublishedAt     time.Time `json:"publishedAt"`
    UpdatedAt       time.Time `json:"updatedAt"`
    IsLatest        bool      `json:"isLatest"`
    VersionSequence int       `json:"versionSequence"`
}

type ResponseMeta struct {
    Official OfficialMeta `json:"io.modelcontextprotocol.registry/official"`
}

type ServerResponse struct {
    Server ServerJSON   `json:"server"`
    Meta   ResponseMeta `json:"_meta"`
}

type ServerListResponse struct {
    Servers  []ServerResponse `json:"servers"`
    Metadata struct{ Count int `json:"count"` } `json:"metadata"`
}

func main() {
    sourcePath := "data/curated_servers.json"
    existingIndexPath := "docs/v0/servers/index.json"

    sourceData, err := os.ReadFile(sourcePath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to read curated source: %v\n", err)
        os.Exit(1)
    }

    var servers []ServerJSON
    if err := json.Unmarshal(sourceData, &servers); err != nil {
        fmt.Fprintf(os.Stderr, "failed to parse curated source: %v\n", err)
        os.Exit(1)
    }

  // Load existing index if present to reuse versionSequence
    existingSequences := map[string]struct{ Version string; Seq int }{}
    if data, err := os.ReadFile(existingIndexPath); err == nil {
        var existing ServerListResponse
        if json.Unmarshal(data, &existing) == nil {
            for _, sr := range existing.Servers {
                existingSequences[sr.Server.Name] = struct{ Version string; Seq int }{Version: sr.Server.Version, Seq: sr.Meta.Official.VersionSequence}
            }
        }
    }

    now := time.Now().UTC()
    out := ServerListResponse{}

    for _, s := range servers {
        prev, ok := existingSequences[s.Name]
        seq := 1
        if ok {
            if prev.Version == s.Version && prev.Seq > 0 {
                seq = prev.Seq // unchanged version, keep sequence
            } else {
                seq = prev.Seq + 1 // version changed or sequence 0 -> increment
            }
        }
        out.Servers = append(out.Servers, ServerResponse{
            Server: s,
            Meta: ResponseMeta{Official: OfficialMeta{
                Status:          "active",
                PublishedAt:     now, // Could preserve previous publishedAt if needed
                UpdatedAt:       now,
                IsLatest:        true,
                VersionSequence: seq,
            }},
        })
    }
    out.Metadata.Count = len(out.Servers)

  // Ensure output directory exists
    if err := os.MkdirAll("docs/v0/servers", 0o755); err != nil {
        fmt.Fprintf(os.Stderr, "failed to create output dir: %v\n", err)
        os.Exit(1)
    }

    jsonBytes, err := json.MarshalIndent(out, "", "  ")
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to marshal output: %v\n", err)
        os.Exit(1)
    }

    if err := os.WriteFile(existingIndexPath, jsonBytes, 0o600); err != nil { // 0600 per gosec recommendation
        fmt.Fprintf(os.Stderr, "failed to write index.json: %v\n", err)
        os.Exit(1)
    }

    fmt.Fprintf(os.Stdout, "Generated docs/v0/servers/index.json with versionSequence numbers.\n")
}

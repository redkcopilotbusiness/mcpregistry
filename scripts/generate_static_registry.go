package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// MCP Registry format following server-list-template.json
// {
//   "servers": [
//     {
//       "_meta": { "io.modelcontextprotocol.registry/official": {...} },
//       "server": { "$schema": "...", "name": "...", ... }
//     }
//   ],
//   "metadata": { "count": 5, "nextCursor": null }
// }

type Transport struct {
	Type string `json:"type"`
	URL  string `json:"url,omitempty"`
}

type Package struct {
	RegistryType string    `json:"registryType"`
	Identifier   string    `json:"identifier"`
	Version      string    `json:"version"`
	RuntimeHint  string    `json:"runtimeHint,omitempty"`
	Transport    Transport `json:"transport"`
}

type Server struct {
	Schema      string      `json:"$schema"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Title       string      `json:"title,omitempty"`
	Version     string      `json:"version"`
	Packages    []Package   `json:"packages,omitempty"`
	Remotes     []Transport `json:"remotes,omitempty"`
}

type OfficialMeta struct {
	Status      string `json:"status"`
	PublishedAt string `json:"publishedAt"`
	IsLatest    bool   `json:"isLatest"`
}

type Meta struct {
	Official OfficialMeta `json:"io.modelcontextprotocol.registry/official"`
}

type ServerEntry struct {
	Meta   Meta   `json:"_meta"`
	Server Server `json:"server"`
}

type Metadata struct {
	Count      int  `json:"count"`
	NextCursor *int `json:"nextCursor"`
}

type RegistryList struct {
	Servers  []ServerEntry `json:"servers"`
	Metadata Metadata      `json:"metadata"`
}

func main() {
	sourcePath := "data/curated_servers.json"
	outputPath := "docs/v0/servers/index.json"

	raw, err := os.ReadFile(sourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read curated_servers.json: %v\n", err)
		os.Exit(1)
	}

	var registryList RegistryList
	if err := json.Unmarshal(raw, &registryList); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse curated_servers.json (expecting MCP registry format): %v\n", err)
		os.Exit(1)
	}

	// Ensure metadata count matches actual server count
	registryList.Metadata.Count = len(registryList.Servers)

	if err := os.MkdirAll("docs/v0/servers", 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to ensure output dir: %v\n", err)
		os.Exit(1)
	}

	outBytes, err := json.MarshalIndent(registryList, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal index.json: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outputPath, outBytes, 0o600); err != nil { // 0600 for secure file permissions
		fmt.Fprintf(os.Stderr, "failed to write %s: %v\n", outputPath, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, "Generated MCP registry docs/v0/servers/index.json")
}

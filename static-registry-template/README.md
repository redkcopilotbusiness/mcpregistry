# Static MCP Registry Template

Minimal template for a read-only curated MCP registry served via GitHub Pages.

## Structure

```
static-registry-template/
  README.md
  docs/
    v0/
      servers/
        index.json   # Static listing of MCP servers
```

## index.json Format

The file `docs/v0/servers/index.json` must follow the registry list response shape:

```jsonc
{
  "servers": [
    {
      "server": {
        "$schema": "https://static.modelcontextprotocol.io/schemas/2025-10-17/server.schema.json",
        "name": "com.example/hello-mcp",
        "description": "Example MCP server from static curated registry template.",
        "version": "1.0.0",
        "packages": [
          {
            "registryType": "npm",
            "identifier": "@example/hello-mcp",
            "version": "1.0.0",
            "environmentVariables": [],
            "transport": { "type": "stdio" }
          }
        ]
      },
      "_meta": {
        "io.modelcontextprotocol.registry/official": {
          "status": "active",
          "publishedAt": "2025-10-20T00:00:00Z",
          "updatedAt": "2025-10-20T00:00:00Z",
          "isLatest": true,
          "versionSequence": 1
        }
      }
    }
  ],
  "metadata": { "count": 1 }
}
```

Add more servers by appending additional objects to the `servers` array. Increase `metadata.count` accordingly.

## Reverse-DNS Naming

Use the pattern `reverse-dns-namespace/name`, e.g.:

- `io.github.usuario/mi-servidor`
- `com.miempresa/ventas`

## Publishing via GitHub Pages

1. Create a new repository and copy the contents of `static-registry-template/` to its root.
2. Enable GitHub Pages: Settings > Pages > Deploy from a branch > `main` / folder `docs`.
3. Access URL: `https://<org>.github.io/<repo>/v0/servers/index.json`.
4. Provide that URL to clients that support custom MCP registries.

## Updating Versions

- Change the `version` field inside the `server` and (if you use it) also inside the package `version`.
- Increment `versionSequence` if the version changed.
- Update `updatedAt` timestamp (ISO 8601 UTC).

## Optional Enhancements

- Add `_meta.io.modelcontextprotocol.registry/publisher-provided` for internal tags.
- Add `icons` for better UI.
- Split into per-server version files (e.g. `docs/v0/servers/com.example%2Fhello-mcp/versions/1.0.0.json`).

## Validation Tips

Basic checks:

- `$schema` URL matches current schema.
- `name` has exactly one `/`.
- `version` is not `latest`.
- `metadata.count` equals length of `servers` array.

You now have the minimal static curated MCP registry.

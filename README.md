# screenjson-export

A free, open-source CLI that converts writer-format screenplays into
[ScreenJSON](https://screenjson.com).

- **Final Draft** — `.fdx`
- **Fountain** — `.fountain`, `.spmd`
- **FadeIn** — `.fadein`

It does one thing: read a writer format, emit a valid ScreenJSON
document on stdout or to a file. It exists so that anyone can convert a
screenplay into the format without a license, without an account, and
without a network call — and so that anyone building their own
converter has a correct, freely-licensed target to compare against.

For the full toolchain — PDF import, exports, validation, AES-256
encryption, REST API server, and storage drivers for MongoDB, DynamoDB,
Elasticsearch, Cassandra, Redis, S3, Azure, and MinIO — see
[**screenjson-cli**](https://screenjson.com/tools/screenjson-cli/).

## Install

### Download a release binary

Builds for macOS, Linux, and Windows are attached to each GitHub release:

<https://github.com/screenjson/screenjson-export/releases>

### Build from source

```bash
git clone https://github.com/screenjson/screenjson-export.git
cd screenjson-export
go build -o screenjson-export ./cmd/screenjson-export
sudo mv screenjson-export /usr/local/bin/
```

Requires Go 1.22+.

### Docker

```bash
docker pull screenjson/export:latest
docker run --rm -v "$PWD:/data" screenjson/export -i /data/screenplay.fdx
```

## Use

```bash
# Convert Final Draft to ScreenJSON (stdout)
screenjson-export -i screenplay.fdx

# Write to a file
screenjson-export -i screenplay.fdx -o screenplay.json

# Fountain / FadeIn — format is auto-detected from the extension
screenjson-export -i screenplay.fountain -o screenplay.json
screenjson-export -i screenplay.fadein   -o screenplay.json

# Force input format on a renamed file
screenjson-export -i script.txt --format fountain -o screenplay.json

# Pipe into jq
screenjson-export -i screenplay.fdx | jq '.title.en'

# Primary language other than English
screenjson-export -i script.fdx --lang fr -o script.json
```

## Options

| Flag | Description |
|------|-------------|
| `-i`, `--input <path>` | Input file. **Required.** |
| `-o`, `--output <path>` | Output file. Defaults to stdout. |
| `-f`, `--format <name>` | Force `fdx`, `fountain`, or `fadein`. Auto-detected otherwise. |
| `--lang <tag>` | BCP 47 primary language tag. Default `en`. |
| `--pretty` | Pretty-print JSON. Default `true`. |
| `-v`, `--version` | Print version. |
| `-h`, `--help` | Show help. |

## What this tool is — and isn't

`screenjson-export` is a minimal subset of the full ScreenJSON CLI,
released free and open-source.

| Capability | `screenjson-export` | `screenjson-cli` |
|---|---|---|
| FDX → ScreenJSON | ✅ | ✅ |
| Fountain → ScreenJSON | ✅ | ✅ |
| FadeIn → ScreenJSON | ✅ | ✅ |
| PDF → ScreenJSON | ❌ | ✅ |
| ScreenJSON → FDX / Fountain / FadeIn / PDF | ❌ | ✅ |
| Schema validation | ❌ | ✅ |
| AES-256 content encryption | ❌ | ✅ |
| REST API server mode | ❌ | ✅ |
| Storage drivers (MongoDB, DynamoDB, ES, S3, …) | ❌ | ✅ |
| MCP mode for AI agents | ❌ | ✅ |

## License

MIT — see [LICENSE](./LICENSE).

The ScreenJSON schema is also open. The full specification is at
<https://screenjson.com/specification/>.

## Related

- [ScreenJSON home](https://screenjson.com)
- [Full CLI](https://screenjson.com/tools/screenjson-cli/) (commercial)
- [Greenlight](https://screenjson.com/tools/greenlight/) — batch microservice
- [screenjson-ui](https://screenjson.com/tools/screenjson-ui/) — browser viewer

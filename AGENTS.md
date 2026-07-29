# AGENTS.md

Go single-binary bookmark homepage (Gin). No database. Runtime data = YAML next to the process cwd.

## Commands

```bash
go run build/build.go          # REQUIRED before go build/run if embed/ or icon vendor changed
go run .                       # default :5005
go build -o flare .
go test ./...
go test ./config/data -run TestFavoriteBookmarks
gofmt -s -w .
```

Docker/release also run `go run build/build.go` first (see `Dockerfile`, `.goreleaser.yaml`).

## Docker

Image: `cyzzxc/flare` (multi-arch: `linux/amd64`, `linux/arm64`). Runtime YAML lives in container cwd `/app`.

```bash
# compose (recommended)
docker compose up -d
# → http://localhost:5005  data in ./data/{apps,bookmarks,config}.yml

# one-shot
docker run -d --name flare -p 5005:5005 -v "$PWD/data:/app" cyzzxc/flare:latest
```

Build/push (needs `docker buildx` docker-container driver for multi-arch):

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t cyzzxc/flare:latest --push .
```

Health: `GET /ping` (JSON `{"message":"pong"}`). Not `/health` (Dockerfile HEALTHCHECK path is stale).

## Source vs generated (easy to get wrong)

| Edit here | Build writes here (go:embed / runtime) |
|-----------|----------------------------------------|
| `embed/templates/` | `internal/resources/templates/html/` (minified) |
| `embed/assets/css/` | inlined into `config/define/style.go` |
| `embed/assets/vendor/remixicon/paths.json` | `internal/resources/mdi/icons.go` + `mdi-cheat-sheets/` |
| `embed/assets/vendor/{guide,editor}-assets/` | `internal/pages/{guide,editor}/*-assets/` |
| `embed/assets/favicon.ico` | `internal/resources/assets/favicon.ico` |

- **Do not hand-edit** generated trees under `internal/resources/`, `internal/pages/*/guide-assets|editor-assets`, or `icons.go` — change source + re-run build.
- Package path is still `internal/resources/mdi` and import alias `FlareMDI`; icon data is **Remix Icon**, not Material Design Icons. Icon names are kebab-case (`home-line`, `mail-fill`). `GetIconByName` keys are lowercased.
- `/icons` is a build-generated search/copy page from the same path map.
- `.gitignore` lists some old paths (`internal/mdi/...`); live paths are under `internal/resources/`.

## Layout

- `main.go` → `cmd.Parse()` → `internal/server.StartDaemon`
- `cmd/` CLI + env (`.env` in workdir); flags override env
- `config/model` structs · `config/define` constants/routes · `config/data` YAML load/save
- `internal/pages/{home,editor,guide}` HTML pages · `internal/settings/*` settings POST handlers
- Runtime files (cwd, gitignored): `apps.yml`, `bookmarks.yml`, `config.yml`

## Icons & bookmarks

- Bookmark/app `icon` field: Remix name, or `http(s)://` URL.
- Empty/unknown name → no icon. **No MDI→Remix mapping** — old names (`homeCircle`) render blank.
- `IconMode` FILLING uses Yandex favicon with Remix path as fallback (`internal/fn/favicon.go`).
- Search helpers: `FlareMDI.SearchIcons` / `IconExists` in `internal/resources/mdi`.

## Conventions (project-specific)

- Prefer plain functions over pointer receivers; care about allocs (`CONVENTIONS-OF-CODE.md`).
- `gofmt -s`; std `go test` only (testify ok).
- Default port **5005** (`config/define/cmd.go`). Health JSON: `GET /ping`.
- Do not add deps for what a few lines + stdlib cover.

## REST API

On by default (`--enable_api` / `FLARE_API`). Optional auth: `--api_key` / `FLARE_API_KEY` (Bearer or `X-API-Key`); empty key = open.

| Method | Path | Body |
|--------|------|------|
| GET | `/api/v1` | index / notes for agents |
| GET/PUT | `/api/v1/apps` | whole `apps.yml` (`Bookmarks` JSON) |
| GET/PUT | `/api/v1/bookmarks` | whole `bookmarks.yml` |
| GET/PUT | `/api/v1/settings` | whole `config.yml` (`Application`) |
| GET | `/api/v1/icons?q=&limit=20` | Remix name search |

- Code: `internal/api/` · registered from `internal/server`.
- PUT replaces the full document. Fields: `name`, `link`, `icon`, `desc`, `category`, `private`.
- Unknown non-URL icons rejected on PUT.
- Agent guide (always, no auth): `GET /llm.txt` ← `internal/api/llm.txt`. Linked from `/help`.

## Gotchas

- Fresh clone / clean tree: if embed targets missing, run `go run build/build.go` before compile.
- Tests under `config/data` may touch workdir YAML paths; prefer package-scoped runs when iterating.
- Login/auth defaults: without `FLARE_USER`/`FLARE_PASS`, username defaults and password is auto-generated (logged at start). `--disable-login` skips auth for local UI work.
- API is independent of session login; empty `api_key` means open if API is enabled.
- Do not commit local binary `flare` or workdir `*.yml` / `.env`.

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
- Empty/unknown name → no icon. Old MDI names (e.g. `homeCircle`) no longer resolve.
- `IconMode` FILLING uses Yandex favicon with Remix path as fallback (`internal/fn/favicon.go`).

## Conventions (project-specific)

- Prefer plain functions over pointer receivers; care about allocs (`CONVENTIONS-OF-CODE.md`).
- `gofmt -s`; std `go test` only (testify ok).
- Default port **5005** (`config/define/cmd.go`). Health: `/health`.
- Do not add deps for what a few lines + stdlib cover.

## Gotchas

- Fresh clone / clean tree: if embed targets missing, run `go run build/build.go` before compile.
- Tests under `config/data` may touch workdir YAML paths; prefer package-scoped runs when iterating.
- Login/auth defaults: without `FLARE_USER`/`FLARE_PASS`, username defaults and password is auto-generated (logged at start). `--disable-login` skips auth for local UI work.

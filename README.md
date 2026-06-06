# winrice

Desktop Rice Manager for Windows 11 (GlazeWM, Zebar, WezTerm, Fastfetch) built using Svelte 5, Go, and Tauri v2.

Follows the G2G (geek-to-geek) minimalist design philosophy with a Zenbones Sand theme.

## Architecture

```
[ Tauri Desktop Shell (Rust) ]
   ├── Spawns & monitors: [ Go Backend (Port 54321) ] ──> [ Active C:\ Configs ]
   └── Displays: [ Svelte 5 Frontend (Port 5173) ]
```

- **Svelte 5**: Desktop web interface for previewing, adopting, editing, and linking configurations.
- **Go 1.22+**: Handles file operations (read, write, copy), junction/symlink creation on Windows, and launching PowerShell scripts (`sync_theme.ps1`).
- **Tauri v2**: Desktop wrapper that orchestrates the UI, spawns the Go server on startup, and terminates it cleanly on window close.

## Supported Configurations

- **WezTerm**: `C:\Users\Timofey\.wezterm.lua` -> `config\wezterm\.wezterm.lua`
- **GlazeWM**: `C:\Users\Timofey\.glzr\glazewm\config.yaml` -> `config\glazewm\config.yaml`
- **Zebar**: `C:\Users\Timofey\.glzr\zebar\` -> `config\zebar`
- **Fastfetch**: `C:\Users\Timofey\AppData\Roaming\fastfetch\config.jsonc` -> `config\fastfetch\config.jsonc`
- **Theme Syncer**: `C:\Users\Timofey\Theme\sync_theme.ps1` -> `config\Theme\sync_theme.ps1`

## Quick Start

### Development

To start the Svelte frontend and Go backend in development mode with HMR:

```powershell
# Set Rust variables and run tauri dev (spawns both Vite and Go backend)
$env:RUSTUP_HOME="E:\Dev\rust\.rustup"
$env:CARGO_HOME="E:\Dev\rust\.cargo"
$env:PATH="E:\Dev\rust\.cargo\bin;" + $env:PATH
npm run tauri dev
```

### Production Build

To compile a production binary (bundled `.exe`):

```powershell
# 1. Compile Go backend first
cd src-go
go build -o winrice-backend.exe main.go
cd ..

# 2. Copy winrice-backend.exe to release target directory
# (Tauri Rust app looks for it in the same directory as the executable)
mkdir src-tauri\target\release -Force
Copy-Item src-go\winrice-backend.exe src-tauri\target\release\winrice-backend.exe

# 3. Build Tauri application
$env:RUSTUP_HOME="E:\Dev\rust\.rustup"
$env:CARGO_HOME="E:\Dev\rust\.cargo"
$env:PATH="E:\Dev\rust\.cargo\bin;" + $env:PATH
npm run tauri build
```

## Structure

```
E:\Dev\projects\winrice\
├── config\             — adopted configuration files (created on first adopt)
├── src-go\             — Go backend API server source
│   ├── main.go
│   └── go.mod
├── src-tauri\          — Tauri v2 configuration & Rust entry point
│   ├── Cargo.toml
│   └── src\
│       ├── main.rs
│       └── lib.rs
├── src\                — Svelte 5 frontend source
│   ├── App.svelte      — dashboard layout & interactions
│   ├── app.css         — G2G reset and typography
│   └── theme.css       — Zenbones Sand theme palette
├── package.json
└── vite.config.ts
```

## License

Proprietary. Developed for Timofey.

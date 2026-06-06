# winrice

![WinRice Preview](preview.png)

Desktop Rice Manager for Windows 11 (GlazeWM, Zebar, WezTerm, Fastfetch, Pinterest Wallpaper) built using Svelte 5 and Tauri v2 (Rust).

Follows the G2G (geek-to-geek) minimalist design philosophy with a Zenbones Sand theme.

## Architecture

```
[ Tauri Desktop Shell (Rust) ] ──> [ Active User Configs ]
         └── Displays: [ Svelte 5 Frontend (Tauri IPC) ]
```

- **Svelte 5**: Desktop web interface for previewing, adopting, editing, and linking configurations.
- **Tauri v2 (Rust)**: Orchestrates the UI window, registers native IPC command handlers for file operations, and manages silent PowerShell script/process execution.

## Supported Configurations

- **WezTerm**: `<UserProfile>\.wezterm.lua` -> `config\wezterm\.wezterm.lua`
- **GlazeWM**: `<UserProfile>\.glzr\glazewm\config.yaml` -> `config\glazewm\config.yaml`
- **Zebar**: `<UserProfile>\.glzr\zebar\` -> `config\zebar`
- **Fastfetch**: `<AppData>\Roaming\fastfetch\config.jsonc` -> `config\fastfetch\config.jsonc`
- **Theme Syncer**: `<UserProfile>\Theme\sync_theme.ps1` -> `config\Theme\sync_theme.ps1`
- **Zed Editor**: `<AppData>\Roaming\Zed\settings.json` -> `config\Zed\settings.json`
- **VS Code**: `<AppData>\Roaming\Code\User\settings.json` -> `config\VSCode\settings.json`
- **Komorebi**: `<UserProfile>\komorebi.json` -> `config\komorebi\komorebi.json`
- **Komorebi Bar**: `<UserProfile>\komorebi.bar.json` -> `config\komorebi\komorebi.bar.json`
- **whkd Shortcuts**: `<UserProfile>\.config\whkdrc` -> `config\whkd\whkdrc`
- **Git Config**: `<UserProfile>\.gitconfig` -> `config\git\.gitconfig`
- **Pinterest Wallpaper**: dynamically located `config.json` -> `config\pinterest-collage\config.json`

## Quick Start

### Development

To start the Svelte frontend and Tauri app in development mode with HMR:

```powershell
# Set Rust variables and run tauri dev
$env:RUSTUP_HOME="E:\Dev\rust\.rustup"
$env:CARGO_HOME="E:\Dev\rust\.cargo"
$env:PATH="E:\Dev\rust\mingw64\bin;E:\Dev\rust\.cargo\bin;" + $env:PATH
npx tauri dev
```

### Production Build

To compile a production installer and standalone binary:

```powershell
# Build Tauri application
$env:RUSTUP_HOME="E:\Dev\rust\.rustup"
$env:CARGO_HOME="E:\Dev\rust\.cargo"
$env:PATH="E:\Dev\rust\mingw64\bin;E:\Dev\rust\.cargo\bin;" + $env:PATH
npx tauri build
```

## Structure

```
E:\Dev\projects\winrice\
├── config\             — adopted configuration files (created on first adopt)
├── src-tauri\          — Tauri v2 configuration & Rust entry point (Current Backend)
│   ├── Cargo.toml
│   └── src\
│       ├── main.rs
│       └── lib.rs
├── src\                — Svelte 5 frontend source (Current Frontend)
│   ├── App.svelte      — dashboard layout & interactions
│   ├── app.css         — G2G reset and typography
│   └── theme.css       — Zenbones Sand theme palette
├── src-go\             — [LEGACY] Go backend source (not used in current Tauri v2 build)
├── package.json
└── vite.config.ts
```

## License

Proprietary. Developed for Timofey.

use std::process::Command;
use std::io::Write;

#[cfg(target_os = "windows")]
use std::os::windows::process::CommandExt;

#[cfg(target_os = "windows")]
const CREATE_NO_WINDOW: u32 = 0x08000000;

#[cfg(target_os = "windows")]
fn configure_command(cmd: &mut Command) {
    cmd.creation_flags(CREATE_NO_WINDOW);
}

#[cfg(not(target_os = "windows"))]
fn configure_command(_cmd: &mut Command) {}

#[derive(serde::Serialize, serde::Deserialize, Clone)]
#[serde(rename_all = "camelCase")]
pub struct ConfigInfo {
    pub key: String,
    pub name: String,
    pub target_path: String,
    pub repo_path: String,
    pub exists: bool,
    pub in_repo: bool,
    pub is_symlink: bool,
    pub symlink_to: String,
    pub is_dir: bool,
    pub skip_link: bool,
}

#[derive(serde::Serialize, serde::Deserialize, Clone)]
pub struct ThemeInfo {
    pub font_family: String,
    pub font_size: f64,
    pub opacity: f64,
    pub border_focused: String,
    pub border_unfocused: String,
    pub bg_color: String,
    pub fg_color: String,
    pub accent_color: String,
    pub lavender: String,
    pub lilac: String,
    pub lavender_grey: String,
    pub pine_blue: String,
    pub jungle_teal: String,
}

#[derive(serde::Serialize, serde::Deserialize, Clone)]
pub struct CustomPreset {
    pub name: String,
    pub bg_color: String,
    pub fg_color: String,
    pub accent_color: String,
    pub border_focused: String,
    pub border_unfocused: String,
    pub lavender: String,
    pub lilac: String,
    pub lavender_grey: String,
    pub pine_blue: String,
    pub jungle_teal: String,
}

fn get_home_dir() -> std::path::PathBuf {
    dirs::home_dir().unwrap_or_else(|| {
        std::env::var("USERPROFILE")
            .map(std::path::PathBuf::from)
            .unwrap_or_else(|_| std::path::PathBuf::from("C:\\Users\\Default"))
    })
}

fn get_appdata_dir() -> std::path::PathBuf {
    dirs::config_dir().unwrap_or_else(|| {
        std::env::var("APPDATA")
            .map(std::path::PathBuf::from)
            .unwrap_or_else(|_| get_home_dir().join("AppData").join("Roaming"))
    })
}

fn find_pinterest_collage_config() -> String {
    let home = get_home_dir();
    let appdata = get_appdata_dir();
    
    let candidates = [
        appdata.join("pinterest-collage").join("config.json"),
        std::path::PathBuf::from("E:\\Dev\\projects\\pinterest-collage\\config.json"),
        home.join("Dev").join("projects").join("pinterest-collage").join("config.json"),
        home.join("projects").join("pinterest-collage").join("config.json"),
    ];

    for path in &candidates {
        if path.exists() {
            return path.to_string_lossy().to_string();
        }
    }

    // Default fallback
    candidates[0].to_string_lossy().to_string()
}

fn get_configs_definition() -> Vec<ConfigInfo> {
    let home = get_home_dir();
    let appdata = get_appdata_dir();
    
    vec![
        ConfigInfo {
            key: "wezterm".to_string(),
            name: "WezTerm".to_string(),
            target_path: home.join(".wezterm.lua").to_string_lossy().to_string(),
            repo_path: "config\\wezterm\\.wezterm.lua".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: false,
        },
        ConfigInfo {
            key: "glazewm".to_string(),
            name: "GlazeWM".to_string(),
            target_path: home.join(".glzr").join("glazewm").join("config.yaml").to_string_lossy().to_string(),
            repo_path: "config\\glazewm\\config.yaml".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: false,
        },
        ConfigInfo {
            key: "zebar".to_string(),
            name: "Zebar Directory".to_string(),
            target_path: home.join(".glzr").join("zebar").to_string_lossy().to_string(),
            repo_path: "config\\zebar".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: true,
            skip_link: false,
        },
        ConfigInfo {
            key: "zebar_settings".to_string(),
            name: "Zebar Settings".to_string(),
            target_path: home.join(".glzr").join("zebar").join("settings.json").to_string_lossy().to_string(),
            repo_path: "config\\zebar\\settings.json".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: true,
        },
        ConfigInfo {
            key: "zebar_zpack".to_string(),
            name: "Zebar Layout (zpack)".to_string(),
            target_path: home.join(".glzr").join("zebar").join("custom-mei").join("zpack.json").to_string_lossy().to_string(),
            repo_path: "config\\zebar\\custom-mei\\zpack.json".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: true,
        },
        ConfigInfo {
            key: "fastfetch".to_string(),
            name: "Fastfetch".to_string(),
            target_path: appdata.join("fastfetch").join("config.jsonc").to_string_lossy().to_string(),
            repo_path: "config\\fastfetch\\config.jsonc".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: false,
        },
        ConfigInfo {
            key: "sync_script".to_string(),
            name: "Theme Syncer".to_string(),
            target_path: home.join("Theme").join("sync_theme.ps1").to_string_lossy().to_string(),
            repo_path: "config\\Theme\\sync_theme.ps1".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: false,
        },
        ConfigInfo {
            key: "zed".to_string(),
            name: "Zed Editor".to_string(),
            target_path: appdata.join("Zed").join("settings.json").to_string_lossy().to_string(),
            repo_path: "config\\Zed\\settings.json".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: false,
        },
        ConfigInfo {
            key: "vscode".to_string(),
            name: "VS Code".to_string(),
            target_path: appdata.join("Code").join("User").join("settings.json").to_string_lossy().to_string(),
            repo_path: "config\\VSCode\\settings.json".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: false,
        },
        ConfigInfo {
            key: "komorebi".to_string(),
            name: "Komorebi".to_string(),
            target_path: home.join("komorebi.json").to_string_lossy().to_string(),
            repo_path: "config\\komorebi\\komorebi.json".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: false,
        },
        ConfigInfo {
            key: "komorebi_bar".to_string(),
            name: "Komorebi Bar".to_string(),
            target_path: home.join("komorebi.bar.json").to_string_lossy().to_string(),
            repo_path: "config\\komorebi\\komorebi.bar.json".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: false,
        },
        ConfigInfo {
            key: "whkd".to_string(),
            name: "whkd Shortcuts".to_string(),
            target_path: home.join(".config").join("whkdrc").to_string_lossy().to_string(),
            repo_path: "config\\whkd\\whkdrc".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: false,
        },
        ConfigInfo {
            key: "gitconfig".to_string(),
            name: "Git Config".to_string(),
            target_path: home.join(".gitconfig").to_string_lossy().to_string(),
            repo_path: "config\\git\\.gitconfig".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: false,
        },
        ConfigInfo {
            key: "pinterest_collage".to_string(),
            name: "Pinterest Wallpaper".to_string(),
            target_path: find_pinterest_collage_config(),
            repo_path: "config\\pinterest-collage\\config.json".to_string(),
            exists: false,
            in_repo: false,
            is_symlink: false,
            symlink_to: "".to_string(),
            is_dir: false,
            skip_link: false,
        },
    ]
}

fn get_project_root() -> std::path::PathBuf {
    if let Ok(mut dir) = std::env::current_dir() {
        if dir.ends_with("src-tauri") {
            dir.pop();
        }
        if dir.join("src-tauri").exists() || dir.join("config").exists() {
            return dir;
        }
    }

    if let Ok(exe_path) = std::env::current_exe() {
        if let Some(exe_dir) = exe_path.parent() {
            let mut dir = exe_dir.to_path_buf();
            while dir.pop() {
                if dir.join("src-tauri").exists() {
                    return dir;
                }
            }
            return exe_dir.to_path_buf();
        }
    }

    let dev_path = std::path::PathBuf::from("E:\\Dev\\projects\\winrice");
    if dev_path.exists() {
        dev_path
    } else {
        std::path::PathBuf::from(".")
    }
}

fn get_info(mut config: ConfigInfo, project_root: &std::path::Path) -> ConfigInfo {
    config.repo_path = project_root.join(&config.repo_path).to_string_lossy().to_string();

    let repo_path = std::path::Path::new(&config.repo_path);
    if let Ok(metadata) = std::fs::symlink_metadata(repo_path) {
        config.in_repo = true;
        config.is_dir = metadata.is_dir();
    }

    let target_path = std::path::Path::new(&config.target_path);
    if let Ok(metadata) = std::fs::symlink_metadata(target_path) {
        config.exists = true;
        let file_type = metadata.file_type();
        if file_type.is_symlink() {
            config.is_symlink = true;
            if let Ok(link_dest) = std::fs::read_link(target_path) {
                config.symlink_to = link_dest.to_string_lossy().to_string();
            }
        }
    }
    config
}

fn copy_dir_all(src: impl AsRef<std::path::Path>, dst: impl AsRef<std::path::Path>) -> std::io::Result<()> {
    std::fs::create_dir_all(&dst)?;
    for entry in std::fs::read_dir(src)? {
        let entry = entry?;
        let ty = entry.file_type()?;
        if ty.is_dir() {
            copy_dir_all(entry.path(), dst.as_ref().join(entry.file_name()))?;
        } else {
            std::fs::copy(entry.path(), dst.as_ref().join(entry.file_name()))?;
        }
    }
    Ok(())
}

fn adopt_config_sync(config: &ConfigInfo, project_root: &std::path::Path) -> Result<(), String> {
    let info = get_info(config.clone(), project_root);
    if !info.exists {
        return Err("Source does not exist".to_string());
    }
    if info.in_repo {
        return Ok(());
    }

    let repo_path = std::path::Path::new(&info.repo_path);
    if let Some(parent) = repo_path.parent() {
        std::fs::create_dir_all(parent).map_err(|e| e.to_string())?;
    }

    if info.is_dir {
        copy_dir_all(&info.target_path, &info.repo_path).map_err(|e| e.to_string())?;
    } else {
        std::fs::copy(&info.target_path, &info.repo_path).map_err(|e| e.to_string())?;
    }
    Ok(())
}

fn link_config_sync(config: &ConfigInfo, project_root: &std::path::Path) -> Result<(), String> {
    let info = get_info(config.clone(), project_root);
    if !info.in_repo {
        return Err("Must be in repo before linking".to_string());
    }

    let mut cmd = if info.is_dir {
        let mut c = Command::new("powershell");
        c.args([
            "-Command",
            &format!("New-Item -ItemType Junction -Path \"{}\" -Target \"{}\" -Force", info.target_path, info.repo_path)
        ]);
        c
    } else {
        let mut c = Command::new("powershell");
        c.args([
            "-Command",
            &format!("New-Item -ItemType SymbolicLink -Path \"{}\" -Target \"{}\" -Force", info.target_path, info.repo_path)
        ]);
        c
    };
    configure_command(&mut cmd);
    let link_result = cmd.output();

    let success = match link_result {
        Ok(output) => output.status.success(),
        Err(_) => false,
    };

    if !success {
        println!("Link failed for {} (likely elevation required). Falling back to direct sync-copy.", config.name);
        let target_path = std::path::Path::new(&info.target_path);
        // Clean up target path if it exists
        if target_path.exists() {
            if info.is_dir {
                let _ = std::fs::remove_dir_all(target_path);
            } else {
                let _ = std::fs::remove_file(target_path);
            }
        }
        if let Some(parent) = target_path.parent() {
            let _ = std::fs::create_dir_all(parent);
        }
        if info.is_dir {
            copy_dir_all(&info.repo_path, target_path).map_err(|e| e.to_string())?;
        } else {
            std::fs::copy(&info.repo_path, target_path).map_err(|e| e.to_string())?;
        }
    }

    Ok(())
}

fn perform_auto_adopt_and_link(project_root: &std::path::Path) {
    println!("Starting automatic adopt and link process...");
    for c in get_configs_definition() {
        let mut info = get_info(c.clone(), project_root);
        // 1. Adopt if exists on C: but not in repo
        if info.exists && !info.in_repo && !c.skip_link {
            println!("Auto-Adopting configuration: {}...", c.name);
            if let Err(e) = adopt_config_sync(&c, project_root) {
                println!("Warning: failed to auto-adopt {}: {}", c.name, e);
            }
        }

        // Refresh info after possible adopt
        info = get_info(c.clone(), project_root);

        // 2. Link if in repo but not symlinked
        if info.in_repo && !info.is_symlink && !c.skip_link {
            println!("Auto-Linking configuration: {}...", c.name);
            if let Err(e) = link_config_sync(&c, project_root) {
                println!("Warning: failed to auto-link {}: {}", c.name, e);
            }
        }
    }
    println!("Auto-adopt and link process completed.");
}

// Helpers to parse lua values without regex crate dependency
fn get_theme_block(lua_content: &str) -> Option<&str> {
    let start_markers = ["-- @theme", "--@theme"];
    let end_markers = ["-- @theme-end", "--@theme-end"];
    
    let mut start_idx = None;
    for marker in &start_markers {
        if let Some(idx) = lua_content.find(marker) {
            start_idx = Some(idx + marker.len());
            break;
        }
    }
    
    let start = start_idx?;
    
    let mut end_idx = None;
    for marker in &end_markers {
        if let Some(idx) = lua_content[start..].find(marker) {
            end_idx = Some(start + idx);
            break;
        }
    }
    
    let end = end_idx?;
    Some(&lua_content[start..end])
}

fn get_string_val(lua_block: &str, key: &str, default_val: &str) -> String {
    let prefix1 = format!("{} = '", key);
    let prefix2 = format!("{} = \"", key);
    
    if let Some(idx) = lua_block.find(&prefix1) {
        let start = idx + prefix1.len();
        if let Some(end) = lua_block[start..].find('\'') {
            return lua_block[start..start+end].to_string();
        }
    }
    if let Some(idx) = lua_block.find(&prefix2) {
        let start = idx + prefix2.len();
        if let Some(end) = lua_block[start..].find('"') {
            return lua_block[start..start+end].to_string();
        }
    }
    
    let prefix3 = format!("{}='", key);
    let prefix4 = format!("{}=\"", key);
    if let Some(idx) = lua_block.find(&prefix3) {
        let start = idx + prefix3.len();
        if let Some(end) = lua_block[start..].find('\'') {
            return lua_block[start..start+end].to_string();
        }
    }
    if let Some(idx) = lua_block.find(&prefix4) {
        let start = idx + prefix4.len();
        if let Some(end) = lua_block[start..].find('"') {
            return lua_block[start..start+end].to_string();
        }
    }
    
    default_val.to_string()
}

fn get_float_val(lua_block: &str, key: &str, default_val: f64) -> f64 {
    let prefix1 = format!("{} = ", key);
    let prefix2 = format!("{}=", key);
    
    let parse_from = |start_str: &str| -> Option<f64> {
        let mut num_str = String::new();
        for c in start_str.chars() {
            if c.is_digit(10) || c == '.' {
                num_str.push(c);
            } else {
                break;
            }
        }
        num_str.parse::<f64>().ok()
    };
    
    if let Some(idx) = lua_block.find(&prefix1) {
        if let Some(val) = parse_from(&lua_block[idx + prefix1.len()..]) {
            return val;
        }
    }
    if let Some(idx) = lua_block.find(&prefix2) {
        if let Some(val) = parse_from(&lua_block[idx + prefix2.len()..]) {
            return val;
        }
    }
    
    default_val
}

// TAURI COMMANDS

#[tauri::command]
fn get_status() -> Result<Vec<ConfigInfo>, String> {
    let project_root = get_project_root();
    let mut result = Vec::new();
    for c in get_configs_definition() {
        result.push(get_info(c, &project_root));
    }
    Ok(result)
}

#[tauri::command]
fn get_config_status(key: String) -> Result<ConfigInfo, String> {
    let project_root = get_project_root();
    for c in get_configs_definition() {
        if c.key == key {
            return Ok(get_info(c, &project_root));
        }
    }
    Err("Config key not found".to_string())
}

#[tauri::command]
fn read_config_file(key: String) -> Result<String, String> {
    let project_root = get_project_root();
    let configs = get_configs_definition();
    let found = configs.iter().find(|c| c.key == key)
        .ok_or_else(|| "Config not found".to_string())?;

    let info = get_info(found.clone(), &project_root);
    if !info.in_repo {
        return Err("Config not in repository yet".to_string());
    }
    if info.is_dir {
        return Err("Target is a directory, cannot read directly".to_string());
    }

    std::fs::read_to_string(&info.repo_path).map_err(|e| e.to_string())
}

fn restart_pinterest_collage() {
    let home = get_home_dir();
    let appdata = get_appdata_dir();
    
    let config_path = find_pinterest_collage_config();
    let mut exe_path = None;
    
    let path = std::path::Path::new(&config_path);
    if let Some(parent) = path.parent() {
        let possible_exe = parent.join("pinterest-collage.exe");
        if possible_exe.exists() {
            exe_path = Some(possible_exe);
        }
    }
    
    if exe_path.is_none() {
        let candidates = [
            std::path::PathBuf::from("E:\\Dev\\projects\\pinterest-collage\\pinterest-collage.exe"),
            home.join("Dev").join("projects").join("pinterest-collage").join("pinterest-collage.exe"),
            home.join("projects").join("pinterest-collage").join("pinterest-collage.exe"),
            appdata.join("pinterest-collage").join("pinterest-collage.exe"),
        ];

        for path in &candidates {
            if path.exists() {
                exe_path = Some(path.clone());
                break;
            }
        }
    }

    if let Some(exe) = exe_path {
        let exe_str = exe.to_string_lossy().to_string();
        let parent_str = exe.parent().map(|p| p.to_string_lossy().to_string()).unwrap_or_default();
        let mut cmd = Command::new("powershell");
        cmd.args([
            "-NoProfile",
            "-Command",
            &format!(
                "Stop-Process -Name pinterest-collage -Force -ErrorAction SilentlyContinue; Start-Sleep -Milliseconds 200; Start-Process '{}' -WorkingDirectory '{}' -WindowStyle Hidden",
                exe_str, parent_str
            )
        ]);
        configure_command(&mut cmd);
        let _ = cmd.output();
    }
}

#[tauri::command]
fn write_config_file(key: String, content: String) -> Result<(), String> {
    let project_root = get_project_root();
    let configs = get_configs_definition();
    let found = configs.iter().find(|c| c.key == key)
        .ok_or_else(|| "Config not found".to_string())?;

    let info = get_info(found.clone(), &project_root);
    if info.is_dir {
        return Err("Target is a directory, cannot write directly".to_string());
    }

    let repo_path = std::path::Path::new(&info.repo_path);
    if let Some(parent) = repo_path.parent() {
        std::fs::create_dir_all(parent).map_err(|e| e.to_string())?;
    }

    std::fs::write(repo_path, &content).map_err(|e| e.to_string())?;

    // Fallback copy: if symlinking is not active, also write to the target on C:
    if !info.is_symlink && info.exists {
        let _ = std::fs::write(&info.target_path, &content);
    }

    if key == "pinterest_collage" {
        restart_pinterest_collage();
    }

    Ok(())
}

#[tauri::command]
fn adopt_config(key: String) -> Result<(), String> {
    let project_root = get_project_root();
    let configs = get_configs_definition();
    let found = configs.iter().find(|c| c.key == key)
        .ok_or_else(|| "Config not found".to_string())?;

    adopt_config_sync(found, &project_root)
}

#[tauri::command]
fn link_config(key: String) -> Result<(), String> {
    let project_root = get_project_root();
    let configs = get_configs_definition();
    let found = configs.iter().find(|c| c.key == key)
        .ok_or_else(|| "Config not found".to_string())?;

    link_config_sync(found, &project_root)
}

#[tauri::command]
fn run_sync() -> Result<String, String> {
    let sync_script_path = get_home_dir().join("Theme").join("sync_theme.ps1");
    let sync_script = sync_script_path.to_string_lossy();
    let mut cmd = Command::new("powershell");
    cmd.args(["-ExecutionPolicy", "Bypass", "-File", &sync_script]);
    configure_command(&mut cmd);
    let output = cmd.output()
        .map_err(|e| format!("Failed to run sync: {}", e))?;

    let out_str = String::from_utf8_lossy(&output.stdout).to_string();
    let err_str = String::from_utf8_lossy(&output.stderr).to_string();

    if !output.status.success() {
        return Err(format!("Sync failed:\nStdout: {}\nStderr: {}", out_str, err_str));
    }
    Ok(out_str)
}

#[tauri::command]
fn reload_glazewm() -> Result<String, String> {
    let mut cmd = Command::new("glazewm");
    cmd.args(["command", "wm-reload-config"]);
    configure_command(&mut cmd);
    let output = cmd.output()
        .map_err(|e| format!("Failed to reload GlazeWM: {}", e))?;

    let out_str = String::from_utf8_lossy(&output.stdout).to_string();
    let err_str = String::from_utf8_lossy(&output.stderr).to_string();

    if !output.status.success() {
        return Err(format!("Reload failed:\nStdout: {}\nStderr: {}", out_str, err_str));
    }
    Ok(out_str)
}

#[tauri::command]
fn get_theme() -> Result<ThemeInfo, String> {
    let project_root = get_project_root();
    let configs = get_configs_definition();
    let wezterm_config = get_info(configs[0].clone(), &project_root);

    let default_theme = ThemeInfo {
        font_family: "JetBrains Mono".to_string(),
        font_size: 14.0,
        opacity: 0.85,
        border_focused: "#2b6f7c".to_string(),
        border_unfocused: "#b2a9a5".to_string(),
        bg_color: "#f0edec".to_string(),
        fg_color: "#2c363c".to_string(),
        accent_color: "#cbd9e3".to_string(),
        lavender: "#dfd9e2".to_string(),
        lilac: "#c3acce".to_string(),
        lavender_grey: "#89909f".to_string(),
        pine_blue: "#538083".to_string(),
        jungle_teal: "#2a7f62".to_string(),
    };

    if !wezterm_config.in_repo {
        return Ok(default_theme);
    }

    let lua_content = std::fs::read_to_string(&wezterm_config.repo_path)
        .map_err(|e| format!("Failed to read WezTerm config: {}", e))?;

    let theme_block = match get_theme_block(&lua_content) {
        Some(block) => block,
        None => return Ok(default_theme),
    };

    let theme = ThemeInfo {
        font_family: get_string_val(theme_block, "font_family", &default_theme.font_family),
        font_size: get_float_val(theme_block, "font_size", default_theme.font_size),
        opacity: get_float_val(theme_block, "opacity", default_theme.opacity),
        border_focused: get_string_val(theme_block, "border_focused", &default_theme.border_focused),
        border_unfocused: get_string_val(theme_block, "border_unfocused", &default_theme.border_unfocused),
        bg_color: get_string_val(theme_block, "bg_color", &default_theme.bg_color),
        fg_color: get_string_val(theme_block, "fg_color", &default_theme.fg_color),
        accent_color: get_string_val(theme_block, "accent_color", &default_theme.accent_color),
        lavender: get_string_val(theme_block, "lavender", &default_theme.lavender),
        lilac: get_string_val(theme_block, "lilac", &default_theme.lilac),
        lavender_grey: get_string_val(theme_block, "lavender_grey", &default_theme.lavender_grey),
        pine_blue: get_string_val(theme_block, "pine_blue", &default_theme.pine_blue),
        jungle_teal: get_string_val(theme_block, "jungle_teal", &default_theme.jungle_teal),
    };

    Ok(theme)
}

#[tauri::command]
fn save_theme(theme: ThemeInfo) -> Result<String, String> {
    let project_root = get_project_root();
    let configs = get_configs_definition();
    let wezterm_config = get_info(configs[0].clone(), &project_root);

    if !wezterm_config.in_repo {
        return Err("WezTerm config must be adopted first".to_string());
    }

    let lua_content = std::fs::read_to_string(&wezterm_config.repo_path)
        .map_err(|e| format!("Failed to read WezTerm lua: {}", e))?;

    let new_theme_block = format!(
        "-- @theme\nlocal theme = {{\n  color_scheme = 'zenbones',\n  font_family = '{}',\n  font_size = {:.1},\n  opacity = {:.2},\n  border_focused = '{}',\n  border_unfocused = '{}',\n  bg_color = '{}',\n  fg_color = '{}',\n  accent_color = '{}',\n  lavender = '{}',\n  lilac = '{}',\n  lavender_grey = '{}',\n  pine_blue = '{}',\n  jungle_teal = '{}',\n}}\n-- @theme-end",
        theme.font_family,
        theme.font_size,
        theme.opacity,
        theme.border_focused,
        theme.border_unfocused,
        theme.bg_color,
        theme.fg_color,
        theme.accent_color,
        theme.lavender,
        theme.lilac,
        theme.lavender_grey,
        theme.pine_blue,
        theme.jungle_teal
    );

    let start_markers = ["-- @theme", "--@theme"];
    let end_markers = ["-- @theme-end", "--@theme-end"];
    
    let mut start_idx = None;
    for marker in &start_markers {
        if let Some(idx) = lua_content.find(marker) {
            start_idx = Some(idx);
            break;
        }
    }
    
    let mut end_idx = None;
    if let Some(start) = start_idx {
        for marker in &end_markers {
            if let Some(idx) = lua_content[start..].find(marker) {
                end_idx = Some(start + idx + marker.len());
                break;
            }
        }
    }

    let updated_lua = match (start_idx, end_idx) {
        (Some(start), Some(end)) => {
            let mut s = lua_content[..start].to_string();
            s.push_str(&new_theme_block);
            s.push_str(&lua_content[end..]);
            s
        }
        _ => {
            let mut s = lua_content.clone();
            s.push_str("\n");
            s.push_str(&new_theme_block);
            s
        }
    };

    std::fs::write(&wezterm_config.repo_path, &updated_lua)
        .map_err(|e| format!("Failed to write WezTerm lua: {}", e))?;

    if !wezterm_config.is_symlink && wezterm_config.exists {
        let _ = std::fs::write(&wezterm_config.target_path, &updated_lua);
    }

    run_sync()
}

#[tauri::command]
fn get_presets() -> Result<Vec<CustomPreset>, String> {
    let project_root = get_project_root();
    let presets_file = project_root.join("config").join("theme_presets.json");

    if !presets_file.exists() {
        return Ok(Vec::new());
    }

    let data = std::fs::read_to_string(&presets_file)
        .map_err(|e| format!("Failed to read presets file: {}", e))?;

    let list: Vec<CustomPreset> = serde_json::from_str(&data)
        .unwrap_or_else(|_| Vec::new());

    Ok(list)
}

#[tauri::command]
fn save_preset(preset: CustomPreset) -> Result<(), String> {
    let project_root = get_project_root();
    let config_dir = project_root.join("config");
    std::fs::create_dir_all(&config_dir).map_err(|e| e.to_string())?;
    
    let presets_file = config_dir.join("theme_presets.json");

    let mut list = get_presets().unwrap_or_else(|_| Vec::new());

    if let Some(idx) = list.iter().position(|p| p.name == preset.name) {
        list[idx] = preset;
    } else {
        list.push(preset);
    }

    let out_bytes = serde_json::to_string_pretty(&list)
        .map_err(|e| format!("Failed to serialize presets: {}", e))?;

    std::fs::write(presets_file, out_bytes)
        .map_err(|e| format!("Failed to write presets file: {}", e))?;

    Ok(())
}

#[tauri::command]
fn delete_preset(name: String) -> Result<(), String> {
    let project_root = get_project_root();
    let presets_file = project_root.join("config").join("theme_presets.json");

    if !presets_file.exists() {
        return Ok(());
    }

    let mut list = get_presets().unwrap_or_else(|_| Vec::new());
    list.retain(|p| p.name != name);

    let out_bytes = serde_json::to_string_pretty(&list)
        .map_err(|e| format!("Failed to serialize presets: {}", e))?;

    std::fs::write(presets_file, out_bytes)
        .map_err(|e| format!("Failed to write presets file: {}", e))?;

    Ok(())
}

#[tauri::command]
fn rename_preset(old_name: String, new_name: String) -> Result<(), String> {
    let project_root = get_project_root();
    let presets_file = project_root.join("config").join("theme_presets.json");

    if !presets_file.exists() {
        return Err("Presets file not found".to_string());
    }

    let mut list = get_presets().unwrap_or_else(|_| Vec::new());
    
    // Check if new_name already exists to prevent duplicate issues
    if list.iter().any(|p| p.name == new_name) {
        return Err("A preset with the new name already exists".to_string());
    }

    if let Some(idx) = list.iter().position(|p| p.name == old_name) {
        list[idx].name = new_name;
    } else {
        return Err("Preset not found".to_string());
    }

    let out_bytes = serde_json::to_string_pretty(&list)
        .map_err(|e| format!("Failed to serialize presets: {}", e))?;

    std::fs::write(presets_file, out_bytes)
        .map_err(|e| format!("Failed to write presets file: {}", e))?;

    Ok(())
}

#[tauri::command]
fn write_log(message: String) -> Result<(), String> {
    let log_path = get_appdata_dir().join("winrice-debug.log");
    let mut file = std::fs::OpenOptions::new()
        .create(true)
        .append(true)
        .open(&log_path)
        .map_err(|e| e.to_string())?;
    let now = std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .unwrap_or_default()
        .as_secs();
    writeln!(file, "[{}] {}", now, message).map_err(|e| e.to_string())?;
    Ok(())
}

#[tauri::command]
fn get_log_path() -> String {
    get_appdata_dir().join("winrice-debug.log").to_string_lossy().to_string()
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .setup(|app| {
            if cfg!(debug_assertions) {
                app.handle().plugin(
                    tauri_plugin_log::Builder::default()
                        .level(log::LevelFilter::Info)
                        .build(),
                )?;
            }

            // Run automatic adopt and link in a background thread on startup so it does not block the GUI window
            let project_root = get_project_root();
            std::thread::spawn(move || {
                perform_auto_adopt_and_link(&project_root);
            });

            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            get_status,
            get_config_status,
            read_config_file,
            write_config_file,
            adopt_config,
            link_config,
            run_sync,
            reload_glazewm,
            get_theme,
            save_theme,
            get_presets,
            save_preset,
            delete_preset,
            rename_preset,
            write_log,
            get_log_path,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

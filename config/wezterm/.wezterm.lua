local wezterm = require 'wezterm'
local config = wezterm.config_builder()

-- @theme
local theme = {
  color_scheme = 'zenbones',
  font_family = 'JetBrains Mono',
  font_size = 14.0,
  opacity = 0.85,
  border_focused = '#2b6f7c',
  border_unfocused = '#b2a9a5',
  bg_color = '#ff4000',
  fg_color = '#2c363c',
  accent_color = '#cbd9e3',
  lavender = '#dfd9e2',
  lilac = '#c3acce',
  lavender_grey = '#89909f',
  pine_blue = '#538083',
  jungle_teal = '#2a7f62',
}
-- @theme-end

-- Create a custom high-contrast variant of the zenbones color scheme
local scheme = wezterm.color.get_builtin_schemes()[theme.color_scheme]
if scheme then
  -- Override pale colors (blue, magenta, cyan) to be darker/high-contrast
  scheme.ansi[5] = '#285c7a' -- darker blue
  scheme.ansi[6] = '#7a4d7a' -- darker magenta
  scheme.ansi[7] = '#1e6a7a' -- darker cyan

  scheme.brights[5] = '#3a7fa8' -- bright blue
  scheme.brights[6] = '#a85fa8' -- bright magenta
  scheme.brights[7] = '#2894a8' -- bright cyan

  config.color_schemes = {
    ['zenbones_custom'] = scheme,
  }
  config.color_scheme = 'zenbones_custom'
else
  config.color_scheme = theme.color_scheme
end

-- Run synchronization script to apply theme to GlazeWM and Zebar
local success, stdout, stderr = wezterm.run_child_process({
  'powershell.exe',
  '-NoProfile',
  '-ExecutionPolicy', 'Bypass',
  '-File', 'C:/Users/Timofey/Theme/sync_theme.ps1'
})
if not success then
  wezterm.log_error("Theme sync failed: " .. (stderr or ""))
end

-- Window Decorations & Tabs (Clean borderless look for GlazeWM)
config.window_decorations = "NONE"
config.enable_tab_bar = false

-- Font Configuration
config.font = wezterm.font(theme.font_family)
config.font_size = theme.font_size

-- Disable ligatures globally
config.harfbuzz_features = { 'calt=0', 'clig=0', 'liga=0' }

-- Color Theme & Custom Color Overrides (for high contrast)
config.colors = {
  foreground = theme.fg_color,
  background = theme.bg_color,
  cursor_bg = theme.border_focused,
  cursor_border = theme.border_focused,
  selection_fg = theme.bg_color,
  selection_bg = theme.border_focused,
}

-- Background Transparency & Blur (Acrylic)
config.window_background_opacity = theme.opacity
config.win32_system_backdrop = 'Acrylic'

-- Custom Keybindings
config.keys = {
  {
    key = 'V',
    mods = 'CTRL|SHIFT',
    action = wezterm.action_callback(function(window, pane)
      -- Run PowerShell helper to detect clipboard formats
      local success, stdout, stderr = wezterm.run_child_process({
        'powershell.exe',
        '-NoProfile',
        '-Command',
        [[
          Add-Type -AssemblyName System.Windows.Forms;
          if ([System.Windows.Forms.Clipboard]::ContainsFileDropList()) {
              $files = [System.Windows.Forms.Clipboard]::GetFileDropList();
              Write-Output $files[0];
          } elseif ([System.Windows.Forms.Clipboard]::ContainsImage()) {
              $latest = Get-ChildItem -Path "C:\Users\Timofey\Pictures\Screenshots" -File | Sort-Object LastWriteTime -Descending | Select-Object -First 1;
              if ($latest) {
                  Write-Output $latest.FullName;
              }
          }
        ]]
      })

      -- Clean up the output path
      local path = ""
      if success and stdout then
        path = stdout:gsub('^%s*(.-)%s*$', '%1')
      end

      -- If we successfully retrieved an image/file path, paste it; otherwise, do normal text paste
      if path ~= "" then
        pane:send_text(path)
      else
        window:perform_action(wezterm.action.PasteFrom 'Clipboard', pane)
      end
    end)
  }
}

-- Default program (run fastfetch on startup)
config.default_prog = { 'C:/Users/Timofey/AppData/Local/Microsoft/WindowsApps/pwsh.exe', '-NoExit', '-Command', 'fastfetch' }

return config

local wezterm = require 'wezterm'
local config = wezterm.config_builder()

-- @theme
local theme = {
  color_scheme = 'zenbones',
  font_family = 'JetBrains Mono',
  font_size = 14.0,
  opacity = 0.85,
  border_focused = '#fe8019',
  border_unfocused = '#7c6f64',
  bg_color = '#282828',
  fg_color = '#ebdbb2',
  accent_color = '#504945',
  lavender = '#d3869b',
  lilac = '#b16286',
  lavender_grey = '#a89984',
  pine_blue = '#83a598',
  jungle_teal = '#8ec07c',
}
-- @theme-end

local bg = theme.bg_color:lower()
local is_dark = bg:match("^#1") or bg:match("^#2") or bg:match("^#3")
local dim_color = is_dark and '#909090' or '#504945'

-- Create a custom high-contrast and readability variant of the selected color scheme
local scheme = wezterm.color.get_builtin_schemes()[theme.color_scheme]
if scheme then
  -- Override pale colors (blue, magenta, cyan) to be darker/high-contrast
  scheme.ansi[5] = '#285c7a' -- darker blue
  scheme.ansi[6] = '#7a4d7a' -- darker magenta
  scheme.ansi[7] = '#1e6a7a' -- darker cyan

  scheme.brights[5] = '#3a7fa8' -- bright blue
  scheme.brights[6] = '#a85fa8' -- bright magenta
  scheme.brights[7] = '#2894a8' -- bright cyan

  -- Readability overrides for gray text (brights[1]) in dark/light modes
  if is_dark then
    scheme.brights[1] = '#909090' -- Make gray text bright enough on dark background
    scheme.ansi[1] = '#505050'
  else
    scheme.brights[1] = '#504945' -- Make gray text dark enough on sand background
    scheme.ansi[1] = '#2c363c'
  end

  config.color_schemes = {
    ['custom_scheme'] = scheme,
  }
  config.color_scheme = 'custom_scheme'
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
local indexed_colors = {}
for i = 232, 255 do
  indexed_colors[i] = dim_color
end

config.colors = {
  foreground = theme.fg_color,
  background = theme.bg_color,
  cursor_bg = theme.border_focused,
  cursor_border = theme.border_focused,
  selection_fg = theme.bg_color,
  selection_bg = theme.border_focused,
  indexed = indexed_colors,
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

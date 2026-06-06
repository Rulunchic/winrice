$weztermPath = "C:\Users\Timofey\.wezterm.lua"
$glazewmPath = "C:\Users\Timofey\.glzr\glazewm\config.yaml"
$zebarThemePath = "C:\Users\Timofey\.glzr\zebar\custom-mei\theme.css"

if (!(Test-Path $weztermPath)) {
    Write-Error "WezTerm configuration not found at $weztermPath"
    exit 1
}

$content = Get-Content $weztermPath -Raw

# Extract values from -- @theme block
if ($content -match '(?s)--\s*@theme\r?\n(.*?)\r?\n--\s*@theme-end') {
    $themeBlock = $Matches[1]
    
    function Get-StringValue($key) {
        if ($themeBlock -match "$key\s*=\s*['`"]([^'`"]+)['`"]") {
            return $Matches[1]
        }
        return $null
    }
    
    function Get-NumericValue($key) {
        if ($themeBlock -match "$key\s*=\s*([\d\.]+)") {
            return [double]$Matches[1]
        }
        return $null
    }
    
    $color_scheme = Get-StringValue "color_scheme"
    $font_family = Get-StringValue "font_family"
    $font_size = Get-NumericValue "font_size"
    $opacity = Get-NumericValue "opacity"
    $border_focused = Get-StringValue "border_focused"
    $border_unfocused = Get-StringValue "border_unfocused"
    $bg_color = Get-StringValue "bg_color"
    $fg_color = Get-StringValue "fg_color"
    $accent_color = Get-StringValue "accent_color"
    
    # 1. Update GlazeWM config.yaml border colors
    if (Test-Path $glazewmPath) {
        $yaml = Get-Content $glazewmPath -Raw
        
        # Replace focused window border color
        if ($border_focused) {
            $yaml = $yaml -replace "(?s)(focused_window:\s*(?:#[^\r\n]*\r?\n|\s)*border:\s*(?:#[^\r\n]*\r?\n|\s)*enabled:\s*\w+\s*(?:#[^\r\n]*\r?\n|\s)*color:\s*)'[^']+'", "`$1'$border_focused'"
        }
        
        # Replace other windows border color
        if ($border_unfocused) {
            $yaml = $yaml -replace "(?s)(other_windows:\s*(?:#[^\r\n]*\r?\n|\s)*border:\s*(?:#[^\r\n]*\r?\n|\s)*enabled:\s*\w+\s*(?:#[^\r\n]*\r?\n|\s)*color:\s*)'[^']+'", "`$1'$border_unfocused'"
        }
        
        Set-Content $glazewmPath $yaml -NoNewline
        Write-Host "Updated GlazeWM borders."
    } else {
        Write-Warning "GlazeWM config not found at $glazewmPath"
    }
    
    # 2. Update Zebar theme.css
    function Convert-HexToRgba($hex, $opacityVal) {
        if ($hex -match '#([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})') {
            $r = [Convert]::ToInt32($Matches[1], 16)
            $g = [Convert]::ToInt32($Matches[2], 16)
            $b = [Convert]::ToInt32($Matches[3], 16)
            return "rgba($r, $g, $b, $opacityVal)"
        }
        return $hex
    }
    
    $rgba_bg = if ($bg_color -and $opacity) { Convert-HexToRgba $bg_color $opacity } else { $bg_color }
    
    $cssContent = @"
/* Automatically generated from WezTerm theme. Do not edit directly. */
:root {
  --font-family: '$font_family', monospace;
  --font-size: ${font_size}px;
  --bg-color: $rgba_bg;
  --fg-color: $fg_color;
  --border-focused: $border_focused;
  --border-unfocused: $border_unfocused;
  --accent-color: $accent_color;
  --opacity: $opacity;

  /* Full user palette references */
  --lavender: #dfd9e2;
  --lilac: #c3acce;
  --lavender-grey: #89909f;
  --pine-blue: #538083;
  --jungle-teal: #2a7f62;
}
"@
    
    $parentDir = Split-Path $zebarThemePath -Parent
    if (!(Test-Path $parentDir)) {
        New-Item -ItemType Directory -Force -Path $parentDir | Out-Null
    }
    
    Set-Content $zebarThemePath $cssContent -NoNewline
    Write-Host "Updated Zebar theme.css."
    
    # 3. Reload GlazeWM config
    & glazewm command wm-reload-config
    Write-Host "Reloaded GlazeWM config."
} else {
    Write-Warning "Could not find a valid -- @theme block in $weztermPath"
}

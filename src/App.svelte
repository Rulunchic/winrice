<script lang="ts">
  import { onMount } from 'svelte';

  interface ConfigInfo {
    key: string;
    name: string;
    targetPath: string;
    repoPath: string;
    exists: boolean;
    inRepo: boolean;
    isSymlink: boolean;
    symlinkTo: string;
    isDir: boolean;
  }

  interface ThemeInfo {
    font_family: string;
    font_size: number;
    opacity: number;
    border_focused: string;
    border_unfocused: string;
    bg_color: string;
    fg_color: string;
    accent_color: string;
    lavender: string;
    lilac: string;
    lavender_grey: string;
    pine_blue: string;
    jungle_teal: string;
  }

  // Active view: 'theme' | 'configs'
  let activeTab = $state<'theme' | 'configs'>('theme');

  // Configurations state
  let configs = $state<ConfigInfo[]>([]);
  let selectedKey = $state<string>('');
  let rawContent = $state<string>('');
  let statusMsg = $state<string>('');
  let isLoading = $state<boolean>(false);
  let editMode = $state<'visual' | 'text'>('visual');

  // Global Theme state
  let theme = $state<ThemeInfo>({
    font_family: 'JetBrains Mono',
    font_size: 14.0,
    opacity: 0.85,
    border_focused: '#2b6f7c',
    border_unfocused: '#b2a9a5',
    bg_color: '#f0edec',
    fg_color: '#2c363c',
    accent_color: '#cbd9e3',
    lavender: '#dfd9e2',
    lilac: '#c3acce',
    lavender_grey: '#89909f',
    pine_blue: '#538083',
    jungle_teal: '#2a7f62',
  });

  // Visual Form state variables (synchronized to config file structures)
  let glazeInnerGap = $state<number>(10);
  let glazeOuterGap = $state<number>(10);
  let glazeFocusFollowsCursor = $state<boolean>(false);

  let zedFontSize = $state<number>(14);
  let zedFontFamily = $state<string>('JetBrains Mono');
  let zedMinimap = $state<boolean>(false);

  let vscodeFontSize = $state<number>(14);
  let vscodeFontFamily = $state<string>('JetBrains Mono');
  let vscodeMinimap = $state<boolean>(false);

  let gitName = $state<string>('');
  let gitEmail = $state<string>('');

  let komorebiBorder = $state<boolean>(false);
  let komorebiBorderWidth = $state<number>(2);
  let komorebiLayout = $state<string>('bsp');

  // Presets definition
  const presets = {
    zenbones: {
      name: 'Zenbones Sand',
      bg_color: '#f0edec',
      fg_color: '#2c363c',
      accent_color: '#cbd9e3',
      border_focused: '#2b6f7c',
      border_unfocused: '#b2a9a5',
      lavender: '#dfd9e2',
      lilac: '#c3acce',
      lavender_grey: '#89909f',
      pine_blue: '#538083',
      jungle_teal: '#2a7f62',
    },
    tokyonight: {
      name: 'Tokyo Night',
      bg_color: '#1a1b26',
      fg_color: '#c0caf5',
      accent_color: '#33467c',
      border_focused: '#7aa2f7',
      border_unfocused: '#3b4261',
      lavender: '#bb9af7',
      lilac: '#9d7cd8',
      lavender_grey: '#565f89',
      pine_blue: '#0db9d7',
      jungle_teal: '#41a6b5',
    },
    gruvbox: {
      name: 'Gruvbox Dark',
      bg_color: '#282828',
      fg_color: '#ebdbb2',
      accent_color: '#504945',
      border_focused: '#fe8019',
      border_unfocused: '#7c6f64',
      lavender: '#d3869b',
      lilac: '#b16286',
      lavender_grey: '#a89984',
      pine_blue: '#83a598',
      jungle_teal: '#8ec07c',
    },
    catppuccin: {
      name: 'Catppuccin Macchiato',
      bg_color: '#24273a',
      fg_color: '#cad3f5',
      accent_color: '#363a4f',
      border_focused: '#8aadf4',
      border_unfocused: '#494d64',
      lavender: '#b7bdf8',
      lilac: '#c6a0f6',
      lavender_grey: '#a5adcb',
      pine_blue: '#8bd5ca',
      jungle_teal: '#a6da95',
    }
  };

  // Computed state
  let selectedConfig = $derived(configs.find(c => c.key === selectedKey));

  async function fetchStatus() {
    try {
      const res = await fetch('/api/status');
      if (res.ok) {
        configs = await res.json();
        if (!selectedKey && configs.length > 0) {
          selectedKey = configs[0].key;
        }
      } else {
        statusMsg = 'Error fetching config status';
      }
    } catch (e) {
      statusMsg = `Connection failed: ${e}`;
    }
  }

  async function fetchTheme() {
    try {
      const res = await fetch('/api/theme');
      if (res.ok) {
        theme = await res.json();
      }
    } catch (e) {
      statusMsg = `Failed to fetch theme: ${e}`;
    }
  }

  async function saveTheme() {
    isLoading = true;
    statusMsg = 'Saving and syncing theme...';
    try {
      const res = await fetch('/api/theme/save', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(theme)
      });
      const output = await res.text();
      if (res.ok) {
        statusMsg = 'Theme saved and applied globally';
        setTimeout(() => { if (statusMsg.includes('saved')) statusMsg = ''; }, 2000);
      } else {
        statusMsg = `Theme save failed: ${output}`;
      }
    } catch (e) {
      statusMsg = `Theme save failed: ${e}`;
    } finally {
      isLoading = false;
    }
  }

  function applyPreset(key: keyof typeof presets) {
    const preset = presets[key];
    theme = {
      ...theme,
      bg_color: preset.bg_color,
      fg_color: preset.fg_color,
      accent_color: preset.accent_color,
      border_focused: preset.border_focused,
      border_unfocused: preset.border_unfocused,
      lavender: preset.lavender,
      lilac: preset.lilac,
      lavender_grey: preset.lavender_grey,
      pine_blue: preset.pine_blue,
      jungle_teal: preset.jungle_teal,
    };
    statusMsg = `Loaded preset: ${preset.name}`;
  }

  // Parse config values into visual variables
  function parseVisualVariables(key: string, content: string) {
    if (key === 'glazewm') {
      const innerMatch = content.match(/inner:\s*(\d+)/);
      if (innerMatch) glazeInnerGap = parseInt(innerMatch[1]);

      const outerMatch = content.match(/outer:\s*(\d+)/);
      if (outerMatch) glazeOuterGap = parseInt(outerMatch[1]);

      const focusMatch = content.match(/focus_follows_cursor:\s*(true|false)/);
      if (focusMatch) glazeFocusFollowsCursor = focusMatch[1] === 'true';
    } else if (key === 'zed') {
      try {
        const cleanJson = content.replace(/\/\/.*|\/\*[\s\S]*?\*\//g, ''); // strip comments
        const parsed = JSON.parse(cleanJson);
        if (parsed.buffer_font_size) zedFontSize = parsed.buffer_font_size;
        if (parsed.buffer_font_family) zedFontFamily = parsed.buffer_font_family;
        if (parsed.show_minimap !== undefined) zedMinimap = parsed.show_minimap;
      } catch (e) {
        console.warn('Failed to parse Zed settings:', e);
      }
    } else if (key === 'vscode') {
      try {
        const cleanJson = content.replace(/\/\/.*|\/\*[\s\S]*?\*\//g, '');
        const parsed = JSON.parse(cleanJson);
        if (parsed['editor.fontSize']) vscodeFontSize = parsed['editor.fontSize'];
        if (parsed['editor.fontFamily']) vscodeFontFamily = parsed['editor.fontFamily'];
        if (parsed['editor.minimap.enabled'] !== undefined) vscodeMinimap = parsed['editor.minimap.enabled'];
      } catch (e) {
        console.warn('Failed to parse VS Code settings:', e);
      }
    } else if (key === 'gitconfig') {
      const nameMatch = content.match(/name\s*=\s*(.+)/);
      if (nameMatch) gitName = nameMatch[1].trim();

      const emailMatch = content.match(/email\s*=\s*(.+)/);
      if (emailMatch) gitEmail = emailMatch[1].trim();
    } else if (key === 'komorebi') {
      try {
        const parsed = JSON.parse(content);
        if (parsed.enable_border !== undefined) komorebiBorder = parsed.enable_border;
        if (parsed.border_width !== undefined) komorebiBorderWidth = parsed.border_width;
        if (parsed.default_layout !== undefined) komorebiLayout = parsed.default_layout;
      } catch (e) {
        console.warn('Failed to parse Komorebi settings:', e);
      }
    }
  }

  // Serialize visual variables back to config format
  function serializeVisualVariables(key: string, content: string): string {
    if (key === 'glazewm') {
      let result = content;
      result = result.replace(/inner:\s*\d+/g, `inner: ${glazeInnerGap}`);
      result = result.replace(/outer:\s*\d+/g, `outer: ${glazeOuterGap}`);
      result = result.replace(/focus_follows_cursor:\s*(true|false)/g, `focus_follows_cursor: ${glazeFocusFollowsCursor}`);
      return result;
    } else if (key === 'zed') {
      try {
        const cleanJson = content.replace(/\/\/.*|\/\*[\s\S]*?\*\//g, '');
        const parsed = JSON.parse(cleanJson);
        parsed.buffer_font_size = zedFontSize;
        parsed.buffer_font_family = zedFontFamily;
        parsed.show_minimap = zedMinimap;
        return JSON.stringify(parsed, null, 2);
      } catch (e) {
        statusMsg = 'JSON formatting error, fallback to text editor';
        return content;
      }
    } else if (key === 'vscode') {
      try {
        const cleanJson = content.replace(/\/\/.*|\/\*[\s\S]*?\*\//g, '');
        const parsed = JSON.parse(cleanJson);
        parsed['editor.fontSize'] = vscodeFontSize;
        parsed['editor.fontFamily'] = vscodeFontFamily;
        parsed['editor.minimap.enabled'] = vscodeMinimap;
        return JSON.stringify(parsed, null, 2);
      } catch (e) {
        statusMsg = 'JSON formatting error';
        return content;
      }
    } else if (key === 'gitconfig') {
      let result = content;
      result = result.replace(/name\s*=\s*(.+)/g, `name = ${gitName}`);
      result = result.replace(/email\s*=\s*(.+)/g, `email = ${gitEmail}`);
      return result;
    } else if (key === 'komorebi') {
      try {
        const parsed = JSON.parse(content);
        parsed.enable_border = komorebiBorder;
        parsed.border_width = komorebiBorderWidth;
        parsed.default_layout = komorebiLayout;
        return JSON.stringify(parsed, null, 2);
      } catch (e) {
        return content;
      }
    }
    return content;
  }

  async function loadFile(key: string) {
    statusMsg = '';
    rawContent = '';
    
    const config = configs.find(c => c.key === key);
    if (!config || config.isDir || !config.inRepo) return;

    try {
      const res = await fetch(`/api/file?key=${key}`);
      if (res.ok) {
        rawContent = await res.text();
        parseVisualVariables(key, rawContent);
      } else {
        statusMsg = 'Failed to load file content';
      }
    } catch (e) {
      statusMsg = `Failed to load: ${e}`;
    }
  }

  async function saveConfig(visual: boolean = true) {
    if (!selectedKey) return;
    statusMsg = 'Saving...';
    
    const contentToSave = visual ? serializeVisualVariables(selectedKey, rawContent) : rawContent;

    try {
      const res = await fetch('/api/file/write', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ key: selectedKey, content: contentToSave })
      });
      if (res.ok) {
        statusMsg = 'Configuration saved';
        rawContent = contentToSave;
        setTimeout(() => { if (statusMsg === 'Configuration saved') statusMsg = ''; }, 2000);
      } else {
        const err = await res.text();
        statusMsg = `Save failed: ${err}`;
      }
    } catch (e) {
      statusMsg = `Save failed: ${e}`;
    }
  }

  async function runSync() {
    isLoading = true;
    statusMsg = 'Syncing theme...';
    try {
      const res = await fetch('/api/action/sync', { method: 'POST' });
      const output = await res.text();
      if (res.ok) {
        statusMsg = 'Theme sync complete';
      } else {
        statusMsg = `Sync failed: ${output}`;
      }
    } catch (e) {
      statusMsg = `Sync failed: ${e}`;
    } finally {
      isLoading = false;
    }
  }

  async function runReloadGlaze() {
    isLoading = true;
    statusMsg = 'Reloading GlazeWM...';
    try {
      const res = await fetch('/api/action/reload-glazewm', { method: 'POST' });
      const output = await res.text();
      if (res.ok) {
        statusMsg = 'GlazeWM config reloaded';
      } else {
        statusMsg = `Reload failed: ${output}`;
      }
    } catch (e) {
      statusMsg = `Reload failed: ${e}`;
    } finally {
      isLoading = false;
    }
  }

  onMount(async () => {
    await fetchStatus();
    await fetchTheme();
  });

  $effect(() => {
    if (selectedKey && activeTab === 'configs') {
      loadFile(selectedKey);
    }
  });
</script>

<div class="layout">
  <!-- Top Bar -->
  <header class="header panel">
    <div class="brand">WinRice Manager</div>
    <div class="tab-controls">
      <button class="tab-btn" class:active={activeTab === 'theme'} onclick={() => activeTab = 'theme'}>Theme Editor</button>
      <button class="tab-btn" class:active={activeTab === 'configs'} onclick={() => activeTab = 'configs'}>Configs ({configs.length})</button>
    </div>
    <div class="actions">
      <button onclick={runSync} disabled={isLoading}>Sync Theme</button>
      <button onclick={runReloadGlaze} disabled={isLoading}>Reload GlazeWM</button>
    </div>
  </header>

  <!-- Left Sidebar -->
  <aside class="sidebar panel">
    {#if activeTab === 'theme'}
      <div class="section-title">THEME PRESETS</div>
      <div class="preset-list">
        {#each Object.keys(presets) as k}
          <button class="preset-item" onclick={() => applyPreset(k as keyof typeof presets)}>
            <div class="preset-preview">
              <span class="preview-dot" style="background: {presets[k as keyof typeof presets].bg_color}"></span>
              <span class="preview-dot" style="background: {presets[k as keyof typeof presets].border_focused}"></span>
              <span class="preview-dot" style="background: {presets[k as keyof typeof presets].accent_color}"></span>
            </div>
            <span class="preset-name">{presets[k as keyof typeof presets].name}</span>
          </button>
        {/each}
      </div>
    {:else}
      <div class="section-title">CONFIGURATIONS</div>
      <div class="config-list">
        {#each configs as c}
          <button 
            class="config-item" 
            class:active={selectedKey === c.key}
            onclick={() => selectedKey = c.key}
          >
            <div class="name-row">
              <span class="cfg-name">{c.name}</span>
              <span class="badge" class:symlinked={c.isSymlink}>
                {c.isSymlink ? 'Linked' : 'Local'}
              </span>
            </div>
          </button>
        {/each}
      </div>
    {/if}
  </aside>

  <!-- Main Content Panel -->
  <main class="main panel">
    {#if activeTab === 'theme'}
      <!-- Visual Theme Editor -->
      <div class="theme-editor">
        <div class="editor-header">
          <h2>Visual Global Theme Settings</h2>
          <button class="primary" onclick={saveTheme} disabled={isLoading}>Save & Apply Theme</button>
        </div>

        <p class="subtitle">Modifying values here updates WezTerm config and runs the theme syncer to update GlazeWM and Zebar status bar globally.</p>

        <div class="theme-grid">
          <!-- Text properties -->
          <div class="theme-section">
            <h3>Terminal & Font settings</h3>
            <div class="editor-row">
              <label for="font-family">Font Family:</label>
              <input type="text" id="font-family" bind:value={theme.font_family} />
            </div>
            <div class="editor-row">
              <label for="font-size">Font Size (px):</label>
              <input type="number" id="font-size" step="0.5" bind:value={theme.font_size} />
            </div>
            <div class="editor-row">
              <label for="opacity">Opacity:</label>
              <input type="range" id="opacity" min="0.1" max="1.0" step="0.05" bind:value={theme.opacity} />
              <span class="range-val">{theme.opacity}</span>
            </div>
          </div>

          <!-- Color Pickers -->
          <div class="theme-section">
            <h3>Global System Colors</h3>
            
            <div class="color-picker-row">
              <div class="color-circle-wrapper" style="background: {theme.bg_color}">
                <input type="color" class="color-input" bind:value={theme.bg_color} />
              </div>
              <span class="color-label">Background Color (Sand)</span>
              <code class="hex-val">{theme.bg_color}</code>
            </div>

            <div class="color-picker-row">
              <div class="color-circle-wrapper" style="background: {theme.fg_color}">
                <input type="color" class="color-input" bind:value={theme.fg_color} />
              </div>
              <span class="color-label">Foreground Color (Stone)</span>
              <code class="hex-val">{theme.fg_color}</code>
            </div>

            <div class="color-picker-row">
              <div class="color-circle-wrapper" style="background: {theme.accent_color}">
                <input type="color" class="color-input" bind:value={theme.accent_color} />
              </div>
              <span class="color-label">Accent Color (Selection)</span>
              <code class="hex-val">{theme.accent_color}</code>
            </div>

            <div class="color-picker-row">
              <div class="color-circle-wrapper" style="background: {theme.border_focused}">
                <input type="color" class="color-input" bind:value={theme.border_focused} />
              </div>
              <span class="color-label">Focused Border (Teal)</span>
              <code class="hex-val">{theme.border_focused}</code>
            </div>

            <div class="color-picker-row">
              <div class="color-circle-wrapper" style="background: {theme.border_unfocused}">
                <input type="color" class="color-input" bind:value={theme.border_unfocused} />
              </div>
              <span class="color-label">Unfocused Border (Muted)</span>
              <code class="hex-val">{theme.border_unfocused}</code>
            </div>
          </div>

          <!-- Accent Palette Colors -->
          <div class="theme-section span-2">
            <h3>Custom Contrast Accent Palette</h3>
            <div class="palette-grid">
              
              <div class="color-picker-row">
                <div class="color-circle-wrapper" style="background: {theme.lavender}">
                  <input type="color" class="color-input" bind:value={theme.lavender} />
                </div>
                <span class="color-label">Lavender</span>
                <code class="hex-val">{theme.lavender}</code>
              </div>

              <div class="color-picker-row">
                <div class="color-circle-wrapper" style="background: {theme.lilac}">
                  <input type="color" class="color-input" bind:value={theme.lilac} />
                </div>
                <span class="color-label">Lilac</span>
                <code class="hex-val">{theme.lilac}</code>
              </div>

              <div class="color-picker-row">
                <div class="color-circle-wrapper" style="background: {theme.lavender_grey}">
                  <input type="color" class="color-input" bind:value={theme.lavender_grey} />
                </div>
                <span class="color-label">Lavender-Grey</span>
                <code class="hex-val">{theme.lavender_grey}</code>
              </div>

              <div class="color-picker-row">
                <div class="color-circle-wrapper" style="background: {theme.pine_blue}">
                  <input type="color" class="color-input" bind:value={theme.pine_blue} />
                </div>
                <span class="color-label">Pine Blue</span>
                <code class="hex-val">{theme.pine_blue}</code>
              </div>

              <div class="color-picker-row">
                <div class="color-circle-wrapper" style="background: {theme.jungle_teal}">
                  <input type="color" class="color-input" bind:value={theme.jungle_teal} />
                </div>
                <span class="color-label">Jungle Teal</span>
                <code class="hex-val">{theme.jungle_teal}</code>
              </div>

            </div>
          </div>
        </div>
      </div>
    {:else}
      <!-- Configurations Tab -->
      {#if selectedConfig}
        <div class="config-view-header">
          <div class="info-grid">
            <h2>{selectedConfig.name} Settings</h2>
            <div class="meta-row">
              <span class="meta-label">Active Path:</span>
              <code class="meta-value">{selectedConfig.targetPath}</code>
            </div>
            <div class="meta-row">
              <span class="meta-label">Sync Mode:</span>
              <span class="meta-value">{selectedConfig.isSymlink ? 'Junction Link' : 'Copy Sync (Fallback)'}</span>
            </div>
          </div>
          <div class="view-controls">
            <button class="tab-btn" class:active={editMode === 'visual'} onclick={() => editMode = 'visual'} disabled={selectedConfig.isDir || selectedConfig.key === 'sync_script'}>Visual Editor</button>
            <button class="tab-btn" class:active={editMode === 'text'} onclick={() => editMode = 'text'} disabled={selectedConfig.isDir}>Code Editor</button>
          </div>
        </div>

        <!-- Custom Visual Form depending on Config Key -->
        {#if editMode === 'visual'}
          <div class="visual-form panel">
            <!-- GlazeWM Visual settings -->
            {#if selectedConfig.key === 'glazewm'}
              <div class="visual-section">
                <h3>GlazeWM Window Tiling Settings</h3>
                <div class="editor-row">
                  <label for="glaze-inner">Inner Window Gap (px):</label>
                  <input type="number" id="glaze-inner" bind:value={glazeInnerGap} />
                </div>
                <div class="editor-row">
                  <label for="glaze-outer">Outer Border Gap (px):</label>
                  <input type="number" id="glaze-outer" bind:value={glazeOuterGap} />
                </div>
                <div class="editor-row">
                  <label class="checkbox-label">
                    <input type="checkbox" bind:checked={glazeFocusFollowsCursor} />
                    Focus Follows Cursor
                  </label>
                </div>
                <button class="primary" onclick={() => saveConfig(true)}>Save Settings</button>
              </div>

            <!-- Zed Editor Visual settings -->
            {:else if selectedConfig.key === 'zed'}
              <div class="visual-section">
                <h3>Zed Editor Global Settings</h3>
                <div class="editor-row">
                  <label for="zed-font">Font Family:</label>
                  <input type="text" id="zed-font" bind:value={zedFontFamily} />
                </div>
                <div class="editor-row">
                  <label for="zed-size">Font Size (px):</label>
                  <input type="number" id="zed-size" bind:value={zedFontSize} />
                </div>
                <div class="editor-row">
                  <label class="checkbox-label">
                    <input type="checkbox" bind:checked={zedMinimap} />
                    Show Code Minimap
                  </label>
                </div>
                <button class="primary" onclick={() => saveConfig(true)}>Save Settings</button>
              </div>

            <!-- VS Code Visual settings -->
            {:else if selectedConfig.key === 'vscode'}
              <div class="visual-section">
                <h3>VS Code Editor Settings</h3>
                <div class="editor-row">
                  <label for="vscode-font">Font Family:</label>
                  <input type="text" id="vscode-font" bind:value={vscodeFontFamily} />
                </div>
                <div class="editor-row">
                  <label for="vscode-size">Font Size (px):</label>
                  <input type="number" id="vscode-size" bind:value={vscodeFontSize} />
                </div>
                <div class="editor-row">
                  <label class="checkbox-label">
                    <input type="checkbox" bind:checked={vscodeMinimap} />
                    Enable Code Minimap
                  </label>
                </div>
                <button class="primary" onclick={() => saveConfig(true)}>Save Settings</button>
              </div>

            <!-- Git Config Visual settings -->
            {:else if selectedConfig.key === 'gitconfig'}
              <div class="visual-section">
                <h3>Git Global User Settings</h3>
                <div class="editor-row">
                  <label for="git-name">Global Username:</label>
                  <input type="text" id="git-name" bind:value={gitName} />
                </div>
                <div class="editor-row">
                  <label for="git-email">Global Email:</label>
                  <input type="text" id="git-email" bind:value={gitEmail} />
                </div>
                <button class="primary" onclick={() => saveConfig(true)}>Save Settings</button>
              </div>

            <!-- Komorebi Visual settings -->
            {:else if selectedConfig.key === 'komorebi'}
              <div class="visual-section">
                <h3>Komorebi Window Manager Settings</h3>
                <div class="editor-row">
                  <label class="checkbox-label">
                    <input type="checkbox" bind:checked={komorebiBorder} />
                    Enable Active Window Border
                  </label>
                </div>
                <div class="editor-row">
                  <label for="komorebi-width">Border Width (px):</label>
                  <input type="number" id="komorebi-width" bind:value={komorebiBorderWidth} />
                </div>
                <div class="editor-row">
                  <label for="komorebi-layout">Default Workspace Layout:</label>
                  <select id="komorebi-layout" bind:value={komorebiLayout}>
                    <option value="bsp">BSP (Binary Space Partitioning)</option>
                    <option value="columns">Columns</option>
                    <option value="rows">Rows</option>
                  </select>
                </div>
                <button class="primary" onclick={() => saveConfig(true)}>Save Settings</button>
              </div>

            <!-- WezTerm is configured via global theme, but show details here -->
            {:else if selectedConfig.key === 'wezterm'}
              <div class="visual-section">
                <h3>WezTerm Settings</h3>
                <p class="desc-text">WezTerm settings are managed directly in the **Theme Editor** tab to serve as the unified source of truth for the entire workspace style pipeline.</p>
                <button onclick={() => activeTab = 'theme'}>Open Theme Editor</button>
              </div>

            {:else if selectedConfig.isDir}
              <div class="editor-placeholder">
                Zebar is a directory containing multiple files. Files are linked automatically, but must be edited via terminal or filesystem inside <code>config/zebar/</code>.
              </div>
            {:else}
              <div class="editor-placeholder">
                Visual editor not available for this configuration. Use the **Code Editor** tab instead.
              </div>
            {/if}
          </div>
        {:else}
          <!-- Text/Code Editor Section -->
          <div class="editor-section">
            <div class="editor-header">
              <span>Plaintext Config Editor (Direct Repo Sync)</span>
              <button class="primary" onclick={() => saveConfig(false)}>Save Config</button>
            </div>
            <textarea bind:value={rawContent} spellcheck="false"></textarea>
          </div>
        {/if}
      {:else}
        <div class="editor-placeholder">
          Select a configuration file to inspect and manage.
        </div>
      {/if}
    {/if}
  </main>

  <!-- Footer / Status messages -->
  <footer class="footer panel">
    <div class="status-bar">
      <span>Status: {statusMsg || 'Idle'}</span>
      <span>API: http://127.0.0.1:54321</span>
    </div>
  </footer>
</div>

<style>
  .layout {
    display: grid;
    grid-template-columns: 240px 1fr;
    grid-template-rows: auto 1fr auto;
    height: 100vh;
    padding: 12px;
    gap: 8px;
    background: var(--bg-base);
  }

  .header {
    grid-column: 1 / span 2;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 16px;
  }

  .brand {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--accent);
  }

  .tab-controls {
    display: flex;
    gap: 4px;
    background: rgba(0, 0, 0, 0.05);
    padding: 2px;
    border-radius: 4px;
    border: 1px solid var(--border);
  }

  .tab-btn {
    background: transparent;
    border: none;
    padding: 4px 12px;
    border-radius: 2px;
    color: var(--fg-muted);
    font-weight: 500;
    cursor: pointer;
  }

  .tab-btn.active {
    background: var(--bg-base);
    color: var(--fg);
    border: 1px solid var(--border);
  }

  .actions {
    display: flex;
    gap: 8px;
  }

  .sidebar {
    display: flex;
    flex-direction: column;
    padding: 12px;
    gap: 12px;
    overflow-y: auto;
  }

  .section-title {
    font-size: 0.8rem;
    font-weight: 600;
    color: var(--fg-muted);
    letter-spacing: 0.05em;
  }

  .config-list, .preset-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .config-item, .preset-item {
    text-align: left;
    width: 100%;
    background: transparent;
    border: 1px solid transparent;
    padding: 8px;
    border-radius: 4px;
    color: var(--fg);
    cursor: pointer;
    font-family: var(--font-mono);
  }

  .config-item:hover, .preset-item:hover {
    background: rgba(44, 54, 60, 0.05);
    border-color: var(--border);
  }

  .config-item.active {
    background: var(--bg-base);
    border-color: var(--border);
    color: var(--accent);
  }

  .preset-item {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .preset-preview {
    display: flex;
    gap: 2px;
  }

  .preview-dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    border: 1px solid rgba(0, 0, 0, 0.15);
  }

  .preset-name {
    font-size: 0.85rem;
  }

  .name-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .badge {
    font-size: 0.75rem;
    padding: 1px 6px;
    border-radius: 4px;
    border: 1px solid var(--border);
    color: var(--accent);
    border-color: var(--accent);
    background: rgba(43, 111, 124, 0.05);
  }

  .main {
    display: flex;
    flex-direction: column;
    padding: 16px;
    gap: 16px;
    overflow-y: auto;
  }

  .config-view-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    border-bottom: 1px solid var(--border);
    padding-bottom: 12px;
  }

  .config-view-header .info-grid {
    border: none;
    padding: 0;
  }

  .view-controls {
    display: flex;
    gap: 4px;
    background: rgba(0, 0, 0, 0.03);
    padding: 2px;
    border-radius: 4px;
    border: 1px solid var(--border);
  }

  .theme-editor {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .subtitle {
    font-size: 0.85rem;
    color: var(--fg-muted);
  }

  .theme-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    margin-top: 8px;
  }

  .theme-section, .visual-form {
    display: flex;
    flex-direction: column;
    gap: 10px;
    background: rgba(0, 0, 0, 0.02);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 12px;
  }

  .visual-form {
    background: var(--bg-panel);
    flex-grow: 1;
  }

  .visual-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
    max-width: 500px;
  }

  .visual-section h3 {
    font-size: 0.95rem;
    color: var(--accent);
    border-bottom: 1px solid var(--border);
    padding-bottom: 4px;
    margin-bottom: 4px;
  }

  .desc-text {
    font-size: 0.875rem;
    color: var(--fg-muted);
    line-height: 1.5;
  }

  .theme-section.span-2 {
    grid-column: 1 / span 2;
  }

  .theme-section h3 {
    font-size: 0.95rem;
    color: var(--accent);
    border-bottom: 1px solid var(--border);
    padding-bottom: 4px;
    margin-bottom: 4px;
  }

  .editor-row {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 0.875rem;
  }

  .editor-row label {
    min-width: 170px;
    color: var(--fg-muted);
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    user-select: none;
  }

  .checkbox-label input {
    cursor: pointer;
  }

  .editor-row input[type="text"], .editor-row input[type="number"], .editor-row select {
    flex-grow: 1;
    padding: 4px 8px;
    border-radius: 4px;
    border: 1px solid var(--border);
    background: var(--bg-base);
    color: var(--fg);
    font-family: var(--font-mono);
  }

  .editor-row input[type="range"] {
    flex-grow: 1;
    cursor: pointer;
  }

  .range-val {
    min-width: 3ch;
    text-align: right;
  }

  .color-picker-row {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 0.875rem;
    padding: 2px 0;
  }

  .color-circle-wrapper {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    border: 1px solid rgba(0, 0, 0, 0.2);
    cursor: pointer;
    position: relative;
  }

  .color-input {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    opacity: 0;
    cursor: pointer;
  }

  .color-label {
    flex-grow: 1;
    color: var(--fg-muted);
  }

  .hex-val {
    font-size: 0.8rem;
    color: var(--fg);
  }

  .palette-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
  }

  .info-grid {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  h2 {
    font-size: 1.15rem;
    margin-bottom: 2px;
  }

  .meta-row {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.8rem;
  }

  .meta-label {
    min-width: 100px;
    color: var(--fg-muted);
  }

  .meta-value {
    background: rgba(0, 0, 0, 0.04);
    padding: 1px 6px;
    border-radius: 4px;
    font-size: 0.75rem;
    word-break: break-all;
  }

  .editor-section {
    display: flex;
    flex-direction: column;
    flex-grow: 1;
    gap: 8px;
    min-height: 250px;
  }

  .editor-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
    font-size: 0.9rem;
    color: var(--fg-muted);
  }

  textarea {
    flex-grow: 1;
    font-family: var(--font-mono);
    font-size: 0.875rem;
    padding: 8px;
    border-radius: 4px;
    border: 1px solid var(--border);
    background: var(--bg-base);
    color: var(--fg);
    resize: none;
    line-height: 1.45;
  }

  textarea:focus {
    outline: none;
    border-color: var(--accent);
  }

  .editor-placeholder {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-grow: 1;
    border: 1px dashed var(--border);
    border-radius: 4px;
    color: var(--fg-muted);
    padding: 24px;
    text-align: center;
  }

  .footer {
    grid-column: 1 / span 2;
    padding: 6px 12px;
  }

  .status-bar {
    display: flex;
    justify-content: space-between;
    font-size: 0.8rem;
    color: var(--fg-muted);
  }

  .panel {
    background: var(--bg-panel);
    border: 1px solid var(--border);
    border-radius: 4px;
  }

  button {
    font-family: var(--font-mono);
    font-size: 0.875rem;
    padding: 4px 12px;
    border-radius: 4px;
    border: 1px solid var(--border);
    background: var(--bg-base);
    color: var(--fg);
    cursor: pointer;
    transition: border-color 80ms ease, color 80ms ease;
  }

  button:hover {
    border-color: var(--accent);
    color: var(--accent);
  }

  button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    border-color: var(--border);
    color: var(--fg-muted);
  }

  button.primary {
    background: var(--accent);
    color: var(--accent-fg);
    border-color: var(--accent);
  }

  button.primary:hover {
    opacity: 0.9;
    color: var(--accent-fg);
  }
</style>

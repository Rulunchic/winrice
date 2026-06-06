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

  // Active view: 'configs' | 'theme'
  let activeTab = $state<'configs' | 'theme'>('theme');

  // Configurations state
  let configs = $state<ConfigInfo[]>([]);
  let selectedKey = $state<string>('');
  let editorContent = $state<string>('');
  let isEditing = $state<boolean>(false);
  let statusMsg = $state<string>('');
  let isLoading = $state<boolean>(false);

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

  async function loadFile(key: string) {
    statusMsg = '';
    isEditing = false;
    editorContent = '';
    
    const config = configs.find(c => c.key === key);
    if (!config || config.isDir || !config.inRepo) return;

    try {
      const res = await fetch(`/api/file?key=${key}`);
      if (res.ok) {
        editorContent = await res.text();
        isEditing = true;
      } else {
        statusMsg = 'Failed to load file content';
      }
    } catch (e) {
      statusMsg = `Failed to load: ${e}`;
    }
  }

  async function saveFile() {
    if (!selectedKey) return;
    statusMsg = 'Saving...';
    try {
      const res = await fetch('/api/file/write', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ key: selectedKey, content: editorContent })
      });
      if (res.ok) {
        statusMsg = 'Saved';
        setTimeout(() => { if (statusMsg === 'Saved') statusMsg = ''; }, 2000);
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

  <!-- Left Sidebar (Configs list or Presets list depending on activeTab) -->
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
        <div class="info-grid">
          <h2>{selectedConfig.name} Configuration</h2>
          
          <div class="meta-row">
            <span class="meta-label">Active Path:</span>
            <code class="meta-value">{selectedConfig.targetPath}</code>
          </div>
          <div class="meta-row">
            <span class="meta-label">Repository Path:</span>
            <code class="meta-value">{selectedConfig.repoPath}</code>
          </div>
          
          <div class="status-row">
            <span class="meta-label">Status:</span>
            <span class="status-value">
              {#if selectedConfig.isSymlink}
                Linked dynamically via symlink to repository.
              {:else}
                Warning: Configuration not symlinked yet.
              {/if}
            </span>
          </div>
        </div>

        <!-- Text Editor Section -->
        {#if isEditing}
          <div class="editor-section">
            <div class="editor-header">
              <span>Plaintext Config Editor (Direct Repo Sync)</span>
              <button class="primary" onclick={saveFile}>Save Config</button>
            </div>
            <textarea bind:value={editorContent} spellcheck="false"></textarea>
          </div>
        {:else if selectedConfig.isDir}
          <div class="editor-placeholder">
            Zebar is a directory containing multiple files. Files are linked automatically, but must be edited via terminal or filesystem inside <code>config/zebar/</code>.
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

  .theme-section {
    display: flex;
    flex-direction: column;
    gap: 10px;
    background: rgba(0, 0, 0, 0.02);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 12px;
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
    min-width: 120px;
    color: var(--fg-muted);
  }

  .editor-row input[type="text"], .editor-row input[type="number"] {
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
    gap: 8px;
    border-bottom: 1px solid var(--border);
    padding-bottom: 16px;
  }

  h2 {
    font-size: 1.15rem;
    margin-bottom: 4px;
  }

  .meta-row, .status-row {
    display: flex;
    align-items: flex-start;
    gap: 8px;
  }

  .meta-label {
    min-width: 130px;
    color: var(--fg-muted);
    font-weight: 500;
  }

  .meta-value {
    background: rgba(0, 0, 0, 0.04);
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 0.85rem;
    word-break: break-all;
  }

  .status-value {
    color: var(--fg);
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

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

  // App State (Svelte 5)
  let configs = $state<ConfigInfo[]>([]);
  let selectedKey = $state<string>('');
  let editorContent = $state<string>('');
  let isEditing = $state<boolean>(false);
  let statusMsg = $state<string>('');
  let isLoading = $state<boolean>(false);

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
        statusMsg = 'Error fetching status';
      }
    } catch (e) {
      statusMsg = `Connection failed: ${e}`;
    }
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

  async function adoptConfig(key: string) {
    statusMsg = 'Adopting...';
    try {
      const res = await fetch(`/api/action/adopt?key=${key}`, { method: 'POST' });
      if (res.ok) {
        statusMsg = 'Adopted successfully';
        await fetchStatus();
        if (selectedKey === key) {
          await loadFile(key);
        }
      } else {
        const err = await res.text();
        statusMsg = `Adopt failed: ${err}`;
      }
    } catch (e) {
      statusMsg = `Adopt failed: ${e}`;
    }
  }

  async function linkConfig(key: string) {
    statusMsg = 'Linking...';
    try {
      const res = await fetch(`/api/action/link?key=${key}`, { method: 'POST' });
      if (res.ok) {
        statusMsg = 'Linked successfully';
        await fetchStatus();
      } else {
        const err = await res.text();
        statusMsg = `Link failed: ${err}`;
      }
    } catch (e) {
      statusMsg = `Link failed: ${e}`;
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
        console.log(output);
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

  onMount(() => {
    fetchStatus();
  });

  $effect(() => {
    if (selectedKey) {
      loadFile(selectedKey);
    }
  });
</script>

<div class="layout">
  <!-- Top Bar -->
  <header class="header panel">
    <div class="brand">WinRice Manager</div>
    <div class="actions">
      <button onclick={runSync} disabled={isLoading}>Sync Theme</button>
      <button onclick={runReloadGlaze} disabled={isLoading}>Reload GlazeWM</button>
    </div>
  </header>

  <!-- Sidebar / File List -->
  <aside class="sidebar panel">
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
            <span class="badge" class:symlinked={c.isSymlink} class:adopted={c.inRepo && !c.isSymlink}>
              {c.isSymlink ? 'Linked' : c.inRepo ? 'Local' : 'Unlinked'}
            </span>
          </div>
        </button>
      {/each}
    </div>
  </aside>

  <!-- Main Content Panel -->
  <main class="main panel">
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
              Active configuration is symlinked to repository.
            {:else if selectedConfig.inRepo}
              Adopted locally. Not linked to active path.
            {:else if selectedConfig.exists}
              Exists at active path only. Ready to adopt.
            {:else}
              Config does not exist at active path.
            {/if}
          </span>
        </div>

        <div class="operations">
          {#if !selectedConfig.inRepo}
            <button onclick={() => adoptConfig(selectedConfig.key)}>Adopt into Repo</button>
          {/if}
          {#if selectedConfig.inRepo && !selectedConfig.isSymlink}
            <button class="primary" onclick={() => linkConfig(selectedConfig.key)}>Link to C:</button>
          {/if}
          {#if selectedConfig.isSymlink}
            <button onclick={() => linkConfig(selectedConfig.key)}>Re-create Symlink</button>
          {/if}
        </div>
      </div>

      <!-- Text Editor Section -->
      {#if isEditing}
        <div class="editor-section">
          <div class="editor-header">
            <span>Editor</span>
            <button class="primary" onclick={saveFile}>Save Changes</button>
          </div>
          <textarea bind:value={editorContent} spellcheck="false"></textarea>
        </div>
      {:else if selectedConfig.isDir}
        <div class="editor-placeholder">
          Zebar is a directory containing multiple files. Use file explorer or terminal to edit directories directly in <code>config/zebar/</code>.
        </div>
      {:else if !selectedConfig.inRepo}
        <div class="editor-placeholder">
          Adopt this configuration to edit it in the dashboard.
        </div>
      {/if}
    {:else}
      <div class="editor-placeholder">
        Select a configuration file to inspect and manage.
      </div>
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

  .config-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .config-item {
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

  .config-item:hover {
    background: rgba(44, 54, 60, 0.05);
    border-color: var(--border);
    color: var(--fg);
  }

  .config-item.active {
    background: var(--bg-base);
    border-color: var(--border);
    color: var(--accent);
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
    color: var(--fg-muted);
    background: rgba(0, 0, 0, 0.03);
  }

  .badge.symlinked {
    color: var(--accent);
    border-color: var(--accent);
    background: rgba(43, 111, 124, 0.05);
  }

  .badge.adopted {
    color: var(--fg);
    border-color: var(--fg-muted);
  }

  .main {
    display: flex;
    flex-direction: column;
    padding: 16px;
    gap: 16px;
    overflow-y: auto;
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

  .operations {
    display: flex;
    gap: 8px;
    margin-top: 8px;
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

  /* G2G Panel helper style */
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

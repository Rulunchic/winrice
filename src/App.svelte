<script lang="ts">
  let count = $state(0);
  let doubled = $derived(count * 2);

  $effect(() => {
    document.title = `count = ${count}`;
  });
</script>

<main>
  <h1>Svelte 5 + Vite + TypeScript</h1>

  <p>Geek-minimal starter. Edit <code>src/App.svelte</code> to start.</p>

  <div class="row">
    <button onclick={() => count--}>-</button>
    <span class="value">{count}</span>
    <button onclick={() => count++}>+</button>
  </div>

  <p class="muted">doubled = {doubled}</p>
</main>

<style>
  main {
    max-width: 720px;
    margin: 0 auto;
    padding: 32px 24px;
  }
  h1 { margin-bottom: 12px; }
  p { margin-bottom: 16px; }
  .row {
    display: flex;
    align-items: center;
    gap: 12px;
    margin: 24px 0;
  }
  button {
    padding: 4px 12px;
    border: 1px solid var(--border);
    border-radius: 4px;
    background: var(--bg-panel);
    color: var(--fg);
    cursor: pointer;
    transition: border-color 80ms ease, color 80ms ease;
  }
  button:hover {
    border-color: var(--accent);
    color: var(--accent);
  }
  .value {
    font-size: 1.25rem;
    min-width: 3ch;
    text-align: center;
    color: var(--accent);
  }
  .muted { color: var(--fg-muted); font-size: 0.875rem; }
  code {
    background: var(--bg-panel);
    border: 1px solid var(--border);
    padding: 1px 4px;
    border-radius: 2px;
  }
</style>

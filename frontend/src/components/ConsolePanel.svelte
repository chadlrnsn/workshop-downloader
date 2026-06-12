<script>
  import { createEventDispatcher } from 'svelte';
  import Card from './common/Card.svelte';
  const dispatch = createEventDispatcher();

  export let activeJobLogs = [];
  export let logContainer;
</script>

<Card title="SteamCMD Live Logs" customClass="console-card">
  <button slot="header-right" class="btn-icon-text" on:click={() => dispatch('clear')}>🗑️ Clear Logs</button>
  
  <div class="console-body" bind:this={logContainer}>
    {#if (activeJobLogs || []).length === 0}
      <div class="console-placeholder">Listening for SteamCMD console events...</div>
    {:else}
      {#each activeJobLogs || [] as log (log.id)}
        <div class="line" class:error-line={log.isError}>
          <span class="timestamp">[{log.time}]</span>
          <span class="msg">{log.text}</span>
        </div>
      {/each}
    {/if}
  </div>
</Card>

<style>
  :global(.console-card) {
    flex: 1.5;
    min-height: 120px;
    background: #090d16 !important;
  }
  :global(.console-card .card-body) {
    padding: 0;
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  .console-body {
    flex: 1;
    background: #020617;
    padding: 12px;
    overflow-y: auto;
    font-family: 'Fira Code', 'Cascadia Code', monospace;
    font-size: 11px;
    line-height: 1.6;
    display: flex;
    flex-direction: column;
    gap: 4px;
    color: #a7f3d0;
  }
  .console-placeholder {
    color: #475569;
    text-align: center;
    margin: auto;
    font-family: sans-serif;
  }
  .line {
    word-break: break-all;
    white-space: pre-wrap;
  }
  .line.error-line {
    color: #fca5a5;
  }
  .timestamp {
    color: #4b5563;
    margin-right: 6px;
  }
  .btn-icon-text {
    background: transparent;
    color: #94a3b8;
    font-size: 12px;
    padding: 4px 8px;
    border: none;
    cursor: pointer;
    font-weight: 600;
  }
  .btn-icon-text:hover {
    color: #f8fafc;
    background: rgba(255, 255, 255, 0.05);
    border-radius: 4px;
  }
</style>

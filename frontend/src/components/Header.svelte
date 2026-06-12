<script>
  import { createEventDispatcher } from 'svelte';
  import { WindowMinimise, WindowToggleMaximise, Quit } from '../../wailsjs/runtime/runtime';

  const dispatch = createEventDispatcher();

  function minimize() {
    WindowMinimise();
  }

  function toggleMaximize() {
    WindowToggleMaximise();
  }

  function close() {
    Quit();
  }
</script>

<header class="top-nav">
  <div class="brand">
    <span class="pulse-icon"></span>
    <span class="brand-title">SteamCMD Workshop Downloader</span>
  </div>
  
  <div class="header-actions">
    <button class="btn-secondary settings-btn" on:click={() => dispatch('toggleSettings')}>
      ⚙️ Settings & Auth
    </button>
    
    <div class="window-controls">
      <button class="win-btn" on:click={minimize} title="Minimize">
        <svg viewBox="0 0 10 10"><path d="M0 5h10v1H0z" fill="currentColor"/></svg>
      </button>
      <button class="win-btn" on:click={toggleMaximize} title="Maximize">
        <svg viewBox="0 0 10 10"><path d="M0 0v10h10V0H0zm1 1h8v8H1V1z" fill="currentColor"/></svg>
      </button>
      <button class="win-btn close" on:click={close} title="Close">
        <svg viewBox="0 0 10 10"><path d="M0 0l10 10M10 0L0 10" stroke="currentColor" stroke-width="1.5" fill="none"/></svg>
      </button>
    </div>
  </div>
</header>

<style>
  .top-nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #111827;
    border-bottom: 1px solid #1f2937;
    padding: 0 0 0 16px;
    height: 40px;
    min-height: 40px;
    box-sizing: border-box;
    /* Enable Wails window dragging */
    --wails-draggable: drag;
    user-select: none;
  }
  
  .brand {
    display: flex;
    align-items: center;
    gap: 10px;
  }
  
  .pulse-icon {
    width: 8px;
    height: 8px;
    background: #10b981;
    border-radius: 50%;
    box-shadow: 0 0 8px #10b981;
  }
  
  .brand-title {
    font-size: 13px;
    font-weight: 700;
    color: #e2e8f0;
  }
  
  .header-actions {
    display: flex;
    align-items: center;
    height: 100%;
    /* Disable drag to keep click events working */
    --wails-draggable: no-drag;
  }
  
  :global(.settings-btn) {
    font-size: 11px !important;
    padding: 4px 10px !important;
    height: 26px !important;
    margin-right: 12px;
  }
  
  .window-controls {
    display: flex;
    height: 100%;
  }
  
  .win-btn {
    background: transparent;
    color: #94a3b8;
    border: none;
    cursor: pointer;
    width: 46px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 0 !important;
    padding: 0;
    transition: background 0.15s, color 0.15s;
  }
  
  .win-btn svg {
    width: 10px;
    height: 10px;
  }
  
  .win-btn:hover {
    background: #1f2937;
    color: #f8fafc;
  }
  
  .win-btn.close:hover {
    background: #ef4444;
    color: #ffffff;
  }
</style>

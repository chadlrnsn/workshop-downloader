<script>
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
</script>

<header class="top-nav">
  <div class="brand">
    <span class="pulse-icon"></span>
    <span class="brand-title">SteamCMD Workshop Downloader</span>
  </div>
  <div class="navbar-actions">
    <button class="btn-secondary" on:click={() => dispatch('toggleSettings')}>
      ⚙️ Settings & Auth
    </button>
  </div>
</header>

<style>
  .top-nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #111827;
    border-bottom: 1px solid #1f2937;
    padding: 12px 20px;
    height: 54px;
    box-sizing: border-box;
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
    font-size: 15px;
    font-weight: 700;
    letter-spacing: 0.5px;
    color: #f8fafc;
  }
</style>

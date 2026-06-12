<script>
  import { createEventDispatcher } from 'svelte';
  import Card from './common/Card.svelte';
  const dispatch = createEventDispatcher();

  let workshopUrl = '';
  let manualAppId = '';

  function handleFormSubmit() {
    dispatch('submit', { workshopUrl, manualAppId });
    workshopUrl = '';
    manualAppId = '';
  }
</script>

<Card title="New Download Task">
  <form on:submit|preventDefault={handleFormSubmit} class="download-form">
    <div class="form-group">
      <label for="urlInput">Workshop Item URL or ID</label>
      <input 
        id="urlInput" 
        type="text" 
        placeholder="https://steamcommunity.com/sharedfiles/..." 
        bind:value={workshopUrl} 
        required 
      />
    </div>
    
    <div class="form-group">
      <label for="appIdInput">Manual AppID (if query appid failed)</label>
      <input 
        id="appIdInput" 
        type="text" 
        placeholder="e.g. 281990" 
        bind:value={manualAppId} 
      />
    </div>

    <button type="submit" class="btn-primary btn-block">Add to Queue</button>
  </form>
</Card>

<style>
  .download-form {
    display: flex;
    flex-direction: column;
  }
  .form-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 14px;
  }
  label {
    font-weight: 600;
    font-size: 12px;
    color: #94a3b8;
    margin-bottom: 2px;
  }
  input[type="text"] {
    background: #0f172a;
    border: 1px solid #334155;
    border-radius: 6px;
    padding: 10px 12px;
    color: #f8fafc;
    font-size: 14px;
    width: 100%;
    box-sizing: border-box;
    transition: all 0.2s;
  }
  input[type="text"]:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
  }
</style>

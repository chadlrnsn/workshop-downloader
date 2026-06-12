<script>
  import { createEventDispatcher } from 'svelte';
  import Card from './common/Card.svelte';
  import { _ } from 'svelte-i18n';
  const dispatch = createEventDispatcher();

  let workshopUrl = '';
  let manualAppId = '';

  function handleFormSubmit() {
    dispatch('submit', { workshopUrl, manualAppId });
    workshopUrl = '';
    manualAppId = '';
  }
</script>

<Card title={$_('new_task.title')}>
  <form on:submit|preventDefault={handleFormSubmit} class="download-form">
    <div class="form-group">
      <label for="urlInput">{$_('new_task.url_label')}</label>
      <input 
        id="urlInput" 
        type="text" 
        placeholder={$_('new_task.url_placeholder')} 
        bind:value={workshopUrl} 
        required 
      />
    </div>
    
    <div class="form-group">
      <label for="appIdInput">{$_('new_task.appid_label')}</label>
      <input 
        id="appIdInput" 
        type="text" 
        placeholder={$_('new_task.appid_placeholder')} 
        bind:value={manualAppId} 
      />
    </div>

    <button type="submit" class="btn-primary btn-block">{$_('new_task.submit')}</button>
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

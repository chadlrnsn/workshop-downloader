<script>
  import { createEventDispatcher } from 'svelte';
  import Modal from './common/Modal.svelte';
  import Checkbox from './common/Checkbox.svelte';
  import { _ } from 'svelte-i18n';
  const dispatch = createEventDispatcher();

  export let config;
  export let steamPassword = '';
  export let steamGuardCode = '';
  export let isSavingSettings = false;
  export let isLoggingIn = false;
  export let isResettingAuth = false;
  export let isCheckingSteamCmd = false;
  export let loginStatusMsg = '';

  function handleSave() {
    dispatch('save');
  }

  function handleCheck() {
    dispatch('check');
  }

  function handleLogin() {
    dispatch('login');
  }

  function handleResetAuth() {
    dispatch('resetAuth');
  }

  function handleCancel() {
    dispatch('close');
  }
</script>

<Modal title={$_('settings.title')} sizeClass="settings-modal" on:close={handleCancel}>
  <div class="modal-body modal-scroll">
    <div class="modal-section-title">📂 {$_('settings.section_paths')}</div>
    
    <div class="form-group">
      <label for="steamPath">{$_('settings.steam_path')}</label>
      <div class="input-with-button">
        <input 
          id="steamPath" 
          type="text" 
          bind:value={config.steamCmdPath} 
          placeholder={$_('settings.steam_path_placeholder')} 
          disabled={isLoggingIn}
        />
        <button 
          class="btn-check" 
          type="button" 
          on:click={handleCheck} 
          disabled={isCheckingSteamCmd || isLoggingIn}
        >
          {isCheckingSteamCmd ? $_('settings.checking') : '🔍 ' + $_('settings.check')}
        </button>
      </div>
    </div>
    
    <div class="form-group">
      <label for="outDir">{$_('settings.downloads_directory')}</label>
      <input 
        id="outDir" 
        type="text" 
        bind:value={config.outputDir} 
        placeholder={$_('settings.downloads_directory_placeholder')} 
        disabled={isLoggingIn}
      />
    </div>

    <div class="modal-section-title">🔑 {$_('settings.section_auth')}</div>
    <p class="description-text">
      {$_('settings.auth_description')}
    </p>

    <div class="auth-status-card" class:authorized={config.username !== 'anonymous'}>
      <div class="auth-status-info">
        <span class="status-dot" class:active={config.username !== 'anonymous'}></span>
        <span class="status-text">
          {#if config.username !== 'anonymous'}
            {$_('settings.status_authorized', { values: { username: config.username } })}
          {:else}
            {$_('settings.status_anonymous')}
          {/if}
        </span>
      </div>
      {#if config.username !== 'anonymous'}
        <button 
          class="btn-reset-auth" 
          type="button" 
          on:click={handleResetAuth} 
          disabled={isLoggingIn || isResettingAuth}
        >
          🗑️ {isResettingAuth ? $_('settings.resetting') : $_('settings.reset_auth')}
        </button>
      {/if}
    </div>

    <div class="form-group">
      <label for="username">{$_('settings.username')}</label>
      <input 
        id="username" 
        type="text" 
        bind:value={config.username} 
        placeholder="anonymous" 
        disabled={isLoggingIn || isResettingAuth} 
      />
      <small class="hint">{$_('settings.username_hint')}</small>
    </div>

    {#if config.username !== 'anonymous'}
      <div class="form-group">
        <label for="password">{$_('settings.password')}</label>
        <input 
          id="password" 
          type="password" 
          bind:value={steamPassword} 
          placeholder="••••••••" 
          disabled={isLoggingIn || isResettingAuth} 
        />
        <small class="hint">{$_('settings.password_hint')}</small>
      </div>

      <Checkbox 
        bind:checked={config.rememberPassword} 
        disabled={isLoggingIn || isResettingAuth} 
        label={$_('settings.remember_password')} 
      />

      <div class="form-group">
        <label for="steamGuardCode">{$_('settings.steam_guard_code')}</label>
        <input 
          id="steamGuardCode" 
          type="text" 
          bind:value={steamGuardCode} 
          placeholder="ABCDE" 
          disabled={isLoggingIn || isResettingAuth} 
        />
        <small class="hint">{$_('settings.steam_guard_code_hint')}</small>
      </div>
    {/if}

    {#if isLoggingIn}
      <div class="login-spinner-container">
        <span class="spinner"></span>
        <p class="spinner-msg">{loginStatusMsg}</p>
      </div>
    {:else if config.username !== 'anonymous'}
      <button class="btn-auth btn-block" on:click={handleLogin} disabled={isResettingAuth}>
        🚀 {$_('settings.authenticate')}
      </button>
    {/if}
  </div>

  <div slot="footer" class="modal-footer border-top">
    <button 
      class="btn-primary" 
      on:click={handleSave} 
      disabled={isLoggingIn || isSavingSettings}
    >
      {isSavingSettings ? $_('settings.saving') : $_('settings.save')}
    </button>
    <button 
      class="btn-secondary" 
      on:click={handleCancel} 
      disabled={isLoggingIn}
    >
      {$_('settings.cancel')}
    </button>
  </div>
</Modal>

<style>
  :global(.settings-modal) {
    width: 520px;
    max-height: 85vh;
  }

  label {
    font-weight: 600;
    font-size: 12px;
    color: #94a3b8;
    margin-bottom: 2px;
  }

  input[type="text"], input[type="password"] {
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

  input[type="text"]:focus, input[type="password"]:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
  }

  input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  button {
    font-size: 13px;
    font-weight: 600;
    border-radius: 6px;
    padding: 8px 16px;
    cursor: pointer;
    transition: all 0.2s;
    border: none;
  }

  button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn-primary {
    background: #2563eb;
    color: #ffffff;
  }

  .btn-primary:hover:not(:disabled) {
    background: #1d4ed8;
  }

  .btn-secondary {
    background: #1e293b;
    border: 1px solid #334155;
    color: #cbd5e1;
  }

  .btn-secondary:hover:not(:disabled) {
    background: #334155;
    color: #f8fafc;
  }

  .input-with-button {
    display: flex;
    gap: 8px;
    width: 100%;
  }

  .btn-check {
    background: #334155;
    color: #f8fafc;
    border: 1px solid #475569;
    white-space: nowrap;
    padding: 10px 16px;
  }

  .btn-check:hover:not(:disabled) {
    background: #475569;
  }

  .auth-status-card {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: #0f172a;
    border: 1px solid #1e293b;
    padding: 10px 14px;
    border-radius: 8px;
    margin-bottom: 20px;
    gap: 12px;
  }
  .auth-status-card.authorized {
    border-color: rgba(16, 185, 129, 0.3);
    background: rgba(16, 185, 129, 0.03);
  }
  .auth-status-info {
    display: inline-flex;
    align-items: center;
    gap: 8px;
  }
  .status-dot {
    width: 8px;
    height: 8px;
    background: #64748b;
    border-radius: 50%;
  }
  .status-dot.active {
    background: #10b981;
    box-shadow: 0 0 8px #10b981;
  }
  .status-text {
    font-size: 13px;
    color: #e2e8f0;
    font-weight: 500;
  }
  .btn-reset-auth {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;
    border: 1px solid rgba(239, 68, 68, 0.2);
    font-size: 11px;
    padding: 6px 10px;
    height: auto;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 600;
  }
  .btn-reset-auth:hover:not(:disabled) {
    background: #ef4444;
    color: #ffffff;
  }

  .btn-auth {
    background: #10b981;
    color: #ffffff;
    padding: 11px;
    margin-top: 10px;
  }

  .btn-auth:hover:not(:disabled) {
    background: #059669;
  }

  .btn-block {
    width: 100%;
    display: block;
  }

  .hint {
    font-size: 11px;
    color: #64748b;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 14px;
  }

  .modal-body {
    padding: 20px;
  }

  .modal-scroll {
    overflow-y: auto;
    max-height: calc(85vh - 120px);
  }

  .modal-section-title {
    font-size: 11px;
    font-weight: 800;
    text-transform: uppercase;
    color: #3b82f6;
    letter-spacing: 0.5px;
    margin: 10px 0 14px 0;
    border-bottom: 1px solid #1f2937;
    padding-bottom: 6px;
  }

  .description-text {
    font-size: 12px;
    line-height: 1.5;
    color: #94a3b8;
    margin-bottom: 16px;
  }

  .login-spinner-container {
    display: flex;
    align-items: center;
    gap: 12px;
    background: #0f172a;
    padding: 12px;
    border-radius: 6px;
    border: 1px dashed #334155;
    margin-top: 14px;
  }

  .spinner {
    width: 20px;
    height: 20px;
    border: 2px solid #3b82f6;
    border-top-color: transparent;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  .spinner-msg {
    margin: 0;
    font-size: 12px;
    color: #e2e8f0;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    padding: 14px 20px;
    background: #0f172a;
  }

  .border-top {
    border-top: 1px solid #1f2937;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>

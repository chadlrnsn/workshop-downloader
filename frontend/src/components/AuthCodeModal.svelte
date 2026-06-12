<script>
  import { createEventDispatcher } from 'svelte';
  import Modal from './common/Modal.svelte';
  import { _ } from 'svelte-i18n';
  const dispatch = createEventDispatcher();

  export let authPromptType = '';
  export let steamCode = '';

  function handleSubmit() {
    dispatch('submit');
  }

  function handleCancel() {
    dispatch('cancel');
  }
</script>

<Modal title={$_('auth.title')} sizeClass="auth-code-modal" zHigh={true} on:close={handleCancel}>
  <div class="modal-body">
    <p>
      {$_('auth.prompt', { values: { type: authPromptType === 'email' ? $_('auth.email') : $_('auth.2fa') } })}
    </p>
    <input 
      type="text" 
      bind:value={steamCode} 
      placeholder={$_('auth.placeholder')} 
      class="input-code" 
      autofocus 
    />
  </div>
  
  <div slot="footer" class="modal-footer">
    <button class="btn-primary" on:click={handleSubmit}>{$_('auth.submit')}</button>
    <button class="btn-secondary" on:click={handleCancel}>{$_('auth.cancel')}</button>
  </div>
</Modal>

<style>
  :global(.auth-code-modal) {
    width: 380px;
    border-color: #3b82f6 !important;
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

  .btn-primary {
    background: #2563eb;
    color: #ffffff;
  }

  .btn-primary:hover {
    background: #1d4ed8;
  }

  .btn-secondary {
    background: #1e293b;
    border: 1px solid #334155;
    color: #cbd5e1;
  }

  .btn-secondary:hover {
    background: #334155;
    color: #f8fafc;
  }

  .modal-body {
    padding: 20px;
  }

  .modal-body p {
    margin: 0 0 12px 0;
    font-size: 13px;
    line-height: 1.5;
    color: #cbd5e1;
  }

  .input-code {
    background: #0f172a;
    border: 2px solid #3b82f6;
    border-radius: 6px;
    padding: 10px;
    width: 100%;
    color: #fff;
    font-size: 20px;
    font-weight: 700;
    letter-spacing: 4px;
    text-align: center;
    box-sizing: border-box;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    padding: 14px 20px;
    background: #0f172a;
  }
</style>

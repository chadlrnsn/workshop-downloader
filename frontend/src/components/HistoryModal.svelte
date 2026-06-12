<!-- HistoryModal.svelte -->
<script>
  import { onMount, createEventDispatcher } from 'svelte';
  import Modal from './common/Modal.svelte';
  import { GetHistory, DeleteHistoryItem } from '../../wailsjs/go/main/App';
  import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
  import { _ } from 'svelte-i18n';

  const dispatch = createEventDispatcher();
  let historyItems = [];
  let isLoading = false;

  function handleCancel() {
    dispatch('close');
  }

  async function loadHistory() {
    isLoading = true;
    try {
      const result = await GetHistory();
      historyItems = result || [];
    } catch (err) {
      console.error('Failed to load download history:', err);
    } finally {
      isLoading = false;
    }
  }

  async function handleDelete(id) {
    try {
      await DeleteHistoryItem(id);
      await loadHistory();
    } catch (err) {
      alert(`Failed to delete history item: ${err}`);
    }
  }

  onMount(() => {
    loadHistory();
    EventsOn('history:updated', loadHistory);

    return () => {
      EventsOff('history:updated');
    };
  });
</script>

<Modal title={$_('history.title')} sizeClass="history-modal" on:close={handleCancel}>
  <div class="modal-body modal-scroll">
    {#if isLoading && historyItems.length === 0}
      <div class="empty-state">
        <p>{$_('history.loading')}</p>
      </div>
    {:else}
      {#if historyItems.length === 0}
        <div class="empty-state">
          <p>{$_('history.empty')}</p>
          <small>{$_('history.empty_hint')}</small>
        </div>
      {:else}
        <div class="history-list">
          {#each historyItems as item (item.id)}
            <div class="history-item" class:is-deleted={!item.folderExists}>
              <div class="history-card-top">
                {#if item.previewUrl}
                  <img src={item.previewUrl} alt={item.title} class="history-preview-img" class:grayscale={!item.folderExists} />
                {:else}
                  <div class="history-preview-placeholder" class:grayscale={!item.folderExists}>📂</div>
                {/if}
                <div class="history-details">
                  <div class="history-title-row">
                    <strong class="history-title-text" title={item.title || `Workshop ID: ${item.workshopId}`}>
                      {item.title || `Workshop ID: ${item.workshopId}`}
                    </strong>
                    {#if !item.folderExists}
                      <span class="history-status-badge deleted">{$_('history.status.deleted')}</span>
                    {:else}
                      <span class="history-status-badge success">{$_('history.status.downloaded')}</span>
                    {/if}
                  </div>
                  <div class="history-ids">
                    <span class="text-secondary">AppID: {item.appId} | ID: {item.workshopId}</span>
                  </div>
                </div>
              </div>

              <div class="history-actions">
                {#if item.folderExists}
                  <button class="btn-open-folder" on:click={() => dispatch('openFolder', { appId: item.appId, workshopId: item.workshopId })}>📁 {$_('history.open_folder')}</button>
                {:else}
                  <span class="folder-deleted-text">⚠️ {$_('history.folder_deleted')}</span>
                {/if}
                <button class="btn-delete-text" on:click={() => handleDelete(item.id)}>🗑️ {$_('history.remove')}</button>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    {/if}
  </div>

  <div slot="footer" class="modal-footer border-top">
    <button class="btn-secondary" on:click={handleCancel}>{$_('history.close')}</button>
  </div>
</Modal>

<style>
  :global(.history-modal) {
    width: 600px;
    max-height: 85vh;
  }

  .modal-body {
    padding: 20px;
  }

  .modal-scroll {
    overflow-y: auto;
    max-height: calc(85vh - 120px);
  }

  .empty-state {
    text-align: center;
    color: #64748b;
    padding: 40px 0;
  }
  .empty-state p {
    margin: 0 0 6px 0;
    font-weight: 600;
    font-size: 14px;
  }
  .empty-state small {
    font-size: 12px;
  }
  
  .history-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
  
  .history-item {
    background: #1e293b;
    border: 1px solid #334155;
    border-left: 4px solid #10b981;
    border-radius: 6px;
    padding: 12px 16px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    transition: all 0.2s ease;
  }

  .history-item.is-deleted {
    background: rgba(30, 41, 59, 0.3);
    border: 1px dashed #475569;
    border-left: 4px solid #64748b;
    opacity: 0.7;
  }

  .history-card-top {
    display: flex;
    gap: 12px;
    align-items: center;
    width: 100%;
  }

  .history-preview-img {
    width: 50px;
    height: 50px;
    object-fit: cover;
    border-radius: 4px;
    border: 1px solid #475569;
    background: #0f172a;
    flex-shrink: 0;
    transition: filter 0.2s ease;
  }

  .history-preview-img.grayscale {
    filter: grayscale(100%) contrast(80%);
    opacity: 0.6;
  }

  .history-preview-placeholder {
    width: 50px;
    height: 50px;
    border-radius: 4px;
    border: 1px solid #475569;
    background: #0f172a;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    flex-shrink: 0;
  }
  
  .history-preview-placeholder.grayscale {
    filter: grayscale(100%);
    opacity: 0.5;
  }

  .history-details {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
  }

  .history-title-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
    width: 100%;
  }

  .history-title-text {
    font-size: 13px;
    font-weight: 700;
    color: #f8fafc;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    display: block;
    max-width: 380px;
  }

  .history-item.is-deleted .history-title-text {
    color: #94a3b8;
  }

  .history-status-badge {
    font-size: 10px;
    text-transform: uppercase;
    font-weight: 700;
    padding: 2px 6px;
    border-radius: 4px;
  }

  .history-status-badge.success {
    background: rgba(16, 185, 129, 0.1);
    color: #10b981;
  }

  .history-status-badge.deleted {
    background: rgba(100, 116, 139, 0.15);
    color: #94a3b8;
  }

  .text-secondary {
    font-size: 11px;
    color: #94a3b8;
  }

  .history-actions {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 12px;
  }

  .folder-deleted-text {
    font-size: 12px;
    color: #94a3b8;
    font-weight: 600;
  }

  button {
    font-size: 12px;
    font-weight: 600;
    border-radius: 4px;
    cursor: pointer;
    background: transparent;
    border: none;
    transition: all 0.2s;
  }

  .btn-open-folder {
    color: #3b82f6;
    padding: 4px 8px;
  }

  .btn-open-folder:hover {
    background: rgba(59, 130, 246, 0.15);
  }

  .btn-delete-text {
    color: #94a3b8;
    padding: 4px 8px;
  }

  .btn-delete-text:hover {
    background: rgba(248, 113, 113, 0.1);
    color: #fca5a5;
  }

  .border-top {
    border-top: 1px solid #1f2937;
  }

  .modal-footer {
    display: flex;
    justify-content: flex-end;
    padding: 14px 20px;
    background: #0f172a;
  }
</style>

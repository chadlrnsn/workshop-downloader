<script>
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();

  export let job;
</script>

<div class="job-item status-{job.status}">
  <div class="job-card-top">
    {#if job.previewUrl}
      <img src={job.previewUrl} alt={job.title} class="job-preview-img" />
    {/if}
    <div class="job-details">
      <div class="job-title-row">
        <strong class="job-title-text" title={job.title || `Workshop ID: ${job.workshopId}`}>
          {job.title || `Workshop ID: ${job.workshopId}`}
        </strong>
        <span class="job-status-badge">{job.status}</span>
      </div>
      <div class="job-ids">
        <span class="text-secondary">AppID: {job.appId} | ID: {job.workshopId}</span>
      </div>
    </div>
  </div>

  <div class="progress-section">
    <div class="progress-bar-bg">
      <div class="progress-bar-fill" style="width: {job.progress}%"></div>
    </div>
    <span class="progress-percentage">{job.progress.toFixed(1)}%</span>
  </div>

  {#if job.errorMsg}
    <div class="job-error-msg">⚠️ {job.errorMsg}</div>
  {/if}

  <div class="job-actions">
    {#if job.status === 'queued' || job.status === 'running'}
      <button class="btn-danger-text" on:click={() => dispatch('cancel', job.id)}>Cancel</button>
    {:else}
      {#if job.status === 'success'}
        <button class="btn-open-folder" on:click={() => dispatch('openFolder', { appId: job.appId, workshopId: job.workshopId })}>📁 Open Folder</button>
      {/if}
      {#if job.status === 'failed'}
        <button class="btn-retry-text" on:click={() => dispatch('retry', job.id)}>🔄 Retry</button>
      {/if}
      <button class="btn-delete-text" on:click={() => dispatch('delete', job.id)}>🗑️ Remove</button>
    {/if}
  </div>
</div>

<style>
  .job-item {
    background: #1e293b;
    border: 1px solid #334155;
    border-left: 4px solid #94a3b8;
    border-radius: 6px;
    padding: 12px 16px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .job-item.status-queued {
    border-left-color: #fbbf24;
  }

  .job-item.status-running {
    border-left-color: #3b82f6;
  }

  .job-item.status-success {
    border-left-color: #10b981;
  }

  .job-item.status-failed {
    border-left-color: #ef4444;
  }

  .job-card-top {
    display: flex;
    gap: 12px;
    align-items: center;
    width: 100%;
  }

  .job-preview-img {
    width: 50px;
    height: 50px;
    object-fit: cover;
    border-radius: 4px;
    border: 1px solid #475569;
    background: #0f172a;
    flex-shrink: 0;
  }

  .job-details {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
  }

  .job-title-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
    width: 100%;
  }

  .job-title-text {
    font-size: 13px;
    font-weight: 700;
    color: #f8fafc;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    display: block;
    max-width: 320px;
  }

  .job-status-badge {
    font-size: 10px;
    text-transform: uppercase;
    font-weight: 700;
    background: rgba(255, 255, 255, 0.05);
    padding: 2px 6px;
    border-radius: 4px;
    color: #cbd5e1;
  }

  .text-secondary {
    font-size: 11px;
    color: #94a3b8;
  }

  .progress-section {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .progress-bar-bg {
    flex: 1;
    background: #0f172a;
    height: 8px;
    border-radius: 4px;
    overflow: hidden;
  }

  .progress-bar-fill {
    background: #3b82f6;
    height: 100%;
    border-radius: 4px;
    transition: width 0.3s ease;
  }

  .job-item.status-success .progress-bar-fill {
    background: #10b981;
  }

  .progress-percentage {
    font-size: 12px;
    font-weight: 700;
    color: #f1f5f9;
    min-width: 36px;
    text-align: right;
  }

  .job-error-msg {
    background: rgba(239, 68, 68, 0.1);
    color: #fca5a5;
    padding: 6px 10px;
    border-radius: 4px;
    font-size: 12px;
  }

  .job-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
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

  .btn-danger-text {
    color: #f87171;
    padding: 4px 8px;
  }

  .btn-danger-text:hover {
    background: rgba(239, 68, 68, 0.15);
  }

  .btn-retry-text {
    color: #10b981;
    padding: 4px 8px;
  }

  .btn-retry-text:hover {
    background: rgba(16, 185, 129, 0.15);
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
</style>

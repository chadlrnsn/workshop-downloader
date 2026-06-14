<script>
  import Card from './common/Card.svelte';
  import { _ } from 'svelte-i18n';
  export let jobs = [];
</script>

<Card customClass="stats-card">
  <div class="status-summary">
    <div class="stat-item total">
      <span class="stat-val total-val">{(jobs || []).length}</span>
      <span class="stat-lbl">{$_('stats.total')}</span>
    </div>
    <div class="stat-item">
      <span class="stat-val queued-val">{(jobs || []).filter(j => j && j.status === 'queued').length}</span>
      <span class="stat-lbl">{$_('stats.queued')}</span>
    </div>
    <div class="stat-item">
      <span class="stat-val active-val">{(jobs || []).filter(j => j && j.status === 'running').length}</span>
      <span class="stat-lbl">{$_('stats.active')}</span>
    </div>
    <div class="stat-item">
      <span class="stat-val success-val">{(jobs || []).filter(j => j && j.status === 'success').length}</span>
      <span class="stat-lbl">{$_('stats.success')}</span>
    </div>
    <div class="stat-item">
      <span class="stat-val failed-val">{(jobs || []).filter(j => j && j.status === 'failed').length}</span>
      <span class="stat-lbl">{$_('stats.failed')}</span>
    </div>
  </div>
</Card>

<style>
  .status-summary {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 8px 12px;
    text-align: center;
  }
  .stat-item {
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .stat-item.total {
    grid-column: span 2;
    border-bottom: 1px solid #1f2937;
    padding-bottom: 6px;
    margin-bottom: 2px;
  }
  .stat-val {
    font-size: 15px;
    font-weight: 800;
  }
  .total-val {
    font-size: 18px;
    color: #3b82f6;
  }
  .queued-val {
    color: #94a3b8;
  }
  .active-val {
    color: #fbbf24;
  }
  .success-val {
    color: #10b981;
  }
  .failed-val {
    color: #f87171;
  }
  .stat-lbl {
    font-size: 9px;
    font-weight: 700;
    color: #64748b;
    margin-top: 1px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    width: 100%;
  }
</style>
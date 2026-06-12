<script>
  import Card from './common/Card.svelte';
  export let jobs = [];
</script>

<Card customClass="stats-card">
  <div class="status-summary">
    <div class="stat-item">
      <span class="stat-val">{(jobs || []).length}</span>
      <span class="stat-lbl">Total Tasks</span>
    </div>
    <div class="stat-item">
      <span class="stat-val">{(jobs || []).filter(j => j && j.status === 'running').length}</span>
      <span class="stat-lbl">Active</span>
    </div>
    <div class="stat-item">
      <span class="stat-val">{(jobs || []).filter(j => j && j.status === 'success').length}</span>
      <span class="stat-lbl">Success</span>
    </div>
  </div>
</Card>

<style>
  .status-summary {
    display: flex;
    justify-content: space-between;
    text-align: center;
  }
  .stat-item {
    display: flex;
    flex-direction: column;
  }
  .stat-val {
    font-size: 20px;
    font-weight: 800;
    color: #e2e8f0;
  }
  .stat-lbl {
    font-size: 11px;
    color: #64748b;
    margin-top: 2px;
  }
</style>

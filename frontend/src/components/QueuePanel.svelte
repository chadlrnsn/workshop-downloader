<script>
  import Card from './common/Card.svelte';
  import JobCard from './JobCard.svelte';
  import { _ } from 'svelte-i18n';
  export let jobs = [];
</script>

<Card title={$_('queue.title')} customClass="queue-card">
  <span slot="header-right" class="badge">{$_('queue.items', { values: { count: (jobs || []).length } })}</span>
  
  <div class="scrollable-content">
    {#if !(jobs && jobs.length)}
      <div class="empty-state">
        <p>{$_('queue.empty')}</p>
        <small>{$_('queue.empty_hint')}</small>
      </div>
    {:else}
      <div class="queue-list">
        {#each jobs || [] as job (job.id)}
          <JobCard 
            {job} 
            on:cancel 
            on:retry 
            on:delete 
            on:openFolder
          />
        {/each}
      </div>
    {/if}
  </div>
</Card>

<style>
  :global(.queue-card) {
    flex: 2;
    min-height: 120px;
    overflow: hidden;
  }
  :global(.queue-card .card-body) {
    padding: 0;
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  .badge {
    background: #1e293b;
    color: #94a3b8;
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 11px;
    font-weight: 600;
  }
  .scrollable-content {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
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
  .queue-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }
</style>

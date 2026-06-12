<script>
  import Card from './common/Card.svelte';
  import { _ } from 'svelte-i18n';
  export let username = 'anonymous';
</script>

<Card title={$_('session.title')} customClass="details-card">
  <div class="session-info">
    <div class="session-row">
      <span class="lbl font-semibold">{$_('session.active_user')}</span>
      <span class="val highlight-text">{username}</span>
    </div>
    <div class="session-row">
      <span class="lbl">{$_('session.engine_status')}</span>
      <span class="val text-green">{$_('session.ready')}</span>
    </div>
  </div>
</Card>

<style>
  .session-info {
    display: flex;
    flex-direction: column;
    gap: 8px;
    font-size: 13px;
  }
  .session-row {
    display: flex;
    justify-content: space-between;
  }
  .lbl {
    color: #94a3b8;
  }
  .font-semibold {
    font-weight: 600;
  }
  .highlight-text {
    color: #3b82f6;
    font-weight: 700;
  }
  .text-green {
    color: #10b981;
    font-weight: 700;
  }
</style>

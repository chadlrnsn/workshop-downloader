<script>
  export let title = '';
  export let customClass = '';
</script>

<div class="card {customClass}">
  {#if title}
    <div class="card-header">
      <h3>{title}</h3>
      <slot name="header-right"></slot>
    </div>
  {/if}
  <div class="card-body">
    <slot></slot>
  </div>
</div>

<style>
  .card {
    background: #111827;
    border: 1px solid #1f2937;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  .card-header {
    background: #151e2e;
    padding: 12px 16px;
    border-bottom: 1px solid #1f2937;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .card-header h3 {
    margin: 0;
    font-size: 13px;
    font-weight: 700;
    color: #f1f5f9;
  }
  .card-body {
    padding: 16px;
  }
</style>

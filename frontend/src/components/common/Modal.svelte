<script>
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();

  export let title = '';
  export let sizeClass = ''; 
  export let zHigh = false;

  function handleClose() {
    dispatch('close');
  }
</script>

<div class="modal-backdrop" class:z-high={zHigh} on:click|self={handleClose}>
  <div class="modal-panel {sizeClass}">
    <div class="modal-header">
      <h3>{title}</h3>
    </div>
    
    <div class="modal-body-wrapper">
      <slot></slot>
    </div>

    <!-- Optional footer slot -->
    <slot name="footer"></slot>
  </div>
</div>

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(4, 6, 15, 0.88);
    backdrop-filter: blur(5px);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .z-high {
    z-index: 1100;
  }

  .modal-panel {
    background: #111827;
    border: 1px solid #1f2937;
    border-radius: 8px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.6);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .modal-header {
    background: #151e2e;
    padding: 12px 16px;
    border-bottom: 1px solid #1f2937;
  }

  .modal-header h3 {
    margin: 0;
    font-size: 14px;
    font-weight: 700;
    color: #f8fafc;
  }

  .modal-body-wrapper {
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
</style>

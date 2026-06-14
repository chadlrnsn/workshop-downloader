<script>
  export let checked = false;
  export let disabled = false;
  export let label = '';
</script>

<div class="form-group checkbox-group">
  <label class="checkbox-container">
    <input 
      type="checkbox" 
      bind:checked={checked} 
      {disabled} 
    />
    <span class="checkbox-label">{label}</span>
  </label>
</div>

<style>
  .checkbox-group {
    display: flex;
    flex-direction: row !important;
    align-items: center;
    margin-bottom: 14px;
  }
  .checkbox-container {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    user-select: none;
  }
  .checkbox-container input[type="checkbox"] {
    width: 16px;
    height: 16px;
    accent-color: #3b82f6;
    cursor: pointer;
  }
  .checkbox-label {
    font-size: 13px;
    color: #cbd5e1;
    font-weight: 500;
  }
</style>

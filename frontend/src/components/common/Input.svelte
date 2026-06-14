<script>
  export let id = '';
  export let type = 'text';
  export let value = '';
  export let placeholder = '';
  export let disabled = false;
  export let min = undefined;
  export let max = undefined;
  export let required = false;
  export let autofocus = false;
</script>

<!-- svelte-ignore a11y_autofocus -->
<input
  {id}
  {type}
  bind:value={value}
  {placeholder}
  {disabled}
  {min}
  {max}
  {required}
  {autofocus}
  {...$$restProps}
/>

<style>
  input {
    background: #0f172a;
    border: 1px solid #334155;
    border-radius: 6px;
    padding: 10px 12px;
    color: #f8fafc;
    font-size: 14px;
    width: 100%;
    box-sizing: border-box;
    transition: all 0.2s;
  }

  input:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
  }

  input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
</style>

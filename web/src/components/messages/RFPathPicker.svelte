<script>
  import { Select } from '@chrissnell/chonky-ui';
  import CallsignAutocomplete from './CallsignAutocomplete.svelte';
  import { RF_PATH_PRESETS, composeRFPath, splitRFPath } from './rf-path.js';
  /** @type {{value?: string | null, valid?: boolean, disabled?: boolean, onChange?: (value: string | null) => void}} */
  let { value = null, valid = $bindable(true), disabled = false, onChange } = $props();
  let mode = $state('inherit');
  let digi = $state('');
  let custom = $state('');
  let marker = $state('');
  let lastValue;
  let initialized = false;
  $effect(() => {
    if (initialized && value === lastValue) return;
    initialized = true;
    lastValue = value;
    const parsed = splitRFPath(value);
    mode = parsed.mode; digi = parsed.digi; custom = parsed.custom; marker = parsed.marker;
  });
  const result = $derived(composeRFPath(mode, digi, custom, marker));
  $effect(() => { valid = !result.error; });
  function update() {
    const next = composeRFPath(mode, digi, custom, marker);
    if (!next.error) { lastValue = next.value; onChange?.(next.value); }
  }
  function updateCustom(event) {
    custom = event.currentTarget.value;
    const parts = custom.toUpperCase().split(',').map(p => p.trim());
    const enteredMarker = parts.find(p => p === 'NOGATE' || p === 'RFONLY');
    if (enteredMarker) {
      marker = enteredMarker;
      custom = parts.filter(p => p !== 'NOGATE' && p !== 'RFONLY').join(',');
      event.currentTarget.value = custom;
    }
    update();
  }
  const markers = [
    { value: '', label: 'No iGate restriction marker' },
    { value: 'NOGATE', label: 'NOGATE (request no iGating)' },
    { value: 'RFONLY', label: 'RFONLY (request no iGating)' },
  ];
</script>
<div class="rf-path-picker">
  <span class="field-label">APRS packet path (RF)</span>
  <Select value={mode} options={RF_PATH_PRESETS} {disabled} onValueChange={v => { mode = String(v); update(); }} aria-label="APRS RF packet path" />
  {#if mode === 'digi'}
    <CallsignAutocomplete bind:value={digi} excludeBots={true} {disabled} autofocus={false} placeholder="Digipeater callsign, e.g. DB0ABC-2" onCommit={v => { digi = v; update(); }} onValueChange={v => { digi = v; update(); }} />
    <p class="hint">Choose a known station or enter a callsign. Check that it operates as a digipeater.</p>
  {:else if mode === 'custom'}
    <input aria-label="Custom RF packet path" value={custom} {disabled} oninput={updateCustom} placeholder="WIDE2-1" />
  {/if}
  {#if mode !== 'inherit'}
    <Select value={marker} options={markers} {disabled} onValueChange={v => { marker = String(v); update(); }} aria-label="Internet forwarding marker" />
  {/if}
  <p class="hint">Path inside the radio packet: requests digipeater hops. NOGATE/RFONLY ask iGates not to forward it to APRS-IS; compliance is not guaranteed.</p>
  <p class="preview">Packet path: <code>{result.error ? '—' : result.value == null ? 'conversation / global default' : result.value || '(direct)'}</code></p>
  {#if result.error}<p class="error" role="status">{result.error}</p>{/if}
</div>
<style>
  .rf-path-picker { display: flex; flex-direction: column; gap: 5px; min-width: 0; }
  .field-label { font-size: 11px; font-weight: 700; color: var(--color-text-muted); }
  .hint, .preview, .error { margin: 0; font-size: 11px; line-height: 1.35; overflow-wrap: anywhere; }
  .hint { color: var(--color-text-muted); }
  .error { color: var(--color-danger); }
  code { font-family: var(--font-mono); }
  input { width: 100%; padding: 6px 8px; border: 1px solid var(--color-border); background: var(--color-bg); color: var(--color-text); border-radius: var(--radius); }
</style>

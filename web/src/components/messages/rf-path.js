export const RF_PATH_PRESETS = [
  { value: 'inherit', label: 'Use default path' },
  { value: 'direct', label: 'Direct (no digipeater)' },
  { value: 'WIDE1-1', label: 'WIDE1-1 (one hop)' },
  { value: 'WIDE2-1', label: 'WIDE2-1 (one hop)' },
  { value: 'WIDE1-1,WIDE2-1', label: 'WIDE1-1,WIDE2-1 (up to two hops)' },
  { value: 'digi', label: 'Specific digipeater…' },
  { value: 'custom', label: 'Custom path…' },
];
const address = /^[A-Z0-9]{1,6}(?:-(?:[0-9]|1[0-5]))?$/;
export function validateRFPath(path) {
  if (!path) return '';
  const parts = path.split(',');
  if (parts.length > 8) return 'Use no more than eight path addresses.';
  if (parts.some(p => !address.test(p) || ['TCPIP', 'TCPXX'].includes(p) || /^QA.$/.test(p))) {
    return 'Enter valid AX.25 callsigns or aliases, separated by commas (SSID 0–15).';
  }
  return '';
}
export function composeRFPath(mode, digi, custom, marker) {
  if (mode === 'inherit') return { value: null, error: '' };
  let base = mode === 'direct' ? '' : mode === 'digi' ? digi : mode === 'custom' ? custom : mode;
  base = (base || '').toUpperCase().split(',').map(p => p.trim()).join(',');
  if (mode === 'digi' && (!base || base.includes(',') || !address.test(base))) {
    return { value: null, error: 'Enter one digipeater callsign with an optional SSID (0–15).' };
  }
  const error = validateRFPath(base);
  if (error) return { value: null, error };
  const hops = base ? base.split(',') : [];
  // The separate selector owns Internet-gating markers, even for custom paths.
  const parts = hops.filter(p => p !== 'RFONLY' && p !== 'NOGATE');
  if (marker) parts.push(marker);
  const value = parts.join(',');
  return { value, error: validateRFPath(value) };
}
export function splitRFPath(value) {
  if (value == null) return { mode: 'inherit', digi: '', custom: '', marker: '' };
  const parts = value.toUpperCase().split(',').map(p => p.trim()).filter(Boolean);
  const marker = parts.find(p => p === 'RFONLY' || p === 'NOGATE') || '';
  const base = parts.filter(p => p !== 'RFONLY' && p !== 'NOGATE').join(',');
  const preset = RF_PATH_PRESETS.find(p => p.value === base);
  const mode = !base ? 'direct' : preset ? base : !base.includes(',') ? 'digi' : 'custom';
  return { mode, digi: mode === 'digi' ? base : '', custom: base, marker };
}
// Save transport and packet path before enqueueing the first message. Preserve
// unrelated conversation settings, including the operator's ACK preference.
export async function prepareComposeRouting(kind, key, options, api) {
  if (options.sendPath == null && !options.pathChanged) return;
  const prefs = await api.get(kind, key);
  await api.put(kind, key, {
    send_path: options.sendPath ?? prefs.send_path ?? '',
    wait_for_ack: prefs.wait_for_ack ?? true,
    rf_path: options.pathChanged ? options.rfPath : prefs.rf_path ?? null,
  });
}

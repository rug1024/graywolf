// Compact, always-visible receive provenance below incoming chat bubbles.
// Source wins over channel metadata: APRS-IS must never look like RF.
export function receiveRouteLabel(msg, channels = []) {
  if (msg?.direction !== 'in') return '';
  if (msg.source === 'is') return 'RX · APRS-IS';
  if (msg.source === 'sim') return 'RX · Simulation';
  if (msg.source !== 'rf') return 'RX · Source unknown';
  const channel = msg.channel;
  if (!Number.isInteger(channel) || channel <= 0) return 'RX · RF';
  const name = channels.find(c => c.id === channel)?.name?.trim();
  return name ? `RX · Ch ${channel} · ${name}` : `RX · Ch ${channel}`;
}

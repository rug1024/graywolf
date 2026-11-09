// Infer a DM reply route from persisted conversation history. Explicit UI
// choices are handled by the caller. Never borrow a tactical participant's
// radio channel for the entire group.
export function replyChannel(messages, kind, key) {
  if (kind !== 'dm') return 0;
  const rows = messages.filter(m => m.thread_kind === kind && m.thread_key === key
    && Number.isInteger(m.channel) && m.channel > 0 && !m.is_ack && !m.is_rej);
  const newest = list => list.reduce((best, m) => {
    const time = row => Date.parse(row.received_at || row.sent_at || row.created_at || '') || 0;
    return !best || time(m) > time(best)
      || (time(m) === time(best) && (m.id || 0) > (best.id || 0)) ? m : best;
  }, null);
  // Incoming RF determines the reply path. Without RF history, keep the
  // channel used to start this conversation (including after a reload).
  return (newest(rows.filter(m => m.direction === 'in' && m.source === 'rf'))
    || newest(rows.filter(m => m.direction === 'out')))?.channel || 0;
}

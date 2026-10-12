// Describe live routing eligibility, never promise that every packet passes filters.
export function channelGateStatus(channel, igate, fresh = true) {
  const detail = 'Channel routing status. APRS rules, packet filters and source exclusions still apply.';
  return ['RF → IS', 'IS → RF'].map((direction, index) => {
    const result = (label, tone) => ({ direction, label, tone, detail });
    if (!fresh || typeof channel.enabled !== 'boolean') return result('Status unavailable', 'unknown');
    if (!channel.enabled || !igate) return result('Disabled', 'disabled');
    // RF receive fan-out submits APRS UI frames from every enabled channel.
    // The legacy rf_channel/gate_rf_to_is config fields do not select this route.
    if (index === 1) {
      if (typeof igate.is_to_rf_enabled !== 'boolean') return result('Status unavailable', 'unknown');
      if (!igate.is_to_rf_enabled || !igate.tx_channel || igate.tx_channel !== channel.id || channel.mode === 'packet') return result('Disabled', 'disabled');
    }
    if (!igate.connected) return result('IS disconnected', 'pending');
    // The backend simulation flag suppresses RF->IS only; IS->RF is not simulated.
    if (index === 0 && igate.simulation_mode) return result('Simulation', 'simulation');
    return result('Active · filtered', 'active');
  });
}

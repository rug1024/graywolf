import test from 'node:test';
import assert from 'node:assert/strict';
import { channelGateStatus as gate } from './igateDashboard.js';
const ch = { id: 4, enabled: true, mode: 'aprs' };
const ig = { rf_to_is_enabled: true, connected: true, tx_channel: 1, is_to_rf_enabled: true };
test('uplink spans enabled channels; downlink uses resolved runtime channel', () => {
  assert.deepEqual(gate(ch, ig).map(x=>x.tone), ['active','disabled']);
  assert.equal(gate(ch, {...ig, tx_channel:4})[1].tone, 'active');
});
test('master downlink off, packet channel and disabled channel are respected', () => {
  assert.equal(gate(ch, {...ig, tx_channel:4, is_to_rf_enabled:false})[1].tone, 'disabled');
  assert.equal(gate({...ch, mode:'packet'}, {...ig, tx_channel:4})[1].tone, 'disabled');
  assert.deepEqual(gate({...ch, enabled:false},ig).map(x=>x.tone), ['disabled','disabled']);
});
test('simulation applies only to uplink and offline is never green', () => {
  assert.deepEqual(gate(ch,{...ig,tx_channel:4,simulation_mode:true}).map(x=>x.tone),['simulation','active']);
  assert.deepEqual(gate(ch,{...ig,tx_channel:4,connected:false}).map(x=>x.tone),['pending','pending']);
});
test('stale and missing status never imply disabled or active routing', () => {
  assert.deepEqual(gate(ch,ig,false).map(x=>x.tone), ['unknown','unknown']);
  assert.equal(gate({id:4},ig)[0].tone,'unknown');
  assert.deepEqual(gate(ch,null).map(x=>x.tone), ['disabled','disabled']);
});

test('RF→IS master off takes precedence over connection and simulation without changing downlink', () => {
  for (const connected of [true, false]) {
    for (const simulation_mode of [true, false]) {
      const state = gate(ch, {...ig, tx_channel: 4, rf_to_is_enabled: false, connected, simulation_mode});
      assert.equal(state[0].tone, 'disabled');
      assert.equal(state[1].tone, connected ? 'active' : 'pending');
    }
  }
});
test('missing RF→IS master state in an older API never appears active', () => {
  const {rf_to_is_enabled, ...oldStatus} = ig;
  assert.equal(gate(ch, oldStatus)[0].tone, 'unknown');
  assert.equal(gate(ch, {...oldStatus, rf_to_is_enabled: null})[0].tone, 'unknown');
});
test('live RF→IS disable and re-enable is reflected without restarting the channel', () => {
  for (const enabled of [true, false, true]) {
    assert.equal(gate(ch, {...ig, rf_to_is_enabled: enabled})[0].tone, enabled ? 'active' : 'disabled');
  }
});

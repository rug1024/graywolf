import test from 'node:test';
import assert from 'node:assert/strict';
import { channelGateStatus as gate } from './igateDashboard.js';
const ch = { id: 4, enabled: true, mode: 'aprs' };
const ig = { connected: true, tx_channel: 1, is_to_rf_enabled: true };
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

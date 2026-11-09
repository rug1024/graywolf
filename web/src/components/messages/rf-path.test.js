import test from 'node:test';
import assert from 'node:assert/strict';
import { composeRFPath, splitRFPath, prepareComposeRouting } from './rf-path.js';

test('direct, inherited and no-gate paths remain distinct', () => {
  assert.deepEqual(composeRFPath('inherit', '', '', ''), { value: null, error: '' });
  assert.deepEqual(composeRFPath('direct', '', '', ''), { value: '', error: '' });
  assert.equal(composeRFPath('direct', '', '', 'RFONLY').value, 'RFONLY');
  assert.equal(composeRFPath('WIDE2-1', '', '', 'NOGATE').value, 'WIDE2-1,NOGATE');
});
test('specific digi normalizes the callsign and appends the selected marker', () => {
  assert.equal(composeRFPath('digi', ' db0abc-2 ', '', 'NOGATE').value, 'DB0ABC-2,NOGATE');
  for (const digi of ['', 'DB0ABC-16', 'DB0ABC*', 'A,B', 'TCPIP']) {
    assert.ok(composeRFPath('digi', digi, '', '').error, digi);
  }
});
test('custom path validation rejects overflow, empty hops and Internet markers', () => {
  for (const custom of ['A,B,C,D,E,F,G,H,I', 'WIDE2-1,', 'qAR', 'TCPXX', 'A\r\nB']) {
    assert.ok(composeRFPath('custom', '', custom, '').error, custom);
  }
  assert.equal(composeRFPath('custom', '', ' WIDE1-1, WIDE2-1 ', 'RFONLY').value, 'WIDE1-1,WIDE2-1,RFONLY');
});
test('stored path can be reopened without losing direct, digi or marker choices', () => {
  for (const path of [null, '', 'RFONLY', 'NOGATE', 'WIDE1-1,WIDE2-1', 'DB0ABC-2,NOGATE', 'DB0ABC-2,DB0DEF-3,RFONLY']) {
    const p = splitRFPath(path);
    assert.equal(composeRFPath(p.mode, p.digi, p.custom, p.marker).value, path);
  }
});
test('compose saves routing before sending and preserves the ACK preference', async () => {
  const events = [];
  await prepareComposeRouting('dm', 'W1ABC', {sendPath:'rf_only',pathChanged:true,rfPath:'RFONLY'}, {
    get: async () => { events.push('get'); return {send_path:'is_only',wait_for_ack:false,rf_path:'WIDE2-1'}; },
    put: async (_kind,_key,prefs) => { events.push('put'); assert.deepEqual(prefs,{send_path:'rf_only',wait_for_ack:false,rf_path:'RFONLY'}); },
  });
  events.push('send');
  assert.deepEqual(events,['get','put','send']);
});
test('unchanged routing performs no preference write and failed saves abort preparation', async () => {
  await prepareComposeRouting('dm','W1ABC',{sendPath:null,pathChanged:false},{get:()=>assert.fail('unnecessary read')});
  await assert.rejects(prepareComposeRouting('dm','W1ABC',{sendPath:'rf_only',pathChanged:false},{get:async()=>({}),put:async()=>{throw Error('offline');}}),/offline/);
});

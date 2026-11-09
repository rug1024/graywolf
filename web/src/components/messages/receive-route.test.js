import test from 'node:test';
import assert from 'node:assert/strict';
import { receiveRouteLabel } from './receive-route.js';

test('incoming RF shows the received channel, not the configured transmit channel', () => {
  const channels = [{ id: 1, name: 'VHF APRS' }, { id: 2, name: 'LoRa APRS' }];
  assert.equal(receiveRouteLabel({ direction: 'in', source: 'rf', channel: 2 }, channels), 'RX · Ch 2 · LoRa APRS');
  assert.equal(receiveRouteLabel({ direction: 'in', source: 'rf', channel: 1 }, channels), 'RX · Ch 1 · VHF APRS');
});
test('missing or removed channel names retain the received channel number', () => {
  assert.equal(receiveRouteLabel({ direction: 'in', source: 'rf', channel: 6 }), 'RX · Ch 6');
  assert.equal(receiveRouteLabel({ direction: 'in', source: 'rf', channel: 6 }, [{ id: 6, name: '  ' }]), 'RX · Ch 6');
});
test('old RF messages without valid channel metadata do not invent a channel', () => {
  for (const channel of [undefined, null, 0, -1, 1.5]) {
    assert.equal(receiveRouteLabel({ direction: 'in', source: 'rf', channel }), 'RX · RF');
  }
});
test('Internet source takes precedence over any stale RF channel field', () => {
  assert.equal(receiveRouteLabel({ direction: 'in', source: 'is', channel: 2 }, [{ id: 2, name: 'LoRa APRS' }]), 'RX · APRS-IS');
  assert.equal(receiveRouteLabel({ direction: 'in', source: 'sim' }), 'RX · Simulation');
  assert.equal(receiveRouteLabel({ direction: 'in' }), 'RX · Source unknown');
});
test('outgoing and missing rows never get a receive label', () => {
  assert.equal(receiveRouteLabel({ direction: 'out', source: 'rf', channel: 2 }), '');
  assert.equal(receiveRouteLabel(null), '');
});

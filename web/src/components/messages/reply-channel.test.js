import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { replyChannel } from './reply-channel.js';

const row = (overrides = {}) => ({
  id: 1, direction: 'in', source: 'rf', channel: 2,
  thread_kind: 'dm', thread_key: 'DA6PG-5',
  received_at: '2026-10-20T10:00:00Z', ...overrides,
});
const route = rows => replyChannel(rows, 'dm', 'DA6PG-5');

test('LoRa incoming message selects channel 2 rather than the VHF default', () => {
  assert.equal(route([row()]), 2);
});
test('newest RF receive wins independently of API ordering and outgoing sends', () => {
  const latest = row({ id: 3, channel: 4, received_at: '2026-10-20T11:00:00Z' });
  const outgoing = row({ id: 4, direction: 'out', channel: 1, received_at: '2026-10-20T12:00:00Z' });
  for (const rows of [[row(), latest, outgoing], [outgoing, latest, row()]]) {
    assert.equal(route(rows), 4);
  }
});
test('APRS-IS, ACK, REJ and unknown channels cannot replace the RF reply route', () => {
  const invalid = [
    row({ source: 'is', channel: 1 }), row({ is_ack: true, channel: 1 }),
    row({ is_rej: true, channel: 1 }), row({ channel: 0 }),
    row({ channel: null }), row({ channel: -1 }),
  ].map(m => ({ ...m, received_at: '2026-10-20T12:00:00Z' }));
  assert.equal(route([row(), ...invalid]), 2);
  assert.equal(route(invalid), 0);
});
test('other conversations and tactical participants cannot leak a channel', () => {
  assert.equal(route([row({ thread_key: 'OTHER' })]), 0);
  assert.equal(replyChannel([row({ thread_kind: 'tactical', thread_key: 'NET' })], 'tactical', 'NET'), 0);
});
test('conversation started on an explicit channel retains it without RF history', () => {
  assert.equal(route([row({ direction: 'out', channel: 4 })]), 4);
  assert.equal(route([]), 0);
});
test('equal receive timestamps use persisted message order', () => {
  assert.equal(route([row({ id: 9, channel: 3 }), row()]), 3);
});

test('chat wiring forwards the chosen channel and freezes it for multipart sends', () => {
  const source = name => readFileSync(new URL(name, import.meta.url), 'utf8');
  const composer = source('./ComposeBar.svelte');
  assert.match(composer, /const batchChannel = channel/);
  assert.match(composer, /onSend\?\.\(tagged, target, batchChannel\)/);
  assert.match(composer, /onSend\?\.\(body, target, batchChannel\)/);
  assert.match(composer, /aria-label="Reply transmit channel"/);
  const thread = source('./MessageThread.svelte');
  assert.match(thread, /channel=\{transmitChannel\}/);
  assert.match(thread, /onCompose\?\.\(body, to, channel\)/);
  assert.match(thread, /store\.messageById\.values\(\)/);
  assert.match(source('../../routes/Messages.svelte'), /async function handleThreadSend\(text, to, channel\)\s*\{\s*await optimisticSend\(text, to, channel\)/);
});

import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import postcss from 'postcss';

const css = postcss.parse(readFileSync(new URL('../../themes/signal-typography.css', import.meta.url), 'utf8'));
const rules = css.nodes.filter(n => n.type === 'rule');
const declarations = rule => Object.fromEntries(rule.nodes.filter(n => n.type === 'decl').map(n => [n.prop, n.value]));

for (const theme of ['signal', 'signal-night', 'neon-signal']) {
  test(`${theme} uses bundled Inter for controls and preserves packet/callsign monospace`, () => {
    const themed = rules.filter(r => r.selector.includes(`[data-theme="${theme}"]`));
    assert.ok(themed.some(r => declarations(r)['--font-sans'] === "'Inter', system-ui, sans-serif"));
    const controls = themed.find(r => r.selector.includes('button, input, select, textarea'));
    assert.equal(declarations(controls)['font-family'], 'var(--font-sans)');
    const packet = themed.find(r => r.selector.includes('.pkt-srcdst, .pkt-raw, .beacon-callsign'));
    assert.equal(declarations(packet)['font-family'], 'var(--font-mono)');
    assert.ok(themed.every(r => !('--font-mono' in declarations(r))));
  });
}

test('Inter font remains bundled and available without an external font service', () => {
  const face = css.nodes.find(n => n.type === 'atrule' && n.name === 'font-face');
  const src = declarations(face).src;
  assert.equal(src, "url('/fonts/inter/InterVariable.woff2') format('woff2')");
  const bytes = readFileSync(new URL('../../public/fonts/inter/InterVariable.woff2', import.meta.url));
  assert.equal(bytes.subarray(0, 4).toString(), 'wOF2');
  assert.ok(bytes.length > 10000);
});

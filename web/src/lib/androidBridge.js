// Single source of truth for the per-launch bearer token injected by
// the Android Service via WebView.addJavascriptInterface(...,
// "GraywolfWebInterface"). Returns null on desktop builds so the
// bearer wiring becomes a no-op there.
//
// The token is cached at first read because the injected JS bridge
// object is stable for the WebView's lifetime; calling the JNI
// getBearerToken on every fetch would cross the JS<->Java boundary
// needlessly.

let cached; // sentinel: undefined = not read; null = absent; string = token

export function getBearerToken() {
  if (cached !== undefined) return cached;
  try {
    const v = globalThis.GraywolfWebInterface?.getBearerToken?.();
    cached = (typeof v === 'string' && v.length > 0) ? v : null;
  } catch {
    cached = null;
  }
  return cached;
}

export function getKeepRunningInBackground() {
  try {
    const bridge = globalThis.GraywolfWebInterface;
    if (!bridge?.getKeepRunningInBackground) return true;
    return Boolean(bridge.getKeepRunningInBackground());
  } catch {
    return true;
  }
}

export function setKeepRunningInBackground(enabled) {
  try {
    const bridge = globalThis.GraywolfWebInterface;
    if (!bridge?.setKeepRunningInBackground) return false;
    bridge.setKeepRunningInBackground(Boolean(enabled));
    return true;
  } catch {
    return false;
  }
}

// Test-only: reset the cache between unit tests. Not part of the
// public surface; do not call from app code.
export function _resetForTests() {
  cached = undefined;
}

# Personal build continuity

Future personal builds must start from the last verified integration state and retain the fixes below. If a fix must be removed or behavior is uncertain, ask the owner before changing it. Do not replace this integration with a vanilla upstream build.

## Android baseline

Baseline: c5347959509a432aee09d964aa11d4efb3b59179 (successful Android run 37971077343).
Preserve Bluetooth Classic and BLE KISS, GPS MSL and review corrections, audio-only microphone permission, Android system-bar insets, notification/launcher reopening, background mode, KISS channel counters and live dashboard status, compact dashboard spacing, configurable RFONLY/NOGATE message path and acknowledgements, and existing APRS server selection.

Keep the existing debug signing certificate for upgrade compatibility. Do not silently substitute a new key.

## Shared iGate fixes approved for Android and Linux

Source: b42f220193971baddbb74587607f23b108de9507.
- Preserve repeated digipeater H-bit (`*`) in RF-to-IS serialization.
- Use qAR/qAO based on actual eligible RF message reachability.
- Truncate outgoing RF-to-IS information at the first CR/LF while preserving the original RF packet and preceding binary bytes (PR #10).
- Persist and live-reload the configurable RF-to-IS source exclusion list; preserve exact SSID semantics. No hardcoded DB0SDA exclusion.
- Retain the approved DAPNET/POCSAG filtering at receive, iGate and transmit boundaries.

Linux build label remains Graywolf/0.14.14-mod. Do not automatically rename Android.

## Build gate

Compare against the previous integration, run frontend tests and production web build, the affected Go regression packages, and Android BLE/GPS unit tests. Check restored signing-key provenance and a successful APK build. Hardware verification remains a separate user test.

## Upstream synchronization 2026-10-12

Merged upstream main 0acbb7f9faf19a57b34b7858fb3230b5489fbd1f (VERSION remains 0.14.14).
Retain the runtime per-channel iGate direction/status display and Signal/Signal Night colors.
Android migration 30 remains the deployed BLE type rename; upstream serial/USB TNC repair runs as migration 31 so existing fork databases receive it.
Serial/USB defaults follow upstream; BLE retains its existing governor-TX choice. OSM, metric units, custom map server, radar defaults and disabled update checks remain intact.

Signal and Signal Night use bundled Inter for UI text, readable channel status sizes and tabular counters; preserve monospace callsigns/raw packets and offline font availability.

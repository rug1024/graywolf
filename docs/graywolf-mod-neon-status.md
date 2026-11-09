# Graywolf mod: Neon Signal typography and RF→IS dashboard status

Baseline: `bc7fc153b6b48168ee5413e6f960fb76527850c8` on `graywolf-mod`, PR #16.

## Changes

- Neon Signal joins Signal and Signal Night in the shared, locally bundled Inter typography. Packet source/destination callsigns, raw APRS packets and beacon callsigns keep `--font-mono`. No palette, font asset or terminal font replacement.
- iGate runtime status exposes `rf_to_is_enabled`, read from the same atomic master switch used by the RF→IS gating pipeline. `/api/status` projects the flag without deriving it from APRS-IS connectivity.
- Dashboard uplink eligibility requires that flag. Off overrides both connected and simulation states; a missing flag is unavailable, never active. Downlink, channel selection, exclusions and APRS filters are unchanged.
- Swagger JSON/YAML and the generated TypeScript client are regenerated together.

## Regression coverage

- Live master-switch enable/disable/re-enable; downlink stays independent.
- Both true and false master-switch values survive status DTO JSON serialization.
- Frontend: off with each combination of connectivity/simulation; old API without the flag; live re-enable.
- All three Signal themes use Inter for controls while preserving callsign/raw-packet monospace, with a valid local WOFF2 asset.

Local validation: 554 frontend tests passed; web production build passed; Go tests passed for APRS, app wiring, iGate, TX governor, web API and messages. Full race/vet and schema drift validation continue in local/CI checks. Final CI results and downloadable artifacts are recorded on PR #16 and in the delivery report.

## Baseline preservation

The implementation delta is limited to shared Signal typography, iGate dashboard status and related tests/schema. No Android/BLE/KISS/GPS code, APRS encoders, gating rules, RF master-switch wiring, theme palettes, database migrations or device settings are removed or altered. Builds use the existing Android and native ARM64 workflows. Both Pi executables must have identical `v0.14.14-mod-<commit>` version output before artifact upload.

No changes to `main`; no branches deleted; no installation or deployment.

## Remaining device validation

Hardware BLE connection, live KISS radio traffic, GPS altitude and visual typography on the user's phone/Pi require a later device test. CI and static preservation checks do not substitute for those tests. Existing Android host-JVM tests outside the BLE/GPS subset have framework-stub limitations documented in the Android workflow.

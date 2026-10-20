# Android update artifacts for graywolf-mod

When delivering an Android APK from this fork as an update to the user's installed Graywolf app:

- Use the existing Android signing certificate. Its SHA-256 fingerprint is:
  `5e8aef2ea71a52a84e723974317c6af39549fb9167d6ff8639cb23d7eae299eb`.
- Before linking an APK, verify its signature with `apksigner verify --verbose --print-certs` and compare that fingerprint exactly. A successful compile alone is insufficient.
- The regular `graywolf-debug-apk` build artifact is an intermediate artifact. Do not present it as a compatible update without verifying its certificate.
- The original debug keystore is retained in the Actions cache `graywolf-debug-keystore-v1` under PR #2's cache scope. Identically named caches in other PR scopes can contain different keys.
- The workflow `.github/workflows/android-legacy-update-sign.yml` on `ble-kiss-tbeam` (PR #2) re-signs a pinned, successfully tested graywolf-mod build with the original key and checks the certificate before uploading `*-compatible-update`.
- Never generate a replacement key when the original cache cannot be restored. Stop delivery and explain the blocker instead.
- Never recommend uninstalling the existing app to work around a signature mismatch; preserve its settings and message history.
- Keep private key material out of repository files, logs, and downloadable artifacts.

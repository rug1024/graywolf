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

# Build cadence

Collect proposed changes first. Only start Pi or Android builds when the user explicitly requests a build. Pushing graywolf-mod updates triggers builds automatically, so hold those pushes until a build is authorized. Read-only investigation and preparing changes remain allowed.

# Installation instructions with build deliveries

Every future build delivery must include a plain-text installation guide as a downloadable `.txt` file. Include the exact artifact names and the Raspberry Pi commands for checksum verification, stopping the service, backing up binaries and data into separate `bin/` and `data/` directories, installing both executables, restarting, and checking the version/service. The installed paths are `/usr/bin/graywolf`, `/usr/bin/graywolf-modem`, and `/var/lib/graywolf`; the service is `graywolf.service`. Include Android update steps when an APK is delivered, preserving the installed app and its settings. Do not trigger a build just to add or update the installation guide. For the 2026-10-25 delivery the user explicitly requested the commands in chat only.

# APRS log corrections included in the next authorized build (2026-10-29)

- Move DETAILS back to the right; the user explicitly withdrew the proposal to move it farther left.
- Shorten `Channel:` to `Ch:` and `Distance:` to `Dist:`.
- Round displayed distances to whole kilometres, without decimal places.
- Keep the compact row spacing and the full-width raw packet line.
- Inter is used across the mod UI and palettes; the user authorized a new Pi/Android build on 2026-10-29, including upstream PR #612 and the collected changes.

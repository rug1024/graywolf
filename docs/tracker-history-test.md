# Android BLE tracker history: hardware acceptance test

Use a Guzman tracker with the authenticated BLE mailbox firmware and a
Graywolf Android version with the matching history reader. Set Graywolf and
the tracker to the same APRS callsign, including the SSID. Use an independent
receiver or the Raspberry Pi log to observe RF traffic.

1. Disconnect Android from the tracker. From another station, send a numbered
   direct message over LoRa using the `RFONLY` path. Confirm that the tracker
   displays the message and sends its normal APRS ACK. Wait for the sender to
   receive the ACK before reconnecting.
2. Reconnect Android with the normal BLE PIN. The message should appear once
   in its conversation with a blue “Vom Tracker” label. Its reception time
   should be unknown. Verify that importing it sends no additional RF ACK.
3. Disconnect and reconnect again. The imported message should still appear
   once, and it should remain readable on the tracker.
4. With BLE connected, send a fresh direct message. Verify that it appears as
   a live message without the history label. Verify normal KISS transmission
   from the app with KISS TX enabled.
5. Connect to a tracker with older firmware that has no history service.
   Verify that live KISS RX and TX still work.

The firmware stores up to 32 direct messages received while Android is
disconnected. It does not backfill messages from before the firmware update.
The mailbox protocol has no reception timestamp or original APRS message
number. Re-importing a mailbox record is deduplicated. An RF retry delivered
live after reconnect can still appear separately.

Record the Android build ID and tracker firmware revision before testing.

Test hardware: T-Beam 1W.

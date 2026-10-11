package aprs

import (
	"bytes"

	"github.com/chrissnell/graywolf/pkg/ax25"
)

// IsDAPNETTraffic identifies DAPNET/POCSAG payloads by their protocol
// labels, independent of the originating callsign or RF/IS transport.
// The user-specific deny policy is applied only at the application edges;
// the ordinary APRS decoder remains protocol-neutral.
//
// This is intentionally a case-insensitive content rule. It may also
// match a human-authored APRS message discussing DAPNET/POCSAG.
func IsDAPNETTraffic(info []byte) bool {
	if len(info) == 0 {
		return false
	}
	upper := bytes.ToUpper(info)
	return bytes.Contains(upper, []byte("DAPNET")) ||
		bytes.Contains(upper, []byte("POCSAG"))
}

// IsDAPNETFrame applies the rule only to AX.25 UI (APRS) frames.
// The complete info field is examined, so third-party ('}') packets
// containing a DAPNET object are blocked too.
func IsDAPNETFrame(frame *ax25.Frame) bool {
	return frame != nil && frame.IsUI() && IsDAPNETTraffic(frame.Info)
}

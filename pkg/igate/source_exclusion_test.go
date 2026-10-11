package igate

import (
	"testing"
	"time"

	"github.com/chrissnell/graywolf/pkg/aprs"
)

func TestDB0SDAExclusionScope(t *testing.T) {
	tests := []struct {
		call string
		want bool
	}{
		{"DB0SDA", true},
		{"db0sda", true},
		{"DB0SDA-2", true},
		{" Db0SdA-10 ", true},
		{"DB0SDAB", false},
		{"DA6PG-7", false},
		{"ON4MAX-2", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.call, func(t *testing.T) {
			if got := blockedRFToISSource(tt.call); got != tt.want {
				t.Errorf("blockedRFToISSource(%q) = %t, want %t", tt.call, got, tt.want)
			}
		})
	}
}

func TestDB0SDANeverBecomesIGateRecipient(t *testing.T) {
	ig := mustIgate(t)
	for _, call := range []string{"DB0SDA", "db0sda-2"} {
		pkt := &aprs.DecodedAPRSPacket{
			Source: call,
			Dest: "APRS",
			Type: aprs.PacketPosition,
			Raw: buildRawFrame(t, "DB0SDA", "APRS", nil, "!5046.57NS00604.17E#PHG2140Aachen"),
		}
		if got := ig.gateRFToIS(pkt); got != GateReasonSourceBlock {
			t.Errorf("%s: outcome = %q, want %q", call, got, GateReasonSourceBlock)
		}
	}
	// Other packets still reach the existing policy pipeline (offline
	// here because the fake iGate has no live APRS-IS connection).
	ordinary := &aprs.DecodedAPRSPacket{
		Source: "ON4MAX-2",
		Dest: "APRS",
		Type: aprs.PacketPosition,
		Raw: buildRawFrame(t, "ON4MAX-2", "APRS", nil, "!5046.57NS00604.17E#PHG2140Aachen"),
	}
	if got := ig.gateRFToIS(ordinary); got != GateReasonOffline {
		t.Errorf("other station outcome = %q, want %q", got, GateReasonOffline)
	}
	if ig.heard.HeardWithin("DB0SDA", 30*time.Minute) {
		t.Error("DB0SDA should not be admitted to local IS→RF heard-direct recipient tracker")
	}
	if !ig.heard.HeardWithin("ON4MAX-2", 30*time.Minute) {
		t.Error("normal station was not recorded in the heard-direct tracker")
	}
}

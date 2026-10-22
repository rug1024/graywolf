package aprs

import (
	"testing"

	"github.com/chrissnell/graywolf/pkg/ax25"
)

func TestIsDAPNETTraffic(t *testing.T) {
	tests := []struct {
		name string
		info string
		want bool
	}{
		{"pocsag object", ";P-DB0SDA *111403z5050.00NP00605.00E#DAPNET POCSAG Transmitter(ONLINE)", true},
		{"mixed case", ";P-X      *111403z5050.00NP00605.00E#DaPnEt", true},
		{"pocsag only", ">POCSAG gateway", true},
		{"direct APRS message mentioning DAPNET is deliberately blocked", ":DA6PG-1  :Do you still use dapnet?", true},
		{"unrelated binary control bytes", ">status\\x1coperational", false},
		{"third party DAPNET", "}DB0SDA>APRS,TCPIP*:;P-XYZ    *111403z5050.00NP00605.00E#DAPNET", true},
		{"third party ordinary", "}ON3HJM-12>APLT00,TCPIP:!5054.61N/00551.97E> _Bat.: 4.11V", false},
		{"same source ordinary", "!5046.57NS00604.17E#PHG2140 Aachen", false},
		{"ordinary status", ">APRS receiver operational", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDAPNETTraffic([]byte(tt.info)); got != tt.want {
				t.Fatalf("IsDAPNETTraffic(%q) = %v, want %v", tt.info, got, tt.want)
			}
		})
	}
}

func TestIsDAPNETFrameNil(t *testing.T) {
	if IsDAPNETFrame(nil) {
		t.Fatal("nil frame should never be blocked")
	}
	if IsDAPNETFrame(&ax25.Frame{}) {
		t.Fatal("empty non-UI frame should never be blocked")
	}
}

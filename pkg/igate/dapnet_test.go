package igate

import (
	"errors"
	"testing"

	"github.com/chrissnell/graywolf/pkg/aprs"
	"github.com/chrissnell/graywolf/pkg/ax25"
)

func TestDAPNETIngressNotRecordedOrForwarded(t *testing.T) {
	calls := 0
	ig, err := New(Config{
		Server: "127.0.0.1:1",
		StationCallsign: "DA6PG-1",
		IsRxHook: func(_ *aprs.DecodedAPRSPacket, _ string) { calls++ },
	})
	if err != nil { t.Fatal(err) }

	ig.handleISLine("db0sda>APRS,TCPIP*,qAC,EIGHTH:;P-TEST   *111403z5050.00NP00605.00E#DAPNET POCSAG Transmitter(ONLINE)")
	if calls != 0 { t.Fatalf("DAPNET packet reached IS receive hook: %d", calls) }

	ig.handleISLine("DB0SDA>APRS,TCPIP*:!5046.57NS00604.17E#PHG2140Aachen")
	if calls != 1 { t.Fatalf("ordinary DB0SDA position should reach hook: %d", calls) }
}

func TestDAPNETRejectedForRFToISAndDirectIS(t *testing.T) {
	ig, err := New(Config{Server: "127.0.0.1:1", StationCallsign: "DA6PG-1"})
	if err != nil { t.Fatal(err) }

	src, err := ax25.ParseAddress("DB0SDA")
	if err != nil { t.Fatal(err) }
	dest, err := ax25.ParseAddress("APRS")
	if err != nil { t.Fatal(err) }
	frame, err := ax25.NewUIFrame(src, dest, nil, []byte(";P-TEST   *111403z5050.00NP00605.00E#DAPNET POCSAG Transmitter"))
	if err != nil { t.Fatal(err) }
	pkt, err := aprs.Parse(frame)
	if err != nil { t.Fatal(err) }
	if got := ig.gateRFToIS(pkt); got != GateReasonDAPNET {
		t.Fatalf("gateRFToIS = %q, want %q", got, GateReasonDAPNET)
	}
	if err := ig.SendLine("DB0SDA>APRS,TCPIP*:;P-TEST   *111403z5050.00NP00605.00E#DAPNET"); err == nil {
		t.Fatal("direct APRS-IS upload of DAPNET content should be blocked")
	}
	if err := ig.SendLine("DB0SDA>APRS,TCPIP*:!5046.57NS00604.17E#PHG2140Aachen"); !errors.Is(err, errors.New("igate: not connected")) && err == nil {
		t.Fatal("ordinary line should only fail due to offline connection")
	}
}

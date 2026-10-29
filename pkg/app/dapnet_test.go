package app

import (
	"context"
	"testing"

	"github.com/chrissnell/graywolf/pkg/app/ingress"
	pb "github.com/chrissnell/graywolf/pkg/ipcproto"
	"github.com/chrissnell/graywolf/pkg/packetlog"
)

func TestDAPNETRFDiscardedBeforeAnyConsumer(t *testing.T) {
	// An otherwise unwired App intentionally panics if the RF packet
	// reaches KISS, AGW, digi or other consumers. The DAPNET guard
	// must short-circuit immediately, leaving packetlog untouched.
	a := &App{plog: packetlog.New(packetlog.Config{})}
	for _, info := range []string{
		";P-TEST   *111403z5050.00NP00605.00E#DAPNET POCSAG Transmitter",
		"}DB0SDA>APRS,TCPIP*:;P-TEST   *111403z5050.00NP00605.00E#DAPNET",
		">pocsag transmitter status",
	} {
		frame := buildUIFrame(t, "DB0SDA", info, nil)
		raw, err := frame.Encode()
		if err != nil { t.Fatal(err) }
		a.dispatchRxFrame(context.Background(), rxFanoutItem{
			rf: &pb.ReceivedFrame{Channel: 1, Data: raw},
			src: ingress.Modem(),
		}, nil)
	}
	if got := len(a.plog.Query(packetlog.Filter{Channel: -1})); got != 0 {
		t.Fatalf("DAPNET RF packets logged: %d", got)
	}
}

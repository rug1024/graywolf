package txgovernor

import (
	"context"
	"errors"
	"testing"
)

func TestDAPNETBlockedAtAllTransmitSources(t *testing.T) {
	for _, src := range []string{"digipeater", "igate", "kiss", "agw", "beacon"} {
		t.Run(src, func(t *testing.T) {
			g := newDeterministic(&captureSender{})
			err := g.Submit(context.Background(), 1, makeFrame(t, ";P-TEST   *111403z5050.00NP00605.00E#DAPNET POCSAG Transmitter"), SubmitSource{Kind: src})
			if !errors.Is(err, ErrDAPNET) {
				t.Fatalf("Submit = %v, want ErrDAPNET", err)
			}
			if got := g.Stats().Enqueued; got != 0 {
				t.Fatalf("blocked frame entered the TX queue: %d", got)
			}
		})
	}
}

func TestNormalAPRSNotBlockedAtTransmitGovernor(t *testing.T) {
	g := newDeterministic(&captureSender{})
	if err := g.Submit(context.Background(), 1, makeFrame(t, "!5046.57NS00604.17E#PHG2140Aachen"), SubmitSource{Kind: "digipeater"}); err != nil {
		t.Fatalf("ordinary APRS frame must not be blocked: %v", err)
	}
	if got := g.Stats().Enqueued; got != 1 {
		t.Fatalf("normal frame Enqueued = %d, want 1", got)
	}
}

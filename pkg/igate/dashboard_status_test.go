package igate

import (
	"log/slog"
	"testing"
)

func TestDashboardStatusTracksLiveRouting(t *testing.T) {
	ig := &Igate{logger: slog.Default()}
	ig.SetTxChannel(4)
	ig.Reconfigure("", nil, &stubGovernor{})
	st := ig.Status()
	if st.TxChannel != 4 || !st.IsToRFEnabled {
		t.Fatalf("missing live route: %+v", st)
	}
	ig.SetTxChannel(7)
	ig.Reconfigure("", nil, nil)
	st = ig.Status()
	if st.TxChannel != 7 || st.IsToRFEnabled {
		t.Fatalf("stale route after reconfigure: %+v", st)
	}
}

func TestDashboardStatusTracksRFToISMasterSwitch(t *testing.T) {
	ig := mustIgate(t)
	ig.SetTxChannel(4)
	ig.Reconfigure("", nil, &stubGovernor{})
	for _, enabled := range []bool{true, false, true} {
		ig.SetRFToISEnabled(enabled)
		st := ig.Status()
		if st.RFToISEnabled != enabled {
			t.Fatalf("RF→IS status=%v, want %v", st.RFToISEnabled, enabled)
		}
		if st.TxChannel != 4 || !st.IsToRFEnabled {
			t.Fatalf("uplink toggle changed downlink: %+v", st)
		}
	}
}

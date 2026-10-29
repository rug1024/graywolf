package igate

import (
 "testing"
 "log/slog"
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

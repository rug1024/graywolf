package messages

import (
	"context"
	"testing"
	"time"

	"github.com/chrissnell/graywolf/pkg/configstore"
)

func TestImportTrackerHistory(t *testing.T) {
	s, cs := newTestStore(t)
	ctx := context.Background()
	received := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)
	m := TrackerHistoryMessage{"tracker-epoch", 1, "DA6PG-1", "DA6PG-7", "Offline test", &received}
	for _, want := range []bool{true, false} {
		got, err := s.ImportTrackerHistory(ctx, "DA6PG-7", m)
		if err != nil || got != want {
			t.Fatalf("import = %v, %v; want %v", got, err, want)
		}
	}
	var rows []configstore.Message
	if err := cs.DB().Find(&rows).Error; err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("stored %d rows; want one", len(rows))
	}
	r := rows[0]
	if r.Source != "hist" || r.Direction != "in" || r.ThreadKey != m.From || r.ReceivedAt == nil || !r.ReceivedAt.Equal(received) || r.NextRetryAt != nil || r.IsAck || r.MsgID != "" {
		t.Fatalf("unexpected history record: %+v", r)
	}
	m.To = "DA6PG-5"
	if _, err := s.ImportTrackerHistory(ctx, "DA6PG-7", m); err == nil {
		t.Fatal("accepted a message addressed to a different station")
	}
}

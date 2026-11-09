package platformsvc

import (
	"context"
	"encoding/binary"
	"errors"
	"testing"
)

func historyTestRecord(seq uint32) []byte {
	r := make([]byte, 14)
	r[0], r[1], r[2] = 1, 1, 42
	binary.LittleEndian.PutUint32(r[10:], seq)
	for _, s := range []string{"DA6PG-1", "DA6PG-7", "Offline test"} {
		r = append(r, byte(len(s)))
		r = append(r, s...)
	}
	return r
}

func TestTrackerMailboxAdvancesOnlyAfterSave(t *testing.T) {
	ctx := context.Background()
	failure := errors.New("database unavailable")
	for _, saveFailure := range []bool{true, false} {
		var reads, saves int
		err := importTrackerMailbox(ctx, func(_ context.Context, cursor []byte) ([]byte, error) {
			reads++
			if reads == 1 {
				if binary.LittleEndian.Uint64(cursor) != 0 || binary.LittleEndian.Uint32(cursor[8:]) != 0 {
					t.Fatal("reconnect did not start with a clean cursor")
				}
				return historyTestRecord(5), nil
			}
			if saves != 1 || binary.LittleEndian.Uint64(cursor) != 42 || binary.LittleEndian.Uint32(cursor[8:]) != 5 {
				t.Fatal("cursor did not follow the committed record")
			}
			end := historyTestRecord(5)[:14]
			end[1] = 0
			return end, nil
		}, func(_ context.Context, r TrackerHistoryRecord) error {
			saves++
			if r.Sequence != 5 || r.To != "DA6PG-7" || r.From != "DA6PG-1" || r.Text != "Offline test" {
				t.Fatalf("decoded unexpected record: %+v", r)
			}
			if saveFailure {
				return failure
			}
			return nil
		})
		if saveFailure {
			if !errors.Is(err, failure) || reads != 1 {
				t.Fatalf("advanced after failure: reads=%d err=%v", reads, err)
			}
		} else if err != nil || reads != 2 {
			t.Fatalf("successful replay did not reach end: reads=%d err=%v", reads, err)
		}
	}
}

func TestTrackerHistoryRejectsTruncatedRecords(t *testing.T) {
	record := historyTestRecord(5)
	for n := range record {
		if _, _, err := decodeTrackerHistory(record[:n]); err == nil {
			t.Fatalf("accepted truncated record at %d bytes", n)
		}
	}
}

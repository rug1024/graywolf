package platformsvc

import (
	"context"
	"encoding/binary"
	"fmt"
	"unicode/utf8"
)

type TrackerHistoryRecord struct {
	Generation     uint64
	Sequence       uint32
	From, To, Text string
}

func decodeTrackerHistory(data []byte) (TrackerHistoryRecord, bool, error) {
	var r TrackerHistoryRecord
	if len(data) < 14 || data[0] != 1 || data[1] > 1 {
		return r, false, fmt.Errorf("invalid tracker history header")
	}
	r.Generation = binary.LittleEndian.Uint64(data[2:10])
	r.Sequence = binary.LittleEndian.Uint32(data[10:14])
	if r.Generation == 0 {
		return r, false, fmt.Errorf("invalid tracker generation")
	}
	if data[1] == 0 && len(data) == 14 {
		return r, false, nil
	}
	if data[1] != 1 || r.Sequence == 0 {
		return r, false, fmt.Errorf("invalid tracker history status")
	}
	data = data[14:]
	for i, field := range []*string{&r.From, &r.To, &r.Text} {
		max := []int{9, 9, 200}[i]
		if len(data) == 0 || data[0] == 0 || int(data[0]) > max || int(data[0]) >= len(data) {
			return r, false, fmt.Errorf("invalid tracker history field")
		}
		n := int(data[0])
		if !utf8.Valid(data[1 : 1+n]) {
			return r, false, fmt.Errorf("invalid tracker history encoding")
		}
		*field = string(data[1 : 1+n])
		data = data[1+n:]
	}
	if len(data) != 0 {
		return r, false, fmt.Errorf("unexpected tracker history trailer")
	}
	return r, true, nil
}

// Replay the bounded mailbox on reconnect. The persistent import key is the
// acknowledgement: if saving fails, do not advance the cursor.
func importTrackerMailbox(ctx context.Context, read func(context.Context, []byte) ([]byte, error), save func(context.Context, TrackerHistoryRecord) error) error {
	cursor := make([]byte, 12)
	var generation uint64
	var previous uint32
	for i := 0; i < 32; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		data, err := read(ctx, cursor)
		if err != nil || data == nil {
			return err
		}
		r, found, err := decodeTrackerHistory(data)
		if err != nil || !found {
			return err
		}
		if i > 0 && (r.Generation != generation || r.Sequence <= previous) {
			return fmt.Errorf("tracker mailbox changed or cursor did not advance")
		}
		if err := save(ctx, r); err != nil {
			return err
		}
		generation, previous = r.Generation, r.Sequence
		binary.LittleEndian.PutUint64(cursor, generation)
		binary.LittleEndian.PutUint32(cursor[8:], previous)
	}
	return nil
}

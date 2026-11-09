//go:build android

package platformsvc

import (
	"context"
	"errors"
	"log/slog"
	"time"

	pb "github.com/chrissnell/graywolf/pkg/platformproto"
)

// SetTrackerHistoryImporter is called during app wiring, before BLE opens.
func (c *clientImpl) SetTrackerHistoryImporter(importer func(context.Context, string, TrackerHistoryRecord) error) {
	c.historyImport = importer
}

func (c *clientImpl) syncBLEHistory(handle uint32, address string, closed <-chan struct{}) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case <-ctx.Done():
		case <-closed:
			cancel()
		case <-c.closeCh:
			cancel()
		}
	}()
	err := importTrackerMailbox(ctx, func(ctx context.Context, cursor []byte) ([]byte, error) {
		return c.readBLEHistory(ctx, handle, cursor)
	}, func(ctx context.Context, record TrackerHistoryRecord) error {
		return c.historyImport(ctx, address, record)
	})
	if err != nil && ctx.Err() == nil {
		slog.Warn("BLE tracker history sync stopped", "error", err)
	}
}

func (c *clientImpl) readBLEHistory(ctx context.Context, handle uint32, cursor []byte) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	id := c.nextSerialHandle()
	ch := make(chan *pb.BleHistoryResponse, 1)
	c.subsMu.Lock()
	if c.historyReplies == nil {
		c.historyReplies = make(map[uint32]chan *pb.BleHistoryResponse)
	}
	c.historyReplies[id] = ch
	c.subsMu.Unlock()
	defer func() {
		c.subsMu.Lock()
		delete(c.historyReplies, id)
		c.subsMu.Unlock()
	}()
	if err := c.send(&pb.PlatformMessage{Body: &pb.PlatformMessage_BleHistoryRequest{
		BleHistoryRequest: &pb.BleHistoryRequest{RequestId: id, Handle: handle, Cursor: cursor},
	}}); err != nil {
		return nil, err
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case reply := <-ch:
		if reply.GetFailed() {
			return nil, errors.New("tracker history GATT read failed")
		}
		if reply.GetUnsupported() {
			return nil, nil
		}
		if len(reply.GetRecord()) == 0 {
			return nil, errors.New("empty tracker history response")
		}
		return reply.GetRecord(), nil
	}
}

package messages

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/chrissnell/graywolf/pkg/configstore"
	"gorm.io/gorm/clause"
)

// TrackerHistoryMessage is a historical direct message from a trusted tracker.
// TrackerID identifies the mailbox generation; Sequence is never reused within it.
type TrackerHistoryMessage struct {
	TrackerID  string
	Sequence   uint64
	From       string
	To         string
	Text       string
	ReceivedAt *time.Time
}

// ImportTrackerHistory persists history without routing a received packet:
// no transmit, ACK, retry or remote action can result from this call.
// A false result with no error means this record was already imported.
func (s *Store) ImportTrackerHistory(ctx context.Context, ownCall string, msg TrackerHistoryMessage) (bool, error) {
	row, err := s.insertTrackerHistory(ctx, ownCall, msg)
	return row != nil, err
}

func (s *Store) insertTrackerHistory(ctx context.Context, ownCall string, msg TrackerHistoryMessage) (*configstore.Message, error) {
	ownCall = strings.ToUpper(strings.TrimSpace(ownCall))
	from := strings.ToUpper(strings.TrimSpace(msg.From))
	to := strings.ToUpper(strings.TrimSpace(msg.To))
	if !rfPacketAddress.MatchString(ownCall) || !rfPacketAddress.MatchString(from) || to != ownCall || from == ownCall {
		return nil, fmt.Errorf("tracker history: expected a direct message addressed to the local station")
	}
	if strings.TrimSpace(msg.TrackerID) == "" || len(msg.TrackerID) > 128 || msg.Sequence == 0 {
		return nil, fmt.Errorf("tracker history: invalid record identity")
	}
	if strings.TrimSpace(msg.Text) == "" || len(msg.Text) > 200 || !utf8.ValidString(msg.Text) {
		return nil, fmt.Errorf("tracker history: invalid message text")
	}

	key := fmt.Sprintf("%s:%d", msg.TrackerID, msg.Sequence)
	row := configstore.Message{
		TrackerImportKey: &key,
		Direction:        "in", OurCall: ownCall, FromCall: from, ToCall: to,
		PeerCall: from, ThreadKind: ThreadKindDM, ThreadKey: from,
		Text: msg.Text, Source: "hist", AckState: AckStateNone, Unread: true,
	}
	if msg.ReceivedAt != nil && !msg.ReceivedAt.IsZero() {
		received := msg.ReceivedAt.UTC()
		row.ReceivedAt = &received
		row.CreatedAt = received
	}

	result := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "tracker_import_key"}}, DoNothing: true,
	}).Create(&row)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, result.Error
	}
	return &row, nil
}

// ImportTrackerHistory publishes a display event only after a new row is saved.
func (s *Service) ImportTrackerHistory(ctx context.Context, ownCall string, msg TrackerHistoryMessage) error {
	row, err := s.cfg.Store.insertTrackerHistory(ctx, ownCall, msg)
	if err != nil || row == nil {
		return err
	}
	s.hub.Publish(Event{Type: EventMessageReceived, MessageID: row.ID,
		ThreadKind: row.ThreadKind, ThreadKey: row.ThreadKey, Timestamp: time.Now().UTC()})
	return nil
}

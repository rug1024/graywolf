package messages

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/chrissnell/graywolf/pkg/configstore"
)

func TestValidateRFPacketPath(t *testing.T) {
	for _, path := range []string{"", "RFONLY", "NOGATE", "WIDE1-1,WIDE2-1", "DB0ABC-2,NOGATE", " wide2-1 , rfonly "} {
		if err := ValidateRFPacketPath(path); err != nil {
			t.Errorf("valid %q: %v", path, err)
		}
	}
	for _, path := range []string{"TCPIP", "TCPXX", "qAR", "DB0ABC*", "DB0ABC-16", "DB0ABC-2junk", "DB0ABC-1\nB", "WIDE2-1,", ",NOGATE", "A,B,C,D,E,F,G,H,I", "A\r\nB"} {
		if err := ValidateRFPacketPath(path); err == nil {
			t.Errorf("accepted invalid path %q", path)
		}
	}
}

func TestServiceConversationRFPathSurvivesReloadAndPreferenceChanges(t *testing.T) {
	for _, selected := range []string{"", "RFONLY", "DB0ABC-2,NOGATE"} {
		t.Run(selected, func(t *testing.T) {
			svc, rig, _, cleanup := buildService(t)
			defer cleanup()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if err := svc.Start(ctx); err != nil {
				t.Fatal(err)
			}
			defer svc.Stop()
			path := selected
			if err := rig.cs.UpsertConversationPrefs(ctx, &configstore.ConversationPrefs{
				ThreadKind: ThreadKindDM, ThreadKey: "W1ABC", SendPath: FallbackPolicyRFOnly, WaitForAck: true, RFPath: &path,
			}); err != nil {
				t.Fatal(err)
			}
			row, err := svc.SendMessage(ctx, SendMessageRequest{OurCall: "N0CALL", To: "W1ABC", Text: "hello", Channel: 2})
			if err != nil {
				t.Fatal(err)
			}
			// Await initial dispatch and retry enrollment before closing this
			// shared-memory SQLite test harness or changing conversation prefs.
			deadline := time.Now().Add(3 * time.Second)
			for {
				current, readErr := rig.store.GetByID(ctx, row.ID)
				if readErr != nil {
					t.Fatal(readErr)
				}
				if current.NextRetryAt != nil {
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("initial dispatch did not finish")
				}
				time.Sleep(5 * time.Millisecond)
			}
			persisted, err := rig.store.GetByID(ctx, row.ID)
			if err != nil {
				t.Fatal(err)
			}
			if persisted.RFPath == nil || *persisted.RFPath != selected || persisted.SendPath != FallbackPolicyRFOnly {
				t.Fatalf("lost routing snapshot: %+v", persisted)
			}
			changed := "WIDE2-2"
			if err := rig.cs.UpsertConversationPrefs(ctx, &configstore.ConversationPrefs{
				ThreadKind: ThreadKindDM, ThreadKey: "W1ABC", WaitForAck: true, RFPath: &changed,
			}); err != nil {
				t.Fatal(err)
			}
			frame, err := rig.sender.buildFrame(persisted)
			if err != nil {
				t.Fatal(err)
			}
			parts := make([]string, 0, len(frame.Path))
			for _, hop := range frame.Path {
				parts = append(parts, hop.String())
			}
			if got := strings.Join(parts, ","); got != selected {
				t.Fatalf("retry path %q, want %q", got, selected)
			}
		})
	}
}

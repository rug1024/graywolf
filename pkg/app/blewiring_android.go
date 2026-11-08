//go:build android

package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/chrissnell/graywolf/pkg/kiss"
	"github.com/chrissnell/graywolf/pkg/messages"
	"github.com/chrissnell/graywolf/pkg/platformsvc"
)

// injectAndroidBLEClient pushes the live platformsvc client into the kiss
// package so ScanBLEMobilinkd and OpenBLEMobilinkd can route through the
// Kotlin BLE bridge. Called once during wireServicesInner after the platform
// client is confirmed non-nil.
func (a *App) injectAndroidBLEClient() {
	if a.platformClient != nil {
		kiss.SetAndroidBLEClient(a.platformClient)
	}
}

func (a *App) wireTrackerHistory() {
	client, ok := a.platformClient.(interface {
		SetTrackerHistoryImporter(func(context.Context, string, platformsvc.TrackerHistoryRecord) error)
	})
	if !ok {
		return
	}
	client.SetTrackerHistoryImporter(func(ctx context.Context, addr string, r platformsvc.TrackerHistoryRecord) error {
		call, err := a.store.ResolveStationCallsign(ctx)
		if err != nil {
			return err
		}
		if !strings.EqualFold(call, r.To) {
			return nil
		}
		return a.msgSvc.ImportTrackerHistory(ctx, call, messages.TrackerHistoryMessage{
			TrackerID: fmt.Sprintf("%s:%016x", strings.ToUpper(addr), r.Generation),
			Sequence:  uint64(r.Sequence), From: r.From, To: r.To, Text: r.Text,
		})
	})
}

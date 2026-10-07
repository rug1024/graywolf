package webapi

import "github.com/chrissnell/graywolf/pkg/configstore"

// channelUsesModemStats is the single source of truth for selecting the
// modembridge stats cache. A nil channel fails closed: only a persisted
// channel with an audio input is modem-backed; KISS-only channels must read
// their counters from kiss.Manager instead.
func channelUsesModemStats(ch *configstore.Channel) bool {
	return ch != nil && ch.InputDeviceID != nil
}

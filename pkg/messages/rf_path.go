package messages

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chrissnell/graywolf/pkg/ax25"
)

var rfPacketAddress = regexp.MustCompile(`^[A-Z0-9]{1,6}(-([0-9]|1[0-5]))?$`)

// ValidateRFPacketPath accepts an explicit AX.25 RF path, including an
// empty (direct) path. Internet markers and already-used hops are not
// valid operator-selected RF routing requests.
func ValidateRFPacketPath(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	parts := strings.Split(path, ",")
	if len(parts) > 8 {
		return fmt.Errorf("RF packet path exceeds 8 addresses")
	}
	for _, part := range parts {
		part = strings.ToUpper(strings.TrimSpace(part))
		if !rfPacketAddress.MatchString(part) || strings.Contains(part, "*") || part == "TCPIP" || part == "TCPXX" || (len(part) == 3 && strings.HasPrefix(part, "QA")) {
			return fmt.Errorf("invalid RF packet path address %q", part)
		}
		if _, err := ax25.ParseAddress(part); err != nil {
			return fmt.Errorf("invalid RF packet path: %w", err)
		}
	}
	return nil
}

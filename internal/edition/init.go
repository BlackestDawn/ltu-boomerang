package edition

import (
	"fmt"

	"github.com/BlackestDawn/ltu-boomerang/internal/edition/australia"
)

func GetEdition(edition string) (Edition, error) {
	switch edition {
	case "AU":
		return australia.New(), nil
	default:
		return nil, fmt.Errorf("unknown edition identifier: '%s'", edition)
	}
}

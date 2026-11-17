package edition

import (
	"github.com/BlackestDawn/ltu-boomerang/internal/core"
)

type Edition interface {
	Cards() []*core.Card
	Rules() map[core.Category]core.ScoringRule
}

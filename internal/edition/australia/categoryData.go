package australia

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

func GetCategoryStrings() map[core.Category]string {
	return map[core.Category]string{
		core.ThrowCatch:  "Throw & Catch",
		core.Site:        "Tourist Site",
		core.RegionBonus: "Regions",
		core.Collection:  "Collections",
		core.Animal:      "Animals",
		core.Activity:    "Activity",
		core.TieBreaker:  "Tie Breaker",
	}
}

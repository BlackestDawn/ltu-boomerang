package usa

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

func GetCategoryStrings() map[core.Category]string {
	return map[core.Category]string{
		core.ThrowCatch:  "Throw & Catch",
		core.Site:        "Tourist Sites",
		core.RegionBonus: "Regions",
		core.Collection:  "Americana",
		core.Animal:      "Animals",
		core.Activity:    "Activities",
		core.TravelBonus: "Coast to Coast",
	}
}

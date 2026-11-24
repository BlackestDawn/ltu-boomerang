package europe

import "github.com/BlackestDawn/ltu-boomerang/internal/core"

func GetCategoryStrings() map[core.Category]string {
	return map[core.Category]string{
		core.ThrowCatch:  "Throw & Catch",
		core.Site:        "Countries",
		core.RegionBonus: "Regions",
		core.Collection:  "Cuisine",
		core.Animal:      "Transport",
		core.Activity:    "National Treasures",
		core.TravelBonus: "Passport",
	}
}

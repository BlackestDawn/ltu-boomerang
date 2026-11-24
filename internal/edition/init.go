package edition

import (
	"fmt"

	"github.com/BlackestDawn/ltu-boomerang/internal/core"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition/australia"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition/europe"
	"github.com/BlackestDawn/ltu-boomerang/internal/edition/usa"
)

func GetEdition(edition string) (Edition, error) {
	switch edition {
	case "AU":
		core.CatStrings = australia.GetCategoryStrings()
		return australia.New(), nil
	case "EU":
		core.CatStrings = europe.GetCategoryStrings()
		return europe.New(), nil
	case "US":
		core.CatStrings = usa.GetCategoryStrings()
		return usa.New(), nil
	default:
		return nil, fmt.Errorf("unknown edition identifier: '%s'", edition)
	}
}

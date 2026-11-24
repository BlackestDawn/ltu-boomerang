package usa

var regionNames = []string{
	"Pacific",
	"Mountain",
	"Southwest",
	"Midwest",
	"Southeast",
	"Mid-Atlantic",
	"Northeast",
}

var regionClusters = [][]byte{
	{'A', 'B', 'C', 'D'},
	{'E', 'F', 'G', 'H'},
	{'I', 'J', 'K', 'L'},
	{'M', 'N', 'O', 'P'},
	{'Q', 'R', 'S', 'T'},
	{'U', 'V', 'W', 'X'},
	{'Y', 'Z', '@', '#'},
}

// Every region is worth 3 points to the first player completing it.
var regionBonus = []int{3, 3, 3, 3, 3, 3, 3}

// Coast towns are circled in blue on the score sheet.
var (
	westCoast = []byte{'A', 'B', 'C', 'D'}
	eastCoast = []byte{'T', 'X', 'Z', '@', '#'}
)

// routes are the pathways between towns on the score sheet map.
//
// NOTE: transcribed from the small score sheet picture in the rulebook; check
// against a real score sheet.
var routes = [][2]byte{
	{'A', 'E'}, {'B', 'G'}, {'C', 'G'}, {'C', 'F'}, {'D', 'F'}, {'D', 'I'},
	{'E', 'G'}, {'E', 'M'}, {'F', 'H'}, {'G', 'H'}, {'H', 'M'}, {'H', 'J'},
	{'I', 'J'}, {'J', 'L'}, {'K', 'L'}, {'L', 'P'}, {'K', 'R'},
	{'M', 'N'}, {'N', 'O'}, {'O', 'U'}, {'P', 'Q'}, {'Q', 'U'}, {'Q', 'S'},
	{'R', 'S'}, {'S', 'T'}, {'S', 'W'},
	{'U', 'V'}, {'V', 'W'}, {'V', 'X'}, {'V', 'Z'}, {'V', 'Y'},
	{'Y', '@'}, {'Y', '#'},
}

package europe

var regionNames = []string{
	"Scandinavia",
	"Baltics",
	"British Isles",
	"Benelux",
	"Central Europe",
	"Eastern Europe",
	"Southern Europe",
}

var regionClusters = [][]byte{
	{'A', 'B', 'C', 'D', 'E'},
	{'F', 'G', 'H'},
	{'I', 'J'},
	{'K', 'L', 'M'},
	{'N', 'O', 'P'},
	{'Q', 'R', 'S', 'T'},
	{'U', 'V', 'W', 'X', 'Y', 'Z', '@', '#'},
}

// Bonus for the first player to visit every country in the region.
//
// NOTE: Benelux, Central, Eastern and Southern Europe are read from the
// rulebook. The values for Scandinavia, the Baltics and the British Isles are
// scribbled over in the only score sheet picture available and are guesses
// within the rulebook's 2-5 range; check against a real score sheet.
var regionBonus = []int{4, 3, 2, 3, 3, 3, 5}

package edition

import "testing"

func TestGetEditionAustralia(t *testing.T) {
	ed, err := GetEdition("AU")
	if err != nil {
		t.Fatalf("GetEdition(AU): %v", err)
	}
	if len(ed.Cards()) != 28 {
		t.Errorf("AU edition has %d cards, want 28", len(ed.Cards()))
	}
}

func TestGetEditionUnknown(t *testing.T) {
	if _, err := GetEdition("XX"); err == nil {
		t.Error("expected error for unknown edition")
	}
}

func TestGetEditionEuropeAndUSA(t *testing.T) {
	for _, id := range []string{"EU", "US"} {
		ed, err := GetEdition(id)
		if err != nil {
			t.Fatalf("GetEdition(%s): %v", id, err)
		}
		if len(ed.Cards()) != 28 {
			t.Errorf("%s edition has %d cards, want 28", id, len(ed.Cards()))
		}
	}
}

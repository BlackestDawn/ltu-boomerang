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

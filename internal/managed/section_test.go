package managed

import (
	"bytes"
	"errors"
	"testing"
)

func TestMergeAppendsSectionWithoutReplacingHumanContent(t *testing.T) {
	human := []byte("# Team instructions\n\nKeep this text.\n")
	desired := []byte(Wrap("generated policy"))

	got, err := Merge(human, desired)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(got, human) {
		t.Fatalf("human content was not preserved: %q", got)
	}
	section, state, err := Extract(got)
	if err != nil {
		t.Fatal(err)
	}
	if state != Present || !bytes.Equal(section, desired) {
		t.Fatalf("managed section = %q, want %q", section, desired)
	}
}

func TestMergeReplacesOnlyExistingSection(t *testing.T) {
	before := []byte("# Before\n\n")
	old := []byte(Wrap("old policy"))
	after := []byte("\n# After\n")
	current := append(append(append([]byte(nil), before...), old...), after...)
	desired := []byte(Wrap("new policy"))

	got, err := Merge(current, desired)
	if err != nil {
		t.Fatal(err)
	}
	want := append(append(append([]byte(nil), before...), desired...), after...)
	if !bytes.Equal(got, want) {
		t.Fatalf("Merge() = %q, want %q", got, want)
	}
}

func TestMergeMigratesLegacyFullyManagedFile(t *testing.T) {
	current := []byte(LegacyMarker + "\n\nlegacy generated content\n")
	desired := []byte(Wrap("new policy"))

	got, err := Merge(current, desired)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, desired) {
		t.Fatalf("Merge() = %q, want %q", got, desired)
	}
}

func TestLocateRejectsDuplicateMarkers(t *testing.T) {
	content := []byte(StartMarker + "\n" + StartMarker + "\n" + EndMarker + "\n")
	if _, _, err := Locate(content); !errors.Is(err, ErrMalformed) {
		t.Fatalf("Locate() error = %v, want ErrMalformed", err)
	}
}

func TestLocateRejectsReversedMarkers(t *testing.T) {
	content := []byte(EndMarker + "\n" + StartMarker + "\n")
	if _, _, err := Locate(content); !errors.Is(err, ErrMalformed) {
		t.Fatalf("Locate() error = %v, want ErrMalformed", err)
	}
}

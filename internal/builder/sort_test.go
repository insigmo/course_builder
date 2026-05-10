package builder_test

import (
	"testing"

	"github.com/insigmo/course_builder/internal/builder"
)

func TestSortNames_NumberedFirst(t *testing.T) {
	names := []string{"Введение", "3. Практика", "1. Основы", "2. Теория"}
	builder.SortNames(names, nil)
	want := []string{"1. Основы", "2. Теория", "3. Практика", "Введение"}
	for i, n := range names {
		if n != want[i] {
			t.Errorf("SortNames[%d] = %q, want %q", i, n, want[i])
		}
	}
}

func TestSortNames_NumericOrder(t *testing.T) {
	names := []string{"10.mp4", "2.mp4", "1.mp4", "9.mp4"}
	builder.SortNames(names, nil)
	want := []string{"1.mp4", "2.mp4", "9.mp4", "10.mp4"}
	for i, n := range names {
		if n != want[i] {
			t.Errorf("SortNames[%d] = %q, want %q", i, n, want[i])
		}
	}
}

func TestSortNames_AlphabeticUnumbered(t *testing.T) {
	names := []string{"Zebra", "Apple", "Мир", "Горутина"}
	builder.SortNames(names, nil)
	// All unnumbered — alphabetical
	if names[0] != "Apple" {
		t.Errorf("expected Apple first, got %q", names[0])
	}
}

func TestSortNames_WithRemovablePrefix(t *testing.T) {
	removable := map[string]struct{}{"[sw]": {}}
	names := []string{"[sw] 3 Практика.mp4", "[sw] 1 Основы.mp4", "[sw] 2 Теория.mp4"}
	builder.SortNames(names, removable)
	want := []string{"[sw] 1 Основы.mp4", "[sw] 2 Теория.mp4", "[sw] 3 Практика.mp4"}
	for i, n := range names {
		if n != want[i] {
			t.Errorf("SortNames[%d] = %q, want %q", i, n, want[i])
		}
	}
}

func TestSortNames_Empty(t *testing.T) {
	var names []string
	builder.SortNames(names, nil) // should not panic
}

func TestSortNames_Single(t *testing.T) {
	names := []string{"one.mp4"}
	builder.SortNames(names, nil)
	if names[0] != "one.mp4" {
		t.Errorf("got %q", names[0])
	}
}

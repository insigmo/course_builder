package builder

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/user/course-builder/internal/prefix"
)

var numRE = regexp.MustCompile(`\d+`)

type sortKey struct {
	num  int
	text string
}

func makeSortKey(name string, removable map[string]struct{}) sortKey {
	s := strings.ToLower(prefix.Strip(name, removable))
	var num int
	if m := numRE.FindString(s); m != "" {
		num, _ = strconv.Atoi(m)
	}
	var runes []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			runes = append(runes, r)
		}
	}
	return sortKey{num: num, text: string(runes)}
}

// SortNames sorts a slice of filenames by numeric-then-alpha order.
func SortNames(names []string, removable map[string]struct{}) {
	type item struct {
		key  sortKey
		name string
	}
	items := make([]item, len(names))
	for i, n := range names {
		items[i] = item{key: makeSortKey(n, removable), name: n}
	}
	sort.SliceStable(items, func(i, j int) bool {
		a, b := items[i].key, items[j].key
		if a.num != b.num {
			return a.num < b.num
		}
		return a.text < b.text
	})
	for i, it := range items {
		names[i] = it.name
	}
}

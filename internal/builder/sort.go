package builder

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/insigmo/course_builder/internal/prefix"
)

var numRE = regexp.MustCompile(`\d+`)

type sortKey struct {
	hasNum bool
	num    int
	text   string
}

func makeSortKey(name string, removable map[string]struct{}) sortKey {
	s := strings.ToLower(prefix.Strip(name, removable))

	var num int
	hasNum := false
	if m := numRE.FindString(s); m != "" {
		hasNum = true
		num, _ = strconv.Atoi(m)
	}

	var runes []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			runes = append(runes, r)
		}
	}
	return sortKey{hasNum: hasNum, num: num, text: string(runes)}
}

// SortNames sorts filenames: numbered entries first (by number), then unnumbered alphabetically.
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
		if a.hasNum != b.hasNum {
			return a.hasNum
		}
		if a.hasNum && a.num != b.num {
			return a.num < b.num
		}
		return a.text < b.text
	})
	for i, it := range items {
		names[i] = it.name
	}
}

package internal

import (
	"errors"
	"fmt"
)

// сравнить два массива со строками
func areStringSlicesSame(a, b []string) error {
	if len(a) != len(b) {
		return fmt.Errorf("Размерности не сходятся. Ожидалось: %d, получено %d", len(a), len(b))
	}
	type Pair struct{ E, R bool }
	var seen = make(map[string]Pair, 0)
	for i := range a {
		s := seen[a[i]]
		s.E = true
		seen[a[i]] = s

		s = seen[b[i]]
		s.R = true
		seen[b[i]] = s
	}
	var (
		errMsg string
	)
	for id, s := range seen {
		if s.E && s.R {
			continue
		}
		if !s.E {
			errMsg += fmt.Sprintf("\nНе хватает элемента: %s", id)
		} else {
			errMsg += fmt.Sprintf("\nЛишний элемент: %s", id)
		}
	}
	if len(errMsg) == 0 {
		return nil
	}
	return errors.New(errMsg)
}

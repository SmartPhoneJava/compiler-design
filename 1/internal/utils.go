package internal

func RemoveBrackets(s string) string {
	if len(s) == 0 || s[0] != '(' {
		return s
	}
	var countBrackets int
	for i, symbol := range s {
		if symbol == '(' {
			countBrackets++
		} else if symbol == ')' {
			countBrackets--
			if countBrackets == 0 {
				if i == len(s)-1 { // последняя скобка первого уровня
					return s[1 : len(s)-1]
				}
				return s
			}
		}
	}

	return s
}

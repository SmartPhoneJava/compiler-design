package internal

/*
	Построить мапу дистанций - откуда куда можно залететь напрямую
	ключ - исходная вершина
	значение -массив вершин, куда ведёт

	O(P)
*/
func (cfr CFR) buildDistMap() map[string][]string {
	var fromTo = make(map[string][]string)
	for _, q := range cfr.P {
		goTo, ok := fromTo[q.From]
		if !ok {
			goTo = cfr.ToNoneTerminals(q.To)
		} else {
			goTo = append(goTo, cfr.ToNoneTerminals(q.To)...)
		}
		fromTo[q.From] = goTo
	}
	return fromTo
}

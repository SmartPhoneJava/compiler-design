package constants

const (
	// Action2NKA По регулярному выражению строит НКА
	Action2NKA = iota + 1
	// Action2DKA По НКА строит эквивалентный ему ДКА.
	Action2DKA
	// ActionMinimize По ДКА строит эквивалентный ему КА
	/*
		имеющий наименьшее возможное количество состояний.
		Указание. Воспользоваться алгоритмом, приведенным по адресу
		http://neerc.ifmo.ru/wiki/index.php?title=Алгоритм_Бржозовского
	*/
	ActionMinimize
)

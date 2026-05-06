package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
	[]string -> срез. Разбираем подробнее:

	УЛУЧШЕНИЕ: синтаксис ниже был неверным — исправленные версии:
	var a [3]rune — массив из 3 rune (нельзя присвоить строку напрямую)
	var a string = "321" — строка, но это НЕ статический массив.
		string в Go — read-only slice байт: структура {ptr *byte, len int}.
		Иммутабельна — нельзя изменить символ по индексу.
	var a [5]int = [5]int{1, 2, 3, 4, 5} — статический массив, литерал через {}, не []

	Итого:
	[N]T  — статический массив, размер фиксирован на этапе компиляции (как int arr[5] в C)
	[]T   — срез (динамический), хранит ptr + len + cap, данные в heap
	string — read-only slice байт, не массив символов

	[]string  — срез, каждый элемент хранит строку (структуру ptr+len)
	[]*string — срез указателей на строки, сами строки живут отдельно в heap
	*[]string — указатель на заголовок среза (нужен редко)

	Аналогия с C:
	char *arr[]  ≈ []*string  (массив указателей на строки)
	char **arr   ≈ []*string  (то же самое, другой синтаксис)
	*(char *[])  — разыменование указателя на массив указателей, в Go прямого аналога нет.
	              *[]string указывает на весь заголовок среза (ptr+len+cap), в C такого нет.
*/

const msg = "Hey there! I am using WhatsApp! WhatsApp"
const cutset = ",.!?:;\"'()`"

func main() {
	if len(msg) == 0 {
		fmt.Printf("msg string is empty")
		return
	}

	counts := countWords(msg)

	keys := sortByFreq(counts)
	for _, v := range keys {
		fmt.Printf("%s : %d\n", v, counts[v])
	}

	// A: без make(map[string]int) переменная типа map == nil (нулевое значение).
	// Читать из nil-map можно (вернёт 0), писать — паника.
	// make аллоцирует hash-таблицу и возвращает указатель на неё.

	// A: stdlib (standard library) — пакеты встроенные в Go: fmt, strings, sort, os, net и ~150 других.
	// Всё что импортируешь без `go get` — это stdlib.
	// Сторонние пакеты устанавливаются через `go get`, прописываются в go.mod.

	/*
		Памятка по fmt.Printf форматам:
		%v  — значение в стандартном виде
		%+v — структура с именами полей
		%#v — Go-синтаксис (как написать литерал в коде)
		%q  — строка в кавычках (или слайс строк в кавычках)
		%s  — строка без кавычек
		%d  — целое число
	*/
}

func countWords(text string) map[string]int {
	if len(text) == 0 {
		return make(map[string]int) // УЛУЧШЕНИЕ: make(map[string]int, 0) == make(map[string]int), второй аргумент лишний
	}

	counts := make(map[string]int)

	fields := strings.Fields(text)
	for _, v := range fields {
		v = strings.Trim(strings.ToLower(v), cutset)
		counts[v]++
	}

	// A: FieldsSeq (Go 1.23+) возвращает итератор iter.Seq[string] вместо слайса.
	// Итератор отдаёт элементы по одному через range, не аллоцирует весь слайс сразу.
	// Для больших строк экономит память. Приставка Seq — соглашение для всех
	// функций-итераторов в stdlib. В нашем случае Fields достаточно.

	return counts
}

func sortByFreq(counts map[string]int) []string {
	if counts == nil {
		return nil
	}

	keys := make([]string, 0, len(counts)) // УЛУЧШЕНИЕ: pre-allocate по размеру map — избегаем лишних реаллокаций
	for k := range counts {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return counts[keys[i]] > counts[keys[j]]
	})
	// A: counts[keys[i]] > counts[keys[j]] — понял правильно.
	// keys[i] — слово по индексу i, counts[keys[i]] — его частота.
	// less(i, j) == true означает "i стоит раньше j".
	// > : элемент с большей частотой идёт раньше → убывание.
	// < : элемент с меньшей частотой идёт раньше → возрастание.
	//
	// A: >= и <= нельзя — sort.Slice требует строгое упорядочение.
	// Если a == b, то a >= b и b >= a оба true — противоречие, алгоритм ломается.
	// != и == тоже не имеют смысла как компаратор порядка.
	// Только строгие > или <.

	return keys
}

/*
	Слайсы

	В Go слайс — структура из 3-х полей ("заголовок"), рантайм Go держит её на стеке:
	    ptr *T   — указатель на массив в heap
	    len int  — сколько элементов используется сейчас
	    cap int  — сколько места выделено в heap

	"Заголовок" — это и есть эта структура {ptr, len, cap}.
	Сам массив данных живёт в heap (как malloc в C).

	Аналог в C:
	    int *ptr = malloc(n * sizeof(int));
	    int len = n;
	    int cap = n;
	Go просто упаковывает это в одну структуру и управляет ей автоматически.

	var s []int        — {ptr: nil, len: 0, cap: 0}, nil slice
	s = append(s, 1)  — Go аллоцирует массив в heap, возвращает новый заголовок {ptr, 1, ?}
	s = append(s, 2, 3) — append вариативная: можно передать несколько элементов сразу.
	                       Если cap хватает — просто len += N, данные записываются в heap.
	                       Если нет — realloc: новый массив ~2x, копирование, новый ptr.

	Рантайм Go — код от команды Go, вшитый в каждый бинарник при компиляции.
	Отвечает за: планировщик горутин, garbage collector, растущий стек, defer, panic/recover.
	Аналогия в C: как если бы malloc, pthread и signal handlers были вшиты автоматически.
	Поэтому Go-бинарник весит ~2–5 МБ даже для Hello World — рантайм уже внутри.

	Этапы компиляции Go (аналогия с C):
	C:   препроцессор → компилятор → ассемблер → линкер → бинарник
	Go:  синтаксический разбор + type check → SSA → оптимизации → машинный код → статическая линковка
	go build делает всё за одну команду. Линковка статическая — всё в одном бинарнике, включая рантайм.
*/

/*
	strings.Fields

	msg := "Hey there! I am using WhatsApp!"
	msgFields := strings.Fields(msg)
	// → ["Hey", "there!", "I", "am", "using", "WhatsApp!"]
	// Разбивает по пробельным символам (пробел, \t, \n). Несколько пробелов подряд — ок.

	strings.ToLower / ToUpper — возвращают новую строку, оригинал не меняется.
	strings.ToTitle — конвертирует все символы в titlecase (≈ uppercase для Unicode edge cases).
	                  Отличается от устаревшего Title, которая делала заглавной первую букву слова.
	                  В обычном коде не нужна.

	strings.Trim(s, cutset) — убирает с обоих концов строки любые символы из cutset.
	Не подстроку — именно любой символ из набора по отдельности.
*/

/*
	пакет sort

	sort.Slice(slice any, less func(i, j int) bool)

	less — функция сравнения двух элементов по индексам.
	Можно передать анонимную функцию (как здесь) или именованную.
	Анонимная функция может захватывать переменные из внешнего scope (замыкание) —
	именно так counts доступен внутри less без передачи параметром.

	Требование к less: строгое упорядочение.
	less(i, j) и less(j, i) не могут оба быть true одновременно.
	Нарушение → непредсказуемый результат сортировки.
*/

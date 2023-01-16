package main

// CyclicQueue циклическая очередь - память выделяем один раз, размер не меняется.
type CyclicQueue[T any] struct {
	items    []T
	capacity int
	first    int
	last     int
}

// GetFirst возвращает первый элемент очереди и признак true вторым параметром, что элемент успешно получен. При попытке
// получить элемент из пустой очереди, второй параметр вернется false.
func (q *CyclicQueue[T]) GetFirst() (T, bool) {
	if q.Len() == 0 {
		var blank T
		return blank, false
	}

	q.first++
	v := q.items[q.first]
	if q.first == q.capacity {
		q.first = -1
	}

	return v, true
}

// AddLast добавляет элемент в конец очереди, возвращает true, если удалось добавить элемент. Вернет false, если пытаемся
// добавить элемент в полную очередь.
func (q *CyclicQueue[T]) AddLast(v T) bool {
	if q.Len() == q.Cap() {
		return false
	}

	q.items[q.last] = v
	q.last++
	if q.last > q.capacity {
		q.last = 0
	}

	return true
}

// Len возвращает размер очереди.
func (q *CyclicQueue[T]) Len() int {
	if q.first >= q.last {
		return q.last - q.first + q.capacity
	}

	return q.last - q.first - 1
}

// Cap возвращает емкость очереди.
func (q *CyclicQueue[T]) Cap() int {
	return q.capacity
}

// NewCyclicQueue возвращает новую циклическую очередь.
func NewCyclicQueue[T any](capacity int) *CyclicQueue[T] {
	return &CyclicQueue[T]{
		items:    make([]T, capacity+1),
		capacity: capacity,
		first:    -1,
	}
}

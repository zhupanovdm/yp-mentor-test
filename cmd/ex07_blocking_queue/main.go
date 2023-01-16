// Демонстрирует шаблон синхронизации по монитору.
// Попробуем реализовать блокирующую очередь. Если очередь пустая, то вызывающий при получении элемента будет блокироваться
// до тех пор, пока в очередь не будет добавлен элемент. Если очередь полная, то вызывающий при попытке добавить элемент
// будет заблокирован до тех пор, пока из очереди не извлекут элемент.
// На практике такое реализовывать не пригодится, т.к. в го уже есть каналы) но ради развлечения, почему бы и нет)
package main

import (
	"fmt"
	"sync"
)

// Блокирующая очередь, обертка над циклической очередью.
type blockingQueue[T any] struct {
	cond  *sync.Cond
	queue *CyclicQueue[T]
}

func main() {
	// Создадим очередь с емкостью 30.
	queue := newBlockingQueue[int](30)

	go func() {
		for i := 0; i < 100; i++ {
			n := i + 1
			queue.add(n)
		}
	}()

	for i := 0; i < 100; i++ {
		fmt.Println(queue.get())
	}
}

// Добавляет элемент в конец очереди.
func (w *blockingQueue[T]) add(v T) {
	w.cond.L.Lock()
	defer w.cond.L.Unlock()

	for w.queue.Len() == w.queue.Cap() {
		w.cond.Wait()
	}

	w.queue.AddLast(v)

	w.cond.Signal()
}

// Получает элемент в начале очереди.
func (w *blockingQueue[T]) get() T {
	w.cond.L.Lock()
	defer w.cond.L.Unlock()

	for w.queue.Len() == 0 {
		w.cond.Wait()
	}

	v, _ := w.queue.GetFirst()

	w.cond.Signal()

	return v
}

// Создать блокирующую очередь.
func newBlockingQueue[T any](capacity int) *blockingQueue[T] {
	var mu sync.Mutex

	return &blockingQueue[T]{
		queue: NewCyclicQueue[T](capacity),
		cond:  sync.NewCond(&mu),
	}
}

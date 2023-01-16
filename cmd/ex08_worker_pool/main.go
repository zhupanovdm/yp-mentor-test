// Демонстрирует шаблон worker pool и работу с контекстов в го-рутинах.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Работа для воркера, при вызове cancel отменяется все.
type job func(ctx context.Context, cancel func())

func main() {
	// После пяти секунд пошлем сигнал завершения.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создадим 4 рабочих.
	pool := make(chan job, 4)
	complete := createWorkerPool(ctx, 4, pool)

	go func() {
		for {
			pool <- generateRandomWordsJob // Добавляем работы работникам.
		}
	}()

	// Ждем завершения всех воркеров.
	<-complete

	log.Println("all workers stopped")
}

//
func createWorkerPool(ctx context.Context, n int, pool chan job) chan any {
	var wg sync.WaitGroup

	// Создадим контекст с отменой, чтобы можно было отменить все.
	ctx, cancel := context.WithCancel(ctx)

	for i := 0; i < n; i++ {
		createWorker(ctx, pool, cancel, &wg)
	}

	// Широковещательный сигнал завершения всех воркеров.
	complete := make(chan any)

	go func() {
		// Освободим ресурсы при завершении.
		defer cancel()
		wg.Wait()

		// Просинализируем, что все воркеры отработали.
		close(complete)
	}()

	return complete
}

// Запускает воркера.
func createWorker(ctx context.Context, jobs <-chan job, cancel func(), wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		log.Println("worker started")

		done := ctx.Done()

		// Цикл активен пока есть что делать и контекст не отменен.
		active := true
		for active {
			select {
			case job, ok := <-jobs:
				// При закрытии канала нечего будет делать.
				if !ok {
					active = false
					log.Println("jobs channel is closed")
				}

				// Передадим контекст в job, чтобы была возможность отменить длительные операции.
				job(ctx, cancel)

			// Следим, не было ли отмены контекста.
			case <-done:
				active = false
				log.Println("worker context is cancelled")
			}
		}

		log.Println("worker stopped")
	}()
}

// Полезная нагрузка - генерация случайных слов.
func generateRandomWordsJob(_ context.Context, cancel func()) {
	resp, err := http.Get("https://random-word-api.herokuapp.com/word?number=3")
	if err != nil {
		cancel()
		return
	}

	defer resp.Body.Close()

	var words []string
	if err = json.NewDecoder(resp.Body).Decode(&words); err != nil {
		cancel()
		return
	}

	fmt.Println(strings.Join(words, " "))
}

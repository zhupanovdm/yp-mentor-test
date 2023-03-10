// Демонстрирует мультиплексирование с помощью select.
// Представим задачу, какой-то источник генерирует события, но мы не хотим вызывать сервис для каждого события, а хотим
// сгруппировать события в пакеты по N штук и передать пакет в сервис. Но события приходят не регулярно, а данные требуют
// оперативной актуализации. Рассмотрим пример программы дебаунсера.
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Получим дебаунсер.
	input := debounce(10, time.Second, func(s ...string) {
		// Просто пощитаем количество событий в пакете.
		fmt.Println(len(s))
	})

	for {
		// Имитация неравномерной генерации событий.
		d := time.Duration(50+rand.Intn(3)*100) * time.Millisecond
		time.Sleep(d)

		// Получили событие.
		input <- "foo"
	}
}

// Группирует входящие события в пакеты заданного размера batch и оповещает о поступлении пакета событий, обратным
// вызовом accept. Если за указанный таймаут timeout указанного размера пакета не набралось, уведомляет обо всех имеющихся
// событиях (размер пакета меньше, чем заданная величина).
func debounce[T any](batch int, timeout time.Duration, accept func(s ...T)) chan<- T {
	input := make(chan T, batch) // Для входящих событий.
	buf := make([]T, 0, batch)   // Временный буфер для еще необработанных событий.

	// Обработчик пакета.
	handle := func() {
		if len(buf) == 0 {
			return
		}

		accept(buf...) // Оповестим о пакете.
		buf = buf[:0]  // Опустошим буфер.
	}

	go func() {
		ticker := time.NewTicker(timeout)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				handle() // Оповестим по таймауту.

			case item, ok := <-input:
				// Обработка прекратится, как только входной канал будет закрыт.
				if !ok {
					return
				}

				// Запомним в буфере.
				buf = append(buf, item)

				// Буфер заполнен?
				if len(buf) == batch {
					handle()
					ticker.Reset(timeout) // Отсчет таймаута начинается заново.
				}
			}
		}
	}()

	return input
}

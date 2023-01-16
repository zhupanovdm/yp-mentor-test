// Демонстрирует шаблон параллельный цикл.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// Исходный слайс, который будем преобразовывать.
	src := make([]int, 50)
	for i := 0; i < len(src); i++ {
		src[i] = i
	}

	// Выполнить параллельное преобразование.
	list := <-parallelMap(src, func(i int) []string {
		list, err := generateRandomWords(i + 10)
		if err != nil {
			log.Fatalln(err)
		}
		return list
	}, 3)

	for _, word := range list {
		fmt.Println(word)
	}
}

// Выполняет преобразование слайса типа A в слайс типа R используя заданное количество параллельно выполняющихся горутин.
// Возвращает канал, в который будет помещен результирующий слайс по окончании преобразования всех элементов.
func parallelMap[A, R any](src []A, transform func(A) R, parallel int) chan []R {
	tokens := make(chan any, parallel)     // Семафор для ограничения количества параллельно выполняющихся горутин.
	results := make(chan func(), parallel) // Для сигналов завершенного преобразования элемента слайса.

	result := make([]R, len(src)) // Результат преобразования.
	resultCh := make(chan []R)    // Сигнал завершения преобразования всех слайсов.

	// Для каждого элемента слайса запустим горутину.
	for i := 0; i < len(src); i++ {
		go func(index int) {
			// Захватим маркер, для ограничения параллельного выполнения.
			tokens <- struct{}{}

			// Выполним преобразование элемента и отправим сигнал окончания преобразования элемента.
			t := transform(src[index])
			results <- func() { result[index] = t }

			// Освободим маркер параллельного выполнения.
			<-tokens
		}(i)
	}

	go func() {
		// Количество преобразованных элементов.
		completed := 0

		// Получаем сигналы об окончании преобразования элементов.
		for apply := range results {
			apply() // Применить к результирующему слайсу.

			// Проверим: все элементы?
			if completed++; completed == len(src) {
				break
			}
		}

		// Тогда отдадим сигнал об окончании работы.
		resultCh <- result
	}()

	return resultCh
}

// Полезная нагрузка: принимает кол-во случайных слов, которые надо сгенерировать, возвращает слайс со случайными
// словами.
func generateRandomWords(number int) ([]string, error) {
	serviceURL := url.URL{
		Scheme:   "https",
		Host:     "random-word-api.herokuapp.com",
		Path:     "/word",
		RawQuery: fmt.Sprintf("number=%d", number),
	}

	resp, err := http.Get(serviceURL.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var words []string
	if err = json.NewDecoder(resp.Body).Decode(&words); err != nil {
		return nil, err
	}

	return words, nil
}

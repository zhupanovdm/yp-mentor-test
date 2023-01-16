// Демонстрирует шаблон pipeline (объединение каналов в конвейер). Выход одного канала может быть "подсоединен"
// ко входу второго и т.д.
// Представим себе, что мы конструируем роботизированную фабрику с конвейерами. На каждой из линий стоит робот (горутина),
// который производит над полуфабрикатом технологическую операцию, затем передает для дальнейшей обработки на следующую
// линию и так далее, пока все технологические операции не будут выполнены. По мере прохождения изделием производственных
// стадий для каждого изделия должен быть сформирован паспорт технологического процесса (какой робот и какую выполнял
// операцию).
package main

import (
	"bytes"
	"fmt"
	"log"
	"time"
)

// Конструктор линии. Принимает вход предыдущей линии, возвращает выход конвейерной линии.
type pipeFactory func(input chan *bytes.Buffer) chan *bytes.Buffer

func main() {
	// Наша фабрика.
	fabric := pipe(
		// Сборочный цех
		pipe(
			robot(op("ROBOT1", "Assembly with detail 1"), 5),
			robot(op("ROBOT2", "Assembly with detail 2"), 3),
		),
		// Сварочный цех
		pipe(
			robot(op("ROBOT3", "Weld detail 1"), 3),
			robot(op("ROBOT4", "Weld detail 2"), 2),
		),
		// Покраска
		robot(op("ROBOT5", "Paint"), 2),
	)

	// Запустим фабрику.
	go accept(fabric(generate(5)))

	// Дадим немного поработать.
	time.Sleep(time.Second)
}

// Возвращает фабрику конвейеров. Каждый созданный конвейер будет состоять из указанных производственных линий.
// Таким образом, линия в цепочке может быть другим конвейером.
func pipe(pipes ...pipeFactory) pipeFactory {
	return func(input chan *bytes.Buffer) chan *bytes.Buffer {
		line := input

		// Выход каждой линии присоединим ко входу следующей линии.
		for _, pipe := range pipes {
			line = pipe(line)
		}

		// Возвращаем выход всей цепочки.
		return line
	}
}

// Возвращает фабрику роботов. Робота тоже можно представить в виде линии: принимает изделие на вход, выполняет
// операцию, подает изделие на выход линии. Для робота должна быть задана работа work, которую он выполняет и количество
// parallel параллельно работающих роботов на участке.
func robot(operation func(*bytes.Buffer), parallel int) pipeFactory {
	return func(input chan *bytes.Buffer) chan *bytes.Buffer {
		// Сюда роботы складывают обработанные изделия, количество мест на ленте должно быть не меньше количества
		// роботов, иначе они будут блокировать друг друга.
		output := make(chan *bytes.Buffer, parallel)

		// Программируем и запускаем каждого робота на выполнение заданной операции.
		for i := 0; i < parallel; i++ {
			go robotOperation(input, output, operation)
		}

		// Возвращаем конец робо-линии.
		return output
	}
}

// Возвращает технологическую операцию, которую робот выполняет на линии.
func op(name, op string) func(*bytes.Buffer) {
	return func(buffer *bytes.Buffer) {
		if _, err := fmt.Fprintf(buffer, "%s: %s\n", name, op); err != nil {
			log.Fatalln(err)
		}
	}
}

// Возвращает канал, в который генерируем изделия для дальнейшей обработки на фабрике.
func generate(prefetch int) chan *bytes.Buffer {
	// Добавим несколько мест, чтобы генератор мог поработать на опережение линии.
	pipe := make(chan *bytes.Buffer, prefetch)

	i := 0

	go func() {
		for {
			i++

			// Назначим каждому изделию уникальный идентификатор.
			item := fmt.Sprintf("ID: %d\n---\n", i)

			// Подаем на обработку.
			pipe <- bytes.NewBufferString(item)

			time.Sleep(50 * time.Millisecond)
		}
	}()

	return pipe
}

// Приемка готовых изделий.
func accept(srcPipe <-chan *bytes.Buffer) {
	for {
		fmt.Println(<-srcPipe)
	}
}

// Робот выполняет работу: берет с линии очередную деталь, выполняет операцию и кладет в конец линии, для дальнейшей
// обработки. Если деталей нет - робот ждет, если следующий участок занят - робот ждет.
func robotOperation(input chan *bytes.Buffer, output chan *bytes.Buffer, operate func(*bytes.Buffer)) {
	for {
		item := <-input
		operate(item)
		output <- item
	}
}

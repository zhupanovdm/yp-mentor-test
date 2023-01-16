// Содержит пример запуска горутин и демнострирует, что может пойти не так, если отсутвует синхронизация.
// Давайте представим себе игру, где есть 26 участников, каждый из которых должен назвать одну букву из алфавита.
// Понаблюдаем, получится ли у них без согласования друг с другом рассказать алфавит в правильной последовательности?
package main

import (
	"fmt"
	"time"
)

func main() {
	// Запустим 26 горутин, каждая назовет свою букву алфавита.
	// Cинхронизации участников нет, поэтому увидим "неправильный" алфавит, что-то вроде: YABIWMQVFOGPTCUHNRKLDJXES.
	for i := 1; i < 26; i++ {
		go sayLetter(i)
	}

	// Главная функция не ждет завершения горутин, поэтому придется немного выждать, этого времени должно хватить.
	time.Sleep(time.Second)
}

// Выводит букву с указанным номером.
func sayLetter(number int) {
	fmt.Printf("%c", 'A'+number-1)
}
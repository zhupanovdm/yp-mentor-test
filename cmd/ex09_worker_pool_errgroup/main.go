// Демонстрирует шаблон worker pool с помощью пакета errgroup.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Установим лимит 2-х одновременных воркеров.
	group, _ := errgroup.WithContext(context.Background())
	group.SetLimit(2)

	// Запустим 100 заданий.
	for i := 0; i < 100; i++ {
		group.Go(generateRandomWordsJob)
	}

	// Ждем завершения всех или пока один из воркеров не "отвалится".
	err := group.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}

// Полезная нагрузка.
func generateRandomWordsJob() error {
	resp, err := http.Get("https://random-word-api.herokuapp.com/word?number=3")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var words []string
	if err = json.NewDecoder(resp.Body).Decode(&words); err != nil {
		return err
	}

	fmt.Println(strings.Join(words, " "))

	return nil
}

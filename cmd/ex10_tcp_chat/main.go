// Демонстрирует работу с горутинами на примере простого чат сервера.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

const serverAddr = ":9001"

type client chan []byte // При записи в канал пишет на сокет клиента.

var clients = make(map[client]any)  // Множество всех подключенных клиентов.
var subscribe = make(chan client)   // Сигнал присоединения клиента.
var unsubscribe = make(chan client) // Сигнал отключения клиента.
var broadcast = make(chan string)   // Канал для широковещательных сообщений.

func main() {
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalln(err)
	}
	defer withErr(listener.Close)

	// Обработчик широковещательных сообщений.
	go publisher()

	// Для каждого нового входящего подключения клиента запускаем горутину обработчик клиента.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go handler(conn)
	}
}

// Привязать соединение к каналу клиента: все что поступает в канал писать в сокет.
func (c client) bind(conn net.Conn) {
	for msg := range c {
		if _, err := conn.Write(msg); err != nil {
			log.Println(conn.RemoteAddr(), err)
		}
	}
}

// Обработчик входящего подключения.
func handler(conn net.Conn) {
	defer withErr(conn.Close)

	remoteAddr := conn.RemoteAddr().String()

	log.Println(remoteAddr, "joins")
	broadcast <- message(remoteAddr, "joins")

	cl := make(client)
	go cl.bind(conn)

	// Для получения сообщений других пользователей.
	subscribe <- cl

	scanner := bufio.NewScanner(conn)
	defer withErr(scanner.Err)

	for scanner.Scan() {
		msg := scanner.Text()

		// Пользователь хочет выйти.
		if msg == "/q" {
			log.Println(remoteAddr, "leaves")
			break
		}

		// Опубликуем сообщения пользователя для всех.
		broadcast <- message(remoteAddr, msg)
	}

	// Не хотим ничего получать больше.
	unsubscribe <- cl

	broadcast <- message(remoteAddr, "leaves")
}

// Обработчик широковещательных сообщений. Синхронизирован с событиями входа и выхода клиентов.
func publisher() {
	for {
		select {
		case msg := <-broadcast:
			// Продублировать сообщение всем активным клиентам.
			for cl := range clients {
				cl <- []byte(msg)
			}
		case cl := <-subscribe:
			// Подписываем клиента на широковещательные сообщения.
			clients[cl] = struct{}{}
		case cl := <-unsubscribe:
			// Отписываем клиента от широковещательных сообщений.
			delete(clients, cl)
			close(cl)
		}
	}
}

// Представление сообщения.
func message(from, text string) string {
	return fmt.Sprintf("%v %s: %s\n", time.Now().Format(time.Kitchen), from, text)
}

// Логирующая ошибку обертка.
func withErr(f func() error) {
	if err := f(); err != nil {
		log.Println(err)
	}
}

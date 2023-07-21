package main

import (
	_ "github.com/lib/pq"
)

func main() {

	startBot()
	createTables()

	for update := range updates {
		if update.Message != nil { // Есть новое сообщение
			msg := processingUpdate(update)

			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}

		}
	}
}

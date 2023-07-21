package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var bot *tgbotapi.BotAPI
var updates tgbotapi.UpdatesChannel
var waitingResponses = make(map[int64]responseCondition)

type responseCondition struct {
	function   string
	stepNumber int16
}

var petsUnsavedData = make(map[int64]pet)

func startBot() {
	var err error
	bot, err = tgbotapi.NewBotAPI("6202309857:AAFqgXhoxno1hXCcmBaKo1EdXINYKJxxqGY")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates = bot.GetUpdatesChan(u)
}

func processingUpdate(update tgbotapi.Update) tgbotapi.MessageConfig {
	text := update.Message.Text      // Текст сообщения
	chatID := update.Message.Chat.ID //  ID чата
	userID := update.Message.From.ID // ID пользователя
	var replyMsg string = fmt.Sprintf("default message recieved from %s user", update.Message.From)

	log.Printf("[%s](%d) %s", update.Message.From.UserName, userID, text)

	// Отправляем ответ
	msg := tgbotapi.NewMessage(chatID, replyMsg) // Создаем новое сообщение
	//msg.ReplyToMessageID = update.Message.MessageID // Указываем сообщение, на которое нужно ответить

	switch update.Message.Text {

	case "/start":
		msg.Text = gainStartMsg()
		msg.ReplyMarkup = mainMenuKeyboard

	case "/help":
		msg.ReplyMarkup = mainMenuKeyboard
		msg.Text = gainHelpMsg()
		msg.ParseMode = "html"

	case "Мои питомцы":
		//addNewTreatmentDay(update.Message.From.UserName, 1)
	case "":

	case "close":
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	default:
		resp, ok := waitingResponses[chatID]
		if ok {
			switch resp.function {
			case "newPet":
				if resp.stepNumber == 1 {
					petsUnsavedData[chatID] = pet{}
				}
				newPet(resp.stepNumber, update.Message.Text, petsUnsavedData[chatID])

			}
		} else {

		}
	}

	return msg
}

func newPet(step int16, msg string, p pet) {
	//TODO возвращать ответное сообщение-запрос информации
	//TODO после каждого этапа получения информации создавать соответствубщий таймер -- удалять его при начале следующего шага
	// если таймер сработает раньше введения сообщения -- удалить запись из словарей responseCondition и petsUnsavedData
	switch step {
	case 1:
		//save pet_name
		p.petName = msg
	case 2:
		//save birthday
		//TODO педобработать строку в формат YYYY-MM-DD для БД
		p.birthday = msg
	case 3:
		//save image and imagePath
		//gain pet_number and save
		//send query to DB
		//TODO сохранить полученную картинку -- пока хз как. получить ее адрес сохранения и сохранить в структуру
		// возможно нужно будет сохранить на предыдущем этапе в БД -- тк изображение необязательно
		//TODO после сохранения в БД -- удалять запись из словарей responseCondition и petsUnsavedData
		p.imagePath = msg

	}
}

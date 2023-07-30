package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
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
		clearWaitingResponsesFromUser(chatID)

	case "/help":
		msg.ReplyMarkup = mainMenuKeyboard
		msg.Text = gainHelpMsg()
		msg.ParseMode = "html"
		clearWaitingResponsesFromUser(chatID)

	case "Мои питомцы":
		clearWaitingResponsesFromUser(chatID)
		msg.ReplyMarkup = myPetsMenu
		//addNewTreatmentDay(update.Message.From.UserName, 1)
	case "Добавить питомца":
		waitingResponses[chatID] = responseCondition{"newPet", 1}
		msg.Text = "Введите имя питомца ответным сообщением"

	case "close":
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	default:
		resp, ok := waitingResponses[chatID]
		if ok {
			switch resp.function {
			case "newPet":
				if resp.stepNumber == 1 {
					petNum := getMaxPetNumber(chatID)
					petsUnsavedData[chatID] = pet{chatId: chatID, petNumber: petNum}
				}
				resp.stepNumber, msg.Text = newPet(resp.stepNumber, update.Message.Text, petsUnsavedData[chatID])
				if resp.stepNumber == 0 {
					addNewPet(petsUnsavedData[chatID])
					clearWaitingResponsesFromUser(chatID)
				}
			}
		} else {
			log.Panic(ok)
		}
	}
	return msg
}

func clearWaitingResponsesFromUser(userId int64) {
	//удалить ожидание ответа
	//удалить сохранения о питомце
	resp, ok := waitingResponses[userId]
	if ok {
		switch resp.function {
		case "newPet":
			delete(petsUnsavedData, userId)
		}
	} else {
	}
	delete(waitingResponses, userId)
}

func newPet(step int16, msg string, p pet) (int16, string) {
	//обработка этапов заполнения информации о питомце
	//TODO возвращать ответное сообщение-запрос информации
	//TODO после каждого этапа получения информации создавать соответствубщий таймер -- удалять его при начале следующего шага
	// если таймер сработает раньше введения сообщения -- удалить запись из словарей responseCondition и petsUnsavedData
	switch step {
	case 1:
		//save pet_name
		p.petName = msg
		return 2, "Имя сохранено! Теперь введите дату рождения в формате день.месяц.год (пример: 01.05.2021) ответным сообщением"
	case 2:
		//save birthday
		//TODO педобработать строку в формат YYYY-MM-DD для БД
		day, err := strconv.Atoi(msg[:2])
		if err != nil {
			return 2, "Данные введены в неверном формате, попробуйте еще раз ввести дату рождения " +
				"в формате день.месяц.год (пример: 01.05.2021). Либо нажмите \"назад\""
		}
		month, err := strconv.Atoi(msg[3:5])
		if err != nil {
			return 2, "Данные введены в неверном формате, попробуйте еще раз ввести дату рождения " +
				"в формате день.месяц.год (пример: 01.05.2021). Либо нажмите \"назад\""
		}
		year, err := strconv.Atoi(msg[6:10])
		if err != nil {
			return 2, "Данные введены в неверном формате, попробуйте еще раз ввести дату рождения " +
				"в формате день.месяц.год (пример: 01.05.2021). Либо нажмите \"назад\""
		}
		p.birthday = strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + strconv.Itoa(day)
		return 3, "Дата рождения сохранена :) При желании: пришлите фото вашего питомца и я сохраню его в карточке питомца!"
	case 3:
		//save image and imagePath
		//gain pet_number and save
		//send query to DB
		//TODO сохранить полученную картинку -- пока не знаю как. получить ее адрес сохранения и сохранить в структуру
		// возможно нужно будет сохранить на предыдущем этапе в БД -- тк изображение необязательно
		//TODO после сохранения в БД -- удалять запись из словарей responseCondition и petsUnsavedData
		p.imagePath = msg
		return 0, "Спасибо, фото добавлено :) Вы будете автоматически перенаправлены в основное меню."
	default:
		//пользователь решил не добавлять фото
		//TODO сделать возможность обрабатывать прислал ли пользователь фото
		return 0, "Спасибо,вы будете автоматически перенаправлены в основное меню."

	}
}

package telegram

import (
	"birthdayBot/internal/configs"
	queries "birthdayBot/internal/db"
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

var mainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("💁‍♂ Жаңасын қосу"),
		tgbotapi.NewKeyboardButton("🗑 Тізімнен жою"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🗓 Тізімді қарау"),
	),
)
var backToMainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🔁 Артқа қайту"),
	),
)

type Bot struct {
	bot  *tgbotapi.BotAPI
	conf *configs.Configuration
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{
		bot: bot,
	}
}

var welcomeString = "Бұл бот сіздің досыңыздың туған күнін енді ешқашан өткізіп алмауыңызды қамтамасыз етеді. Тізімге барлық туған күндерді қосыңыз, хабарландырулар орнатыңыз және тойлау уақыты келгенде бот сізге пинг береді!"

func (b *Bot) Start(db *sql.DB, conf *configs.Configuration) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	b.conf = conf

	updates := b.bot.GetUpdatesChan(u)
	var uId int

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				cmdText := update.Message.Command()
				if cmdText == "start" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Salem "+update.Message.Chat.FirstName+"! "+welcomeString)
					msg.ReplyMarkup = mainMenu
					b.bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Өкінішке орай мен start-тан басқа команда білмеймін :(")
					msg.ReplyMarkup = mainMenu
					b.bot.Send(msg)
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Менен не керек түсінбедім 🤷")

				switch update.Message.Text {

				case mainMenu.Keyboard[0][0].Text:
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ол адамның атын енгізіңіз: ")
					uId, _ = queries.InsertIntoBirthdayList(db, "", "1996-12-02", 0)
					msg.ReplyMarkup = backToMainMenu
					b.bot.Send(msg)
				case mainMenu.Keyboard[0][1].Text:
					objs, err := queries.GetBirthdayList(db)
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
						b.bot.Send(msg)
					} else {
						if len(objs) > 0 {
							for _, obj := range objs {
								text := fmt.Sprintf("*%s* - *%s*\n", obj.BirthDate, obj.Name)
								keyboard := tgbotapi.InlineKeyboardMarkup{}
								var row []tgbotapi.InlineKeyboardButton
								btn := tgbotapi.NewInlineKeyboardButtonData("Жою", obj.Name)
								row = append(row, btn)
								keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
								msg.ReplyMarkup = keyboard
								msg.ParseMode = "markdown"
								b.bot.Send(msg)
							}
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Әзірге тізім бос...")
							msg.ReplyMarkup = backToMainMenu
							b.bot.Send(msg)
						}
					}
				case mainMenu.Keyboard[1][0].Text:
					objs, err := queries.GetBirthdayList(db)
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
						b.bot.Send(msg)
					} else {
						if len(objs) > 0 {
							var text string
							for _, obj := range objs {
								text += fmt.Sprintf("*%s* - *%s*\n", obj.BirthDate, obj.Name)
								msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
							}
							msg.ReplyMarkup = backToMainMenu
							msg.ParseMode = "markdown"
							b.bot.Send(msg)
						} else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Әзірге тізім бос...")
							msg.ReplyMarkup = backToMainMenu
							b.bot.Send(msg)
						}

					}
				case backToMainMenu.Keyboard[0][0].Text:
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "So what..")
					msg.ReplyMarkup = mainMenu
					b.bot.Send(msg)
				default:
					//b.bot.Send(msg)
					userStatus, err := queries.GetUserStatus(db, uId)
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
						b.bot.Send(msg)
					} else {
						switch userStatus {
						case 0:
							name := update.Message.Text
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Адамның туған күнін енгізіңіз. \nФормат осындай DD/MM/YYYY болуы керек.")
							err = queries.UpdateName(db, uId, name)
							if err != nil {
								msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
								b.bot.Send(msg)
							}
							err := queries.UpdateUserStatus(db, uId, 1)
							if err != nil {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
								b.bot.Send(msg)
							}
							msg.ReplyMarkup = backToMainMenu
							b.bot.Send(msg)

						case 1:
							birthDate := update.Message.Text
							_, err := time.Parse("01/02/2006", birthDate)
							if err != nil {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Формат дұрыс емес енгізілді.")
								msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
								b.bot.Send(msg)
								err := queries.UpdateUserStatus(db, uId, 1)
								if err != nil {
									msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
									b.bot.Send(msg)
								}
							} else {
								err := queries.UpdateBirthDate(db, uId, birthDate)
								if err != nil {
									msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
									b.bot.Send(msg)
								}
								err = queries.UpdateUserStatus(db, uId, 2)
								if err != nil {
									msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
									b.bot.Send(msg)
								}
								msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Сәтті түрде енгізілді! 🎉")
								msg.ReplyMarkup = mainMenu
								b.bot.Send(msg)
							}

						}

					}
				}

			}
		} else if update.CallbackQuery != nil {
			fmt.Println(update.CallbackQuery.Data)
			isDeleted, err := queries.DeleteFromBirthdayList(db, update.CallbackQuery.Data)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Database error.")
				b.bot.Send(msg)
			} else {
				if isDeleted {
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Сәтті түрде жойылды.")
					msg.ReplyMarkup = mainMenu
					b.bot.Send(msg)
				}
			}
		}
	}
	return nil

}

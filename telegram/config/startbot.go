package config

import (
	"github.com/MamushevArup/telegram-bot-krisha/internal/texts"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Bot struct {
	TgBot *tgbotapi.BotAPI
}
type Info struct {
	Username   string   `bson:"telegram_tag"`
	Region     []string `bson:"region"`
	StartPrice int      `bson:"start_price"`
	EndPrice   int      `bson:"end_price"`
	IsOwner    string   `bson:"is_owner"`
	URL        string   `bson:"url"`
}

func New(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{TgBot: bot}
}

type Client struct {
	UserID int64 `json:"user_id"`
}

func (c *Client) СhatID(val int64) *Client {
	return &Client{UserID: val}
}
func (c Client) GetID() int64 {
	return c.UserID
}

var Chats int64

func (b *Bot) Start() (*Info, error) {

	var username string
	var regions []string
	var result string
	log.Printf("Authorized on account %s", b.TgBot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := b.TgBot.GetUpdatesChan(u)
	if err != nil {
		log.Println("Error when get updates from channel", err.Error())
		return &Info{}, err
	}
	start := 0
	for update := range updates {
		var c *Client
		c = c.СhatID(update.Message.Chat.ID)
		Chats = c.GetID()
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message == nil {
			continue
		}
		if start == 0 && !update.Message.IsCommand() {
			msg.Text = "Для начала работы введи команду -> \n/start"
			continue
		} else if start == 0 && update.Message.Command() != texts.StartCommand() {
			msg.Text = "Пожалуйста начните с команды /start"
		} else if update.Message.Command() == texts.StartCommand() {
			msg.Text = texts.HelloMessage()
			start = 1
		} else if start == 1 && update.Message.Command() == texts.RegionCommand() {
			msg.Text = "Выберите один или несколько районов разделите их запятой (,)\n" + strings.Join(texts.AllRegions, "\n")
			start = 2
		} else if start == 2 {
			texts.UsersCity = update.Message.Text
			texts.UsersCity = strings.ReplaceAll(texts.UsersCity, " ", "")
			regions = strings.Split(texts.UsersCity, ",")
			sort.Strings(regions)
			l := 0
			i := 0
			allCity := texts.AllRegionsFunc()
			for l < len(regions) {
				if regions[l] == allCity[i] {
					msg.Text = "Отлично!\nВведите пожалуйста начальную цену\nДля разделения используйте пробел"
					l++
				}
				i++
			}
			start = 3
		} else if start == 3 {
			var converter string
			converter = update.Message.Text
			converter = strings.ReplaceAll(converter, " ", "")
			texts.StartP, err = strconv.Atoi(converter)
			if err != nil {
				log.Fatal("Error while convert price to int", err, converter)
			}
			msg.Text = "Пожалуйста введие максимальную цену\nВ случае не надобности максимальной цены введите слово <Нет>"
			start = 4
		} else if start == 4 {
			var converter string
			converter = update.Message.Text
			converter = strings.ReplaceAll(converter, " ", "")
			if strings.EqualFold("Нет", converter) {
				texts.EndP = math.MaxInt32
			} else {
				texts.EndP, err = strconv.Atoi(converter)
				if err != nil {
					log.Fatal("Error while get final price", err)
				}
				if texts.EndP < texts.StartP {
					msg.Text = "Максимальная цена не может быть меньше минимальной повторите попытку"
					start = 4
					continue
				}
			}
			msg.Text = "Отлично!\nТеперь выберите один вариант\n1) Крыша Агент\n2) Специалист\n3) Хозяин недвижимость\n4) Все выше перечисленное"
			start = 5
		} else if start == 5 {
			value, err := strconv.Atoi(update.Message.Text)
			if err != nil {
				log.Fatal("Error while convert to int", err)
			}
			if value == 1 {
				result = "Крыша Агент"
			} else if value == 2 {
				result = "Специалист"
			} else if value == 3 {
				result = "Хозяин недвижимости"
			} else if value == 4 {
				result = "All"
			} else {
				msg.Text = "Пожалуйста введите число от 1 до 4"
				start = 5
			}
			start = 6
		}
		username = update.Message.From.UserName
		b.TgBot.Send(msg)
		if start == 6 {
			msg.Text = "Отлично теперь на основе информации мы будем выдавать вам ссылки на квартиры"
			break
		}
	}
	return &Info{Username: username, Region: regions, StartPrice: texts.StartP, EndPrice: texts.EndP, IsOwner: result, URL: ""}, nil
}
func NewInfo(tag string, region []string, startPrice, endPrice int, isOwner string, URL string) *Info {
	return &Info{Username: tag, Region: region, StartPrice: startPrice, EndPrice: endPrice, IsOwner: isOwner, URL: URL}
}

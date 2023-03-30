package scraping

import (
	"fmt"
	"github.com/MamushevArup/telegram-bot-krisha/telegram/config"
	"github.com/MamushevArup/telegram-bot-krisha/telegram/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
)

type Houses struct {
	Region string `bson:"region"`
	Intro  string
	URL    string `bson:"url"`
	Price  string `bson:"price"`
	Owner  string
	Desc   string
}

func ScrapKrisha(b *tgbotapi.BotAPI, i *config.Info) {
	tick := time.NewTicker(500 * time.Millisecond)
	for range tick.C {
		c := colly.NewCollector()
		var h Houses
		// Find and visit all links
		c.OnHTML("div[class=a-card__descr]", func(e *colly.HTMLElement) {
			intro := e.ChildText("div.a-card__text-preview")
			textLink := e.ChildText("a.a-card__title")
			regions := e.ChildText("div.a-card__subtitle")
			owner := e.ChildText("div.owners__label")
			urls := "https://krisha.kz" + e.ChildAttr("a[href]", "href")
			price := e.ChildText("div.a-card__price")
			h.Region = regions
			h.URL = urls
			h.Price = price
			h.Intro = intro
			h.Owner = owner
			h.Desc = textLink
		})
		c.Visit("https://krisha.kz/prodazha/kvartiry/almaty/")
		if h.URL == i.URL {
			continue
		}
		fmt.Println(h.Owner, "===========================")
		i.URL = h.URL
		booleanRegion, booleanPrice, booleanOwner := validateRegion(i, h), validatePrice(i, h), validateOwner(i, h)
		if booleanRegion && booleanPrice && booleanOwner {
			err := database.UpdateURL(i.Username, h.URL)
			if err != nil {
				log.Println("Error while updating the URLS", err)
			}
			msg := tgbotapi.NewMessage(config.Chats, h.URL+"\n"+h.Region+"\n"+h.Price+"\n"+h.Desc+"\n"+h.Owner)
			_, err = b.Send(msg)
			if err != nil {
				log.Println("Err", err, "923329329932392329932239")
			}
		}
	}
}
func validateRegion(i *config.Info, h Houses) bool {
	fmt.Println(h.Region, "+++++++++++++++++++++++++++++")
	UsersPrefer := i.Region
	var houseUser string
	for _, v := range h.Region {
		if v == ' ' {
			break
		}
		houseUser += string(v)
	}
	res := houseUser
	fmt.Println(res, "_---------------------------------")
	for j := 0; j < len(UsersPrefer); j++ {
		if res == UsersPrefer[j] {
			return true
		}
	}
	return false
}
func validatePrice(i *config.Info, h Houses) bool {
	usersStartPrice := i.StartPrice
	usersEndPrice := i.EndPrice
	fromKrisha := h.Price
	fromKrisha = strings.ReplaceAll(fromKrisha, "〒", "")
	fromKrisha = strings.ReplaceAll(fromKrisha, "от", "")
	fromKrisha = strings.ReplaceAll(fromKrisha, "\u00a0", "")
	fromKrisha = strings.ReplaceAll(fromKrisha, " ", "")
	var result int
	result, err := strconv.Atoi(fromKrisha)
	if err != nil {
		log.Println("Error while converting to the int", err, result)
		h.Owner = "Новостройка"
		return true
	}
	if usersEndPrice > result && usersStartPrice < result {
		return true
	}
	return false
}
func validateOwner(i *config.Info, h Houses) bool {
	if i.IsOwner == "All" {
		return true
	}
	if h.Owner == " " {
		return true
	}
	if i.IsOwner == h.URL {
		return true
	}
	return false
}

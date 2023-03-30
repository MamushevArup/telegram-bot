package texts

type Config struct {
	Token string `json:"telegram_bot_token"`
}

var PathToJSON = "./internal/config.json"
var AllRegions = []string{"Алатауский", "Алмалинский", "Ауэзовский", "Бостандыкский", "Жетысуский", "Медеуский", "Наурызбайский", "Турксибский"}

func AllRegionsFunc() []string {
	return []string{"Алатауский", "Алмалинский", "Ауэзовский", "Бостандыкский", "Жетысуский", "Медеуский", "Наурызбайский", "Турксибский"}
}

var StartCommand = func() string {
	return "start"
}
var RegionCommand = func() string {
	return "region"
}
var UsersCity string
var StartP int
var EndP int

var HelloMessage = func() string {
	return "Приветствуем вас в телеграм боте. Перед началом необходимо поделиться информацией о некоторых категориях\n" +
		"/region выбор района"
}

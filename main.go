package main

import (
	"flag"
	"log"

	tgClient "tg-bot-training/clients/telegram" // че за странная запись первый раз такое вижу
	"tg-bot-training/consumer/event-consumer"
	"tg-bot-training/events/telegram"
	"tg-bot-training/storage/files"
)

/*
с
каркас проекта
token = flags.Get(token)-это вроде бы готово
tgClient = telegram.New(token)					//этот готов
fetcher = fetcher.New(tgClient)					//эти на подходе
processor = processor.New(tgClient)				//эти на подходе
consumer.Start(fetcher, processor)
*/
const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	//новый раздел был разблокирован //tgBotHost="api.telegram.org"
	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("server started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize) //в качестве Фетчера передаем eventsProcessor, в качестве процессора передаем eventsProcessor
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stoped", err) //напишет сообщение об ошибке и остановит программу
	}
}

// достает флаг токена который по сути адрес на токен который лежит на компе
func mustToken() string {
	token := flag.String( //туда сохраняется не само значение с ссылка на значение //и берется оно из флага хз как
		"tg-bot-token",                     //имя флага, во ввремя запуска программы надо указать его имя
		"",                                 //дефолтное значение флага, если флаг не указан
		"token for access to telegram bot", //подсказка к флагу, увидим после компиляции программы
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified") //вызываем панику и аварийно завершаем, программу если пустое значение в токене
	}
	return *token
}

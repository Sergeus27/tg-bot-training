package consumer

import (
	"log"
	"tg-bot-training/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize, //сколько событий обрабатываем за раз
	}
}

func (c Consumer) Start() error {
	for { //вечный цикл будет постоянно ждать новые события и обрабатывать их //ничего нового уже не используем, используем все что написали раннее
		gotEvents, err := c.fetcher.Fetch(c.batchSize) //получим события с помощью фетчера
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue //continue выйдет из if или из цикла
		}

		if len(gotEvents) == 0 { //добавляем задержку в 1 сек. если не получили ивентов
			time.Sleep(1 * time.Second)

			continue
		}
		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

/*
проблемы этой функции:
 1. потеря необработанных из-за ошибок событий событий
    ретраи, возвращение в хранилища, фоллбек, подтверждение Fetcher
 2. обработка всей пачки
    остановка после первой ошибки, вести счетчик остановка после количества
 3. паралельная обработка
    sinc.WateGroup{}  // почитать можно сделать на такой херне, ДЗ типо
*/
func (c *Consumer) handleEvents(events []events.Event) error { //событий может быть несколько по жтому выносим код в отдельную функцию
	for _, event := range events { //перебираем события
		log.Printf("got new events: %s", event.Text) //отчитались когда получили событие

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event:%s", err.Error()) //отчитались когда обработали

			continue //прикол для других разработчиков, которые заходят после этого ифа че то свое дописать еще
		}
	}
	return nil
}

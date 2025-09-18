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
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new events: %s", event.Text)
	}
}

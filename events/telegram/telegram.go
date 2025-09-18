package telegram

import (
	"tg-bot-training/client/telegram"
	"tg-bot-training/events"
	"tg-bot-training/lib/e"
	"tg-bot-training/storage"
)

type Processor struct { //реализовывать оба интерфеса будет единственный тип данных Processor
	tg      *telegram.Client //надо клиент, делали совсем недавно в 3 ролике
	offset  int              //внутренний параметр offset которым будет пользоваться самостоятельно
	storage storage.Storage  //в storage будут сохраняться ссылки //тут мы используем абстрактный интерфейс Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknown event type") //переменные с ошибкой создаем здесь видимо чтобы нас сюда выкинуло при ошибке, мы увидели ее в части кода???
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor { //создавать процессор будет
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

// разница updates и events, updates-термин телеграма относятся только к нему, в других месенджерах возмоно и нет, events-более общая сущность можем преобразовывать все что получаем от других месенджеров в независимости от формата предоставляемых от месенджеров данных
func (p *Processor) Fetch(limit int) ([]events.Event, error) { //Fetcher Processor у нас есть в events>types// возвращаем слайс из ивентов
	updates, err := p.tg.Updates(p.offset, limit) //с помощью клиента нужно получить все апдейты, offset внутренний
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil //возвращаем нулевой результат если не нашли апдейтов
	}

	res := make([]events.Event, 0, len(updates)) //алацируем память под результат, т.к известно сколько будет значений
	for _, u := range updates {                  //обходим апдейты превращая их в ивенты
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1 //обновляем поле offset чтобы при следующем запросе получить апдейты у которых ID больше чем у последнего полученного

	return res, nil
}

func (p *Processor) Process(event events.Event) error { //метод будет выполнять различные действия в зависимости от типа ивента
	switch event.Type {
	case events.Message: //когда работаем с сообщением
		return p.processMessage(event)
	default: //когда не знаем с чем работаем
		return e.Wrap("can't process message", ErrUnknownEventType)
		//если надо будет работать с другими апдейтами от тг просто добавь case
	}
}

func (p *Processor) processMessage(event events.Event) {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) {
	res, ok := event.Meta(Meta) //для поля meta мы попытаемся сделать type assersion	 если здесь будет что-то другое то вторым параметром "ok" придется false
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
}

func event(upd telegram.Update) events.Event { //функция для преобразования объектов в ивенты
	updType := fetchType(upd) //fetchType нужен еще кое-где по этому вынесли в переменную

	res := events.Event{
		Type: updType,        //создаем 2 функции 1 получает из объекта тип...
		Type: fetchText(upd), //...а другая текст
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil { //Message может быть нулевым
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

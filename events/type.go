package events //абстрактный пакет events не привязанный ни к 1 месенджеру

type Fetcher interface {
	Fetch(limit int) ([]Event, error) //хотели прикрутить параметр offset но у нас будет универсальный интерфейс, если его использовать не для tg не факт что там будет offset, offset мы уже реализуем внутри Fetcher
}

type Processor interface {
	Process(e Event) error
}

type Type int //наш кастомный тип, нахрена не знаю, чтобы не запутаться говорит, браво

const (
	Unknown Type = iota //iota-счетчик для констант, первой присваивает 0
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{} //в интерфейс можно положить все что угодно, Events не будет понимать что здесь будет находиться
}

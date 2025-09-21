package telegram //комманды которые будет понимать телеграм бот
import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
	"tg-bot-training/lib/e"
	"tg-bot-training/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error { //основная функция является методом типа Processor, это будет что-то типа API роутера //будем смотреть на текст сообщения и по его формату и содержанию будем понимать какая это команда
	text = strings.TrimSpace(text) //удаляем лишние пробелы
	log.Printf("got new command '%s' from '%s", text, username)
	// add page: 	https://...				сохраниение ссылки
	// rnd page: 	/rnd					отправка рандомной страницы
	// help:		/help 					вывод информации как работать с ботом
	// start:		/start: hi+help			автоматичкски будет выполняться когда пользователь напишет первое сообщение боту, + выводит help

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text { //комманжы с ключевыми словами будем определьть с помощью switch
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExist(context.Background(), page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool { //проверяет является ли сообщение ссылкой
	return isURL(text)
}

func isURL(text string) bool { //проверяет что текст является ссылкой, зачем нам 2 абсолютоно одинаковые функции я хз: проверку можно выполнять различными способами, и поменять способ найдя нужную функция(как по мне это лишнее усложнение)
	u, err := url.Parse(text) //распарсим текст считая его ссылкой функцией Parse пакета url, если ошибка нудевая то текст=ссылка и при этом указан хост. Тут есть минус: ссылка формата google.com за ссылки восприниматься не будут, но мне пофиг. В ссылке обязательно должно быть http в начале

	return err == nil && u.Host != ""
}

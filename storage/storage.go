package storage //походу тут тоже все на интерфейсах будет
import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"tg-bot-training/lib/e"
)

type Storage interface {
	//создаем методы
	Save(ctx context.Context, p *Page) error                        //принимает страницу на вход чтобы сохранить //передаем страницу по ссылке так как тип может расширяться(помимо созранения ссылки бот сможет переходить по ним, сохранял превью или всю статью целиком например)
	PickRandom(ctx context.Context, userName string) (*Page, error) //принимает имя пользователя возвращает рандомную страницу для него
	Remove(ctx context.Context, p *Page) error
	IsExist(ctx context.Context, p *Page) (bool, error) //сообщает существует ли такая страница
}

var ErrNoSavedPages = errors.New("no saved pages") //вынесли ошибку отдельно в переменную пакета чтобы ее можно было проверить снаружи//перенесли ошибку в storage из files, потому что она более общая

type Page struct { //основной тип данных с которым будет работоть Storage
	URL      string //ссылка на страницу которую мы скинули боту
	UserName string //имя пользователя который скинул ссылку
	//Created time.Time //если вместо случайных страниц захотим скидывать самые старые или самые новые то можно это будет использовать
}

func (p Page) Hash() (string, error) { //будем генерировать хэш с помощью алгоритма sha1 уже есть библиотека //Хэш нам нужен для уникальности, вот только чего?
	h := sha1.New() //возвращает Hash который интерфейс и в него встроен интерфейс Writer, чтобы пользоваться просто передай через Write инфу

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil //итоговый хеш будет создан методом суммации, nil пихаем потому что уже записали все что нужно в хеш "h", Sum возвращает байты по этому используем Sprintf чтобы переобразовать их в текст //мы складываем не только адрес ссылки но и имя пользователя так как у нас все может содержаться в 1 папке, а разные пользователи в теории могут загружать одинаковые ссылки
}

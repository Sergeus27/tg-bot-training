package files

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
)

type Storage struct {
	basePath string //хранит информацию в какой папке мы это все это будет хранить, не показывайте мой код филологам. P.S. сеньорам тоже
}

const defaultPerm = 0774 //дает всем пользователям права на чтение и

func New(basePath string) Storage {
	return Storage{basePath: basePath} // вот тут я вообще нихрена не понял

}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }() //определяем способ обработки ошибок

	fPath := filepath.Join(s.basePath, page.UserName) //путь до дирректории куда будем сохранять наш файл //filepath то же что и path но на windows разделяет обратным слешом '\' //а сохраняться все будет в папку UserName ну понятно

	if err := os.MkdirAll(fPath, defaultPerm); err != nil { //Mkdir создаст все дирректории которые входят в переданный путь + defaultPerm(параметры доступа 8-ричн СС)
		return err //благодаря defer просто возвращаем err уже после обертки
	}

	fName, err := fileName(page) //сюда имя файла формируем

	if err != nil {
		return err
	}
	fPath = filepath.Join(fPath, fName) //полный путь до файла(дописываем имя файла к пути)

	file, err := os.Create(fPath) //создаем файл
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }() //можно закрыть файл только в функции если она в defer

	if err := gob.NewEncoder(file).Encode(page); err != nil { //записываем в файл страницу в нужном формате //сериализуем страницу-привести к формату который можно записать в файл и по нему восстановить исходную структуру
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random", err) }()

	fPath := filepath.Join(s.basePath, userName) //получаем путь до дирктории с файлами

	files, err := os.ReadDir(path) //получаем список файлов

	if err != nil {
		return nil, err
	}

	if len(files) == 0 { //ноль файлов? возвращай ошибку
		return nil, storage.ErrNoSavedPages
	}

	// 0-№ last file
	rand.Seed(time.Now().UnixNano()) //псевдорандом, с динамическим семнем(время)
	n := rand.Intn(len(files))       //рандомное число от 0 до порядкового нормера последнего файла
	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name())) //в качестве результата возвращаем вызов этой функции
}

func (s Storage) Remove(p *Storage.Page) error { //метод Remove так же передается странице, все что он будет делать мы уже видели в прошлых методах
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)
	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)
		return e.Wrap(msg, err) //в ошибку добавим путь до файла чтобы знать на каком сломалось
	}

	return nil
}
func (s Storage) IsExist(p *storage.Page) (bool, error) { //метод проверяет существует страница или нет сохранял ли ее польщователь раннее
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); { //возвращает множественные ошибки и только ненаход файла будет означать 0 ошибку
	case error.Is(err, os.ErrNotExist): //ХЗ ЧЕ ТУТ С КЕЙСАМИ НАДО БУДЕТ РАЗОБРАТЬСЯ ПОЧЕМУ FALSE ЕСЛИ НЕ НАШЕЛ
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)
		return false, e.Wrap(msg, err)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) { //осталось декодировать файл и вернуть его содержимое, надо открыть файл и вызвать для него декодер
	f, err := os.Open(filePath) //открыли файл
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }() //закрываем файл

	var p storage.Page //переменная в которую файл будет декодирован

	if err := gob.NewDecoder(f).Decode(&p); err != nil { //осуществляем декодирование файла с помощью пакета gob //&p-ссылка на страницу передаем ее в метод Decode
		return nil, e.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storagePage) (string, error) { //функция для определения имени //для того чтобы названия папок были уникальны используем хеш URL+Username
	return p.Hash() //в функции только 1 строка чтобы код был гибким
}

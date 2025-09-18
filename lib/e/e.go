package e //мы создали целую библиотеку для оборачивания ошибок, т.к они у нас часто могут поторяться

import "fmt"

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err) //msg-текст с подсказкой, err-сама ошибка
}

func WrapIfErr(msg string, err error) error { //эту более длинную функцию мы используем когда у нас может и не быть никакой ошибки
	if err == nil {
		return nil
	}
	return Wrap(msg, err)
}

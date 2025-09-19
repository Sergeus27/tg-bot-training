package telegram //клиент будет общаться с телеграмом

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"tg-bot-training/lib/e"
)

type Client struct {
	host     string      //хост API сервиса телеграма tg-bot.com/...
	basePath string      //базовый путь .../bot<token> только без скобок "<>" токен потом получим вместе с хостом будет путь: tg-bot.com/bot<token>
	client   http.Client //храним еще http клиент чтобы не создавать его отдельно для каждого запроса
}

const (
	getUpdatesMethod  = "getUpdates"
	SendMessageMethod = "sendMessage"
)

func New(host string, token string) *Client { //будет создавать клиентов
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{}, //здесь уже нам изветные поля, только токен решили создавать отдельной функцией
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) { //это метод для получения апдейтов т. новых сообщений     //не ебу что за запись такая c *Client это че тип данных?
	q := url.Values{}                     //а вот какая полезная инфа там будет находиться  https://core.telegram.org/bots/api#getting-updates //для получения запросов мы будет отправлять запрос getUpdates
	q.Add("offset", strconv.Itoa(offset)) //q(query) это параметры запроса //offset(смещение) это с какого запроса нам API рассылает апдейты
	q.Add("limit", strconv.Itoa(limit))   //limit количество апдейтов которые мы получаем за 1 запрос

	data, err := c.doRequest(getUpdatesMethod, q) //getUpdatesMethod="getUpdates"
	if err != nil {
		return nil, err
	}

	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil { //распаршиваем json ну и указываем ссылку на значения куда результат
		return nil, err
	}

	return res.Result, nil //возвращает некую структуру которая будет содержать все что нам нужно знать об апдейте//не ебу что это значит
}

func (c *Client) SendMessage(chatID int, text string) error { //метод для отправки сообщения пользователю //мы перенесли все экспортируемые методы наверх чтобы их было проще искать людям
	q := url.Values{} //подготовим параметры запроса, делаем это уже привычно? КАК?
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	//do request <- getUpdates
	_, err := c.doRequest(SendMessageMethod, q) //опять константа SendMessageMethod=SendMessage
	if err != nil {
		return e.Wrap("can't send message", err) //тут Wrap потому что ошибка не нулевая
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) { //отправляет запрос query это q, method
	const errMsg = "can't do request"

	defer func() { err = e.WrapIfErr(errMsg, err) }()

	u := url.URL{ //url на который будет отправляться запрос
		Scheme: "https",                       //протокол
		Host:   c.host,                        //хост из клиента
		Path:   path.Join(c.basePath, method), //путь из двух частей базова+метод//method берется из документации, например для получения апдейтов метод называется getUpdates//Join просто склеивает 2 части рассатвляя / между ними
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil) //формируем/подготавливаем объект запроса //http.MethodGet="GET" что? Зачем? непонятно //u.String() по сути url  в текстовом виде //body пустое потому что у нас уже все есть в параметрах
	if err != nil {
		return nil, err //тут мы оборачивали ошибку так как при логировании из за того что она находится в пакете http мы не поймем в каком месте программы ошибка возникла
	}

	if err != nil {
		//error.Is() и error.As() как то связаны с fmt.Errorf()но я пока не могу понять как
		return nil, e.Wrap(errMsg, err)
	}
	req.URL.RawQuery = query.Encode() //передаем в объект request параметры которые получили из аргумента, encode приводит параметры к такому виду который мы можем отправить на сервер

	resp, err := c.client.Do(req) // для отправки мы используем тот клиент который заранее подготовили
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	defer func() { _ = resp.Body.Close() }() //в конце закрываем тело запроса игнорируем ошибку

	body, err := io.ReadAll(resp.Body) //получим содержимое

	if err != nil {
		return nil, err
	}
	return body, nil //возвращаем результат
}

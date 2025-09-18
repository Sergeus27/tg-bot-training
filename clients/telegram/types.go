package telegram //тут определяем все типы с которыми будет работать наш клиент

type UpdateResponse struct { //это гавно написано именно так потому что контракт в телеграме такой, так надо писать а думать не надо
	Ok     bool     `json:"ok"`     //но тут не все так просто на самом деле от функции getUpdates мы получаем не голые апдейты а структуру, а некоторый объект ответов где апдейт будет содержаться кроме прочей информации "ok bool", а наши объекты при ok=True будут содержаться в result
	Result []Update `json:"result"` //это все
}

type Update struct { //это гавно написано именно так потому что контракт в телеграме такой, так надо писать а думать не надо
	ID      int              `json:"update_id"` //https://core.telegram.org/bots/api#getting-updates вот тут в json теге и написано что нам приходит в апдейте
	Message *IcommingMessage `json:"message"`   // здесь булет не string а отдельная струкура, смотри ниже //поле message может отстутвовать по этому оставили его ссылкой на структуру
}

//https://core.telegram.org/bots/api#getting-updates //message это объект и тоже состоит из множетсва полей from-от кого кришло сообщение chat-как отправить сообщение обратно text-из него будем получать команжы и ссылки
type IcommingMessage struct { //входяшие сообщения
	Text string `json:"text"` //text-из него будем получать команжы и ссылки
	From From   `json:"from"` //chat-как отправить сообщение обратно //это так же объекты их тоже разложем ниже
	Chat Chat   `json:"chat"` //from-от кого кришло сообщение 		//это так же объекты их тоже разложем ниже
}

type From struct {
	Username string `json:"username"` //из поля From, нам потребуется только имя пользователя
}

type Chat struct {
	ID int `json:"id"` //потребуется только id чата
}

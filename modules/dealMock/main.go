package dealMock

import (
	"math/rand"
	"time"
)

type deal struct {
	expired bool
}

var process = []string{
	"Пока все отлично",
	"Ты в плюсе",
	"Красава",
	"Тебе везет!",
	"Все идет хорошо",
	"Ты в минусе",
	"Еще не все потеряно",
	"Ушел в минус",
	"Это нормально",
	"Все еще может измениться",
}

var finish = []string{
	"Так держать! Ты заработал 40р! Твой баланс %d",
	"К сожалению ты проиграл. кнопка \"попробовать еще раз\" кнопка \"открыть такую же сделку x2\"",
}

// Инициализация моков для сделок
func NewDeal() *deal {
	return &deal{expired:false}
}

func (d *deal) Start(timeOut int64)  {
	// c := make(chan bool, 1)

	go func(dl *deal) {
		timer1 := time.NewTimer(time.Second * time.Duration(timeOut))
		for {
			select {
			case <-timer1.C:
				dl.expired = true
			}
		}
	}(d)


}

// Процесс сделки, рандомные события, по завершении времени выдает рандомный резудьтат сделки
func (d *deal) Process() (string, bool) {
	if d.expired {
		return finish[rand.Intn(len(finish))], true
	}
	return process[rand.Intn(len(process))], false
}
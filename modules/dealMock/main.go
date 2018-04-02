package dealMock

import (
	"math/rand"
	"time"
	"fmt"
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
	timer1 := time.NewTimer(time.Duration(timeOut))

	go func() {
		<-timer1.C
		fmt.Println("Expiret timer")
		d.expired = true
	}()
}

// Процесс сделки, рандомные события, по завершении времени выдает рандомный резудьтат сделки
func (d *deal) Process() string {
	if d.expired {
		return finish[rand.Intn(len(finish))]
	}
	return process[rand.Intn(len(process))]
}
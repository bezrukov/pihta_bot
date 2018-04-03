package dealMock

import (
	"math/rand"
	"time"
	"pihta_bot/modules/account"
)

type deal struct {
	expired bool
}

type finishResult struct {
	msg string
	isWin bool
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

var finish = []finishResult{
	{msg: "Так держать! Ты заработал 40р!", isWin: true},
    {msg:"К сожалению ты проиграл.", isWin: false},
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
func (d *deal) Process(balance *account.Balance, dealCost int) (string, bool) {
	if d.expired {
		result := finish[rand.Intn(len(finish))]

		if result.isWin {
			balance.Inc(dealCost)
		} else {
			balance.Dec(dealCost)
		}

		return result.msg, true
	}
	return process[rand.Intn(len(process))], false
}
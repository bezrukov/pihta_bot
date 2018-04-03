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
	"\xF0\x9F\x92\x94 Ты в минусе",
	"\xE2\x9C\x85 Ты в плюсе",
}

var finish = []finishResult{
	{msg: "Так держать! Ты заработал 40р! \xF0\x9F\x8E\x89\xF0\x9F\x8E\x89\xF0\x9F\x8E\x89", isWin: true},
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
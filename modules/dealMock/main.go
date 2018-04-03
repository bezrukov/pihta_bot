package dealMock

import (
	"math/rand"
	"time"
	"github.com/bezrukov/pihta_bot/modules/account"
)

type deal struct {
	expired bool
}

type finishResult struct {
	msg string
	isWin bool
}

var process = []string{
	"\xE2\x9C\x85 Ты в плюсе",
	"\xE2\x9C\x85 Пока все отлично",
	"\xE2\x9C\x85 Красава",
	"\xE2\x9C\x85 Тебе везет!",
	"\xE2\x9C\x85 Все идет хорошо",

	"\xF0\x9F\x92\x94 Еще не все потеряно",
	"\xF0\x9F\x92\x94 Ушел в минус",
	"\xF0\x9F\x92\x94 Все еще может измениться",
	"\xF0\x9F\x92\x94 Потерпи чуть-чуть",
}

var finish = []finishResult{
	{msg: "\xF0\x9F\x8E\x89\xF0\x9F\x8E\x89\xF0\x9F\x8E\x89 Так держать! Ты заработал 40р!", isWin: true},
    {msg:"\xF0\x9F\x98\xBF К сожалению ты проиграл.", isWin: false},
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
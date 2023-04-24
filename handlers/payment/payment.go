package payment

import (
	"fmt"
	"sync"
	"time"
)

type credit float64

type Account struct {
	AccountID    string
	Balance      credit
	CreatedAt    time.Time
	LastModified time.Time
}

var (
	mutex sync.RWMutex
)

func makePayment(a Account, amount credit) {
	ch := make(chan credit)
	go a.addCredit(amount)
	ch <- a.Balance
	v := <-ch
	fmt.Print("new balance is", v)
}

func (a *Account) addCredit(c credit) {
	mutex.Lock()
	defer mutex.Unlock()
	a.Balance += c
	a.LastModified = time.Now()
}

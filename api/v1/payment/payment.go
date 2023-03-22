package payment

import (
	"fmt"
	"sync"
)

type credit float64

var (
	account credit
	mutex   sync.Mutex
)

func makePayment(account, amount credit) {
	ch := make(chan credit)
	go account.addCredit(amount)
	ch <- account
	v := <-ch
	fmt.Print("new balance is", v)
}

func (a credit) addCredit(value credit) {
	mutex.Lock()
	defer mutex.Unlock()
	a += value
}

// based on this: https://dev.to/joashxu/go-limit-service-time-per-user-pc
package main

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	AllowedSecondsToProcess int64
}

type User struct {
	ID        string
	IsPremium bool
	TimeUsed  int64
	Mx        *sync.Mutex
}

func (s *Service) HandleRequest(process func(), u *User) bool {
	done := make(chan bool)

	go func() {
		process()
		done <- true
	}()

	for {
		select {
		case <-done:
			return true
		case <-time.Tick(time.Second * 1):
			u.Mx.Lock()
			u.TimeUsed++

			if !u.IsPremium && u.TimeUsed > s.AllowedSecondsToProcess {
				u.Mx.Unlock()
				return false
			}

			u.Mx.Unlock()
		}
	}
}

func sampleProcess(seconds int64) {
	start := time.Now()
	time.Sleep(time.Duration(seconds) * time.Second)
	log.Printf("Process finished after: %v", time.Since(start))
}

func sampleCustomizedProcess() {
	sampleProcess(5)
}

func main() {
	srvc := &Service{
		AllowedSecondsToProcess: 10,
	}
	user := &User{
		ID:        uuid.NewString(),
		IsPremium: false,
		Mx:        &sync.Mutex{},
	}

	successful := srvc.HandleRequest(sampleCustomizedProcess, user)
	log.Printf("finished short process with success: %v, user premium: %v", successful, user.IsPremium)
	successful = srvc.HandleRequest(sampleCustomizedProcess, user)
	log.Printf("finished short process with success: %v, user premium: %v", successful, user.IsPremium)
	successful = srvc.HandleRequest(sampleCustomizedProcess, user)
	log.Printf("finished short process with success: %v, user premium: %v", successful, user.IsPremium)

	user.IsPremium = true
	successful = srvc.HandleRequest(sampleCustomizedProcess, user)
	log.Printf("finished short process with success: %v, user premium: %v", successful, user.IsPremium)
	successful = srvc.HandleRequest(sampleCustomizedProcess, user)
	log.Printf("finished short process with success: %v, user premium: %v", successful, user.IsPremium)
	successful = srvc.HandleRequest(sampleCustomizedProcess, user)
	log.Printf("finished short process with success: %v, user premium: %v", successful, user.IsPremium)
	successful = srvc.HandleRequest(sampleCustomizedProcess, user)
	log.Printf("finished short process with success: %v, user premium: %v", successful, user.IsPremium)
}

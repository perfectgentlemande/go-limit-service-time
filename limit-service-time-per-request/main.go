// based on this: https://dev.to/joashxu/go-limit-service-time-per-request-49e1
package main

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	AllowedSecondsToProcess int64
}

type User struct {
	ID        string
	IsPremium bool
}

func (s *Service) HandleRequest(process func(), u *User) bool {
	done := make(chan bool)

	go func() {
		process()
		done <- true
	}()

	success := true
	select {
	case <-done:
		return success
	case <-time.After(time.Second * time.Duration(s.AllowedSecondsToProcess)):
		if !u.IsPremium {
			return false
		}
	}

	return success
}

func sampleProcess(seconds int64) {
	start := time.Now()
	time.Sleep(time.Duration(seconds) * time.Second)
	log.Printf("Process finished after: %v", time.Since(start))
}

func sampleLongProcess() {
	sampleProcess(15)
}
func sampleShortProcess() {
	sampleProcess(5)
}

func main() {
	srvc := &Service{
		AllowedSecondsToProcess: 10,
	}
	user := &User{
		ID:        uuid.NewString(),
		IsPremium: false,
	}

	successful := srvc.HandleRequest(sampleShortProcess, user)
	log.Printf("finished short process with success: %v, user premium: %v", successful, user.IsPremium)

	successful = srvc.HandleRequest(sampleLongProcess, user)
	log.Printf("finished long process with success: %v, user premium: %v", successful, user.IsPremium)

	user.IsPremium = true
	successful = srvc.HandleRequest(sampleShortProcess, user)
	log.Printf("finished short process with success: %v, user premium: %v", successful, user.IsPremium)

	successful = srvc.HandleRequest(sampleLongProcess, user)
	log.Printf("finished long process with success: %v, user premium: %v", successful, user.IsPremium)
}

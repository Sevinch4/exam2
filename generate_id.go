package check

import (
	"fmt"
	"sync"
	"time"
)

const (
	incrementStep = 1
	initialValue  = 11
)

func Generate(stop <-chan struct{}, wg *sync.WaitGroup) <-chan string {
	idCh := make(chan string)

	go func() {
		defer close(idCh)
		defer wg.Done()

		var mu sync.Mutex
		var prefix string = "T"
		var increasingID = initialValue

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				mu.Lock()
				increasingID += incrementStep
				id := fmt.Sprintf("%s-%d", prefix, increasingID)
				mu.Unlock()
				idCh <- id
			case <-stop:
				return
			}
		}
	}()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(<-idCh)
		}()
		time.Sleep(500 * time.Millisecond) // Sleep for half a second between generations
	}
	return idCh
}

func GenerateID() {
	var wg sync.WaitGroup

	stopCh := make(chan struct{})
	defer close(stopCh)

	// Create a function that generates an ID with auto-increment every second
	idCh := Generate(stopCh, &wg)

	// Example: Print 5 IDs without using an explicit for loop
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(<-idCh)
		}()
		time.Sleep(500 * time.Millisecond) // Sleep for half a second between generations
	}

	wg.Wait() // Wait for all goroutines to finish
}

//const (
//	increasingID = iota + 11
//)
//
//func GenerateID() func() string {
//	var mu sync.Mutex
//	var prefix string = "T"
//	return func() string {
//		mu.Lock()
//		defer mu.Unlock()
//
//		id := fmt.Sprintf("%s-%d", prefix, increasingID)
//		//increasingID += 1
//		return id
//	}
//}

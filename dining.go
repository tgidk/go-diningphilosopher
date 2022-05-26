package main

import (
	"fmt"
	"sync"
	"time"
)

type ChopStick struct {
	sync.Mutex
}

type Philosopher struct {
	id, eatingCount int
	leftCS, rightCS *ChopStick
}

func (p Philosopher) Eat(sg *sync.WaitGroup, ch chan int, numberOfMeals int) {
	defer sg.Done()

	// Each philosopher should eat only 3 times
	for i := 0; i < numberOfMeals; i++ {
		ch <- p.id // send request to get permission
		<-ch       // get response from host, but ignore value

		// get chop sticks
		p.leftCS.Mutex.Lock()
		p.rightCS.Lock()
		fmt.Printf("*starting to eat (meal %d of %d) %d  \n", i+1, numberOfMeals, p.id)

		time.Sleep(1 * time.Second) // eating...

		fmt.Printf("finishing eating (meal %d of %d) %d \n", i+1, numberOfMeals, p.id)
		p.leftCS.Unlock()
		p.rightCS.Unlock()
	}
}

func GetPhilosophers(numberOfPhilosophers int) []*Philosopher {
	chopSticks := make([]*ChopStick, numberOfPhilosophers)
	philosophers := make([]*Philosopher, numberOfPhilosophers)

	for i := 0; i < numberOfPhilosophers; i++ {
		chopSticks[i] = new(ChopStick)
	}
	for i := 0; i < numberOfPhilosophers; i++ {
		philosophers[i] = &Philosopher{i + 1, 0, chopSticks[i], chopSticks[(i+1)%numberOfPhilosophers]}
	}
	return philosophers
}

func Host(ch chan int) {
	// host allows no more than 2 philosophers to eat concurrently,
	// this is controlled by the channel capacity
	for x := range ch { // get request
		ch <- x // send response back/gives permission to eat
	}
}

func main() {
	const numberOfGuests int = 5
	const numberOfMeals int = 3
	const channelCapacity int = 2

	ch := make(chan int, channelCapacity)
	defer close(ch)

	sg := sync.WaitGroup{}
	sg.Add(numberOfGuests)
	philosophers := GetPhilosophers(numberOfGuests)

	go Host(ch) // start host task

	for i := 0; i < numberOfGuests; i++ {
		go philosophers[i].Eat(&sg, ch, numberOfMeals)
	}

	sg.Wait()
}

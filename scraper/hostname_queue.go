package scraper

import "sync"

type null struct{}

type HostnameQueue struct {
	hostname                string
	URLPathQueue            []string
	NumReqsPerMinuteAllowed int
	NumReqsSent             int
	Visited                 map[string]null
	Mu                      sync.Mutex
}

/*
	API for using slice as a queue

	queue := make([]int, 0)
	// Push to the queue
	queue = append(queue, 1)
	// Top (just get next element, don't remove it)
	x = queue[0]
	// Discard top element
	queue = queue[1:]
	// Is empty ?
	if len(queue) == 0 {
	fmt.Println("Queue is empty !")
	}
*/

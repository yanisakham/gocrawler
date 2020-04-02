package scraper

import "time"

type MasterQueue struct {
	HostnameQueueMap   map[string]*HostnameQueue
	HostnameQueueQueue []*HostnameQueue
	CurrentQueue       *HostnameQueue
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


func (mq *MasterQueue) nextURL() (string, bool) {
	// TODO
	return "", false
}

// Takes in a HostnameQueue that's already preprocessed and ready to
// replace mq.CurrentQueue
func (mq *MasterQueue) addNewHostname(queue *HostnameQueue) {
	// TODO

}

// Sets all HostnameQueue's numReqsSent to 0
func (mq *MasterQueue) ClearNumReqsSent() {
	// TODO
}

// Runs ClearNumReqsSent every minute,
// thread sleeps upon completion of ClearNumReqsSent
func (mq *MasterQueue) RateReset() {
	nextTime := time.Now().Truncate(time.Minute)
	for {
		nextTime = nextTime.Add(time.Minute)
		mq.ClearNumReqsSent()
		time.Sleep(time.Until(nextTime))
	}
}

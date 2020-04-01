package scraper

type MasterQueue struct {
	HostnameQueueMap   map[string]*HostnameQueue
	HostnameQueueQueue []*HostnameQueue
	CurrentQueue       *HostnameQueue
}

func (mq *MasterQueue) nextURL() (string, bool) {
	return "", false
}

/*
	I debated for a long time on whether paths should be a slice or a map,
	in the end, went, I went with a map, because since this source of data comes from
	the HostnameCoordinatorServer, which directly will use a map regardless to avoid duplicate
	urls, and since we can only have one hostname coordinator, it is best to lighten the load
	as much as possible.
*/
func (mq *MasterQueue) addNewHostname(hostname string, paths map[string]null ) {
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

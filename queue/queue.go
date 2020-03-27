package queue

type Queue interface {
	Enqueue(url string)
	Dequeue() string
}

type Channel chan string

func (c Channel) Enqueue(url string) {
	c <- url
}

func (c Channel) Dequeue() string {
	return <-c
}

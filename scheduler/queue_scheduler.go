package scheduler

import (
	"container/list"
	"crypto/md5"
	"sync"
)

// 队列调度器，管理urls
type QueueScheduler struct {
	locker  *sync.Mutex
	queue   *list.List
	listKey map[[md5.Size]byte]*list.Element
}

func NewQueueScheduler() *QueueScheduler {
	queue := list.New()
	locker := new(sync.Mutex)
	listKey := make(map[[md5.Size]byte]*list.Element)

	return &QueueScheduler{
		queue:   queue,
		locker:  locker,
		listKey: listKey,
	}
}

func (qs *QueueScheduler) Pop() (string, bool) {
	qs.locker.Lock()
	if qs.queue.Len() <= 0 {
		qs.locker.Unlock()
		return "", false
	}
	e := qs.queue.Front()
	url := e.Value.(string)
	key := md5.Sum([]byte(url))
	delete(qs.listKey, key)
	qs.queue.Remove(e)
	qs.locker.Unlock()
	return url, true
}

func (qs *QueueScheduler) Push(url string) {
	qs.locker.Lock()
	key := md5.Sum([]byte(url))
	if _, ok := qs.listKey[key]; ok {
		qs.locker.Unlock()
		return
	}
	e := qs.queue.PushBack(url)
	qs.listKey[key] = e
	qs.locker.Unlock()
}

func (qs *QueueScheduler) getUrlsNum(url string) int {
	return len(qs.listKey)
}

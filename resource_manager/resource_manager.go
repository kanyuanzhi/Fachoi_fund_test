package resource_manager

// 进程资源管理器
type ResourceManager struct {
	ch chan int
}

func NewResourceManager(num int) *ResourceManager {
	return &ResourceManager{ch: make(chan int, num)}
}

func (rm *ResourceManager) GetOne() {
	rm.ch <- 1
}

func (rm *ResourceManager) FreeOne() {
	<-rm.ch
}

func (rm *ResourceManager) Cap() int {
	return cap(rm.ch)
}

func (rm *ResourceManager) Has() int {
	return len(rm.ch)
}

func (rm *ResourceManager) Remain() int {
	return cap(rm.ch) - len(rm.ch)
}

func (rm *ResourceManager) FillToTheFull() {
	for i := 0; i < rm.Cap(); i++ {
		rm.ch <- 1
	}
}

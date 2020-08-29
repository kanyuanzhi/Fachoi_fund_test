package spider

type ResourceManagerChan struct {
	tc chan uint8
}

func NewResourceManagerChan(num uint8) *ResourceManagerChan {
	tc := make(chan uint8, num)
	return &ResourceManagerChan{tc: tc}
}

func (r *ResourceManagerChan) GetOne() {
	r.tc <- 1
}

func (r *ResourceManagerChan) FreeOne() {
	<-r.tc
}

func (r *ResourceManagerChan) Cap() int {
	return cap(r.tc)
}

func (r *ResourceManagerChan) Has() int {
	return len(r.tc)
}

func (r *ResourceManagerChan) Remain() int {
	return cap(r.tc) - len(r.tc)
}

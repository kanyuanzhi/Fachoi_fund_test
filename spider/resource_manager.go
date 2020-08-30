package spider

type ResourceManager struct {
	tc chan uint8
}

func NewResourceManager(num uint8) *ResourceManager {
	tc := make(chan uint8, num)
	return &ResourceManager{tc: tc}
}

func (r *ResourceManager) GetOne() {
	r.tc <- 1
}

func (r *ResourceManager) FreeOne() {
	<-r.tc
}

func (r *ResourceManager) Cap() int {
	return cap(r.tc)
}

func (r *ResourceManager) Has() int {
	return len(r.tc)
}

func (r *ResourceManager) Remain() int {
	return cap(r.tc) - len(r.tc)
}

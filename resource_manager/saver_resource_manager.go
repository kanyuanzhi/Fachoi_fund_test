package resource_manager

// 存储资源管理器，不限制存储并发数，只要数据管道中有数据就可以一直开线程存储数据
type SaverResourceManager struct {
	*ResourceManager
}

func NewSaverResourceManager(num int) *SaverResourceManager {
	crm := new(SaverResourceManager)
	crm.ResourceManager = NewResourceManager(num)
	return crm
}

// 填满存储资源管理器，用于Run()中每存储一个数据，减少一位，以控制并发存储的结束
func (srm *SaverResourceManager) FillToTheFull() {
	for i := 0; i < srm.Cap(); i++ {
		srm.ch <- 1
	}
}

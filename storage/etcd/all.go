package etcd

// 加载所有数据
func (self *Etcd) LoadAll() (err error) {
	if err = self.LoadAllService(); err != nil {
		return
	}
	return
}

// 存储所有数据
func (self *Etcd) SaveAll() (err error) {
	if err = self.SaveAllService(); err != nil {
		return
	}
	return
}

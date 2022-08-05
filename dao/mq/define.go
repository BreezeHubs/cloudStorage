package mq

// 存储类型(表示文件存到哪里)
type StoreType int

const (
	_ StoreType = iota
	// StoreLocal : 节点本地
	StoreLocal
	// StoreCeph : Ceph集群
	StoreCeph
	// StoreOSS : 阿里OSS
	StoreOSS
	// StoreMix : 混合(Ceph及OSS)
	StoreMix
	// StoreAll : 所有类型的存储都存一份数据
	StoreAll
)

type TransferData struct {
	FileHash      string
	CurLocation   string //当前文件的存储路径
	DestLocation  string //目标文件的存储路径
	DestStoreType StoreType
}

//ceph

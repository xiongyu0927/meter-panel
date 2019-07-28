package store

// NewNodeCache is used create a NodeCache that can get somedata
func NewNodeCache() *NodeCache {
	return &NodeCache{
		StoreAllClusterNodeList: StoreAllClusterNodeList,
	}
}

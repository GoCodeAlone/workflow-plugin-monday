package internal

import "sync"

var (
	clientMu       sync.RWMutex
	clientRegistry = make(map[string]*MondayClient)
)

func RegisterClient(name string, c *MondayClient) {
	clientMu.Lock()
	defer clientMu.Unlock()
	clientRegistry[name] = c
}

func GetClient(name string) (*MondayClient, bool) {
	clientMu.RLock()
	defer clientMu.RUnlock()
	c, ok := clientRegistry[name]
	return c, ok
}

func UnregisterClient(name string) {
	clientMu.Lock()
	defer clientMu.Unlock()
	delete(clientRegistry, name)
}

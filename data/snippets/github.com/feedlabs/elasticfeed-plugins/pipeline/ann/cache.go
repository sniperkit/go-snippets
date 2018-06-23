package ann

type Cache struct {
	data  interface{}
}

func (this *Cache) Lock(s string) string {
	return ""
}

func (this *Cache) Unlock(s string) {
}

func (this *Cache) RLock(s string) (string, bool) {
	return "", false
}

func (this *Cache) RUnlock(s string) {
}

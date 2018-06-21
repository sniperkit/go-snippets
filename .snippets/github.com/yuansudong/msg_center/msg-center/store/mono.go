package store

var (
	storeInstance *Store
)

func init() {
     
}

func GetInstance() *Store {

	return storeInstance
}

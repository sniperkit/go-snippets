## concurrent lock-free skiplist ##

### Current Version: 1.0beta
##### Put, PutOnlyIfAbsent, Update, Get, Contains, Remove, Iterator #####

```Go

//set max level
sl := skiplist.NewSkiplist(32)

key := []byte("test")
val := 123

//put
previou, err := sl.Put(key, unsafe.Pointer(&val)) // return nil, nil

//PutOnlyIfAbsent
previou, err = sl.PutOnlyIfAbsent(key, unsafe.Pointer(&val)) // return nil, key already exists

//update
previou, err = sl.Update(key, func(old unsafe.Pointer) unsafe.Pointer {
	return *(*int)(old)++
}) //return 123, nil

//Get
previou, err = sl.Get(key) //return 124, nil

//Contains
ok, err := sl.Contains(key) //return true, nil

//Remove
previou, err = sl.Remove(key) //return 124, nil

//CompareAndRemove
previou, err = sl.CompareAndRemove(key, value) //return nil, nil

//Get
previou, err = sl.Get(key) //return nil, nil

//Contains
ok, err := sl.Contains(key) //return false, nil

//iterator
it := skiplist.NewIterator(sl, nil, nil)

//Next
for it.Next() {
	key, value := it.NextNode()
}
```

##### performance #####
tested 1000 grounite, each grounite put && get 100000 key-value pairs, result shows follow:
![image](https://github.com/HearingFish/skiplist/blob/master/performance.png)

result shows that put is almost 70% faster than lock map and get is almost 10% faster than lock map


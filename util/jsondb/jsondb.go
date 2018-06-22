package jsondb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/distributed-vision/go-resources/util"
)

var BLOCK_SIZE uint = 32
var TAB byte = 9
var NEWLINE byte = 10

type row struct {
	index uint
	key   string
	val   interface{}
}

func (this *row) MarshalJSON() ([]byte, error) {
	val, err := json.Marshal(this.val)

	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("[ %d,\"%s\",%s]", this.index, this.key, string(val))), nil
}

func (this *row) UnmarshalJSON(b []byte) error {
	var array []interface{}
	err := json.Unmarshal(b, &array)
	if err != nil {
		return err
	}
	this.index = uint(array[0].(float64))
	this.key = array[1].(string)
	this.val = array[2]
	return nil
}

type entry struct {
	position uint
	block    uint
	row      *row
}

type JsonDb struct {
	storage     util.Storage
	head        uint
	lastIndex   uint
	entries     map[string]*entry
	freelists   [][]uint //= [[], [], [], [], [], []];
	entryMutex  sync.Mutex
	allocMutex  sync.Mutex
	taskQueue   chan func()
	syncOnWrite bool
}

func NewJsonDb(storage util.Storage, syncOnWrite bool) *JsonDb {

	return &JsonDb{
		storage:     storage,
		head:        0,
		lastIndex:   0,
		entries:     make(map[string]*entry),
		freelists:   [][]uint{},
		syncOnWrite: syncOnWrite}
}

func (this *JsonDb) Open() chan error {
	cerr := make(chan error, 1)

	this.taskQueue = make(chan func(), 10)

	go func() {
		for task := range this.taskQueue {
			task()
		}
	}()

	this.taskQueue <- func() {
		defer close(cerr)
		buf, err := this.storage.Open(true)

		if err != nil {
			cerr <- err
			return
		}

		err = this.parseDatabase(buf)

		if err != nil {
			cerr <- err
		}
	}

	return cerr
}

func (this *JsonDb) IsOpen() bool {
	return this.storage.IsOpen()
}

func (this *JsonDb) Set(key string, val interface{}) chan error {
	cerr := make(chan error, 1)

	if !this.IsOpen() {
		cerr <- fmt.Errorf("database is not open")
		close(cerr)
		return cerr
	}

	this.entryMutex.Lock()
	this.allocMutex.Lock()

	var ent = this.entries[key]
	var oldFreelist *[]uint
	var oldPointer uint

	if ent != nil {
		oldPointer = ent.position
		oldFreelist = &this.freelists[ent.block]
		ent.row.val = val
	} else {
		ent = &entry{
			position: 0,
			block:    0,
			row:      &row{0, key, val}}
		this.entries[key] = ent
	}

	this.lastIndex++
	ent.row.index = this.lastIndex
	if val == nil {
		delete(this.entries, key)
	}

	this.entryMutex.Unlock()
	this.allocMutex.Unlock()

	this.taskQueue <- func() {
		this.allocMutex.Lock()
		defer close(cerr) //fmt.Printf("row=%v\n", ent.row)

		row, err := json.Marshal(&ent.row)
		if err == nil {
			buf := bytes.Join([][]byte{[]byte{TAB}, row, []byte{NEWLINE}}, []byte{})
			//fmt.Printf("buf=%v\n", buf)

			if uint(len(buf)) > (BLOCK_SIZE << ent.block) {
				ent.block = nextBlockSize(uint(len(buf)))
			}

			ent.position = this.alloc(ent.block)

			err = this.write(buf, ent)
			if oldFreelist != nil {
				*oldFreelist = append(*oldFreelist, oldPointer)
			}
		} else {
			cerr <- err
		}

		this.allocMutex.Unlock()
	}

	return cerr
}

func (this *JsonDb) Delete(key string) chan error {
	return this.Set(key, nil)
}

func (this *JsonDb) Clear() chan error {
	cErrOut := make(chan error, 1)

	if !this.IsOpen() {
		cErrOut <- fmt.Errorf("database is not open")
		close(cErrOut)
		return cErrOut
	}

	var wg sync.WaitGroup
	errors := make([]error, 0)

	delete := func(key string) {
		defer wg.Done()

		cerr := this.Delete(key)

		if cerr == nil {
			errors = append(errors, fmt.Errorf("Clear Failed: Error channel is undefined"))
		}

		if err, ok := <-cerr; ok {
			errors = append(errors, err)
			return
		}
	}

	wg.Add(len(this.entries))
	keys := make([]string, 0, len(this.entries))
	for key, _ := range this.entries {
		keys = append(keys, key)
	}

	for _, key := range keys {
		go delete(key)
	}

	go func() {
		wg.Wait()
		if len(errors) > 0 {
			cErrOut <- fmt.Errorf("Clear failed with the following errors %v", errors)
		}
		close(cErrOut)
	}()

	return cErrOut
}

func (this *JsonDb) Get(key string) (interface{}, bool) {
	//fmt.Printf("entries=%v\n", this.entries)
	if !this.IsOpen() {
		return nil, false
	}

	this.entryMutex.Lock()
	defer this.entryMutex.Unlock()
	if entry, ok := this.entries[key]; ok {
		//fmt.Printf("val=%v\n", entry.row.val)
		return entry.row.val, ok
	}
	return nil, false
}

func (this *JsonDb) Has(key string) bool {
	if !this.IsOpen() {
		return false
	}
	this.entryMutex.Lock()
	_, ok := this.entries[key]
	this.entryMutex.Unlock()
	return ok
}

func (this *JsonDb) Len() int {
	this.entryMutex.Lock()
	defer this.entryMutex.Unlock()
	return len(this.entries)
}

func (this *JsonDb) ForEach(callback func(key string, val interface{})) {
	type entry struct {
		key   string
		value interface{}
	}

	this.entryMutex.Lock()

	entries := make([]entry, 0, len(this.entries))

	for k, e := range this.entries {
		entries = append(entries, entry{k, e.row.val})
	}

	this.entryMutex.Unlock()

	for _, entry := range entries {
		callback(entry.key, entry.value)
	}
}

func (this *JsonDb) Close() chan error {
	cErrOut := make(chan error, 1)

	cerr := this.Sync()

	this.taskQueue <- func() {
		defer close(cErrOut)

		if err, ok := <-cerr; ok {
			cErrOut <- err
			return
		}

		err := this.storage.Close()

		if err != nil {
			cErrOut <- err
		}

		this.head = 0
		this.lastIndex = 0
		this.entries = make(map[string]*entry)
		this.freelists = [][]uint{}
	}

	return cErrOut
}

func (this *JsonDb) Sync() chan error {
	cErrOut := make(chan error, 1)

	this.taskQueue <- func() {
		err := this.storage.Sync()
		if err != nil {
			cErrOut <- err
		}
		close(cErrOut)
	}

	return cErrOut
}

func (this *JsonDb) alloc(block uint) uint {
	for uint(len(this.freelists)) <= block {
		this.freelists = append(this.freelists, []uint{})
	}

	var freelist = this.freelists[block]

	if len(freelist) == 0 {
		freelist = append(freelist, this.head)
		this.head += BLOCK_SIZE << block
	}

	allocated := freelist[len(freelist)-1]
	this.freelists[block] = freelist[:len(freelist)-1]

	return allocated
}

func (this *JsonDb) parseDatabase(data []byte) error {
	this.entryMutex.Lock()
	defer this.entryMutex.Unlock()
	this.allocMutex.Lock()
	defer this.allocMutex.Unlock()

	var pointer uint = 0
	var entries = []*entry{}
	var latest = this.entries
	var lastIndex = this.lastIndex

	for i := 0; i < len(data); i++ {
		if data[i] == '\t' {
			pointer = uint(i)
		}

		if data[i] == '\n' {
			var buf = data[pointer:i]
			var row = tryParse(buf)

			if row != nil {
				var entry = &entry{
					position: uint(pointer),
					block:    nextBlockSize(uint(len(buf))),
					row:      row}
				entries = append(entries, entry)
			}
			pointer = uint(i + 1)
		}
	}

	for _, entry := range entries {
		var key = entry.row.key
		if lentry, _ := latest[key]; lentry == nil || lentry.row.index < entry.row.index {
			if entry.row.val != nil {
				latest[key] = entry
			}
		}

		if entry.row.index > lastIndex {
			lastIndex = entry.row.index
		}
	}

	filtered := make([]*entry, 0, len(entries))

	for _, entry := range entries {
		if latest[entry.row.key] != entry {
			filtered = append(filtered, entry)
		}
	}

	entries = filtered

	for _, entry := range entries {
		//fmt.Printf("filt [%v]=%v,latest []=%v\n", entry.row.val, entry, latest[entry.row.key])
		if entry.row.val == nil {
			delete(latest, entry.row.key)
		}
	}

	this.lastIndex = lastIndex
	this.populateFreelist(entries)
	return nil
}

func (this *JsonDb) populateFreelist(entries []*entry) {

	var free = func(from, to, block uint) uint {
		var size = BLOCK_SIZE << block

		for to-from >= size {
			this.freelists[block] = append(this.freelists[block], from)
			from += size
		}

		return from
	}

	var maxBlock uint = 0

	for _, entry := range entries {
		if entry.block > maxBlock {
			maxBlock = entry.block
		}
	}

	for uint(len(this.freelists)) <= maxBlock {
		this.freelists = append(this.freelists, []uint{})
	}

	for _, entry := range entries {
		var from = this.head
		from = free(from, entry.position, maxBlock)
		from = free(from, entry.position, 0)
		this.head = entry.position + (BLOCK_SIZE << entry.block)
	}
}

func (this *JsonDb) write(buf []byte, entry *entry) error {

	//fmt.Printf("write=%s\n", string(buf))
	err := this.storage.Write(buf, entry.position)

	if err != nil {
		return err
	}

	if this.syncOnWrite {
		err = this.storage.Sync()
		if err != nil {
			return err
		}
	}

	return nil
}

func nextBlockSize(length uint) uint {
	var i uint = 0
	for (BLOCK_SIZE << i) < length {
		i++
	}
	return i
}

func tryParse(data []byte) *row {
	var r row
	err := json.Unmarshal(data, &r)

	if err != nil {
		return nil
	}

	return &r
}

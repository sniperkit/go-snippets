package cache

import (
	"github.com/golang/mock/gomock"
	"testing"
	"github.com/stretchr/testify/assert"
	"httpcache/mocks/httpcache"
	"strconv"
	"sync/atomic"
	"time"
	"fmt"
)

//go:generate go get github.com/golang/mock/gomock
//go:generate go get github.com/golang/mock/mockgen
//go:generate mockgen -destination mocks/http/mock_http.go net/http ResponseWriter
//go:generate mockgen -destination mocks/httpcache/mock_cache.go httpcache Cacheable

type testAddData struct {
	testCase string
	hashes   []string
	expected []string
	size     int
	limit    int
}

var addTestCases = []testAddData{
	{"single entry, fitting", []string{"abc"}, []string{"abc"}, 100, 100},
	{"multiple entries, fitting", []string{"a", "b", "c"}, []string{"a", "b", "c"}, 33, 99},
	{"single entry, too large", []string{"a"}, []string{}, 101, 100},
	{"multiple entries, sum too large", []string{"a", "b", "c", "d"}, []string{"b", "c", "d"}, 15, 50},
}

// Data driven test for updates
func Test_Add(t *testing.T) {
	for _, test := range addTestCases {
		t.Run(test.testCase, func(t *testing.T) {
			add_TestFunction(t, test)
		})
	}
}

func add_TestFunction(t *testing.T, data testAddData) {
	//given
	ctrl := gomock.NewController(t)
	cache := NewLruCache(data.limit)
	cMocks := []Cacheable{}


	//expect
	for _, hash := range data.hashes {
		mockCacheable := mock_httpcache.NewMockCacheable(ctrl)
		mockCacheable.EXPECT().Key().Return(hash).AnyTimes()
		mockCacheable.EXPECT().Size().Return(data.size).AnyTimes()
		cMocks =append(cMocks, mockCacheable)
	}


	//when
	for _, cMock := range cMocks {
		cache.Add(cMock)
	}

	//then
	for _, expected := range data.expected {
		assert.NotNil(t,cache.(*lruCache).cache[expected])
	}

	// Check if unexpected hashes are not there
	for _, hash := range data.hashes {
		if !contains(data.expected, hash) {
			println("Not expecting "+hash)
			assert.Nil(t, cache.(*lruCache).cache[hash])
		}

	}

}

func Test_Concurreny(t *testing.T) {
	ctrl := gomock.NewController(t)
	hashes := []string{}
	expectedHashes := []string{}

	cache := NewLruCache(10000)
	var done int32
	for i := 0; i < 10000; i++ {
		hash := strconv.Itoa(i)
		go func() {
			mock := mock_httpcache.NewMockCacheable(ctrl)
			mock.EXPECT().Key().AnyTimes().Return(hash)
			mock.EXPECT().Size().AnyTimes().Return(1)
			cache.Add(mock)
			atomic.AddInt32(&done, 1)
		}()
		hashes = append(hashes, hash)
	}

	for done < 10000{
		time.Sleep(time.Millisecond)
	}

	for i := 10000; i < 11000; i++ {
		hash := strconv.Itoa(i)
		go func() {
			mock := mock_httpcache.NewMockCacheable(ctrl)
			mock.EXPECT().Key().AnyTimes().Return(hash)
			mock.EXPECT().Size().AnyTimes().Return(1)
			cache.Add(mock)
			atomic.AddInt32(&done, 1)
		}()
		expectedHashes = append(expectedHashes, strconv.Itoa(i))
	}

	for done < 11000{
		time.Sleep(time.Millisecond)
	}


	for _, expected := range expectedHashes {
		assert.NotNil(t,cache.(*lruCache).cache[expected])
	}

	assert.Equal(t, 10000, cache.(*lruCache).size)
	assert.Equal(t, 10000, len(cache.(*lruCache).cache))



	println(fmt.Sprint(done))



}

func contains(s []string, e string) bool{
	for _, se := range s {
		if se == e { return true}
	}
	return false
}

package linear_hashing

import (
	"hash/fnv"
	"math"
)

// match background
/*
key = [5, 9, 13] n = 4
key % n = 1

n = n * 2
5 % 8 = 5
9 % 8 = 1
13 % 8 = 5

so =>
key % n = M
key % 2n = M or M+n

rehash just bother current bucket and new bucket
*/

// a every simple implement
// linear hashing site: https://en.wikipedia.org/wiki/Linear_hashing
type LinearHashing struct {
	N int // initial number of buckets
	L int // The current level which is an integer that indicates on a logarithmic scale approximately how many buckets the table has grown by.
	S int // The split pointer which points to a bucket. It initially points to the first bucket in the table.

	maxBucketSlots int
	buckets        []*bucket
}

func NewLinearHashing(n, maxBucketSlots int) *LinearHashing {
	if n < 0 {
		panic("initial number of buckets must greater than 0")
	}

	l := new(LinearHashing)
	l.N = n
	l.L = 0
	l.S = 0

	l.maxBucketSlots = maxBucketSlots
	l.buckets = make([]*bucket, l.N)

	// initialize buckets
	for i := 0; i < l.N; i++ {
		l.buckets[i] = NewBucket()
	}

	return l
}

func (l *LinearHashing) currBucketSize() int {
	return l.N * int(math.Pow(2, float64(l.L)))
}

func (l *LinearHashing) hash(key string, le int) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	keyHash := int(h.Sum32())
	return keyHash % (l.N * int(math.Pow(2, float64(le))))
}

func (l *LinearHashing) splitHash() {
	newBucket := NewBucket()
	for key, val := range l.buckets[l.S].slots {
		idx := l.hash(key, l.L+1)
		if idx >= len(l.buckets) {
			l.buckets[l.S].Remove(key)
			newBucket.Add(key, val)
		}
	}
	l.buckets = append(l.buckets, newBucket)
	l.S++
	if l.S >= l.currBucketSize() {
		// one round over
		l.L++
		l.S = 0
	}
}

func (l *LinearHashing) Get(key string) (val string, ok bool) {
	bucketIdx := l.hash(key, l.L)
	if bucketIdx < l.S {
		bucketIdx = l.hash(key, l.L+1)
	}
	return l.buckets[bucketIdx].Get(key)
}

func (l *LinearHashing) Set(key string, val string) {
	bucketIdx := l.hash(key, l.L)
	if bucketIdx < l.S {
		bucketIdx = l.hash(key, l.L+1)
	}

	l.buckets[bucketIdx].Add(key, val)

	// split hash if slots greater than maxBucketSlots
	if l.buckets[bucketIdx].Size() >= l.maxBucketSlots {
		l.splitHash()
	}
}

// bucket
type bucket struct {
	slots map[string]string
}

func NewBucket() *bucket {
	b := new(bucket)
	b.slots = make(map[string]string)
	return b
}

func (b *bucket) Add(key string, value string) {
	b.slots[key] = value
}

func (b *bucket) Remove(key string) {
	delete(b.slots, key)
}

func (b *bucket) Get(key string) (value string, exists bool) {
	if val, ok := b.slots[key]; ok {
		return val, true
	}
	return "", false
}

func (b *bucket) Size() int {
	return len(b.slots)
}

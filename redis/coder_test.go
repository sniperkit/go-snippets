package redis

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	kmsgpack "github.com/kdar/cachei/msgpack"
	"gopkg.in/vmihailenco/msgpack.v2"
)

var (
	m = map[string]interface{}{"innovator": "False", "brandgenericstatusid": 4, "cp_num": 397, "repackaged": "False", "desistatusid": 6, "legendstatusid": 2, "productid": 47709, "marketerid": 1148, "deaclassificationid": 6, "replacedbyproductid": interface{}(nil), "modifiedaction": interface{}(nil), "offmarket": interface{}(nil), "prescribingname": "Metoprolol 50mg Tab", "alchemymarketedproductid": 231, "productnamelong": "Metoprolol Tartrate 50mg Tablet", "replacedbydate": interface{}(nil), "productnameshort": "Metoprolol 50mg Tab", "onmarket": time.Time{}, "productnametypeid": 2, "licensetypeid": 2, "privatelabel": "False"}
)

func conn() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	return c
}

func msgpackm(c redis.Conn, m map[string]interface{}) {
	b, err := msgpack.Marshal(m)
	if err != nil {
		panic(err)
	}

	c.Send("MULTI")
	c.Send("SET", "drugs:1", string(b))
	c.Send("EXPIRE", "drugs:1", "5")
	r, err := c.Do("EXEC")
	if err != nil {
		panic(err)
	}

	r, err = c.Do("GET", "drugs:1")
	v, err := redis.Bytes(r, err)
	if err != nil {
		panic(err)
	}

	m2 := make(map[string]interface{})
	err = msgpack.Unmarshal(v, &m2)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%#+v", m2)
}

func kmsgpackm(c redis.Conn, m map[string]interface{}) {
	mp := &kmsgpack.Coder{}

	b, err := mp.Marshal(m)
	if err != nil {
		panic(err)
	}

	c.Send("MULTI")
	c.Send("SET", "drugs:1", string(b))
	c.Send("EXPIRE", "drugs:1", "5")
	r, err := c.Do("EXEC")
	if err != nil {
		panic(err)
	}

	r, err = c.Do("GET", "drugs:1")
	v, err := redis.Bytes(r, err)
	if err != nil {
		panic(err)
	}

	m2 := make(map[string]interface{})
	err = mp.Unmarshal(v, &m2)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%#+v", m2)
}

func gobm(c redis.Conn, m map[string]interface{}) {
	gob.Register(time.Time{})

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(m)
	if err != nil {
		panic(err)
	}

	c.Send("MULTI")
	c.Send("SET", "drugs:1", buffer.String())
	c.Send("EXPIRE", "drugs:1", "5")
	r, err := c.Do("EXEC")
	if err != nil {
		panic(err)
	}

	r, err = c.Do("GET", "drugs:1")
	v, err := redis.Bytes(r, err)
	if err != nil {
		panic(err)
	}

	m2 := make(map[string]interface{})
	buffer2 := bytes.NewBuffer(v)
	decoder := gob.NewDecoder(buffer2)
	err = decoder.Decode(&m2)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%#+v", m2)
}

func jsonm(c redis.Conn, m map[string]interface{}) {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	c.Send("MULTI")
	c.Send("SET", "drugs:1", string(b))
	c.Send("EXPIRE", "drugs:1", "5")
	r, err := c.Do("EXEC")
	if err != nil {
		panic(err)
	}

	r, err = c.Do("GET", "drugs:1")
	v, err := redis.Bytes(r, err)
	if err != nil {
		panic(err)
	}

	m2 := make(map[string]interface{})
	err = json.Unmarshal(v, &m2)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%#+v", m2)
}

func hashm(c redis.Conn, m map[string]interface{}) {
	args := []interface{}{"drugs:1"}
	for key, value := range m {
		args = append(args, key, value)
	}

	c.Send("MULTI")
	c.Send("HMSET", args...)
	c.Send("EXPIRE", "drugs:1", "5")
	r, err := c.Do("EXEC")
	if err != nil {
		panic(err)
	}

	r, err = c.Do("HGETALL", "drugs:1")
	v, err := redis.Values(r, err)
	if err != nil {
		panic(err)
	}

	m2 := make(map[string]interface{})
	for x := 0; x < len(v); x += 2 {
		key := fmt.Sprintf("%s", v[x])
		value := v[x+1]

		switch t := value.(type) {
		case []byte:
			m2[key] = string(t)
		}
	}

	//fmt.Printf("%#+v", m2)
}

func BenchmarkRedis_map(b *testing.B) {
	c := conn()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		hashm(c, m)
	}
}

func BenchmarkRedis_json(b *testing.B) {
	c := conn()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsonm(c, m)
	}
}

func BenchmarkRedis_gob(b *testing.B) {
	c := conn()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gobm(c, m)
	}
}

func BenchmarkRedis_msgpack(b *testing.B) {
	c := conn()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		msgpackm(c, m)
	}
}

func BenchmarkRedis_kmsgpack(b *testing.B) {
	c := conn()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		kmsgpackm(c, m)
	}
}

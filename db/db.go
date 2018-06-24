package db

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"project_c9f7_i5l8_o0p4_p0j8/vertices"
)

// ****************************************************************************
// Each job has two db collections. We alternate which one
// we are using after checkpoints. If we only had one db collection,
// then if there is a failure during a checkpoint save, we might end up
// with an inconsistent checkpoint, which would be useless. So we save one
// at a time.

// Use the NewAccess constructor to ensure the primary key is active first

// TODO: Question: is the Primary where we are returning the completed graph? <- yes
var DB_IP string = "";
var DB_PORT string = "";
// Access contains the primary and secondary collection names for a job
type Access struct {
	Primary        string // Key to a collection of vertices
	Secondary      string // Key to the alternate collection of vertices
	PrimaryCurrent bool   // Currently-active state: true = Primary; false = Secondary
}

// NewAccess produces new Access, setting the primary key as current
func NewAccess(primaryKey string, secondaryKey string) Access {
	return Access{primaryKey, secondaryKey, true}
}

// Key produces the currently-active key
func (a Access) Key() string {
	if a.PrimaryCurrent {
		return a.Primary
	}
	return a.Secondary
}

// OtherKey produces the currently-active key
func (a Access) OtherKey() string {
	if a.PrimaryCurrent {
		return a.Secondary
	}
	return a.Primary
}

// PrimaryKey produces the primary key
func (a Access) PrimaryKey() string {
	return a.Primary
}

// SwapKeys switches primary and secondary as current
func (a *Access) SwapKeys() {
	a.PrimaryCurrent = !a.PrimaryCurrent
}

// ****************************************************************************

var dialInfo = mgo.DialInfo{
	Addrs:     []string{DB_IP},
	Direct:    true,
	Timeout:   10 * time.Second,
	Service:   DB_IP + ":" + DB_PORT,
	PoolLimit: 0,
}

const dbName = "test"

// DbVertex is the vertex representation in the db
type DbVertex struct {
	VertexID  int      `bson:"vertex_id"`
	Value     float64  `bson:"value"`
	Adjacent  []int    `bson:"adjacent"`
	Messages  []string `bson:"msgs"`
	Active    bool     `bson:"active"`
	Superstep int      `bson:"step"`
}

type byVertexID []DbVertex

func (v byVertexID) Len() int           { return len(v) }
func (v byVertexID) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byVertexID) Less(i, j int) bool { return v[i].VertexID < v[j].VertexID }

// CreateNewJob parses the file with the graph info and uploads it
// to the DB specified by the jobname. It also sets the initial value
// for each vertex to the specified initVal.
func CreateNewJob(jobname string, filename string, initVal float64) (Access, error) {
	fmt.Println("opening file")
	inFile, err := os.Open(filename)
	emptyAccess := NewAccess("", "")
	if err != nil {
		fmt.Println(err)
		return emptyAccess, err
	}
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	m := make(map[int][]int)
	fmt.Println("Scanning")
	for scanner.Scan() {
		str := scanner.Text()
		if len(str) != 0 && str[0:1] != "#" {
			vals := strings.Fields(str)
			fromVertexString := vals[0]
			toVertexString := vals[1]
			fromVertex, _ := strconv.Atoi(fromVertexString)
			toVertex, _ := strconv.Atoi(toVertexString)
			_, exists := m[toVertex]
			if !exists {
				m[toVertex] = make([]int, 0)
			}
			m[fromVertex] = append(m[fromVertex], toVertex)
		}
	}

	inFile.Close()

	fmt.Println("file closed")

	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		fmt.Println(err)
		return emptyAccess, err
	}
	defer session.Close()

	primKey := jobname
	secKey := primKey + "-secondary"
	prim := session.DB(dbName).C(primKey)
	sec := session.DB(dbName).C(secKey)
	index := mgo.Index{
		Key:        []string{"vertex_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = prim.EnsureIndex(index)
	if err != nil {
		fmt.Println(err)
		return emptyAccess, err
	}

	err = sec.EnsureIndex(index)
	if err != nil {
		fmt.Println(err)
		return emptyAccess, err
	}

	bulkprim := prim.Bulk()
	bulksec := sec.Bulk()

	fmt.Println("map size: ", len(m))

	for key, list := range m {
		var msgs []string

		writeReq := DbVertex{
			VertexID:  key,
			Value:     initVal,
			Adjacent:  list,
			Messages:  msgs,
			Active:    true,
			Superstep: 0,
		}

		bulkprim.Insert(writeReq)
		bulksec.Insert(writeReq)
	}
	_, err = bulkprim.Run()

	if err != nil {
		fmt.Println(err)
		return emptyAccess, err
	}

	_, err = bulksec.Run()

	if err != nil {
		fmt.Println(err)
		return emptyAccess, err
	}

	fmt.Println("Completed puts of all vertices")
	return NewAccess(jobname, secKey), err
}

// BatchUpdate updates several vertices in the DB at once. It only updates one
// collection specified by jobname, and updates the vertices in the slice
func BatchUpdate(jobname string, vertices map[int]vertices.Vertex) error {
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer session.Close()

	c := session.DB(dbName).C(jobname)
	bulkc := c.Bulk()

	for _, vertex := range vertices {
		dbvertex := vertexToDBVertex(vertex)
		query := bson.M{"vertex_id": dbvertex.VertexID}
		change := bson.M{"$set": bson.M{"value": dbvertex.Value, "msgs": dbvertex.Messages, "active": dbvertex.Active, "step": dbvertex.Superstep}}
		bulkc.Update(query, change)
	}
	_, err = bulkc.Run()

	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

// UpdateOne updates one vertex in the given collection
func UpdateOne(jobname string, vertex vertices.Vertex) error {
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer session.Close()

	c := session.DB(dbName).C(jobname)
	dbvertex := vertexToDBVertex(vertex)
	query := bson.M{"vertex_id": dbvertex.VertexID}
	change := bson.M{"$set": bson.M{"value": dbvertex.Value, "msgs": dbvertex.Messages, "active": dbvertex.Active, "step": dbvertex.Superstep}}
	err = c.Update(query, change)
	if err != nil {
		fmt.Println("Couldn't update", err)
		return err
	}
	return err
}

// GetOne gets one vertex from the collection specified by jobname
func GetOne(jobname string, vid int) (vertices.BaseVertex, error) {
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		fmt.Println(err)
		emptyBVertex := vertices.BaseVertex{}
		return emptyBVertex, err
	}
	defer session.Close()

	c := session.DB(dbName).C(jobname)
	result := DbVertex{}
	err = c.Find(bson.M{"vertex_id": vid}).One(&result)
	bvertex := vertexToBaseVertex(result)
	if err != nil {
		fmt.Println("Could not get item", err)
		return bvertex, err
	}
	return bvertex, err
}

// BatchGet gets a batch of vertices from the specified collection (jobname)
// with ids [minVal, maxVal) including minVal but not maxVal
func BatchGet(jobname string, minVal int, maxVal int) (map[int]vertices.BaseVertex, error) {
	var results []DbVertex
	bvertices := make(map[int]vertices.BaseVertex)

	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		fmt.Println(err)
		return bvertices, err
	}
	defer session.Close()

	c := session.DB(dbName).C(jobname)

	err = c.Find(bson.M{"vertex_id": bson.M{"$gt": minVal - 1, "$lt": maxVal}}).All(&results)
	if err != nil {
		fmt.Println("Could not get items in range specified", err)
		return bvertices, err
	}

	for _, vertex := range results {
		bvertex := vertexToBaseVertex(vertex)
		bvertices[vertex.VertexID] = bvertex
	}
	return bvertices, err
}

// NumVertices gets the number of items in the colleciton specified by jobname
func NumVertices(jobname string) (int, error) {
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer session.Close()

	c := session.DB(dbName).C(jobname)
	count, err := c.Count()
	if err != nil {
		fmt.Println("Couldn't get number of vertices", err)
		return 0, err
	}
	return count, err
}

// PrintToFile outpides the vertex id and the value of the
// vertex into the outfile, one vertex per row
func PrintToFile(jobname string, outfile string) error {
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer session.Close()

	c := session.DB(dbName).C(jobname)
	var results []DbVertex
	err = c.Find(nil).All(&results)
	if err != nil {
		fmt.Println(err)
		return err
	}

	sort.Sort(byVertexID(results))

	file, err := os.Create(outfile)
	defer file.Close()

	for _, v := range results {
		str := fmt.Sprintf("%d %g\n", v.VertexID, v.Value)
		_, err := file.WriteString(str)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return err
}

// DeleteJob deletes the collections associated with a job
func DeleteJob(a Access) error {
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer session.Close()
	primKey := a.Primary
	secKey := a.Secondary
	prim := session.DB(dbName).C(primKey)
	sec := session.DB(dbName).C(secKey)
	err = prim.DropCollection()
	if err != nil {
		fmt.Println("Collection couldn't be dropped", err)
		return err
	}
	err = sec.DropCollection()
	if err != nil {
		fmt.Println("Collection couldn't be dropped", err)
		return err
	}
	return err
}

func createStringArray(vertexMsgs []vertices.VertexMessage) []string {
	var msgs []string
	for _, vm := range vertexMsgs {
		out, _ := json.Marshal(vm)
		str := string(out[:])
		msgs = append(msgs, str)
	}

	return msgs
}

func createMessageArray(msgs []string) []vertices.VertexMessage {
	var vmsgs []vertices.VertexMessage
	for _, str := range msgs {
		bs := []byte(str)
		var vm vertices.VertexMessage
		_ = json.Unmarshal(bs, &vm)
		vmsgs = append(vmsgs, vm)
	}
	return vmsgs
}

func vertexToBaseVertex(v DbVertex) vertices.BaseVertex {
	vmsgs := createMessageArray(v.Messages)
	bvertex := vertices.BaseVertex{
		ID:          v.VertexID,
		Value:       v.Value,
		OutVertices: v.Adjacent,
		IncMsgs:     vmsgs,
		Active:      v.Active,
		Superstep:   v.Superstep,
	}
	return bvertex
}

func vertexToDBVertex(v vertices.Vertex) DbVertex {
	msgs := createStringArray(v.GetMessages())
	dbvertex := DbVertex{
		VertexID:  v.GetID(),
		Value:     v.GetValue(),
		Adjacent:  v.GetOutVertices(),
		Messages:  msgs,
		Active:    v.GetActive(),
		Superstep: v.GetSuperstep(),
	}
	return dbvertex
}

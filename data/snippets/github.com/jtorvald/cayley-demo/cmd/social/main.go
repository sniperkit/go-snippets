package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/quad"
)

func main() {
	// Some globally applicable things
	graph.IgnoreMissing = true
	graph.IgnoreDuplicates = true

	// File for your new BoltDB. Use path to regular file and not temporary in the real world
	t := getTempfileName()
	fmt.Printf("%v\n", t)

	defer os.Remove(t)                 // clean up
	store := initializeAndOpenGraph(t) // initialize the graph

	addQuads(store) // add quads to the graph

	countOuts(store, "robertmeta")
	lookAtOuts(store, "robertmeta")
	lookAtIns(store, "robertmeta")

	lookAtOuts(store, "jorgent")
	lookAtIns(store, "jorgent")

	lookAtFriendsOfFriends(store, "barakmich")

}

func lookAtFriendsOfFriends(store *cayley.Handle, to string) {
	fmt.Printf("\nlookAtFriendsOfFriends for subject (%s):\n", to)
	fmt.Printf("============================================\n")

	p := cayley.StartPath(store, quad.Raw(to))

	// find everybody that knows TO
	p = p.Tag("subject").OutWithTags([]string{"predicate"}, quad.Raw("knows")).Tag("friend")

	// display everybody that TO knows
	p.Iterate(nil).TagValues(nil, func(m map[string]quad.Value) {
		fmt.Printf("%s `%s`-> %s\n", m["subject"], m["predicate"], m["friend"])
	})

	// and from there all 'friends of friends'
	p = p.Tag("friend").OutWithTags([]string{"predicate"}, quad.Raw("knows")).Tag("friend_of_friend")

	p.Iterate(nil).TagValues(nil, func(m map[string]quad.Value) {
		fmt.Printf("%s `%s`-> %s\n", m["friend"], m["predicate"], m["friend_of_friend"])
	})
}

// countOuts ... well, counts Outs
func countOuts(store *cayley.Handle, to string) {
	p := cayley.StartPath(store, quad.Raw(to)).Out().Count()
	fmt.Printf("\n\ncountOuts for %s: ", to)
	p.Iterate(nil).EachValue(store, func(v quad.Value) {
		fmt.Printf("%d\n", quad.NativeOf(v))
	})
	fmt.Printf("============================================\n")
}

// lookAtOuts looks at the outbound links from the "to" node
func lookAtOuts(store *cayley.Handle, to string) {
	p := cayley.StartPath(store, quad.Raw(to)) // start from a single node, but we could start from multiple

	// this gives us a path with all the output predicates from our starting point
	p = p.Tag("subject").OutWithTags([]string{"predicate"}).Tag("object")

	fmt.Printf("\nlookAtOuts: subject (%s) -predicate-> object\n", to)
	fmt.Printf("============================================\n")

	p.Iterate(nil).TagValues(nil, func(m map[string]quad.Value) {
		fmt.Printf("%s `%s`-> %s\n", m["subject"], m["predicate"], m["object"])
		if m["predicate"] == quad.Raw("follows") {

			p = cayley.StartPath(store, m["object"]).Tag("subject").OutWithTags([]string{"predicate"}).Tag("object")

			p.Iterate(nil).TagValues(nil, func(m map[string]quad.Value) {
				fmt.Printf("%s `%s`-> %s\n", m["subject"], m["predicate"], m["object"])
			})
		}
	})
}

// lookAtIns looks at the inbound links to the "to" node
func lookAtIns(store *cayley.Handle, to string) {
	fmt.Printf("\nlookAtIns: object <-predicate- subject (%s)\n", to)
	fmt.Printf("=============================================\n")

	cayley.StartPath(store, quad.Raw(to)).Tag("object").InWithTags([]string{"predicate"}).Tag("subject").Iterate(nil).TagValues(nil, func(m map[string]quad.Value) {
		fmt.Printf("%s <-`%s` %s\n", m["object"], m["predicate"], m["subject"])
	})

}

func initializeAndOpenGraph(atLoc string) *cayley.Handle {
	// Initialize the database
	graph.InitQuadStore("bolt", atLoc, nil)

	// Open and use the database
	store, err := cayley.NewGraph("bolt", atLoc, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return store
}

func getTempfileName() string {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	return tmpfile.Name()
}

func addQuads(store *cayley.Handle) {
	store.AddQuad(quad.MakeRaw("barakmich", "drinks_with", "robertmeta", "demo graph"))
	store.AddQuad(quad.MakeRaw("barakmich", "is_a", "cayley creator", "demo graph"))
	store.AddQuad(quad.MakeRaw("barakmich", "knows", "robertmeta", "demo graph"))
	store.AddQuad(quad.MakeRaw("barakmich", "knows", "jorgent", "demo graph"))

	store.AddQuad(quad.MakeRaw("betawaffle", "knows", "robertmeta", "demo graph"))
	store.AddQuad(quad.MakeRaw("betawaffle", "is_a", "cayley advocate", "demo graph"))
	store.AddQuad(quad.MakeRaw("dennwc", "is_a", "cayley coding machine", "demo graph"))
	store.AddQuad(quad.MakeRaw("dennwc", "knows", "robertmeta", "demo graph"))
	store.AddQuad(quad.MakeRaw("henrocdotnet", "is_a", "cayley doubter", "demo graph"))
	store.AddQuad(quad.MakeRaw("henrocdotnet", "knows", "robertmeta", "demo graph"))
	store.AddQuad(quad.MakeRaw("henrocdotnet", "works_with", "robertmeta", "demo graph"))

	store.AddQuad(quad.MakeRaw("oren", "is_a", "cayley advocate", "demo graph"))
	store.AddQuad(quad.MakeRaw("oren", "knows", "robertmeta", "demo graph"))
	store.AddQuad(quad.MakeRaw("oren", "makes_talks_with", "robertmeta", "demo graph"))

	store.AddQuad(quad.MakeRaw("robertmeta", "is_a", "cayley advocate", "demo graph"))
	store.AddQuad(quad.MakeRaw("robertmeta", "knows", "barakmich", "demo graph"))
	store.AddQuad(quad.MakeRaw("robertmeta", "knows", "betawaffle", "demo graph"))
	store.AddQuad(quad.MakeRaw("robertmeta", "knows", "dennwc", "demo graph"))
	store.AddQuad(quad.MakeRaw("robertmeta", "knows", "dennwc", "demo graph")) // purposeful dup, will be ignored
	store.AddQuad(quad.MakeRaw("robertmeta", "knows", "dennwc", "demo graph")) // purposeful dup, will be ignored
	store.AddQuad(quad.MakeRaw("robertmeta", "knows", "dennwc", "demo graph")) // purposeful dup, will be ignored
	store.AddQuad(quad.MakeRaw("robertmeta", "knows", "henrocdotnet", "demo graph"))
	store.AddQuad(quad.MakeRaw("robertmeta", "knows", "oren", "demo graph"))

	store.AddQuad(quad.MakeRaw("jorgent", "knows", "oren", "demo graph"))
	store.AddQuad(quad.MakeRaw("jorgent", "knows", "dennwc", "demo graph"))
	store.AddQuad(quad.MakeRaw("jorgent", "drinks_with", "robertmeta", "demo graph"))

	// store meta data for a relation
	//store.AddQuad(quad.MakeRaw("jorgent", "follows", "632372a5-1085-4e63-a06c-79a6a46fdcea", "demo graph"))
	//store.AddQuad(quad.MakeRaw("632372a5-1085-4e63-a06c-79a6a46fdcea", "follows_from", "jorgent", "demo graph"))
	//store.AddQuad(quad.MakeRaw("632372a5-1085-4e63-a06c-79a6a46fdcea", "follows_to", "robertmeta", "demo graph"))
	//store.AddQuad(quad.MakeRaw("632372a5-1085-4e63-a06c-79a6a46fdcea", "follows_created_at", "2017-03-20 20:58:00", "demo graph"))

	store.AddQuad(quad.MakeRaw("cayley advocate", "is_a", "hard job without docs", "demo graph"))
	store.AddQuad(quad.MakeRaw("cayley codign machine", "is_a", "", "demo graph"))
}

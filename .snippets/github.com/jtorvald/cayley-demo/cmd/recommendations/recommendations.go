package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"math/rand"
	"time"

	"sort"

	"flag"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	_ "github.com/cayleygraph/cayley/graph/bolt"
	"github.com/cayleygraph/cayley/quad"
	"github.com/satori/go.uuid"
)

func main() {
	// Some globally applicable things
	graph.IgnoreMissing = true
	graph.IgnoreDuplicates = true

	// File for your new BoltDB. Use path to regular file and not temporary in the real world
	var t string

	file := flag.String("file", "", "File for the database")
	randomCustomers := flag.Int("customers", 10, "Number of random customers to generate")
	flag.Parse()

	fileExisted := false
	if *file != "" {
		t = *file
		if _, err := os.Stat(*file); err == nil {
			fileExisted = true
		}

	} else {
		t = getTempfileName()
		defer os.Remove(t) // clean up
	}

	fmt.Printf("Using database file: %v\n", t)

	store := initializeAndOpenGraph(t) // initialize the graph

	if !fileExisted {
		fmt.Println("Adding test data")
		addQuads(store, *randomCustomers) // add quads to the graph
	}

	// John Doe
	id := quad.IRI("3117979d-516a-4bac-a55e-b71g4dcb2351")

	// list products that John bought
	findProductsForCustomer(store, id)

	// find product recommendations for John
	findProductRecommendationsForCustomer(store, id)

	// find product recommendations for trackball
	findProductRecommendationsForProduct(store, quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2364"))

}

type ProductRecommendation struct {
	ProductID string
	Name      string
	Count     int32
}

type ProductRecommendations []ProductRecommendation

func (slice ProductRecommendations) Len() int {
	return len(slice)
}

func (slice ProductRecommendations) Less(i, j int) bool {
	return slice[i].Count > slice[j].Count
}

func (slice ProductRecommendations) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func findProductsForCustomer(store *cayley.Handle, to quad.Value) {
	fmt.Printf("\nFind products bought by customer (%s):\n", to)
	fmt.Printf("============================================\n")

	// start from the custoemr
	current_customer := cayley.StartPath(store, to)

	p := current_customer.Out(quad.IRI("bought")).Tag("product").Save(quad.IRI("label"), "name")
	// display all the product recommendations
	p.Iterate(nil).TagValues(nil, func(m map[string]quad.Value) {
		fmt.Printf("%s %s %s\n", to.String(), m["product"], m["name"])

	})
}

// c1 -> products1 -> group <- products2 (- products1) <- c2
func findProductRecommendationsForCustomer(store *cayley.Handle, to quad.Value) {
	fmt.Printf("\nFind product recommendations for customer (%s):\n", to)
	fmt.Printf("============================================\n")

	// start from the custoemr
	current_customer := cayley.StartPath(store, to)

	// find everybody that knows TO
	pred_bought := quad.IRI("bought")
	pred_product_group := quad.IRI("in_group")
	pred_label := quad.IRI("label")
	pred_firstname := quad.IRI("firstname")

	// find the articles the customer bought (to exclude later)
	customer_articles := current_customer.Out(pred_bought) //.Unique() // this one makes it a bit slower
	product_groups := customer_articles.Out(pred_product_group).Unique()

	// who else bought these articles?
	p := product_groups.In(pred_product_group).Except(customer_articles).Tag("product").Save(pred_label, "name")

	p = p.InWithTags([]string{"predicate"}, pred_bought).Tag("customer").Save(pred_firstname, "client_name") // c2

	recmap := make(map[string]ProductRecommendation)

	// display all the product recommendations
	p.Iterate(nil).TagValues(nil, func(m map[string]quad.Value) {
		//fmt.Printf("%s %s `%s`-> %s %s\n", m["customer"], m["client_name"], m["predicate"], m["product"], m["name"])
		if _, ok := recmap[m["product"].String()]; !ok {
			r := ProductRecommendation{}
			r.Name = m["name"].String()
			r.ProductID = m["product"].String()
			recmap[m["product"].String()] = r
		}
		obj := recmap[m["product"].String()]
		obj.Count++
		recmap[m["product"].String()] = obj
	})
	recommendations := ProductRecommendations{}
	for _, r := range recmap {
		recommendations = append(recommendations, r)
	}
	sort.Sort(recommendations)
	fmt.Printf("%v\n", recommendations)
}

// product1 -> group <- products2 (- product1) <- c2 (+ product_id)
func findProductRecommendationsForProduct(store *cayley.Handle, product_id quad.Value) {
	fmt.Printf("\nFind product recommendations for product (%s):\n", product_id)
	fmt.Printf("============================================\n")

	// start from the custoemr
	current_product := cayley.StartPath(store, product_id)

	// find everybody that knows TO
	pred_bought := quad.IRI("bought")
	pred_product_group := quad.IRI("in_group")
	pred_label := quad.IRI("label")
	pred_firstname := quad.IRI("firstname")

	// find the articles the customer bought (to exclude later)
	customer_articles := current_product
	product_groups := customer_articles.Out(pred_product_group).Unique()

	// what did they buy?
	p := product_groups.In(pred_product_group).Except(customer_articles).Tag("product").Save(pred_label, "name")

	// who bought it that also bought the same product
	//.Filter(iterator.CompareGT, quad.Value(time.Date(2017, 01, 23, 0, 0, 0, 0, nil))
	// Tag().Out().Filter().Back()
	p = p.InWithTags([]string{"predicate"}, pred_bought).Has(pred_bought, product_id).Tag("customer").Save(pred_firstname, "client_name") // c2

	recmap := make(map[string]ProductRecommendation)

	// display all the product recommendations
	p.Iterate(nil).TagValues(nil, func(m map[string]quad.Value) {
		//fmt.Printf("%s %s `%s`-> %s %s\n", m["customer"], m["client_name"], m["predicate"], m["product"], m["name"])

		if _, ok := recmap[m["product"].String()]; !ok {
			r := ProductRecommendation{}
			r.Name = m["name"].String()
			r.ProductID = m["product"].String()
			recmap[m["product"].String()] = r
		}
		obj := recmap[m["product"].String()]
		obj.Count++
		recmap[m["product"].String()] = obj
	})

	recommendations := ProductRecommendations{}
	for _, r := range recmap {
		recommendations = append(recommendations, r)
	}
	sort.Sort(recommendations)
	fmt.Printf("%v\n", recommendations)
}

func lookAtFriendsOfFriends(store *cayley.Handle, to quad.Value) {
	fmt.Printf("\nlookAtFriendsOfFriends for subject (%s):\n", to)
	fmt.Printf("============================================\n")

	p := cayley.StartPath(store, to)

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
func countOuts(store *cayley.Handle, to quad.Value) {
	p := cayley.StartPath(store, to).Out().Count()
	fmt.Printf("\n\ncountOuts for %s: ", to)
	p.Iterate(nil).EachValue(store, func(v quad.Value) {
		fmt.Printf("%d\n", quad.NativeOf(v))
	})
	fmt.Printf("============================================\n")
}

// countIns... well, counts Ins
func countIns(store *cayley.Handle, to quad.Value) {
	p := cayley.StartPath(store, to).In().Count()
	fmt.Printf("\n\ncountIns for %s: ", to)
	p.Iterate(nil).EachValue(store, func(v quad.Value) {
		fmt.Printf("%d\n", quad.NativeOf(v))
	})
	fmt.Printf("============================================\n")
}

// lookAtOuts looks at the outbound links from the "to" node
func lookAtOuts(store *cayley.Handle, to quad.Value) {
	p := cayley.StartPath(store, to) // start from a single node, but we could start from multiple

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
func lookAtIns(store *cayley.Handle, to quad.Value) {
	fmt.Printf("\nlookAtIns: object <-predicate- subject (%s)\n", to)
	fmt.Printf("=============================================\n")

	cayley.StartPath(store, to).Tag("object").InWithTags([]string{"predicate"}).Tag("subject").Iterate(nil).TagValues(nil, func(m map[string]quad.Value) {
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

func addQuads(store *cayley.Handle, randomCustomers int) {

	tr := graph.NewWriter(store)

	// register type product
	tr.WriteQuad(quad.Make(quad.IRI("product"), quad.IRI("type"), quad.IRI("class"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("product"), quad.IRI("label"), "Product", "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("product"), quad.IRI("desc"), "A store product", "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("product"), quad.IRI("hasProperty"), quad.IRI("label"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("product"), quad.IRI("hasProperty"), quad.IRI("desc"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("product"), quad.IRI("hasProperty"), quad.IRI("price"), "catalog"))

	// register type product_group
	tr.WriteQuad(quad.Make(quad.IRI("product_group"), quad.IRI("type"), quad.IRI("class"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("product_group"), quad.IRI("label"), "Product Group", "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("product_group"), quad.IRI("desc"), "A product group", "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("product_group"), quad.IRI("hasProperty"), quad.IRI("label"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("product_group"), quad.IRI("hasProperty"), quad.IRI("desc"), "catalog"))

	// add product group: electronics
	tr.WriteQuad(quad.Make(quad.IRI("electronics"), quad.IRI("type"), quad.IRI("product_group"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("electronics"), quad.IRI("label"), quad.IRI("Electronics"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("electronics"), quad.IRI("desc"), quad.IRI("Electronics for in and around the house"), "catalog"))

	// add product group: bed
	tr.WriteQuad(quad.Make(quad.IRI("bedroom"), quad.IRI("type"), quad.IRI("product_group"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("bedroom"), quad.IRI("label"), quad.IRI("Bedroom"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("bedroom"), quad.IRI("desc"), quad.IRI("Everything for in the bedroom"), "catalog"))

	// household
	tr.WriteQuad(quad.Make(quad.IRI("utensils"), quad.IRI("type"), quad.IRI("product_group"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("utensils"), quad.IRI("label"), quad.IRI("Utensils"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("utensils"), quad.IRI("desc"), quad.IRI("Every tool you need in house"), "catalog"))

	// household
	tr.WriteQuad(quad.Make(quad.IRI("household"), quad.IRI("type"), quad.IRI("product_group"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("household"), quad.IRI("label"), quad.IRI("Household"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("household"), quad.IRI("desc"), quad.IRI("Everything for your household"), "catalog"))

	// register type client
	tr.WriteQuad(quad.Make(quad.IRI("client"), quad.IRI("type"), quad.IRI("class"), "crm"))
	tr.WriteQuad(quad.Make(quad.IRI("client"), quad.IRI("label"), "Client", "crm"))
	tr.WriteQuad(quad.Make(quad.IRI("client"), quad.IRI("desc"), "A person who placed an order", "crm"))
	tr.WriteQuad(quad.Make(quad.IRI("client"), quad.IRI("hasProperty"), quad.IRI("firstname"), "crm"))
	tr.WriteQuad(quad.Make(quad.IRI("client"), quad.IRI("hasProperty"), quad.IRI("lastname"), "crm"))

	// add clients
	tr.WriteQuads(generateClientQuads("3117979d-516a-4bac-a55e-b71g4dcb2351", "John", "Doe"))
	tr.WriteQuads(generateClientQuads("3217979d-516a-4bac-a55e-b71f4dcb2352", "Alice", "Blue"))
	tr.WriteQuads(generateClientQuads("3317979d-516a-4bac-a55e-b71e4dcb2353", "Jase", "Folli"))
	tr.WriteQuads(generateClientQuads("3417979d-516a-4bac-a55e-b71d4dcb2355", "Casper", "Walden"))

	// products
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2351", "Walkman", "This is a description", 12.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2352", "Discman", "This is a description", 34.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2353", "Walky talky", "This is a description", 12.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2354", "Pencil", "This is a description", 76.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2355", "Pen", "This is a description", 2.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2356", "Pillow", "This is a description", 54.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2357", "Blanket", "This is a description", 34.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2358", "Sheets", "This is a description", 52.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2359", "Bucket", "This is a description", 21.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2360", "Monitor", "This is a description", 321.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2361", "Laptop", "This is a description", 426.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2362", "Keyboard", "This is a description", 34.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2363", "Mouse", "This is a description", 12.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2364", "Trackball", "This is a description", 23.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2365", "Harddrive", "This is a description", 202.3))
	tr.WriteQuads(generateProductQuads("2017979d-516a-4bac-a55e-b71c4dcb2366", "MagSafe Adapter", "This is a description", 86.3))

	// put products in their groups
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2351"), quad.IRI("in_group"), quad.IRI("electronics"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2352"), quad.IRI("in_group"), quad.IRI("electronics"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2353"), quad.IRI("in_group"), quad.IRI("electronics"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2360"), quad.IRI("in_group"), quad.IRI("electronics"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2361"), quad.IRI("in_group"), quad.IRI("electronics"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2362"), quad.IRI("in_group"), quad.IRI("electronics"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2363"), quad.IRI("in_group"), quad.IRI("electronics"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2364"), quad.IRI("in_group"), quad.IRI("electronics"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2365"), quad.IRI("in_group"), quad.IRI("electronics"), "catalog"))

	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2354"), quad.IRI("in_group"), quad.IRI("utensils"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2355"), quad.IRI("in_group"), quad.IRI("utensils"), "catalog"))

	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2359"), quad.IRI("in_group"), quad.IRI("household"), "catalog"))

	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2356"), quad.IRI("in_group"), quad.IRI("bedroom"), "catalog"))
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2357"), quad.IRI("in_group"), quad.IRI("bedroom"), "catalog"))

	// add purchases
	// john bought a walkman
	tr.WriteQuad(quad.Make(quad.IRI("3117979d-516a-4bac-a55e-b71g4dcb2351"), quad.IRI("bought"), quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2351"), "sales"))
	// and a monitor
	tr.WriteQuad(quad.Make(quad.IRI("3117979d-516a-4bac-a55e-b71g4dcb2351"), quad.IRI("bought"), quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2360"), "sales"))

	// alice bought a pencil and a walkman
	tr.WriteQuad(quad.Make(quad.IRI("3217979d-516a-4bac-a55e-b71f4dcb2352"), quad.IRI("bought"), quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2354"), "sales"))
	tr.WriteQuad(quad.Make(quad.IRI("3217979d-516a-4bac-a55e-b71f4dcb2352"), quad.IRI("bought"), quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2351"), "sales"))
	tr.WriteQuad(quad.Make(quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2364"), quad.IRI("bought"), quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2351"), "sales"))

	// casper bought a harddrive pencil and a walkman
	tr.WriteQuad(quad.Make(quad.IRI("3417979d-516a-4bac-a55e-b71d4dcb2355"), quad.IRI("bought"), quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2365"), "sales"))
	tr.WriteQuad(quad.Make(quad.IRI("3417979d-516a-4bac-a55e-b71d4dcb2355"), quad.IRI("bought"), quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2354"), "sales"))
	tr.WriteQuad(quad.Make(quad.IRI("3417979d-516a-4bac-a55e-b71d4dcb2355"), quad.IRI("bought"), quad.IRI("2017979d-516a-4bac-a55e-b71c4dcb2351"), "sales"))

	randomproducts := []string{
		"2017979d-516a-4bac-a55e-b71c4dcb2351",
		"2017979d-516a-4bac-a55e-b71c4dcb2352",
		"2017979d-516a-4bac-a55e-b71c4dcb2353",
		"2017979d-516a-4bac-a55e-b71c4dcb2354",
		"2017979d-516a-4bac-a55e-b71c4dcb2355",
		"2017979d-516a-4bac-a55e-b71c4dcb2356",
		"2017979d-516a-4bac-a55e-b71c4dcb2357",
		"2017979d-516a-4bac-a55e-b71c4dcb2358",
		"2017979d-516a-4bac-a55e-b71c4dcb2359",
		"2017979d-516a-4bac-a55e-b71c4dcb2360",
		"2017979d-516a-4bac-a55e-b71c4dcb2361",
		"2017979d-516a-4bac-a55e-b71c4dcb2362",
		"2017979d-516a-4bac-a55e-b71c4dcb2363",
		"2017979d-516a-4bac-a55e-b71c4dcb2364",
		"2017979d-516a-4bac-a55e-b71c4dcb2365",
		"2017979d-516a-4bac-a55e-b71c4dcb2366",
	}
	// now create some random data
	for i := 0; i < randomCustomers; i++ {
		userId := uuid.NewV4().String()
		tr.WriteQuads(generateClientQuads(userId, fmt.Sprintf("User %d", i), fmt.Sprintf("Lastname %d", i)))

		rand.Seed(time.Now().UnixNano())
		n := rand.Int() % len(randomproducts)

		// bought two random products
		tr.WriteQuad(quad.Make(quad.IRI(userId), quad.IRI("bought"), quad.IRI(randomproducts[n]), "sales"))

		rand.Seed(time.Now().UnixNano())
		n = rand.Int() % len(randomproducts)
		tr.WriteQuad(quad.Make(quad.IRI(userId), quad.IRI("bought"), quad.IRI(randomproducts[n]), "sales"))
	}

	err := tr.Close()
	if err != nil {
		panic(err)
	}

}

func generateClientQuads(id, firstname, lastname string) []quad.Quad {
	quads := make([]quad.Quad, 0, 3)
	quads = append(quads, quad.Make(quad.IRI(id), quad.IRI("firstname"), firstname, "crm"))
	quads = append(quads, quad.Make(quad.IRI(id), quad.IRI("lastname"), lastname, "crm"))
	quads = append(quads, quad.Make(quad.IRI(id), quad.IRI("type"), quad.IRI("client"), "crm"))
	return quads
}
func generateProductQuads(id, title, description string, price float32) []quad.Quad {
	quads := make([]quad.Quad, 0, 4)

	quads = append(quads, quad.Make(quad.IRI(id), quad.IRI("label"), title, "catalog"))
	quads = append(quads, quad.Make(quad.IRI(id), quad.IRI("desc"), description, "catalog"))
	quads = append(quads, quad.Make(quad.IRI(id), quad.IRI("price"), price, "catalog"))
	quads = append(quads, quad.Make(quad.IRI(id), quad.IRI("type"), quad.IRI("product"), "catalog"))

	return quads
}

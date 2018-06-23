Snippets of Go code for Sublime Text
====================================


csv.Read:
```go

r := csv.NewReader(strings.NewReader(${1:in}))

for {
	record, err := r.Read()
	if err == io.EOF {
		break
	}
	if err != nil {
		log.Fatal(err)
	}

	${2:fmt.Println(record)}
}

```

csv.Write:
```go

w := csv.NewWriter(${1:os.Stdout})

for _, record := range records {
	if err := w.Write(record); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
}

w.Flush()

if err := w.Error(); err != nil {
	log.Fatal(err)
}

```

http.Basic:
```go

func handler${1:route}(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "...")
}

func main() {
	http.HandleFunc("/${1:route}", handler${1:route})
	http.HandleFunc("/", http.NotFound)
	http.ListenAndServe(":8080", nil)
}

```

http.Get:
```go
resp, err := http.Get("${1:http://example.com/}")
if err != nil {
	${2:log.Fatal(err)}
}
defer resp.Body.Close()
body, err := ioutil.ReadAll(resp.Body)
if err != nil {
	${3:log.Fatal(err)}
}
${4:// code...}

```

http.Route:
```go

func handler${1:route}(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "...")
}

func init() {
	http.HandleFunc("/${1:route}", handler${1:route})
}

```

json.Marshal:
```go

b, err := json.Marshal(${1:data})
if err != nil {
	log.Fatal(err)
}

```

json.NewDecoder:
```go

var ${1:data} ${2:interface{}}
if err := json.NewDecoder(${3:r}).Decode(&${1:data}); err != nil {
	${4:log.Fatal(err)}
}

```

json.NewDecoder:
```go

var ${1:data} ${2:interface{}}
if err := json.NewEncoder(${3:w}).Encode(&${1:data}); err != nil {
	${4:log.Fatal(err)}
}

```

json.Unmarshal:
```go

var ${2:dst} ${3:map[string]string}
err := json.Unmarshal(${1:jsonBlob}, &${2:dst})
if err != nil {
	log.Fatal(err)
}

```

strconv.Atoi:
```go

${1:i}, err := strconv.Atoi(${2:str})
if err != nil {
	log.Fatal(err)
}

```

strconv.ParseBool:
```go

${1:v}, err := strconv.ParseBool(${2:str})
if err != nil {
	log.Fatal(err)
}

```

strconv.ParseFloat:
```go

${1:f64}, err := strconv.ParseFloat(${2:str}, ${3:64})
if err != nil {
	log.Fatal(err)
}

```

strconv.ParseInt:
```go

${1:i64}, err := strconv.ParseInt(${2:str}, ${3:10}, ${4:64})
if err != nil {
	log.Fatal(err)
}

```

strconv.ParseUint:
```go

${1:u64}, err := strconv.ParseUint(${2:str}, ${3:10}, ${4:64})
if err != nil {
	log.Fatal(err)
}

```

tabwriter:
```go

w := new(tabwriter.Writer)

// Format in tab-separated columns with a tab stop of 8.
w.Init(${1:os.Stdout}, ${2:0}, ${3:8}, ${4:0}, ${5:'\t'}, ${6:0})
fmt.Fprintln(w, "${7:...\t...}")
fmt.Fprintln(w)

if err := w.Flush(); err != nil {
	log.Fatal(err)
}

```

time.After:
```go

select {
case <-time.After(${1:5 * time.Minute}):
	${2:fmt.Println("timed out")}
${3:case m := <-c:
	handle(m)
}
}

```

time.Date:
```go

${1:t} := time.Date(${2:2009}, time.${3:November}, ${4:10}, ${5:23}, ${6:0}, ${7:0}, ${8:0}, ${9:time.UTC})

```

time.Parse:
```go

${1:t}, err := time.Parse(${2:`Mon Jan 2 15:04:05 -0700 MST 2006`}, ${3:str})
if err != nil {
	log.Fatal(err)
}

```

time.ParseDuration:
```go

${1:d}, err := time.ParseDuration(${2:str})
if err != nil {
	log.Fatal(err)
}

```

time.Tick:
```go

${1:c} := time.Tick(${2:1 * time.Minute})
for now := range ${1:c} {
	${3:fmt.Printf("%v\n", now)}
}

```


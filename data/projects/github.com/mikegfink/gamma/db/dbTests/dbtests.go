package main

import (
	"fmt"
	"time"
	"os"
	"project_c9f7_i5l8_o0p4_p0j8/db"
	"project_c9f7_i5l8_o0p4_p0j8/pregel"
)

func main() {
	fmt.Println("CREATE NEW JOB")
	jobname := "job_wiki_vote"
	filename := "/Users/dorothyordogh/Downloads/wiki-Vote.txt"

	start := time.Now()
	response := db.CreateNewJob(jobname, filename, 1)
	elapsed := time.Since(start)
	fmt.Println(response, " took ", elapsed)

	fmt.Println("GET VERTEX WITH ID 3")
	start = time.Now()
	vertex1 := db.GetOne(jobname, 3)
	fmt.Println(vertex1)
	elapsed = time.Since(start)
	fmt.Println("took", elapsed)

	vertex1.Active = false

	fmt.Println("UPDATING VERTEX 3's MESSAGE")
	start = time.Now()
	bval := db.UpdateOne(jobname, vertex1)
	fmt.Println(bval)
	elapsed = time.Since(start)
	fmt.Println("took", elapsed)

	fmt.Println("CHECKING THAT THE UPDATE WORKED")
	start = time.Now()
	vertex2 := db.GetOne(jobname, 3)
	fmt.Println(vertex2)
	elapsed = time.Since(start)
	fmt.Println("took", elapsed)

	fmt.Println("GET VERTICES in the range 3-6 (ie 3, 4, 5, & 6)")
	start = time.Now()
	vertices := dbmanager.BatchGet(jobname, 1, 6)
	fmt.Println(vertices)
	elapsed = time.Since(start)
	fmt.Println("took", elapsed)

	fmt.Println("CHANGING MESSAGES")
	for k, val := range vertices {
		val.Active = false
		vertices[k] = val
	}

	fmt.Println("CHECKING THAT IT WORKED")
	fmt.Println(vertices)

	fmt.Println("BATCH UPDATE 3, 4, 5, & 6")
	start = time.Now()
	bval = db.BatchUpdate(jobname, vertices)
	fmt.Println(bval)
	elapsed = time.Since(start)
	fmt.Println("took", elapsed)

	fmt.Println("CHECKING THAT BATCH UPDATE WORKED 3, 4, 5, & 6")
	start = time.Now()
	vertices = db.BatchGet(jobname, 1, 6)
	fmt.Println(vertices)
	elapsed = time.Since(start)
	fmt.Println("took", elapsed)

	fmt.Println("DELETING JOB")
	start = time.Now()
	bval = db.DeleteJob(jobname)
	fmt.Println(bval)
	elapsed = time.Since(start)
	fmt.Println("took", elapsed)
}

package main

import "gorlim"
import "fmt"
import "strconv"
import "os"

func testAddIssues() gorlim.IssueRepositoryInterface {
   repo := gorlim.CreateRepo(os.Getenv("HOME") + "/issues0")

   issue1 := gorlim.Issue{Id:1, Opened: true, Assignee:"gark", Milestone:"mile", Title:"mytitle", 
                          Description:"mydescription", Labels: []string{"l1", "l2"}, Comments: []string{"a", "b", "c"}}
   issue2 := gorlim.Issue{Id:2, Opened: false, Assignee:"gark", Milestone:"mile2", Title:"mytitle",
                          Description:"mydescription", Comments: []string{"a", "b", "c", "d"}}
   issues := []gorlim.Issue{issue1, issue2}
   repo.Update("comment", issues)	


   issue3 := gorlim.Issue{Id:3, Opened: true, Assignee:"gark", Milestone:"mile2", Title:"mytitle",
                          Description:"mydescription", Comments: []string{"a", "b", "c", "d"}}
   issues = []gorlim.Issue{issue3}
   repo.Update("comment2", issues)	
   return repo
}

func testGetIssues(repo gorlim.IssueRepositoryInterface) {
    issues, _ := repo.GetIssues()
   	for _, issue := range issues {
   		fmt.Println("Issue")
   		fmt.Println("id: " + strconv.Itoa(issue.Id))
   		fmt.Println("Assignee: " + issue.Assignee)
   		fmt.Println("Milestone: " + issue.Milestone)
   		fmt.Println("Title: "  + issue.Title)
   		fmt.Println("Description: "  + issue.Description)
   		if issue.Opened {
   			fmt.Print("Opened: true\n")
   		} else {
   			fmt.Print("Opened: false\n ")
   		}
   		fmt.Print("Labels: ")
   		for _, label := range issue.Labels {
   			fmt.Print(label + ",")
   		}
   		fmt.Println("")
   		fmt.Print("Comments: ")
   		for _, comment := range issue.Comments {
   			fmt.Print(comment + " ")
   		}
   		fmt.Println("\n\n\n")
   	}
}

func main() {
  r := testAddIssues()
  testGetIssues(r)
}
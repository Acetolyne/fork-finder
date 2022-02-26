package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/google/go-github/v42/github"
)

type Repo struct {
	Name string
	Date *time.Time
}

func main() {
	var r []Repo
	//@todo get this from arguement
	var repo = "qt"
	client := github.NewClient(nil)
	//@todo get the repo from an argument
	opt, res, err := client.Repositories.ListForks(context.Background(), "therecipe", repo, nil)
	commits, res, err := client.Repositories.ListCommits(context.Background(), "therecipe", repo, nil)
	//Last commit date of base repo
	fmt.Println("Last Commit to base repository:", commits[0].Commit.Committer.Date)
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode == 200 {
		for i := range opt {
			owner := string(opt[i].Owner.GetLogin())
			commits, res, err = client.Repositories.ListCommits(context.Background(), owner, repo, nil)
			if len(commits) > 0 {
				r = append(r, Repo{opt[i].GetFullName(), commits[0].Commit.Committer.Date})
			}
			//fmt.Println(opt[i].Comm)
		}
		//sort by date reversed
		sort.Slice(r[:], func(i, j int) bool {
			a := r[i].Date.Unix()
			b := r[j].Date.Unix()
			return a > b
		})
		//@todo if argument --all is used return all results in reverse date order
		for i := range r {
			fmt.Println(r[i].Name, r[i].Date)
		}
		//@todo else we just return the first result
	} else {
		fmt.Println("Github API returned status code:", res.StatusCode)
	}
}

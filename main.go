package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v42/github"
)

type Repo struct {
	Name string
	Date *time.Time
}

func main() {

	var user, reponame string
	allFlag := flag.Bool("a", false, "Print all forked repositories instead of the one with the latest commits")
	flag.Parse()

	if len(os.Args) > 1 {
		s := os.Args[1]
		//repo := s[strings.LastIndex(s, "/")+1:]
		//user := s[strings.LastIndex(s, "/")+1:]
		ss := strings.Split(s, "/")
		reponame = ss[len(ss)-1]
		user = ss[len(ss)-2]
	}
	if user == "" || reponame == "" {
		fmt.Printf("please provide the base repository to search in the format fork-finder REPO where repo is the url of the base repository, can be any url format or just the name of the repository like Acetolyne/fork-finder")
		os.Exit(1)
	}
	var r []Repo
	client := github.NewClient(nil)
	//@todo get the repo from an argument
	opt, res, err := client.Repositories.ListForks(context.Background(), user, reponame, nil)
	if err != nil {
		fmt.Println(err)
	}
	if res.StatusCode == 200 {
		commits, res, err := client.Repositories.ListCommits(context.Background(), user, reponame, nil)
		//Last commit date of base repo
		fmt.Println("Last Commit to base repository:", commits[0].Commit.Committer.Date)
		if err != nil {
			fmt.Println(err)
		}
		if res.StatusCode == 200 {
			for i := range opt {
				owner := string(opt[i].Owner.GetLogin())
				commits, _, _ = client.Repositories.ListCommits(context.Background(), owner, reponame, nil)
				if len(commits) > 0 {
					r = append(r, Repo{opt[i].GetFullName(), commits[0].Commit.Committer.Date})
				}
			}
			//sort by date reversed
			sort.Slice(r[:], func(i, j int) bool {
				a := r[i].Date.Unix()
				b := r[j].Date.Unix()
				return a > b
			})
			if *allFlag {
				for i := range r {
					fmt.Println(r[i].Name, r[i].Date)
				}
			} else {
				fmt.Println(r[0].Name, r[0].Date)
			}
		} else {
			fmt.Println("Github API returned status code:", res.StatusCode, err)
		}
	} else {
		fmt.Println("Github API returned status code:", res.StatusCode, err)
	}
}

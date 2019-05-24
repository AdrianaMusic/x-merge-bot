package main

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

func getMergeOpt(request *github.PullRequest) *github.PullRequestOptions {
	defaultBranch := os.Getenv("DEFAULT_BRANCH")
	commitMergeOpt := github.PullRequestOptions{MergeMethod: "merge"}
	squashMergeOpt := github.PullRequestOptions{MergeMethod: "squash"}
	if request.Base.GetRef() == defaultBranch {
		return &squashMergeOpt
	} else {
		return &commitMergeOpt
	}
}

func reviewRepo(client *github.Client, repo string) {
	owner := os.Getenv("OWNER")

	prs, _, err := client.PullRequests.List(context.Background(), owner, repo, nil)
	if err != nil {
		log.Print(err)
		return
	}

	for i, pr := range prs {
		pullRequest, res, err := client.PullRequests.Get(context.Background(), owner, repo, pr.GetNumber())
		if err != nil {
			log.Print(i, " fetch pull request ", pr.GetNumber(), "failure")
			log.Print(res)
			return
		}

		if pullRequest.GetMergeable() && pullRequest.GetMergeableState() == "clean" {
			opt := getMergeOpt(pullRequest)
			p, res, err := client.PullRequests.Merge(context.Background(), owner, repo, pullRequest.GetNumber(), pullRequest.GetTitle(), opt)
			if err != nil {
				log.Print(res)
				return
			}

			log.Print(p.GetMessage())
		}
	}
}

func initClient() *github.Client {
	log.Print("token:", os.Getenv("TOKEN"))
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func startReview() {
	client := initClient()
	repos := strings.Fields(os.Getenv("REPO"))
	log.Print("Start review repos")
	for i, repo := range repos {
		log.Print(i+1, ".", repo)
		reviewRepo(client, repo)
	}
	log.Print("Finish review repos")
	time.AfterFunc(180*time.Second, startReview)
}

func blockForever() {
foo:
	runtime.Gosched()
	goto foo
}

func main() {
	startReview()
	blockForever()
}

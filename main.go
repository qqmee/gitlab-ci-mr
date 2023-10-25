package main

import (
	"github.com/xanzy/go-gitlab"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Env struct {
	BaseUrl       string
	Token         string
	Pid           int
	UserId        int
	MainBranch    string
	SourceBranch  string
	TargetBranch  string
	CommitMessage string
}

func getTarget(taskName, sourceBranchShortName, mainBranch string) string {
	if sourceBranchShortName == "feature" {
		return "review/" + taskName
	}

	return mainBranch
}

func parseEnv() *Env {
	pid, _ := strconv.Atoi(os.Getenv("CI_PROJECT_ID"))
	userId, _ := strconv.Atoi(os.Getenv("GITLAB_USER_ID"))
	sourceBranch := os.Getenv("CI_COMMIT_REF_NAME")
	mainBranch := os.Getenv("CI_DEFAULT_BRANCH")

	regexTarget := regexp.MustCompile("^(feature|review|hotfix)/(.*)")
	sourceBranchShortName := regexTarget.ReplaceAllString(sourceBranch, "$1")
	taskName := regexTarget.ReplaceAllString(sourceBranch, "$2")
	targetBranch := getTarget(taskName, sourceBranchShortName, mainBranch)

	return &Env{
		BaseUrl:       os.Getenv("CI_SERVER_URL"),
		Token:         os.Getenv("PAT"),
		Pid:           pid,
		UserId:        userId,
		MainBranch:    mainBranch,
		SourceBranch:  sourceBranch,
		TargetBranch:  targetBranch,
		CommitMessage: os.Getenv("CI_COMMIT_MESSAGE"),
	}
}

func main() {
	env := parseEnv()

	api, err := gitlab.NewClient(env.Token, gitlab.WithBaseURL(env.BaseUrl))
	if err != nil {
		log.Fatal(err)
	}

	if env.TargetBranch != env.MainBranch {
		// MainBranch всегда существует
		b, _, _ := api.Branches.GetBranch(env.Pid, env.TargetBranch)

		if b == nil {
			b, _, err := api.Branches.CreateBranch(env.Pid, &gitlab.CreateBranchOptions{
				Branch: &env.TargetBranch,
				Ref:    &env.MainBranch,
			})

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Branch '%s' created. %s", b.Name, b.Commit.WebURL)
		}
	}

	mrState := "opened"
	mergeRequests, _, err := api.MergeRequests.ListProjectMergeRequests(env.Pid, &gitlab.ListProjectMergeRequestsOptions{
		State:        &mrState,
		SourceBranch: &env.SourceBranch,
		TargetBranch: &env.TargetBranch,
	})

	if len(mergeRequests) > 0 {
		log.Printf("MR already exists")
		return
	}

	mr, _, err := api.MergeRequests.CreateMergeRequest(env.Pid, &gitlab.CreateMergeRequestOptions{
		Title:           &env.CommitMessage,
		SourceBranch:    &env.SourceBranch,
		TargetBranch:    &env.TargetBranch,
		AssigneeID:      &env.UserId,
		TargetProjectID: &env.Pid,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("MR '%s' (#%d) created. %s", mr.Title, mr.ID, mr.WebURL)
}

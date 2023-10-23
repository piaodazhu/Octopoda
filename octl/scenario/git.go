package scenario

import (
	"fmt"
	"octl/config"
	"octl/output"

	"github.com/go-git/go-git/v5"
	gitconf "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func GitClone(repo, user string) {
	remoteUrl := makeRemoteUrl(repo, user)
	_, err := git.PlainClone(repo, false, &git.CloneOptions{
		URL: remoteUrl,
		Auth: &http.BasicAuth{
			Username: config.GlobalConfig.Gitinfo.Username,
			Password: config.GlobalConfig.Gitinfo.Password,
		},
	})
	if err != nil {
		output.PrintFatalln(err)
	} else {
		output.PrintInfoln("CLONE DONE!")
	}
}

func GitPush(repo, user string) {
	remoteUrl := makeRemoteUrl(repo, user)
	repository, err := git.PlainOpen(repo)
	if err != nil {
		repository, err = git.PlainInit(repo, false)
		if err != nil {
			output.PrintFatalln(err)
		}
	}
	workTree, err := repository.Worktree()
	if err != nil {
		output.PrintFatalln(err)
	}
	workTree.Add(".")
	workTree.Commit("default msg", &git.CommitOptions{
		AllowEmptyCommits: false,
	})
	remote, err := repository.CreateRemote(&gitconf.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteUrl},
	})
	if err != nil {
		remote, err = repository.Remote("origin")
		if err != nil {
			output.PrintFatalln(err)
		}
	}

	err = remote.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: config.GlobalConfig.Gitinfo.Username,
			Password: config.GlobalConfig.Gitinfo.Password,
		},
	})
	if err != nil {
		output.PrintFatalln(err)
	} else {
		output.PrintInfoln("PUSH REMOTE DONE!")
	}
}

func makeRemoteUrl(repo, user string) string {
	if user == "" {
		user = config.GlobalConfig.Gitinfo.Username
	}
	return fmt.Sprintf("%s/%s/%s.git", config.GlobalConfig.Gitinfo.ServeUrl, user, repo)
}

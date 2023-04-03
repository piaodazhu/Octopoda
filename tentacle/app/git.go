package app

import (
	"os"
	"tentacle/logger"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// must support rollback
func GitReset(path string, hash string, mode string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		logger.Server.Print("git.PlainOpen")
		return err
	}
	// record current hash for back up
	ref, _ := repo.Head()
	hashbackup := ref.Hash()

	wt, err := repo.Worktree()
	if err != nil {
		logger.Server.Print("repo.Worktree")
		return err
	}

	err = wt.Reset(&git.ResetOptions{
		Commit: plumbing.NewHash(hash),
		Mode: git.HardReset,
	})
	if err != nil {
		// failed, rollback
		wt.Reset(&git.ResetOptions{
			Commit: hashbackup,
			Mode: git.HardReset,
		})
		return err
	}
	return nil
}

type EmptyCommitError struct{}

func (e EmptyCommitError) Error() string { return "EmptyCommitError" }
func GitCommit(path string, message string) (Version, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		logger.Server.Print("git.PlainOpen")
		return Version{}, err
	}
	wt, err := repo.Worktree()
	if err != nil {
		logger.Server.Print("repo.Worktree")
		return Version{}, err
	}

	wt.Add("./")
	state, err := wt.Status()
	if err != nil {
		logger.Server.Print("wt.Status")
		return Version{}, err
	}
	if state.IsClean() {
		logger.Server.Print("state.IsClean")
		return Version{}, EmptyCommitError{}
	}

	commitTime := time.Now()
	hash, err := wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name: "Octopoda",
			When: commitTime,
		},
	})
	if err != nil {
		logger.Server.Print("wt.Commit")
		return Version{}, err
	}
	return Version{commitTime.Unix(), hash.String(), message}, nil
}

func GitCreate(path string) bool {
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		logger.Server.Print("os.Mkdir")
		return false
	}
	if _, err := git.PlainInit(path, true); err != nil {
		logger.Server.Print("git.PlainInit")
		return false
	}
	return true
}

// for Fix and FixAll
func gitLogs(path string, N int) ([]Version, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		logger.Server.Print("git.PlainOpen")
		return nil, err
	}
	
	iter, err := repo.Log(&git.LogOptions{
		All: true,
	})
	if err != nil {
		logger.Server.Print("repo.Log")
		return nil, err
	}
	defer iter.Close()
	
	v := []Version{}
	for i := 0; i < N; i++ {
		cmt, err := iter.Next()
		if err != nil {
			break
		}
		v = append(v, Version{
			Time: cmt.Committer.When.Unix(),
			Hash: cmt.Hash.String(),
			Msg: cmt.Message,
		})
	}
	return v, nil
}
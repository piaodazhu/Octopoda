package app

import (
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/piaodazhu/Octopoda/protocols/san"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
)

// must support rollback
func GitReset(app string, hash string, mode string) error {
	path := config.GlobalConfig.Workspace.Root + app
	repo, err := git.PlainOpen(path)
	if err != nil {
		logger.Exceptions.Print("git.PlainOpen")
		return err
	}
	// record current hash for back up
	ref, _ := repo.Head()
	hashbackup := ref.Hash()

	wt, err := repo.Worktree()
	if err != nil {
		logger.Exceptions.Print("repo.Worktree")
		return err
	}

	err = wt.Reset(&git.ResetOptions{
		Commit: plumbing.NewHash(hash),
		Mode:   git.HardReset,
	})
	if err != nil {
		// failed, rollback
		wt.Reset(&git.ResetOptions{
			Commit: hashbackup,
			Mode:   git.HardReset,
		})
		return err
	}
	return nil
}

type EmptyCommitError struct{}

func (e EmptyCommitError) Error() string { return "EmptyCommitError" }
func GitCommit(app string, message string) (san.Version, error) {
	path := config.GlobalConfig.Workspace.Root + app
	repo, err := git.PlainOpen(path)
	if err != nil {
		logger.Exceptions.Print("git.PlainOpen")
		return san.Version{}, err
	}
	wt, err := repo.Worktree()
	if err != nil {
		logger.Exceptions.Print("repo.Worktree")
		return san.Version{}, err
	}

	wt.Add("./")
	state, err := wt.Status()
	if err != nil {
		logger.Exceptions.Print("wt.Status")
		return san.Version{}, err
	}
	if state.IsClean() {
		// logger.Exceptions.Print("state.IsClean")
		return san.Version{}, EmptyCommitError{}
	}

	commitTime := time.Now()
	hash, err := wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name: "Octopoda",
			When: commitTime,
		},
		AllowEmptyCommits: true,
	})
	if err != nil {
		logger.Exceptions.Print("wt.Commit")
		return san.Version{}, err
	}
	return san.Version{Time: san.TsInt64(commitTime.Unix()), Hash: hash.String(), Msg: message}, nil
}

func GitStatus(app string) (san.Version, bool, error) {
	path := config.GlobalConfig.Workspace.Root + app
	isClean := false
	repo, err := git.PlainOpen(path)
	if err != nil {
		logger.Exceptions.Print("git.PlainOpen")
		return san.Version{}, isClean, err
	}
	wt, err := repo.Worktree()
	if err != nil {
		logger.Exceptions.Print("repo.Worktree")
		return san.Version{}, isClean, err
	}

	st, err := wt.Status()
	if err != nil {
		logger.Exceptions.Print("wt.Status")
		return san.Version{}, isClean, err
	}

	isClean = st.IsClean()
	head, err := repo.Head()
	if err != nil {
		logger.Exceptions.Print("repo.Head")
		return san.Version{}, isClean, err
	}

	cmt, err := repo.CommitObject(head.Hash())
	if err != nil {
		logger.Exceptions.Print("repo.CommitObject")
		return san.Version{}, isClean, err
	}

	return san.Version{
		Time: san.TsInt64(cmt.Committer.When.Unix()),
		Hash: cmt.Hash.String(),
		Msg:  cmt.Message,
	}, isClean, nil
}

func GitCreate(app string) bool {
	path := config.GlobalConfig.Workspace.Root + app
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		logger.Exceptions.Print("os.Mkdir: will overwrite")
		err := os.RemoveAll(path)
		if err != nil {
			logger.Exceptions.Print("os.RemoveAll: overwrite failed")
			return false
		}
		if os.Mkdir(path, os.ModePerm) != nil {
			logger.Exceptions.Print("os.Mkdir")
			return false
		}
	}
	if _, err := git.PlainInit(path, false); err != nil {
		logger.Exceptions.Print("git.PlainInit")
		return false
	}
	return true
}

// for Fix and FixAll
func gitLogs(app string, N int) ([]san.Version, error) {
	path := config.GlobalConfig.Workspace.Root + app
	repo, err := git.PlainOpen(path)
	if err != nil {
		logger.Exceptions.Print("git.PlainOpen:", err.Error())
		return nil, err
	}

	iter, err := repo.Log(&git.LogOptions{
		All: true,
	})
	if err != nil {
		logger.Exceptions.Print("repo.Log")
		return nil, err
	}
	defer iter.Close()

	v := []san.Version{}
	for i := 0; i < N; i++ {
		cmt, err := iter.Next()
		if err != nil {
			break
		}
		v = append(v, san.Version{
			Time: san.TsInt64(cmt.Committer.When.Unix()),
			Hash: cmt.Hash.String(),
			Msg:  cmt.Message,
		})
	}

	// reverse logs: commit time asc
	i, j := 0, len(v)-1
	for i < j {
		v[i], v[j] = v[j], v[i]
		i++
		j--
	}
	return v, nil
}

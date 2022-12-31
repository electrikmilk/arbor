/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"strings"

	"github.com/electrikmilk/args-parser"
	"github.com/electrikmilk/ttuy"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var repo *git.Repository

var branches []ttuy.Option
var copyBranch string

// Check to make sure that the current working directory is a git repository, if so, set repo and get branches
func checkForGit() {
	r, repoErr := git.PlainOpen(".")
	handleGit(repoErr)
	repo = r
	getBranches()
}

// Gets branches for this repo, either local or remote and loads them into memory
func getBranches() {
	if args.Using("remote") {
		remote, remoteErr := repo.Remote("origin")
		handleGit(remoteErr)
		refList, listRemoteErr := remote.List(&git.ListOptions{})
		handleGit(listRemoteErr)
		refPrefix := "refs/heads/"
		for _, ref := range refList {
			refName := ref.Name().String()
			if !strings.HasPrefix(refName, refPrefix) {
				continue
			}
			branchName := refName[len(refPrefix):]
			branches = append(branches, ttuy.Option{
				Label: branchName,
				Callback: func() {
					copyBranch = branchName
				},
			})
		}
	} else {
		iter, branchesErr := repo.Branches()
		handleGit(branchesErr)
		iterErr := iter.ForEach(func(r *plumbing.Reference) error {
			branch := strings.ReplaceAll(string(r.Name()), "refs/heads/", "")
			branches = append(branches, ttuy.Option{
				Label: branch,
				Callback: func() {
					copyBranch = branch
				},
			})
			return nil
		})
		handleGit(iterErr)
	}
}

func checkout(branch *string) {
	var branchPath string = fmt.Sprintf("refs/heads/%s", *branch)
	wt, err := repo.Worktree()
	handleGit(err)
	go ttuy.Spinner("Checking out branch "+ttuy.Style(*branch, ttuy.CyanText)+"...", ttuy.Throbber)
	err = wt.Checkout(&git.CheckoutOptions{
		Create: false,
		Force:  false,
		Branch: plumbing.ReferenceName(branchPath),
	})
	ttuy.StopSpinner()
	handleGit(err)
	ttuy.Success("Checked out branch " + ttuy.Style(*branch, ttuy.Bold))
	if args.Using("remote") {
		// ttuy.Menu("Does this pull need an SSH key?", []ttuy.Option{
		// 	{
		// 		Label: "Yes",
		// 		Callback: func() {
		// 			sshSign()
		// 		},
		// 	},
		// 	{
		// 		Label: "No",
		// 		Callback: func() {
		//
		// 		},
		// 	},
		// })
		go ttuy.Spinner("Pulling branch "+ttuy.Style(*branch, ttuy.CyanText)+"...", ttuy.Throbber)
		pullErr := wt.Pull(&git.PullOptions{
			RemoteName:        "origin",
			ReferenceName:     plumbing.ReferenceName(branchPath),
			SingleBranch:      true,
			Depth:             0,
			Auth:              nil,
			RecurseSubmodules: 0,
			Progress:          nil,
			Force:             false,
			InsecureSkipTLS:   false,
			CABundle:          nil,
		})
		ttuy.StopSpinner()
		if pullErr != nil {
			ttuy.Warnf("Failed to pull %s: %s", *branch, pullErr)
		} else {
			ttuy.Success("Pulled branch " + ttuy.Style(*branch, ttuy.Bold))
		}
	}
}

// func sshSign() {
// var publicKey *ssh.PublicKeys
// sshKeyPath := fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME"))
// sshKey, readKeyErr := ioutil.ReadFile(sshKeyPath)
// handleGit(readKeyErr)
// publicKey, keyError := ssh.NewPublicKey(sshKey)
// handleGit(keyError)
//
// ssh.NewPublicKey({})
//
// http.BasicAuth{
// 	Username: "git",
// 	Password: "",
// }
// auth = &gitssh.PublicKeys{User: "git", Signer: signer}
// }

func startBranch() {
	if args.Using("initials") {
		saveInitials()
	} else {
		getInitials()
	}
	var branchType string
	ttuy.Menu("Type of Branch", []ttuy.Option{
		{
			Label: "Hotfix",
			Callback: func() {
				branchType = "hotfix"
			},
		},
		{
			Label: "Bug",
			Callback: func() {
				branchType = "bug"
			},
		},
		{
			Label: "Enhancement",
			Callback: func() {
				branchType = "enhancement"
			},
		},
		{
			Label: "Feature",
			Callback: func() {
				branchType = "feature"
			},
		},
	})
	var reference string
	fmt.Println(ttuy.Style("Enter a ticket number, or dash seperated string describing the branch.", ttuy.Dim))
	ttuy.Ask("Reference", &reference)
	var name = fmt.Sprintf("%s/%s/%s", branchType, initials, reference)
	createBranch(&name)
}

func createBranch(branch *string) {
	ttuy.Menu("Based on Branch", branches)
	checkout(&copyBranch)
	var branchPath string = fmt.Sprintf("refs/heads/%s", *branch)
	wt, err := repo.Worktree()
	handleGit(err)
	go ttuy.Spinner("Creating branch "+ttuy.Style(*branch, ttuy.CyanText)+"...", ttuy.Throbber)
	err = wt.Checkout(&git.CheckoutOptions{
		Create: true,
		Force:  false,
		Branch: plumbing.ReferenceName(branchPath),
	})
	ttuy.StopSpinner()
	handleGit(err)
	ttuy.Success("Created branch " + ttuy.Style(*branch, ttuy.Bold))
}

func createCommit() {
	go ttuy.Spinner("Checking staged...", ttuy.Throbber)
	wt, err := repo.Worktree()
	handleGit(err)
	var stage, statusErr = wt.Status()
	handleGit(statusErr)
	var fileName string
	var stagedFile git.FileStatus
	if len(stage) > 1 {
		ttuy.StopSpinner()
		ttuy.Fail("Unable to auto-commit, more than one file is staged.")
	}
	for name, file := range stage {
		if file.Staging == git.Unmodified || file.Staging == git.Untracked || file.Worktree == git.Untracked {
			continue
		}
		if file.Staging == git.Added || file.Staging == git.Modified || file.Staging == git.Renamed || file.Staging == git.Copied || file.Staging == git.Deleted {
			fileName = name
			stagedFile = *file
			break
		}
	}
	if fileName == "" {
		ttuy.StopSpinner()
		ttuy.Fail("Unable to auto-commit, stage is empty")
	}
	ttuy.StopSpinner()
	ttuy.Success("Single staged file: " + ttuy.Style(fileName, ttuy.Bold))
	go ttuy.Spinner("Committing "+ttuy.Style(fileName, ttuy.CyanText)+"...", ttuy.Throbber)
	var _, addErr = wt.Add(fileName)
	handleGit(addErr)
	var prefix string
	switch stagedFile.Staging {
	case git.Added:
		prefix = "Add"
	case git.Deleted:
		prefix = "Delete"
	case git.Modified:
		prefix = "Update"
	case git.Renamed:
		prefix = "Rename"
	case git.Copied:
		prefix = "Copy"
	}
	var message = prefix + " " + fileName
	var _, commitErr = wt.Commit(message, &git.CommitOptions{})
	ttuy.StopSpinner()
	handleGit(commitErr)
	ttuy.Successf("Committed %s", fileName)
}

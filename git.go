/*
 * Copyright (c) 2022 Brandon Jordan
 */

package main

import (
	"fmt"
	"strings"

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
	if arg("r") == true {
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
	go ttuy.Spinner("Checking out branch "+ttuy.Style(*branch, ttuy.CyanText)+"...", ttuy.Throbber)
	var branchPath string = fmt.Sprintf("refs/heads/%s", *branch)
	wt, err := repo.Worktree()
	if err != nil {
		return
	}
	handleGit(err)
	err = wt.Checkout(&git.CheckoutOptions{
		Create: false,
		Force:  false,
		Branch: plumbing.ReferenceName(branchPath),
	})
	handleGit(err)
	ttuy.StopSpinner()
	ttuy.Success("Checked out branch " + ttuy.Style(*branch, ttuy.Bold))
	if arg("remote") {
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

func create(branch *string) {
	ttuy.Menu("Based on Branch", branches)
	checkout(&copyBranch)
	go ttuy.Spinner("Creating branch "+ttuy.Style(*branch, ttuy.CyanText)+"...", ttuy.Throbber)
	var branchPath string = fmt.Sprintf("refs/heads/%s", *branch)
	wt, err := repo.Worktree()
	handleGit(err)
	err = wt.Checkout(&git.CheckoutOptions{
		Create: true,
		Force:  false,
		Branch: plumbing.ReferenceName(branchPath),
	})
	ttuy.StopSpinner()
	ttuy.Success("Created branch " + ttuy.Style(*branch, ttuy.Bold))
}

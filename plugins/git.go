package plugins

import (
	"errors"
	"os/exec"
	"strings"
)

type Git struct {
	Remote string
}

func (g *Git) New() {
	_, err := exec.LookPath("git")
	if err != nil {
		panic("git is not installed")
	}

	// Get the remote URL
	remote, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		panic("unable to get remote URL")
	}

	g.Remote = string(remote)

	branch, err := g.Branch()
	if err != nil {
		panic("unable to get current branch")
	}

	// Pull the latest changes
	_, err = exec.Command("git", "pull", "origin", branch).Output()
	if err != nil {
		panic("unable to pull latest changes")
	}
}

func (g *Git) Status() (string, error) {
	status, err := exec.Command("git", "status").Output()
	if err != nil {
		return "", err
	}

	return string(status), nil
}

func (g *Git) Add(file string) error {
	_, err := exec.Command("git", "add", file).Output()
	if err != nil {
		return errors.New("unable to add file: " + file)
	}

	return nil
}

func (g *Git) AddAll() error {
	_, err := exec.Command("git", "add", ".").Output()
	if err != nil {
		return errors.New("unable to add all files")
	}

	return nil
}

func (g *Git) Branch() (string, error) {
	branch, err := exec.Command("git", "branch").Output()
	if err != nil {
		return "", errors.New("unable to get current branch")
	}

	lines := strings.Split(string(branch), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "*") {
			return strings.TrimSpace(strings.TrimPrefix(line, "* ")), nil
		}
	}

	return "", nil
}

func (g *Git) Commit(message string) error {
	// Check if there are any changes to commit
	status, err := g.Status()
	if err != nil {
		return err
	}

	if !strings.Contains(status, "Changes to be committed") {
		return errors.New("no changes to commit")
	}

	// Commit the changes
	_, err = exec.Command("git", "commit", "-m", message).Output()
	if err != nil {
		return errors.New("unable to commit changes")
	}

	return nil
}

func (g *Git) PushCurrentBranch() error {
	// Get the current branch
	branch, err := g.Branch()
	if err != nil {
		return err
	}

	// Push the changes
	_, err = exec.Command("git", "push", "origin", branch).Output()
	if err != nil {
		return errors.New("unable to push changes to branch: " + branch)
	}

	return nil
}

func (g *Git) PullCurrentBranch() error {
	// Get the current branch
	branch, err := g.Branch()
	if err != nil {
		return err
	}

	// Pull the changes
	_, err = exec.Command("git", "pull", "origin", branch).Output()
	if err != nil {
		return errors.New("unable to pull changes from branch: " + branch)
	}

	return nil
}

func (g *Git) Checkout(branch string) error {
	// Check if there are any changes to stash
	status, err := g.Status()
	if err != nil {
		return err
	}

	if strings.Contains(status, "Changes not staged for commit") {
		return errors.New("unable to checkout branch: " + branch + " due to uncommitted changes")
	}

	// Checkout the branch
	_, err = exec.Command("git", "checkout", branch).Output()
	if err != nil {
		return errors.New("unable to checkout branch: " + branch)
	}

	return nil
}

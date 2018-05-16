package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var (
	srcRepo   = flag.String("src", envString("GIT_SRC_REPO", ""), "the git repository to clone")
	srcBranch = flag.String("branch", envString("GIT_SRC_BRANCH", "master"), "the git repository branch")
	destRepo  = flag.String("dest", envString("GIT_DEST_REPO", ""), "the git repository to push")
	gitRoot   = flag.String("root", envString("GIT_ROOT_PATH", ""), "the git saved directory path")
)

func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}

func main() {
	flag.Parse()
	inputValidation()

}

func inputValidation() {
	if *srcRepo == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --src or $GIT_SRC_REPO must be provided\n")
		flag.Usage()
		os.Exit(1)
	}

	if *destRepo == "" {
		fmt.Fprintf(os.Stderr, "ERROR: --dest or $GIT_DEST_REPO must be provided\n")
		flag.Usage()
		os.Exit(1)
	}

	if _, err := exec.LookPath("git"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: git executable not found: %v\n", err)
		os.Exit(1)
	}
}

func runCommand(cwd, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	if cwd != "" {
		cmd.Dir = cwd
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running command: %v: %q", err, string(output))
	}
	return string(output), nil
}

func cloneRepo(repo, branch, gitRoot string) error {
	args := []string{"clone", "-b", branch, repo, gitRoot}
	_, err := runCommand("", "git", args...)
	if err != nil {
		return err
	}
	log.Printf("Cloned repo %s", repo)
	return nil
}

func pullRepo(branch, gitRoot string) error {
	args := []string{"pull", "origin", branch}
	_, err := runCommand(gitRoot, "git", args...)
	if err != nil {
		return err
	}
	return nil
}

func addDest(destRepo, gitRoot string) error {
	args := []string{"remote", "show", "dest"}
	_, err := runCommand(gitRoot, "git", args...)
	if err == nil {
		return nil
	}

	args = []string{"remote", "add", "dest", destRepo}
	_, err = runCommand(gitRoot, "git", args...)
	if err != nil {
		return fmt.Errorf("failure to add remote repository: %v", err)
	}
	return nil
}

func pushRepo(branch, destRepo, gitRoot string) error {
	args := []string{"push", "dest", branch}
	_, err := runCommand(gitRoot, "git", args...)
	if err != nil {
		return fmt.Errorf("failure to push to %s: %v", destRepo, err)
	}
	return nil
}

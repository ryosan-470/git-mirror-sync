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
	//	fmt.Printf("srcRepo = %s\nsrcBranch = %s\ndestRepo = %s\n", *srcRepo, *srcBranch, *destRepo)
	err := cloneRepo("https://github.com/ryosan-470/git-mirror-sync", "master", "/tmp/git/")
	if err != nil {
		fmt.Println(err)
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

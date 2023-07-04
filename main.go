package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	const Green = "\033[32m%s\033[0m"
	if len(os.Args) > 1 {
		var commitArg string
		commitArg = os.Args[1]

		if randomMessage, err := getRandomCommitMessage(); verifyCommitMessage(commitArg) {
			if err != nil {
				log.Fatalf("Failed to get random commit message: %v", err)
			}
			fmt.Printf(Green, randomMessage)
			if err := askUser(randomMessage); err != nil {
				log.Fatalf("Failed to ask user: %v", err)
			}
			return
		} else {
			if err := rewriteCommitMessage(randomMessage); err != nil {
				log.Fatalf("Failed to rewrite commit message: %v", err)
			}
		}
	}
}

func getRandomCommitMessage() (string, error) {
	resp, err := http.Get("https://whatthecommit.com/index.txt")
	if err != nil {
		return "", err
	}
	defer func() {
		if resp.Body.Close() != nil {
			return
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func rewriteCommitMessage(message string) error {
	cmd := exec.Command("git", "rebase", "-i", "HEAD~1")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run 'git rebase -i HEAD~1': %v", err)
	}

	cmd = exec.Command("git", "commit", "--amend", "-m", message)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run 'git commit --amend': %v", err)
	}

	return nil
}

func askUser(message string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nDo you want to modify the commit message? (y/n): ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	text = strings.TrimSpace(text)
	if text == "y" {
		message = strings.TrimSpace(message)
		if err := rewriteCommitMessage(message); err != nil {
			return err
		}
	} else if text != "n" {
		fmt.Println("Please enter 'y' or 'n'.")
		if err := askUser(message); err != nil {
			return err
		}
	}
	return nil
}

func verifyCommitMessage(str string) bool {
	return regexp.MustCompile(`^(feat|enhancement|fix|docs|style|refactor|perf|test|chore|ci|build)\(([A-Za-z0-9_-]+)\) ?: ?.+$`).MatchString(str)
}

package main

import (
	"bufio"
	"os"
	"strings"
)

var blacklistMap map[string]bool

func loadBlackList(filename string) (map[string]bool, error) {
	list := make(map[string]bool)
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" && !strings.HasPrefix(domain, "#") {
			list[domain] = true
		}
	}
	return list, scanner.Err()
}

package picker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

func expandPath(path string) string {
	expanded := path
	if strings.Contains(path, "~") {
		home, _ := os.UserHomeDir()
		expanded = filepath.Join(home, path[1:])
	}
	return expanded 
}

// TODO: pull config and ignore any dirs that dont need to be pulled
func findDirectories(searchPaths []string) []list.Item {
	dirList := []list.Item{}
	for _, root := range searchPaths {
		expandedRoot:= expandPath(root)
		filepath.WalkDir(expandedRoot, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("error while walking dir: %v\n", err)
			}

			if !d.IsDir() {
				return nil
			}

			item := Item{
				Name: d.Name(), // use utils to get directory name
				Path: path, 
			}
			dirList = append(dirList, item)

			// Needs to be after directory has been added to keep the search to a max depth of 1
			if d.IsDir() && path != expandedRoot {
				return filepath.SkipDir
			}

			return nil
		})
	}
	return dirList
}

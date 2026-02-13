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

func findDirectories(searchPaths []string) []list.Item {
	m := make(map[string]Item)
	dirList := []list.Item{}
	for _, root := range searchPaths {
		expandedRoot := expandPath(root)
		filepath.WalkDir(expandedRoot, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("error while walking dir: %v\n", err)
				return err
			}

			if !d.IsDir() {
				return nil
			}

			sessionName := d.Name()
			if strings.Contains(d.Name(), ".") {
				sessionName = "_" + sessionName[1:]
			}
			if ok := m[sessionName]; ok == (Item{}) {
				item := Item{
					SessionName: sessionName,
					Path:        path,
				}
				dirList = append(dirList, item)
				m[sessionName] = item
			}

			// Needs to be after directory has been added to keep the search to a max depth of 1
			if d.IsDir() && path != expandedRoot {
				return filepath.SkipDir
			}

			return nil
		})
	}
	return dirList
}

package files

import "os"

func CreateFolder(path string) {
	// Create folder if not exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}

// createFile creates a file given a path and a content
func CreateFile(path string, content string) error {
	// Create file if not exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()
		f.WriteString(content)
	}
	return nil
}

package path

import (
	"os"
)

// ReWriteFile renew the file
func ReWriteFile(fileName string, content []byte) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0644)
	if err == nil {
		_, err = f.Write(content)
		if err != nil {
			return err
		}
		return f.Close()
	}
	return err
}

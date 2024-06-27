package command

import "os"

func getFilePathInSSHFolder(file string) string {
	return os.ExpandEnv("$HOME/.ssh/" + file)
}

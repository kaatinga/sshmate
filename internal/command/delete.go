package command

import (
	"bytes"
	"fmt"
	"os"
)

func DeleteKeyPair(selectedItem KeyPair) error {
	// delete the key pair
	if err := os.Remove(selectedItem.PrivateFile); err != nil {
		fmt.Printf("The private key file '%s' could not be deleted: %v\n", selectedItem.PrivateFile, err)
	}

	if selectedItem.PublicFile != "" {
		if err := os.Remove(selectedItem.PublicFile); err != nil {
			fmt.Printf("The public key file '%s' could not be deleted: %v\n", selectedItem.PublicFile, err)
		}
	}

	configData, err := os.ReadFile(getFilePathInSSHFolder("config"))
	if err != nil {
		return fmt.Errorf("failed to read ~/.ssh/config file: %w", err)
	}

	buffer := bytes.NewBuffer(nil)
	lines := bytes.Split(configData, []byte("\n"))
	hostFound := false
	for _, line := range lines {
		if after, found := bytes.CutPrefix(line, hostPrefix); found {
			if bytes.Equal(after, []byte(selectedItem.Host)) {
				hostFound = true
				continue
			}
			hostFound = false
		}
		if hostFound {
			continue
		}

		buffer.Write(line)
		buffer.WriteByte('\n')
	}

	if err = os.WriteFile(getFilePathInSSHFolder("config"), buffer.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write ~/.ssh/config file: %w", err)
	}

	return nil
}

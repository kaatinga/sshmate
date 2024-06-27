package command

import (
	"bufio"
	"bytes"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

type KeyType byte

func (kt KeyType) String() string {
	return [...]string{"unknown", "ssh-rsa", "ssh-dss", "sh-ed25519", "ecdsa-sha2-nistp256"}[kt]
}

const (
	Unknown KeyType = iota
	RSA
	DSA
	ED25519
	ECDSA
)

type KeyPair struct {
	Type        KeyType
	PrivateFile string
	PublicFile  string
	Host        string
	Available   bool
}

func (k *KeyPair) CheckPublicKey() error {
	publicFilePath := getFilePathInSSHFolder(k.PrivateFile + ".pub")
	file, err := os.Open(publicFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to open public key file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return fmt.Errorf("failed to read public key file")
	}

	lines := strings.Split(scanner.Text(), " ")
	if len(lines) < 2 {
		return fmt.Errorf("invalid public key file")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	switch lines[0] {
	case ssh.KeyAlgoRSA:
		k.Type = RSA
	case ssh.KeyAlgoED25519:
		k.Type = ED25519
	case ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521:
		k.Type = ECDSA
	case ssh.KeyAlgoDSA:
		k.Type = DSA
	default:
		return fmt.Errorf("unsupported key type: %s", lines[0])
	}
	k.PublicFile = k.PrivateFile + ".pub"

	defer file.Close()

	return nil
}

func GetKeyPairs() ([]KeyPair, error) {
	keyPairs, err := getConfigKeyPairs()
	if err != nil {
		return nil, err
	}

	// check availability of the private keys
	for i := range keyPairs {
		getKeyTypeAndCheckAvailability(keyPairs, i)
		if !keyPairs[i].Available {
			continue
		}

		if err := keyPairs[i].CheckPublicKey(); err != nil {
			return nil, err
		}
	}

	return keyPairs, nil
}

func getKeyTypeAndCheckAvailability(keyPairs []KeyPair, i int) {
	filePath := getFilePathInSSHFolder(keyPairs[i].PrivateFile)
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		fmt.Println("failed to open private key file:", err)
		return
	}

	keyPairs[i].Available = true

	key, err := ssh.ParseRawPrivateKey(keyData)
	if err != nil {
		fmt.Println("failed to parse private key", filePath, ":", err)
		return
	}

	switch key.(type) {
	case *rsa.PrivateKey:
		keyPairs[i].Type = RSA
	case *ecdsa.PrivateKey:
		keyPairs[i].Type = ECDSA
	case *ed25519.PrivateKey:
		keyPairs[i].Type = ED25519
	case *dsa.PrivateKey:
		keyPairs[i].Type = DSA
	default:
		fmt.Println("unsupported parsed key type")
		return
	}
}

var (
	hostPrefix         = []byte("Host ")
	identityFilePrefix = []byte("IdentityFile ")
)

func getConfigKeyPairs() ([]KeyPair, error) {
	file, err := os.Open(getFilePathInSSHFolder("config"))
	if err != nil {
		return nil, fmt.Errorf("failed to read ~/.ssh/config file: %w", err)
	}
	defer file.Close()

	keyPairs := make([]KeyPair, 0)
	scanner := bufio.NewScanner(file)
	keyPair := KeyPair{}
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		trimmedLine := bytes.TrimSpace(bytes.TrimPrefix(line, hostPrefix))

		// we detected new HOST section
		if !bytes.Equal(line, trimmedLine) {
			keyPair = KeyPair{
				Host: string(trimmedLine),
			}
			continue
		}

		// we are in the middle of a key pair
		identitySetting := bytes.TrimSpace(bytes.TrimPrefix(line, identityFilePrefix))
		if bytes.Equal(line, identitySetting) {
			continue
		}

		// we detected a new key pair
		keyPair.PrivateFile = filepath.Base(string(identitySetting))
		keyPairs = append(keyPairs, keyPair)
		continue
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return keyPairs, nil
}

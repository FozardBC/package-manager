package ssh

import (
	"errors"
	"fmt"
	"io"
	"os"
	verManager "pm/internal/versionManager"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

var (
	RemotePath = "/opt/"

	ErrNotFound = errors.New("no packages found for %s with version constraint")
)

type Client struct {
	sshClient *ssh.Client
	sftp      *sftp.Client
}

func New(host, user, keyPath string) (*Client, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshClient, err := ssh.Dial("tcp", host+":2222", config)
	if err != nil {
		return nil, err
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		sshClient.Close()
		return nil, err
	}

	return &Client{
		sshClient: sshClient,
		sftp:      sftpClient,
	}, nil
}

func (c *Client) ListDir(RemotePath string) ([]string, error) {
	files, err := c.sftp.ReadDir(RemotePath)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, f := range files {
		if !f.IsDir() {
			names = append(names, f.Name())
		}
	}
	return names, nil
}

func (c *Client) Upload(local, remote string) error {
	dst, err := c.sftp.Create(remote)
	if err != nil {
		return err
	}
	defer dst.Close()

	src, err := os.Open(local)
	if err != nil {
		return err
	}
	defer src.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err

	}
	return nil
}

func (c *Client) Download(remote, local string) error {
	src, err := c.sftp.Open(remote)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(local)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (c *Client) Search(packageName, versionConstraint string) (string, error) {

	files, err := c.ListDir(RemotePath)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(fmt.Sprintf(`^%s[_-]([0-9]+(?:\.[0-9]+)*)\.tar\.gz$`, regexp.QuoteMeta(packageName)))

	var candidates []struct {
		Filename string
		Version  verManager.Version
	}

	for _, f := range files {
		matches := re.FindStringSubmatch(f)
		if len(matches) < 2 {
			continue
		}
		verStr := matches[1]
		parsedVer, err := verManager.ParseVersion(verStr)
		if err != nil {
			continue
		}
		candidates = append(candidates, struct {
			Filename string
			Version  verManager.Version
		}{
			Filename: f,
			Version:  parsedVer,
		})
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("%w:%s'%s'", ErrNotFound, packageName, versionConstraint)
	}

	// sorting by version less to greater
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Version[0] < candidates[j].Version[0] ||
			(candidates[i].Version[0] == candidates[j].Version[0] && candidates[i].Version[1] < candidates[j].Version[1]) ||
			(candidates[i].Version[0] == candidates[j].Version[0] && candidates[i].Version[1] == candidates[j].Version[1])
	})

	// try to parse condition: >=1.10, <=2.0, 1.10 и т.д.
	targetVerStr := strings.TrimPrefix(versionConstraint, ">=")
	if strings.HasPrefix(versionConstraint, ">=") {
		ver, err := verManager.ParseVersion(targetVerStr)
		if err != nil {
			return "", err
		}
		// try to find first ver which >= ver
		for i := len(candidates) - 1; i >= 0; i-- {
			if candidates[i].Version.GreaterEqual(ver) {
				return candidates[i].Filename, nil
			}
		}
		return "", fmt.Errorf("no version >= %s found", targetVerStr)
	}

	targetVerStr = strings.TrimPrefix(versionConstraint, "<=")
	if strings.HasPrefix(versionConstraint, "<=") {
		ver, err := verManager.ParseVersion(targetVerStr)
		if err != nil {
			return "", err
		}

		for i := len(candidates) - 1; i >= 0; i-- {
			if candidates[i].Version.LessEqual(ver) {
				return candidates[i].Filename, nil
			}
		}
		return "", fmt.Errorf("no version <= %s found", targetVerStr)
	}

	if versionConstraint != "" {
		ver, err := verManager.ParseVersion(versionConstraint)
		if err != nil {
			return "", err
		}
		for _, c := range candidates {
			if c.Version == ver {
				return c.Filename, nil
			}
		}
		return "", fmt.Errorf("exact version %s not found", versionConstraint)
	}

	return candidates[len(candidates)-1].Filename, nil
}

func (c *Client) Close() {
	c.sftp.Close()
	c.sshClient.Close()
}

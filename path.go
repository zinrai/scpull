package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// localPath maps a remote path to its local destination under the host
// directory, preserving the remote directory structure. A remote path of
// "/var/log/messages" fetched from host "192.168.2.5" maps to
// "192.168.2.5/var/log/messages".
//
// The same host always maps to the same directory regardless of which
// profile requested it, so files collected across profiles accumulate
// under one directory per host.
//
// localPath rejects remote paths that are not absolute or that contain a
// ".." element, since either could direct output outside the intended
// host directory.
func localPath(host, remote string) (string, error) {
	if !strings.HasPrefix(remote, "/") {
		return "", fmt.Errorf("remote path %q is not absolute", remote)
	}

	// Reject ".." on the raw path, before Clean. filepath.Clean would
	// resolve a path like "/var/../../etc/passwd" to "/etc/passwd",
	// erasing the ".." and silently changing which file is fetched.
	for _, elem := range strings.Split(remote, "/") {
		if elem == ".." {
			return "", fmt.Errorf("remote path %q contains '..'", remote)
		}
	}

	clean := filepath.Clean(remote)

	// Strip the leading slash so the remote structure nests under the
	// host directory rather than resolving to the filesystem root.
	rel := strings.TrimPrefix(clean, "/")

	return filepath.Join(host, rel), nil
}

// ensureDir creates the parent directory of dst if it does not exist.
// Existing directories are left in place so that files fetched earlier
// (possibly via a different profile) are preserved.
func ensureDir(dst string) error {
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create directory %q: %w", dir, err)
	}
	return nil
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// scpArgs builds the argument list for a single scp invocation. It is
// separated from the exec call so the argument construction can be
// tested by comparing slices.
//
// The remote source is "user@host:remote" and the local destination is
// dst. Arguments are passed to scp as a slice, never through a shell, so
// values containing shell metacharacters cannot be interpreted.
func scpArgs(user, host, remote, dst string) []string {
	src := fmt.Sprintf("%s@%s:%s", user, host, remote)
	return []string{"-p", src, dst}
}

// fetchFile fetches a single remote path from host into the local
// destination derived from the host and remote path. It creates the
// destination directory, then runs scp with stdout and stderr connected
// to the process so progress and errors are visible.
func fetchFile(user, host, remote string) error {
	dst, err := localPath(host, remote)
	if err != nil {
		return err
	}

	if err := ensureDir(dst); err != nil {
		return err
	}

	args := scpArgs(user, host, remote, dst)

	// Print the exact command being run. args drives both this line and
	// the exec call below, so the log always matches what is executed.
	fmt.Fprintf(os.Stderr, "scp %s\n", strings.Join(args, " "))

	cmd := exec.Command("scp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scp %s@%s:%s: %w", user, host, remote, err)
	}

	return nil
}

package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/nats-io/jsm.go/natscontext"
	"github.com/google/shlex"
)

func loadContext(softFail bool) error {
	opts := DefaultOptions

	ctxOpts := []natscontext.Option{
		natscontext.WithServerURL(opts.Servers),
		natscontext.WithCreds(opts.Creds),
		natscontext.WithNKey(opts.Nkey),
		natscontext.WithCertificate(opts.TlsCert),
		natscontext.WithKey(opts.TlsKey),
		natscontext.WithCA(opts.TlsCA),
	}

	if opts.TlsFirst {
		ctxOpts = append(ctxOpts, natscontext.WithTLSHandshakeFirst())
	}

	if opts.Username != "" && opts.Password == "" {
		ctxOpts = append(ctxOpts, natscontext.WithToken(opts.Username))
	} else {
		ctxOpts = append(ctxOpts, natscontext.WithUser(opts.Username), natscontext.WithPassword(opts.Password))
	}

	var err error

	exist, _ := fileAccessible(opts.CfgCtx)

	if exist && strings.HasSuffix(opts.CfgCtx, ".json") {
		opts.Config, err = natscontext.NewFromFile(opts.CfgCtx, ctxOpts...)
	} else {
		opts.Config, err = natscontext.New(opts.CfgCtx, !SkipContexts, ctxOpts...)
	}

	if err != nil && softFail {
		opts.Config, err = natscontext.New(opts.CfgCtx, false, ctxOpts...)
	}

	return err
}

func fileAccessible(f string) (bool, error) {
	stat, err := os.Stat(f)
	if err != nil {
		return false, err
	}

	if stat.IsDir() {
		return false, fmt.Errorf("is a directory")
	}

	file, err := os.Open(f)
	if err != nil {
		return false, err
	}
	file.Close()

	return true, nil
}

// Split the string into a command and its arguments.
func splitCommand(s string) (string, []string, error) {
	cmdAndArgs, err := shlex.Split(s)
	if err != nil {
		return "", nil, err
	}

	cmd := cmdAndArgs[0]
	args := cmdAndArgs[1:]
	return cmd, args, nil
}

// Edit the file at filepath f using the environment variable EDITOR command.
func editFile(f string) error {
	rawEditor := os.Getenv("EDITOR")
	if rawEditor == "" {
		return fmt.Errorf("set EDITOR environment variable to your chosen editor")
	}

	editor, args, err := splitCommand(rawEditor)
	if err != nil {
		return fmt.Errorf("could not parse EDITOR: %v", rawEditor)
	}

	args = append(args, f)
	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not edit file %v: %s", f, err)
	}

	return nil
}
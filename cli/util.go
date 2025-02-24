package cli

import (
    "fmt"
    "github.com/google/shlex"
    "os"
    "os/exec"
)

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

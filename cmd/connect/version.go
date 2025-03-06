package main

import (
    "fmt"
    "github.com/choria-io/fisk"
)

var (
    Version        = "0.0.0"
    CommitHash     = "n/a"
    BuildTimestamp = "n/a"
)

func BuildVersion() string {
    return fmt.Sprintf("%s-%s (%s)", Version, CommitHash, BuildTimestamp)
}

func configureVersionCommand(parentCmd *fisk.Application) {
    parentCmd.Command("version", "Show version information").Action(func(pc *fisk.ParseContext) error {
        fmt.Println(BuildVersion())
        return nil
    })
}

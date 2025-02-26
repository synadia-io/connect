package cli

import (
    "github.com/choria-io/fisk"
)

type command struct {
    Name    string
    Order   int
    Command func(app commandHost)
}

type commandHost interface {
    Command(name string, help string) *fisk.CmdClause
}

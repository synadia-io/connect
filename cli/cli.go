package cli

import (
	"github.com/choria-io/fisk"
)

type commandHost interface {
	Command(name string, help string) *fisk.CmdClause
}

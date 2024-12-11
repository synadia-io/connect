package cli

import (
	"context"
	"github.com/choria-io/fisk"
	glog "log"
	"sort"
	"sync"
	"time"
)

type command struct {
	Name    string
	Order   int
	Command func(app commandHost)
}

type commandHost interface {
	Command(name string, help string) *fisk.CmdClause
}

// Logger provides a pluggable logger implementation
type Logger interface {
	Printf(format string, a ...any)
	Fatalf(format string, a ...any)
	Print(a ...any)
	Fatal(a ...any)
	Println(a ...any)
}

var (
	commands []*command
	mu       sync.Mutex
	Version  = "development"
	log      Logger
	ctx      context.Context
)

func registerCommand(name string, order int, c func(app commandHost)) {
	mu.Lock()
	commands = append(commands, &command{name, order, c})
	mu.Unlock()
}

// SkipContexts used during tests
var SkipContexts bool

func SetVersion(v string) {
	mu.Lock()
	defer mu.Unlock()

	Version = v
}

// SetLogger sets a custom logger to use
func SetLogger(l Logger) {
	mu.Lock()
	defer mu.Unlock()

	log = l
}

// SetContext sets the context to use
func SetContext(c context.Context) {
	mu.Lock()
	defer mu.Unlock()

	ctx = c
}

func commonConfigure(cmd commandHost, cliOpts *Options) error {
	if cliOpts != nil {
		DefaultOptions = cliOpts
	} else {
		DefaultOptions = &Options{
			Timeout: 5 * time.Second,
		}
	}

	ctx = context.Background()
	log = goLogger{}

	sort.Slice(commands, func(i int, j int) bool {
		if commands[i].Order != commands[j].Order {
			return commands[i].Order < commands[j].Order
		} else {
			return commands[i].Name < commands[j].Name
		}
	})

	for _, c := range commands {
		c.Command(cmd)
	}

	return nil
}

// ConfigureInApp attaches the cli commands to app, prepare will load the context on demand and should be true unless override nats,
// manager and js context is given in a custom PreAction in the caller.  Disable is a list of command names to skip.
func ConfigureInApp(app *fisk.Application, cliOpts *Options, prepare bool) (*Options, error) {
	err := commonConfigure(app, cliOpts)
	if err != nil {
		return nil, err
	}

	if prepare {
		app.PreAction(preAction)
	}

	return DefaultOptions, nil
}

func preAction(_ *fisk.ParseContext) (err error) {
	loadContext(true)
	return nil
}

type goLogger struct{}

func (goLogger) Fatalf(format string, a ...any) { glog.Fatalf(format, a...) }
func (goLogger) Printf(format string, a ...any) { glog.Printf(format, a...) }
func (goLogger) Print(a ...any)                 { glog.Print(a...) }
func (goLogger) Println(a ...any)               { glog.Println(a...) }
func (goLogger) Fatal(a ...any)                 { glog.Fatal(a...) }

func opts() *Options {
	return DefaultOptions
}

package cli

import (
    "fmt"
    "github.com/choria-io/fisk"
    "github.com/nats-io/nats.go"
    "os"
    "os/signal"
    "strings"
    "syscall"
)

type logsCommand struct {
    opts *Options
}

func ConfigureLogsCommand(parentCmd commandHost, opts *Options) {
    c := &logsCommand{
        opts: opts,
    }

    parentCmd.Command("logs", "View Logs").Alias("log").Action(c.logs)
}

func (c *logsCommand) logs(pc *fisk.ParseContext) error {
    appCtx, err := LoadOptions(c.opts)
    fisk.FatalIfError(err, "failed to load options")
    defer appCtx.Close()

    fmt.Println("Capturing logs for all connectors. Press Ctrl+C to stop.")
    fmt.Println(fmt.Sprintf("$NEX.logs.%s.>", appCtx.Client.Account()))
    sub, err := appCtx.Nc.Subscribe(fmt.Sprintf("$NEX.logs.%s.>", appCtx.Client.Account()), func(msg *nats.Msg) {
        sp := strings.Split(msg.Subject, ".")
        if len(sp) > 5 {
            return
        }

        instanceId := sp[3]
        level := sp[4]
        line := string(msg.Data)

        // Skip metrics
        if level == "metrics" {
            return
        }

        fmt.Printf("%s %s\n", instanceId, line)
    })
    defer sub.Unsubscribe()

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
    <-sigs

    return nil
}

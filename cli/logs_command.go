package cli

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/choria-io/fisk"
	"github.com/nats-io/nats.go"
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
	subject := fmt.Sprintf("$NEX.FEED.%s.logs.>", appCtx.Client.Account())
	fmt.Println(subject)
	sub, err := appCtx.Nc.Subscribe(subject, func(msg *nats.Msg) {
		sp := strings.Split(msg.Subject, ".")
		if len(sp) > 6 {
			return
		}

		instanceId := sp[4]
		level := sp[5]
		line := string(msg.Data)

		// Skip metrics
		if level == "metrics" {
			return
		}

		fmt.Printf("%s %s\n", instanceId, line)
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %w", err)
	}
	defer func() { _ = sub.Unsubscribe() }()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs

	return nil
}

package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log/slog"
	"strings"
)

const (
	GroupsHeader = "X-NATS-LOG-GROUPS"
	LevelHeader  = "X-NATS-LOG-LEVEL"
)

func NewNatsHandler(nc *nats.Conn, subjectPrefix string, levelOffset slog.Level) slog.Handler {
	return &natsHandler{
		levelOffset:   levelOffset,
		subjectPrefix: subjectPrefix,
		nc:            nc,
		attributes:    nil,
		groups:        nil,
	}
}

type natsHandler struct {
	nc            *nats.Conn
	levelOffset   slog.Level
	attributes    []slog.Attr
	subjectPrefix string
	groups        []string
}

func (n *natsHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= n.levelOffset
}

func (n *natsHandler) Handle(ctx context.Context, record slog.Record) error {
	data := map[string]any{
		"level":     record.Level.String(),
		"message":   record.Message,
		"timestamp": record.Time.UnixMilli(),
		"groups":    n.groups,
	}

	// -- add the attributes
	for _, attr := range n.attributes {
		data[attr.Key] = attr.Value
	}

	// -- construct the subject
	subject := fmt.Sprintf("%s.%s", n.subjectPrefix, record.Level.String())

	// -- encode the data
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := nats.NewMsg(subject)
	msg.Data = b
	msg.Header.Add(GroupsHeader, strings.Join(n.groups, ","))
	msg.Header.Add(LevelHeader, record.Level.String())

	// -- publish the message
	return n.nc.PublishMsg(msg)
}

func (n *natsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newl := *n
	newAttrs := make([]slog.Attr, len(newl.attributes)+len(attrs))
	copy(newAttrs, newl.attributes)
	copy(newAttrs[len(newl.attributes):], attrs)
	newl.attributes = newAttrs
	return &newl
}

func (n *natsHandler) WithGroup(name string) slog.Handler {
	newl := *n
	newGroups := make([]string, len(newl.groups)+1)
	copy(newGroups, newl.groups)
	newGroups[len(newl.groups)] = name
	newl.groups = newGroups
	return &newl
}

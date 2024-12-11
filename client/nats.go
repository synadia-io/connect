package client

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"time"
)

func getUserInfo(nc *nats.Conn) (*server.UserInfo, error) {
	resp, err := nc.Request("$SYS.REQ.USER.INFO", nil, time.Second)
	if err != nil {
		return nil, fmt.Errorf("could not get user info: %s", err)
	}

	var res = struct {
		Data   *server.UserInfo  `json:"data"`
		Server server.ServerInfo `json:"server"`
		Error  *server.ApiError  `json:"error"`
	}{}

	if err = json.Unmarshal(resp.Data, &res); err != nil {
		return nil, fmt.Errorf("could not parse user info: %s", err)
	}

	if res.Error != nil {
		return nil, fmt.Errorf("could not get user info: %s", res.Error)
	}

	return res.Data, nil
}

// slackService
package service

import (
	"bytes"
	"crypto-trading-bot-go/core"
	"fmt"
	"io"
	"net/http"
)

func SendSlack(message string) {
	if !core.Config.Slack.Enable {
		return
	}

	jsonData := []byte("{'channel': '" + core.Config.Slack.Channel + "', 'text': '" + message + "'}")

	response, error := http.Post(core.Config.Slack.Webhook, "application/json", bytes.NewBuffer(jsonData))
	if error != nil {
		panic(error)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	Logger.Debug(fmt.Sprintf("result: %s", string(body)))
}

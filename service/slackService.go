// slackService
package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func SendSlack(message string) {
	jsonData := []byte("{'channel': '" + Config.Slack.Channel + "', 'text': '" + message + "'}")

	response, error := http.Post(Config.Slack.Webhook, "application/json", bytes.NewBuffer(jsonData))
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

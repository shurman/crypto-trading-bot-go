// slackService
package core

import (
	"bytes"
	"net/http"
)

func sendSlack(message string) {
	url := Config.Slack.Webhook
	channel := Config.Slack.Channel

	jsonData := []byte("payload={'channel': '" + channel + "', 'username': 'Trading Signal Bot', 'text': '" + message + "'}")

	request, error := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

}

// slackService
package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendSlack(message string) {
	jsonData := []byte("{'channel': '" + Config.Slack.Channel + "', 'username': 'Signal Bot', 'text': '" + message + "'}")

	response, error := http.Post(Config.Slack.Webhook, "application/json", bytes.NewBuffer(jsonData))
	if error != nil {
		panic(error)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	Logger.Debug(fmt.Sprintf("result: %s", string(body)))
}

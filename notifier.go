package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Notifier is responsible for interacting with workplace API
type Notifier struct {
	token  string
	thread string
}

// Notify will send a message to workchat
func (n *Notifier) Notify(message string) (err error) {
	type payloadMessage struct {
		Recipient struct {
			ThreadKey string `json:"thread_key,omitempty"`
		} `json:"recipient,omitempty"`
		Message struct {
			Text string `json:"text,omitempty"`
		} `json:"message,omitempty"`
	}

	url := fmt.Sprintf("https://graph.facebook.com/v3.2/me/messages?access_token=%s", n.token)

	payload := payloadMessage{}
	payload.Recipient.ThreadKey = n.thread
	payload.Message.Text = message

	buffer, err := json.Marshal(&payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(buffer))
	if err != nil {
		return fmt.Errorf("failed to create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to post event to %s with status code %d", url, resp.StatusCode)
	}

	return nil
}

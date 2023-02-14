package notification

import (
	"context"

	"github.com/slack-go/slack"

	"github.com/ONSdigital/dp-integrity-checker/checker"
)

//go:generate moq -out mock/slack.go -pkg mock . SlackClient

// Notifier represents a interface for a generic notifier
type Notifier interface {
	SendCheckerResult(context.Context, *checker.Result) error
}

// SlackClient is an interface to enable mocking of the slack-go/slack.Client
type SlackClient interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}

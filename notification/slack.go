package notification

import (
	"context"
	"fmt"
	"strings"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/slack-go/slack"
	
	"github.com/ONSdigital/dp-integrity-checker/checker"
	"github.com/ONSdigital/dp-integrity-checker/config"
)

// SlackNotifier is a Notifier that uses the slack-go/slack.Client to send notifications
type SlackNotifier struct {
	Config config.Slack
	Client SlackClient
}

// GetClient returns the underlying slack client of this notifier, creating a new one if necessary
func (n *SlackNotifier) GetClient() SlackClient {
	if n.Client == nil {
		n.Client = slack.New(n.Config.ApiToken)
	}
	return n.Client
}

// SendCheckerResult sends a new Slack message for the supplied checker.Result.  It uses GetClient for the underlying
// slack client.
func (n *SlackNotifier) SendCheckerResult(ctx context.Context, result *checker.Result) error {
	logData := log.Data{"config": n.Config}
	log.Info(ctx, "sending slack notification for result", logData)

	client := n.GetClient()

	attachmentText := strings.Builder{}
	for _, inc := range result.Inconsistencies {
		attachmentText.WriteString("- ")
		attachmentText.WriteString(inc)
		attachmentText.WriteRune('\n')
	}

	attachment := slack.Attachment{
		Pretext: fmt.Sprintf("Found %d inconsistencies during integrity check\n", len(result.Inconsistencies)),
		Text:    attachmentText.String(),
		Color:   "danger",
	}

	_, _, err := client.PostMessage(
		n.Config.AlarmChannel,
		slack.MsgOptionAsUser(false),
		slack.MsgOptionUsername(n.Config.UserName),
		slack.MsgOptionIconEmoji(n.Config.AlarmEmoji),
		slack.MsgOptionAttachments(attachment),
	)
	if err != nil {
		log.Error(ctx, "unable to send message to slack api", err, logData)
		return err
	}

	log.Info(ctx, "slack notification successful", logData)

	return nil
}

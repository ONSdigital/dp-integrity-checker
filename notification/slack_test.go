package notification_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-integrity-checker/checker"
	"github.com/ONSdigital/dp-integrity-checker/config"
	"github.com/ONSdigital/dp-integrity-checker/notification"
	"github.com/ONSdigital/dp-integrity-checker/notification/mock"
)

const (
	inc1          = "inconsistency1"
	inc2          = "inconsistencyNUMBER2!@£$"
	channel       = "#some-channel"
	username      = "Some Username"
	emoji         = ":some_emoji:"
	dummyApiToken = "dummydummydummydummy"
)

var slackError = errors.New("some slack error")

var inconsistentResult = checker.Result{
	Success:         false,
	Inconsistencies: []string{inc1, inc2},
}

func TestGetClient(t *testing.T) {
	os.Clearenv()
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests

	Convey("Given a slack notifier", t, func() {

		notifier := notification.SlackNotifier{
			Config: config.Slack{
				ApiToken: dummyApiToken,
			},
		}

		Convey("When a slack client is requested", func() {
			client := notifier.GetClient()

			Convey("Then a new slack-go/slack client is created", func() {
				So(client, ShouldNotBeNil)
				So(fmt.Sprintf("%T", client), ShouldEqual, "*slack.Client")
			})
		})

		Convey("When two slack clients are requested", func() {
			client1 := notifier.GetClient()
			client2 := notifier.GetClient()

			Convey("Then the same slack-go/slack client is returned", func() {
				So(client1, ShouldNotBeNil)
				So(client2, ShouldNotBeNil)
				So(fmt.Sprintf("%T", client1), ShouldEqual, "*slack.Client")
				So(fmt.Sprintf("%T", client2), ShouldEqual, "*slack.Client")
				So(client2, ShouldEqual, client1)
			})
		})
	})

}

func TestSendCheckerResult_Slack(t *testing.T) {
	os.Clearenv()
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests

	Convey("Given a slack notifier with a non-erroring client", t, func() {
		slackClientMock := &mock.SlackClientMock{
			PostMessageFunc: func(channelID string, options ...slack.MsgOption) (string, string, error) {
				return "", "", nil
			},
		}

		notifier := notification.SlackNotifier{
			Config: config.Slack{
				UserName:     username,
				AlarmChannel: channel,
				AlarmEmoji:   emoji,
			},
			Client: slackClientMock,
		}

		Convey("When a result is sent", func() {
			err := notifier.SendCheckerResult(context.Background(), &inconsistentResult)

			Convey("Then the client should be called successfully", func() {
				So(err, ShouldBeNil)
				So(slackClientMock.PostMessageCalls(), ShouldHaveLength, 1)
				postCall := slackClientMock.PostMessageCalls()[0]
				So(postCall.ChannelID, ShouldEqual, channel)
				So(postCall.Options, ShouldHaveLength, 4)
			})
		})
	})

	Convey("Given a slack notifier with an erroring client", t, func() {
		slackClientMock := &mock.SlackClientMock{
			PostMessageFunc: func(channelID string, options ...slack.MsgOption) (string, string, error) {
				return "", "", slackError
			},
		}

		notifier := notification.SlackNotifier{
			Config: config.Slack{
				UserName:     username,
				AlarmChannel: channel,
				AlarmEmoji:   emoji,
			},
			Client: slackClientMock,
		}

		Convey("When a result is sent", func() {
			err := notifier.SendCheckerResult(context.Background(), &inconsistentResult)

			Convey("Then the client should be called successfully and an error returned from the notifier", func() {
				So(err, ShouldNotBeNil)
				So(err, ShouldResemble, slackError)
				So(slackClientMock.PostMessageCalls(), ShouldHaveLength, 1)
				postCall := slackClientMock.PostMessageCalls()[0]
				So(postCall.ChannelID, ShouldEqual, channel)
				So(postCall.Options, ShouldHaveLength, 4)
			})
		})
	})
}

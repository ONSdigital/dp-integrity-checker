// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/ONSdigital/dp-integrity-checker/notification"
	"github.com/slack-go/slack"
	"sync"
)

// Ensure, that SlackClientMock does implement notification.SlackClient.
// If this is not the case, regenerate this file with moq.
var _ notification.SlackClient = &SlackClientMock{}

// SlackClientMock is a mock implementation of notification.SlackClient.
//
//	func TestSomethingThatUsesSlackClient(t *testing.T) {
//
//		// make and configure a mocked notification.SlackClient
//		mockedSlackClient := &SlackClientMock{
//			PostMessageFunc: func(channelID string, options ...slack.MsgOption) (string, string, error) {
//				panic("mock out the PostMessage method")
//			},
//		}
//
//		// use mockedSlackClient in code that requires notification.SlackClient
//		// and then make assertions.
//
//	}
type SlackClientMock struct {
	// PostMessageFunc mocks the PostMessage method.
	PostMessageFunc func(channelID string, options ...slack.MsgOption) (string, string, error)

	// calls tracks calls to the methods.
	calls struct {
		// PostMessage holds details about calls to the PostMessage method.
		PostMessage []struct {
			// ChannelID is the channelID argument value.
			ChannelID string
			// Options is the options argument value.
			Options []slack.MsgOption
		}
	}
	lockPostMessage sync.RWMutex
}

// PostMessage calls PostMessageFunc.
func (mock *SlackClientMock) PostMessage(channelID string, options ...slack.MsgOption) (string, string, error) {
	if mock.PostMessageFunc == nil {
		panic("SlackClientMock.PostMessageFunc: method is nil but SlackClient.PostMessage was just called")
	}
	callInfo := struct {
		ChannelID string
		Options   []slack.MsgOption
	}{
		ChannelID: channelID,
		Options:   options,
	}
	mock.lockPostMessage.Lock()
	mock.calls.PostMessage = append(mock.calls.PostMessage, callInfo)
	mock.lockPostMessage.Unlock()
	return mock.PostMessageFunc(channelID, options...)
}

// PostMessageCalls gets all the calls that were made to PostMessage.
// Check the length with:
//
//	len(mockedSlackClient.PostMessageCalls())
func (mock *SlackClientMock) PostMessageCalls() []struct {
	ChannelID string
	Options   []slack.MsgOption
} {
	var calls []struct {
		ChannelID string
		Options   []slack.MsgOption
	}
	mock.lockPostMessage.RLock()
	calls = mock.calls.PostMessage
	mock.lockPostMessage.RUnlock()
	return calls
}

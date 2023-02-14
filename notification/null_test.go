package notification_test

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/ONSdigital/log.go/v2/log"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-integrity-checker/notification"
)

func TestSendCheckerResult_NullNotifier(t *testing.T) {
	os.Clearenv()
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests

	Convey("Given a null notifier", t, func() {

		notifier := notification.NullNotifier{}

		Convey("When a result is sent", func() {
			err := notifier.SendCheckerResult(context.Background(), &inconsistentResult)

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

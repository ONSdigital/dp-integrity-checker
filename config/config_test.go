package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	os.Clearenv()
	var err error
	var configuration *Config

	Convey("Given an environment with no environment variables set", t, func() {
		Convey("Then cfg should be nil", func() {
			So(cfg, ShouldBeNil)
		})

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned, and values are as expected", func() {
				configuration, err = Get() // This Get() is only called once, when inside this function
				So(err, ShouldBeNil)
				So(configuration, ShouldResemble, &Config{
					ZebedeeRoot:                "content",
					CheckPublishedPreviousDays: 1,
					SlackEnabled:               false,
					SlackConfig: Slack{
						ApiToken:     "",
						UserName:     "Integrity Checker",
						AlarmChannel: "#sandbox-alarm",
						AlarmEmoji:   ":rotating_light:",
					},
				},
				)
			})

			Convey("Then a second call to config should return the same config", func() {
				// This achieves code coverage of the first return in the Get() function.
				newCfg, newErr := Get()
				So(newErr, ShouldBeNil)
				So(newCfg, ShouldResemble, cfg)
			})
		})
	})
}

package checker_test

import (
	"context"
	"github.com/ONSdigital/dp-integrity-checker/checker"
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRun_Unreadable(t *testing.T) {
	os.Clearenv()

	Convey("Given a checker on an unreadable directory", t, func() {
		tempZebedeeRoot, err := os.CreateTemp("", "checkertest")
		So(err, ShouldBeNil)
		defer os.Remove(tempZebedeeRoot.Name())

		chk := checker.New(context.Background(), tempZebedeeRoot.Name())

		Convey("When the the checker is run", func() {
			res, err := chk.Run(context.Background())

			Convey("Then the checker should return an error", func() {
				So(res, ShouldBeNil)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "/zebedee/master: not a directory")
			})
		})
	})
}

func TestRun_EmptyWorkspace(t *testing.T) {
	os.Clearenv()

	Convey("Given a checker on an empty directory", t, func() {
		tempZebedeeRoot, err := os.MkdirTemp("", "checkertest")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempZebedeeRoot)

		chk := checker.New(context.Background(), tempZebedeeRoot)

		Convey("When the the checker is run", func() {
			res, err := chk.Run(context.Background())
			So(err, ShouldBeNil)

			Convey("Then the results should show missing dirs", func() {
				So(res, ShouldNotBeNil)
				So(res.Success, ShouldBeFalse)
				So(res.Inconsistencies, ShouldHaveLength, 2)
				So(res.Inconsistencies[0], ShouldResemble, "'zebedee/master' dir missing from zebedee root")
				So(res.Inconsistencies[1], ShouldResemble, "'zebedee/publish-log' dir missing from zebedee root")
			})
		})
	})
}

func TestRun_GoodWorkspaceNoCollections(t *testing.T) {
	os.Clearenv()

	Convey("Given a checker on an valid but empty zebedee workspace", t, func() {
		tempZebedeeRoot, err := os.MkdirTemp("", "checkertest")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempZebedeeRoot)

		addDirs(tempZebedeeRoot, "zebedee/master", "zebedee/publish-log")

		chk := checker.New(context.Background(), tempZebedeeRoot)

		Convey("When the the checker is run", func() {
			res, err := chk.Run(context.Background())
			So(err, ShouldBeNil)

			Convey("Then the results should show missing dirs", func() {
				So(res, ShouldNotBeNil)
				So(res.Success, ShouldBeTrue)
				So(res.Inconsistencies, ShouldHaveLength, 0)
			})
		})
	})
}

func addDirs(ws string, dirs ...string) {
	for _, dir := range dirs {
		os.MkdirAll(path.Join(ws, dir), 0750)
	}
}

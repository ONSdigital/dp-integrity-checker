package checker_test

import (
	"context"
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ONSdigital/dp-integrity-checker/checker"
	"github.com/ONSdigital/log.go/v2/log"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRun_Unreadable(t *testing.T) {
	os.Clearenv()
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests

	Convey("Given a checker on an unreadable directory", t, func() {
		tempZebedeeRoot, err := os.CreateTemp("", "checkertest")
		So(err, ShouldBeNil)
		defer os.Remove(tempZebedeeRoot.Name())

		chk := checker.Checker{
			ZebedeeRoot: tempZebedeeRoot.Name(),
		}

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
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests

	Convey("Given a checker on an empty directory", t, func() {
		tempZebedeeRoot, err := os.MkdirTemp("", "checkertest")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempZebedeeRoot)

		chk := checker.Checker{
			ZebedeeRoot: tempZebedeeRoot,
		}

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
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests

	Convey("Given a checker on an valid but empty zebedee workspace", t, func() {
		tempZebedeeRoot, err := os.MkdirTemp("", "checkertest")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempZebedeeRoot)

		addDirs(tempZebedeeRoot, "zebedee/master", "zebedee/publish-log")

		chk := checker.Checker{
			ZebedeeRoot: tempZebedeeRoot,
		}

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

func TestRun_MissingDirInMasterForPublishedCollection(t *testing.T) {
	os.Clearenv()
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests

	Convey("Given a checker on a workspace with a published dir missing from master", t, func() {
		tempZebedeeRoot, err := os.MkdirTemp("", "checkertest")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempZebedeeRoot)

		addDirs(tempZebedeeRoot,
			"zebedee/master/somepage/v1",
			"zebedee/publish-log/2023-02-08-08-50-col1test/somepage/v1",
			"zebedee/publish-log/2023-02-09-12-13-collection2/somepage/v1",
			"zebedee/publish-log/2023-02-09-12-13-collection2/somepage/v2",
		)

		chk := checker.Checker{
			ZebedeeRoot: tempZebedeeRoot,
		}

		// Override current time in checker package
		checker.Now = func() time.Time {
			return time.Date(2023, 2, 9, 11, 0, 0, 0, time.UTC)
		}

		Convey("When the the checker is run", func() {
			res, err := chk.Run(context.Background())
			So(err, ShouldBeNil)

			Convey("Then the results should show missing dirs", func() {
				So(res, ShouldNotBeNil)
				So(res.Success, ShouldBeFalse)
				So(res.Inconsistencies, ShouldHaveLength, 1)
				So(res.Inconsistencies[0], ShouldEqual, "dirs from collection '2023-02-09-12-13-collection2' missing from publishing master")
			})
		})
	})
}

func addDirs(ws string, dirs ...string) {
	for _, dir := range dirs {
		os.MkdirAll(path.Join(ws, dir), 0750)
	}
}

func addFile(ws, fpath string, body []byte) error {
	dir := path.Dir(fpath)
	os.MkdirAll(path.Join(ws, dir), 0750)
	return os.WriteFile(path.Join(ws, fpath), body, 0644)
}

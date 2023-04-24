package checker_test

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-integrity-checker/checker"
)

func TestCheckPublishedCollections(t *testing.T) {
	os.Clearenv()
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests

	Convey("Given a checker on a consistent workspace", t, func() {
		tempZebedeeRoot, err := os.MkdirTemp("", "checkertest")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempZebedeeRoot)

		addDirs(tempZebedeeRoot,
			"zebedee/master/somepage/v1",
			"zebedee/master/somepage/v2",
			"zebedee/master/somepage/v3",
			"zebedee/publish-log/2023-02-08-08-50-col1test/somepage/v1",
			"zebedee/publish-log/2023-02-09-12-13-collection2/somepage/v1",
			"zebedee/publish-log/2023-02-09-12-13-collection2/somepage/v2",
		)

		err = addFile(tempZebedeeRoot, "zebedee/publish-log/2023-02-08-08-50-col1test.json", []byte("{}"))
		So(err, ShouldBeNil)
		err = addFile(tempZebedeeRoot, "zebedee/publish-log/2023-02-09-12-13-collection2.json", []byte("{}"))
		So(err, ShouldBeNil)

		chk := checker.Checker{
			ZebedeeRoot:                tempZebedeeRoot,
			CheckPublishedPreviousDays: 1,
		}

		// Override current time in checker package
		checker.Now = func() time.Time {
			return time.Date(2023, 2, 9, 11, 0, 0, 0, time.UTC)
		}

		Convey("When the published collections checker is run", func() {
			valid, err := chk.CheckPublishedCollections(context.Background())

			Convey("Then the checker should return true", func() {
				So(err, ShouldBeNil)
				So(valid, ShouldBeTrue)
			})
		})
	})
}

func TestCheckPublishedCollections_FalsePositiveDeletedInLaterCollection(t *testing.T) {
	os.Clearenv()
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests

	Convey("Given a checker on a workspace with content modified then deleted in a later collection", t, func() {
		tempZebedeeRoot, err := os.MkdirTemp("", "checkertest")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempZebedeeRoot)

		addDirs(tempZebedeeRoot,
			"zebedee/master/somepage/v1",
			"zebedee/master/somepage/v2",
			"zebedee/master/somepage/v3",
			"zebedee/publish-log/2023-02-08-08-50-col1test/somepage/v1",
			"zebedee/publish-log/2023-02-08-08-50-col1test/someotherpage",
			"zebedee/publish-log/2023-02-09-12-13-collection2/somepage/v1",
			"zebedee/publish-log/2023-02-09-12-13-collection2/somepage/v2",
		)

		err = addFile(tempZebedeeRoot, "zebedee/publish-log/2023-02-08-08-50-col1test.json", []byte("{}"))
		So(err, ShouldBeNil)
		err = addFile(tempZebedeeRoot, "zebedee/publish-log/2023-02-09-12-13-collection2.json",
			[]byte(`{"pendingDeletes": [{"root": {"uri": "/someotherpage"}}]}`))
		So(err, ShouldBeNil)

		chk := checker.Checker{
			ZebedeeRoot:                tempZebedeeRoot,
			CheckPublishedPreviousDays: 1,
		}

		// Override current time in checker package
		checker.Now = func() time.Time {
			return time.Date(2023, 2, 9, 11, 0, 0, 0, time.UTC)
		}

		Convey("When the published collections checker is run", func() {
			valid, err := chk.CheckPublishedCollections(context.Background())

			Convey("Then the checker should return true", func() {
				So(err, ShouldBeNil)
				So(valid, ShouldBeTrue)
			})
		})
	})
}

func TestCheckPublishedCollections_MissingDir(t *testing.T) {
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
		err = addFile(tempZebedeeRoot, "zebedee/publish-log/2023-02-08-08-50-col1test.json", []byte("{}"))
		So(err, ShouldBeNil)
		err = addFile(tempZebedeeRoot, "zebedee/publish-log/2023-02-09-12-13-collection2.json", []byte("{}"))
		So(err, ShouldBeNil)

		chk := checker.Checker{
			ZebedeeRoot:                tempZebedeeRoot,
			CheckPublishedPreviousDays: 1,
		}

		// Override current time in checker package
		checker.Now = func() time.Time {
			return time.Date(2023, 2, 9, 11, 0, 0, 0, time.UTC)
		}

		Convey("When the published collections checker is run", func() {
			valid, err := chk.CheckPublishedCollections(context.Background())

			Convey("Then the checker should return true", func() {
				So(err, ShouldBeNil)
				So(valid, ShouldBeFalse)
			})
		})
	})
}

func TestCheckPublishedCollections_UndefinedZebedeeRoot(t *testing.T) {
	os.Clearenv()
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests

	Convey("Given a checker with no root defined", t, func() {

		chk := checker.Checker{}

		Convey("When the published collections checker is run", func() {
			valid, err := chk.CheckPublishedCollections(context.Background())

			Convey("Then the checker should return true", func() {
				So(valid, ShouldBeFalse)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "zebedee root undefined in checker")
			})
		})
	})
}

func TestGetPublishedCollections(t *testing.T) {
	os.Clearenv()
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests
	Convey("Given a checker covering three days on a workspace with collections from five days", t, func() {
		tempZebedeeRoot, err := os.MkdirTemp("", "checkertest")
		So(err, ShouldBeNil)
		defer os.RemoveAll(tempZebedeeRoot)

		addDirs(tempZebedeeRoot,
			"zebedee/publish-log/2023-02-06-08-50-col0test",
			"zebedee/publish-log/2023-02-07-08-50-col1test",
			"zebedee/publish-log/2023-02-08-08-50-col2test",
			"zebedee/publish-log/2023-02-08-11-17-col3test",
			"zebedee/publish-log/2023-02-09-12-13-collection4",
			"zebedee/publish-log/2023-02-10-12-13-collection5",
		)

		chk := checker.Checker{
			ZebedeeRoot:                tempZebedeeRoot,
			CheckPublishedPreviousDays: 2,
		}

		// Override current time in checker package
		checker.Now = func() time.Time {
			return time.Date(2023, 2, 9, 11, 0, 0, 0, time.UTC)
		}

		Convey("When GetPublishedCollections is run", func() {
			cols, err := chk.GetPublishedCollections(context.Background())

			Convey("Then only the collections from the expected date range should be returned", func() {
				So(err, ShouldBeNil)
				So(cols, ShouldNotBeNil)
				So(cols, ShouldHaveLength, 4)
				So(cols[0], ShouldEqual, "2023-02-07-08-50-col1test")
				So(cols[1], ShouldEqual, "2023-02-08-08-50-col2test")
				So(cols[2], ShouldEqual, "2023-02-08-11-17-col3test")
				So(cols[3], ShouldEqual, "2023-02-09-12-13-collection4")
			})
		})
	})
}

func TestCheckDirsInPublishedCollection_UndefinedZebedeeRoot(t *testing.T) {
	os.Clearenv()
	log.SetDestination(io.Discard, io.Discard) // Suppress logs for tests

	Convey("Given a checker with no root defined", t, func() {

		chk := checker.Checker{}

		Convey("When CheckDirsInPublishedCollection is run", func() {
			valid, err := chk.CheckDirsInPublishedCollection(context.Background(), "collection1", []string{})

			Convey("Then the function should return false with an error", func() {
				So(valid, ShouldBeFalse)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "zebedee root undefined in checker")
			})
		})
	})
}

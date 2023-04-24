package checker

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-integrity-checker/zebedee"
)

// Replacable function to allow for unit testing
var Now = func() time.Time { return time.Now().UTC() }

type allDeleted []string

func (a allDeleted) includes(s string) bool {
	for _, d := range a {
		if d == s {
			return true
		}
	}
	return false
}

func (c *Checker) CheckPublishedCollections(ctx context.Context) (bool, error) {
	log.Info(ctx, "checking consistency of published collections")
	collections, err := c.GetPublishedCollections(ctx)
	if err != nil {
		return false, err
	}

	valid := true
	allDel := make([]string, 0)
	for i := len(collections) - 1; i >= 0; i-- { // loop over from recent to oldest in case content deleted later
		deletedContent, err := c.GetDeletedContent(ctx, collections[i])
		if err != nil {
			return false, err
		}
		allDel = append(allDel, deletedContent...)

		colvalid, err := c.CheckPublishedCollection(ctx, collections[i], allDel)
		if err != nil {
			return false, err
		}
		valid = colvalid && valid
	}
	return valid, nil
}

func (c *Checker) GetDeletedContent(ctx context.Context, collection string) ([]string, error) {
	filename := path.Join(c.ZebedeeRoot, publish_log, collection+".json")
	deleted := make([]string, 0)
	col, err := zebedee.GetCollectionFromFile(filename)
	if err != nil {
		return nil, err
	}

	for _, pd := range col.PendingDeletes {
		deleted = append(deleted, pd.Root.URI)
	}

	return deleted, nil
}

func (c *Checker) CheckPublishedCollection(ctx context.Context, collection string, allDeleted allDeleted) (bool, error) {
	logdata := log.Data{"collection": collection}
	log.Info(ctx, "checking published collection", logdata)

	valid, err := c.CheckDirsInPublishedCollection(ctx, collection, allDeleted)
	if err != nil {
		log.Error(ctx, "error while checking dirs in published collection", err, logdata)
		return false, err
	}

	if !valid {
		log.Info(ctx, "inconsistency in published collection", logdata)
		return false, nil
	}
	log.Info(ctx, "published collection consistent", logdata)
	return true, nil
}

func (c *Checker) GetPublishedCollections(ctx context.Context) ([]string, error) {
	root, err := c.ensureZebedeeRoot()
	if err != nil {
		return nil, err
	}

	collections := make([]string, 0)

	startDate := Now().AddDate(0, 0, -c.CheckPublishedPreviousDays)
	log.Info(ctx, "getting list of published collections", log.Data{"start_date": startDate.Format("2006-01-02")})

	for checkDate := startDate; !checkDate.After(Now()); checkDate = checkDate.AddDate(0, 0, 1) {
		glob := path.Join(root, publish_log, checkDate.Format("2006-01-02")+"*")
		matches, err := filepath.Glob(glob)
		if err != nil {
			return nil, errors.Wrap(err, "unexpected error searching for published collections")
		}
		for _, match := range matches {
			col := path.Base(match)
			if !strings.HasSuffix(col, ".json") {
				collections = append(collections, col)
			}
		}
	}

	log.Info(ctx, "found published collections", log.Data{"collection_count": len(collections)})
	return collections, nil
}

func (c *Checker) CheckDirsInPublishedCollection(ctx context.Context, collection string, allDeleted allDeleted) (bool, error) {
	root, err := c.ensureZebedeeRoot()
	if err != nil {
		return false, err
	}

	logData := log.Data{"collection": collection}

	missingDirs := make([]string, 0)

	coldir := path.Join(root, publish_log, collection)
	err = filepath.WalkDir(coldir, func(subdir string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			relativePath := subdir[len(coldir):]

			// skip root dir of collection
			if relativePath == "" {
				return nil
			}

			inMaster, err := c.IsDirInMaster(ctx, relativePath)
			if err != nil {
				return err
			}
			if !inMaster && !allDeleted.includes(relativePath) {
				missingDirs = append(missingDirs, relativePath)
				return filepath.SkipDir // don't bother checking subdirs of missing dirs
			}
		}
		return nil
	})
	if err != nil {
		return false, errors.Wrap(err, "error walking dir tree")
	}

	if len(missingDirs) > 0 {
		logData["missing_dirs"] = missingDirs
		log.Info(ctx, "dirs from collection missing from publishing master", logData)
		c.AddInconsistency(fmt.Sprintf("dirs from collection '%s' missing from publishing master", collection))
		return false, nil
	}
	return true, nil
}

func (c *Checker) IsDirInMaster(ctx context.Context, subdir string) (bool, error) {
	root, err := c.ensureZebedeeRoot()
	if err != nil {
		return false, err
	}
	fulldir := path.Join(root, master, subdir)

	if _, err := os.Stat(fulldir); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

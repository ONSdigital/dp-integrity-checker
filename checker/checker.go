package checker

import (
	"context"
	"fmt"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
	"os"
	"path"
	"sync"
)

const (
	master      = "zebedee/master"
	publish_log = "zebedee/publish-log"
)

// Checker defines a runnable integrity checker
type Checker struct {
	ZebedeeRoot                string
	CheckPublishedPreviousDays int
	mu                         sync.Mutex
	inconsistencies            []string
}

// Result holds final results of an integrity checker run
type Result struct {
	Success         bool
	Inconsistencies []string
}

// Run runs the integrity checker
func (c *Checker) Run(ctx context.Context) (*Result, error) {
	validMaster, err := c.validateDir(ctx, master)
	if err != nil {
		return nil, err
	}

	validPublishLog, err := c.validateDir(ctx, publish_log)
	if err != nil {
		return nil, err
	}

	valid := validMaster && validPublishLog

	if valid {
		validPublishedCols, err := c.CheckPublishedCollections(ctx)
		if err != nil {
			return nil, err
		}
		valid = validPublishedCols && valid
	}

	return &Result{
		Success:         valid,
		Inconsistencies: c.inconsistencies,
	}, nil

}

func (c *Checker) validateDir(ctx context.Context, dir string) (bool, error) {
	root, err := c.ensureZebedeeRoot()
	if err != nil {
		return false, err
	}
	fulldir := path.Join(root, dir)
	logData := log.Data{
		"dir": fulldir,
	}
	if _, err := os.Stat(fulldir); err != nil {
		if !os.IsNotExist(err) {
			log.Error(ctx, "error reading dir in zebedee root", err, logData)
			return false, err
		}
		log.Info(ctx, "dir does not exist in zebedee root", logData)
		c.AddInconsistency(fmt.Sprintf("'%s' dir missing from zebedee root", dir))
		return false, nil
	}
	log.Info(ctx, "dir exists in zebedee root", logData)
	return true, nil
}

func (c *Checker) AddInconsistency(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.inconsistencies == nil {
		c.inconsistencies = make([]string, 0)
	}
	c.inconsistencies = append(c.inconsistencies, msg)
}

func (c *Checker) ensureZebedeeRoot() (string, error) {
	root := c.ZebedeeRoot
	if root == "" {
		return "", errors.New("zebedee root undefined in checker")
	}
	return root, nil
}

package checker

import (
	"context"
	"fmt"
	"github.com/ONSdigital/log.go/v2/log"
	"os"
	"path"
	"sync"
)

// Checker defines a runnable integrity checker
type Checker struct {
	mu              sync.Mutex
	ZebedeeRoot     string
	Inconsistencies []string
}

// Result holds final results of an integrity checker run
type Result struct {
	Success         bool
	Inconsistencies []string
}

// New returns an integrity checker based at the supplied zebedee root
func New(ctx context.Context, zebedeeRoot string) *Checker {
	return &Checker{
		ZebedeeRoot: zebedeeRoot,
	}
}

// Run runs the integrity checker
func (c *Checker) Run(ctx context.Context) (Result, error) {
	valid := true

	valid = c.validateDir(ctx, "zebedee/master") && valid
	valid = c.validateDir(ctx, "zebedee/publish-log") && valid

	return Result{
		Success:         valid,
		Inconsistencies: c.Inconsistencies,
	}, nil

}

func (c *Checker) validateDir(ctx context.Context, dir string) bool {
	root := c.ZebedeeRoot
	fulldir := path.Join(root, dir)
	logData := log.Data{
		"dir": fulldir,
	}
	if _, err := os.Stat(fulldir); os.IsNotExist(err) {
		log.Info(ctx, "dir does not exist in zebedee root", logData)
		c.AddInconsistency(fmt.Sprintf("'%s' dir missing from zebedee root", dir))
		return false
	}
	log.Info(ctx, "dir exists in zebedee root", logData)
	return true
}

func (c *Checker) AddInconsistency(msg string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Inconsistencies == nil {
		c.Inconsistencies = make([]string, 0)
	}
	c.Inconsistencies = append(c.Inconsistencies, msg)
}

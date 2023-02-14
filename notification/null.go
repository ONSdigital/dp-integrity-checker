package notification

import (
	"context"

	"github.com/ONSdigital/dp-integrity-checker/checker"
)

// NullNotifier is a Notifier that does nothing when SendCheckerResult is called
type NullNotifier struct{}

// SendCheckerResult does nothing
func (n *NullNotifier) SendCheckerResult(ctx context.Context, checkerResult *checker.Result) error {
	return nil
}

package engine

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

const (
	// DefaultThresholdSize is the default size threshold in bytes.
	DefaultThresholdSize = 1000000

	// DefaultThresholdDuration is the default duration threshold.
	DefaultThresholdDuration = 7 * 24 * time.Hour
)

// RunOption is the signature function for artifact cleanup feature options.
type RunOption func(RunOptions) RunOptions

// WithLogger sets the logger in options.
func WithLogger(logger Logger) RunOption {
	return func(o RunOptions) RunOptions {
		o.Logger = logger
		return o
	}
}

// WithDryRun sets the dry-run mode in options.
func WithDryRun(dryRun bool) RunOption {
	return func(o RunOptions) RunOptions {
		o.DryRun = dryRun
		return o
	}
}

// WithPaths sets the paths (regexps or raw paths) in options.
func WithPaths(paths ...string) RunOption {
	return func(o RunOptions) RunOptions {
		o.Paths = paths
		return o
	}
}

// WithThresholdSize sets the size threshold in options.
func WithThresholdSize(thresholdSize int) RunOption {
	return func(o RunOptions) RunOptions {
		o.ThresholdSize = thresholdSize
		return o
	}
}

// WithThresholdDuration sets the duration threshold in options.
func WithThresholdDuration(thresholdDuration time.Duration) RunOption {
	return func(o RunOptions) RunOptions {
		o.ThresholdDuration = thresholdDuration
		return o
	}
}

// RunOptions contains all available options for artifact cleanup feature.
type RunOptions struct {
	Logger Logger // should be unexported

	DryRun            bool
	PathRegexps       []*regexp.Regexp
	Paths             []string
	ThresholdDuration time.Duration
	ThresholdSize     int

	thresholdTime time.Time // should be unexported
}

// NewRunOptions creates a new RunOptions instance with the given options.
func NewRunOptions(opts ...RunOption) (RunOptions, error) {
	var ro RunOptions
	for _, opt := range opts {
		ro = opt(ro)
	}

	var errs []error
	for _, path := range ro.Paths {
		reg, err := regexp.Compile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid regexp '%s': %w", path, err))
		}
		ro.PathRegexps = append(ro.PathRegexps, reg)
	}

	if ro.Logger == nil {
		ro.Logger = &noopLogger{}
	}
	if ro.ThresholdSize <= 0 {
		ro.ThresholdSize = DefaultThresholdSize
	}
	if ro.ThresholdDuration <= 0 {
		ro.ThresholdDuration = DefaultThresholdDuration
	}
	ro.thresholdTime = time.Now().Add(ro.ThresholdDuration)

	return ro, errors.Join(errs...)
}

// ThresholdTime returns the time after which artifacts should be deleted.
func (ro RunOptions) ThresholdTime() time.Time {
	return ro.thresholdTime
}

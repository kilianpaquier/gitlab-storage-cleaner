package engine

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
)

// RunOption is the signature function for artifact cleanup feature options.
type RunOption func(RunOptions) RunOptions

// WithLogger sets the logger in run options.
//
// This logger can be accessed later with engine.GetLogger(context.Context) function.
//
// Logger is saved in context with runOptions.Context method to an easier access in sub-functions.
func WithLogger(logger Logger) RunOption {
	return func(o RunOptions) RunOptions {
		o.logger = logger
		return o
	}
}

// WithDryRun sets the dry-run mode in run options.
//
// When running in dry run, no actual cleaning of artifacts will be performed.
func WithDryRun(dryRun bool) RunOption {
	return func(o RunOptions) RunOptions {
		o.DryRun = dryRun
		return o
	}
}

// WithPaths sets the paths (regexps or raw paths) in run options.
//
// A path must be a valid regexp (or else NewRunOptions will return an error).
//
// Paths will be used to filter projects to clean during artifacts.Run function.
func WithPaths(paths ...string) RunOption {
	return func(o RunOptions) RunOptions {
		o.Paths = paths
		return o
	}
}

// WithThresholdDuration sets the duration threshold in run options.
//
// If the creation date of a job is less than current execution time
// minus that threshold, its artifacts will be deleted (and not already deleted by GitLab).
//
// Examples:
//
//	Given a job without a creation date
//	Then its artifacts will be deleted since it's considered a misconfiguration
//
//	Given a job created on 2025-01-10 10:00:00
//	And the threshold duration is 7 days
//	And the current time is 2025-01-17 00:00:00
//	Then its artifacts will not be deleted, because 2025-01-10 10:00:00 > 2025-01-10 00:00:00
//
//	Given an artifact expiring on 2025-01-10 00:00:00
//	And the threshold duration is 7 days
//	And the current time is 2025-01-17 12:00:00
//	Then its artifacts will not be deleted, because 2025-01-10 00:00:00 < 2025-01-10 12:00:00
//
// Default threshold is 7 days.
func WithThresholdDuration(thresholdDuration time.Duration) RunOption {
	return func(o RunOptions) RunOptions {
		o.ThresholdDuration = thresholdDuration
		return o
	}
}

// RunOptions contains all available options for artifact cleanup feature.
type RunOptions struct {
	// DryRun is a flag to enable dry-run mode.
	DryRun bool

	// Paths is a list of paths (regexps or raw paths) to filter projects to clean.
	//
	// It can be useful to only clean specific projects
	// in case given token / developer is maintainer of a lot of projects.
	Paths []string

	// ThresholdDuration is the duration threshold.
	//
	// See WithThresholdDuration option for more information.
	ThresholdDuration time.Duration

	logger  Logger
	regexps []*regexp.Regexp
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
		ro.regexps = append(ro.regexps, reg)
	}

	if ro.logger == nil {
		ro.logger = &noopLogger{}
	}
	if ro.ThresholdDuration <= 0 {
		errs = append(errs, fmt.Errorf("invalid threshold duration '%d'", ro.ThresholdDuration))
	}

	return ro, errors.Join(errs...)
}

// Context returns a new context with various options saved in it.
func (ro RunOptions) Context(parent context.Context) context.Context {
	ctx := context.WithValue(parent, LoggerKey, ro.logger)
	return ctx
}

// Regexps returns the compiled regexps from options paths.
func (ro RunOptions) Regexps() []*regexp.Regexp {
	return ro.regexps
}

package artifacts

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

//go:generate go-builder-generator generate -f options.go -s Options -d tests

// Options is the struct containing all available options in artifacts command.
type Options struct {
	DryRun            bool
	PathRegexps       []*regexp.Regexp `builder:"append"`
	Paths             []string         `builder:"append" validate:"required,dive,required"`
	ThresholdSize     uint64           `                 validate:"required,gt=0"`
	ThresholdDuration time.Duration    `                 validate:"required,gt=0"`
	ThresholdTime     time.Time
}

// EnsureDefaults ensures that all options in CleanOpts are valid
// or valorized with their default values.
func (c *Options) EnsureDefaults() error {
	var errs []error

	// validate regexps
	for _, path := range c.Paths {
		reg, err := regexp.Compile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("invalid regexp '%s': %w", path, err))
		}
		c.PathRegexps = append(c.PathRegexps, reg)
	}

	// setup threshold time
	c.ThresholdTime = time.Now().Add(c.ThresholdDuration)

	// validate overall struct
	if err := validator.New().Struct(c); err != nil {
		errs = append(errs, fmt.Errorf("failed to validate artifacts clean options: %w", err))
	}

	return errors.Join(errs...)
}

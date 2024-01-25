package config

import (
	"fmt"
	"strings"
)

type Environment string

const (
	Local      Environment = "local"
	Staging    Environment = "staging"
	Production Environment = "production"
)

func (e Environment) IsValid() error {
	switch e {
	case Local, Staging, Production:
		return nil
	default:
		return fmt.Errorf("unknown environment: %s", e)
	}
}

func (e Environment) String() string { return strings.ToLower(string(e)) }

func (e Environment) IsNot(environment Environment) bool {
	return e != environment
}

func (e Environment) Is(environment Environment) bool {
	return e == environment
}

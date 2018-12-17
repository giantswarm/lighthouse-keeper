package audit

import "github.com/giantswarm/microerror"

// invalidFlagsError is used when an attempt to write some file fails
var invalidFlagsError = &microerror.Error{
	Kind: "invalidFlagsError",
}

// IsInvalidFlagsError asserts invalidFlagsError
func IsInvalidFlagsError(err error) bool {
	return microerror.Cause(err) == invalidFlagsError
}

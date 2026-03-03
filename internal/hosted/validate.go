package hosted

import (
	"fmt"
	"regexp"
	"strings"
)

// DoltHub naming rules: alphanumeric, hyphens, underscores, 1-64 chars.
var slugRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,63}$`)

// validateSlug checks a single name component (org, db, handle) against DoltHub rules.
func validateSlug(field, value string) error {
	if !slugRe.MatchString(value) {
		return fmt.Errorf("%s must be 1-64 alphanumeric characters, hyphens, or underscores", field)
	}
	return nil
}

// validateUpstream checks that upstream is "org/db" with valid components.
func validateUpstream(value string) error {
	org, db, ok := strings.Cut(value, "/")
	if !ok || org == "" || db == "" {
		return fmt.Errorf("upstream must be in org/db format")
	}
	if err := validateSlug("upstream org", org); err != nil {
		return err
	}
	return validateSlug("upstream db", db)
}

// validateConnectFields validates all fields in a connect request.
func validateConnectFields(rigHandle, forkOrg, forkDB, upstream string) error {
	if err := validateSlug("rig_handle", rigHandle); err != nil {
		return err
	}
	if err := validateSlug("fork_org", forkOrg); err != nil {
		return err
	}
	if err := validateSlug("fork_db", forkDB); err != nil {
		return err
	}
	return validateUpstream(upstream)
}

// validateJoinFields validates all fields in a join request.
func validateJoinFields(forkOrg, forkDB, upstream string) error {
	if err := validateSlug("fork_org", forkOrg); err != nil {
		return err
	}
	if err := validateSlug("fork_db", forkDB); err != nil {
		return err
	}
	return validateUpstream(upstream)
}

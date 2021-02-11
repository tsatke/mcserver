package packet

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Validator describes a serverbound  packet that can validate its own fields
// for invalid values. Validate will be called by Decode if the decoded packet
// implements this interface.
type Validator interface {
	Serverbound
	Validate() error
}

func multiValidate(errs ...error) error {
	actualErrs := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			actualErrs = append(actualErrs, err)
		}
	}

	if len(actualErrs) == 0 {
		return nil
	} else if len(actualErrs) == 1 {
		return actualErrs[0]
	}

	var buf bytes.Buffer
	for _, err := range actualErrs {
		buf.WriteString("\t" + err.Error() + "\n")
	}
	return fmt.Errorf("multiple validation errors:\n%s", buf.String())
}

func stringMaxLength(fieldName string, maxLen int, str string) error {
	count := utf8.RuneCount([]byte(str))
	if count > maxLen {
		return fmt.Errorf("%s is too long (%d > %d)", fieldName, count, maxLen)
	}
	return nil
}

func stringNotEmpty(fieldName, str string) error {
	if str == "" {
		return fmt.Errorf("%s must not be empty", fieldName)
	}
	return nil
}

func stringNotContains(fieldName, str, forbidden string) error {
	if strings.Contains(str, forbidden) {
		return fmt.Errorf("%s must not contain %q", fieldName, forbidden)
	}
	return nil
}

func intWithinRange(fieldName string, lowerInclusive, upperInclusive, val int) error {
	if val < lowerInclusive || val > upperInclusive {
		return fmt.Errorf("%s must be within %d and %d (both inclusive), but was %d", fieldName, lowerInclusive, upperInclusive, val)
	}
	return nil
}

package core_http_utils

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
)

type dataTest struct {
	nameTest  string
	queryData string
	queryKey  string
	expected  error
}

func TestGetIntQueryParm(t *testing.T) {
	testCase := []dataTest{
		{
			nameTest:  "dot valid query",
			queryData: "abc",
			queryKey:  "limit",
			expected:  core_errors.ErrorValidation,
		},
	}

	for _, test := range testCase {
		url := fmt.Sprintf("http://localhost:8080/test?%s=%s", test.queryKey, test.queryData)
		r := httptest.NewRequest(test.queryKey, url, nil)
		_, err := GetIntQueryParm(r, test.queryKey)

		if !errors.Is(err, test.expected) {
			t.Errorf("Test %v: Expected error: %v, but got %v", test.nameTest, test.expected, err)
		}
	}
}

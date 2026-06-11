package core_domain

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestingStruct interface {
	Validate() error
}

type DataTest struct {
	name     string
	item     TestingStruct
	expected error
}

func TestUserValidate(t *testing.T) {
	testCase := []DataTest{
		{
			"UserCorrectValue",
			&User{
				ID:           1,
				Full_name:    "John Doe",
				Email:        "example@gmail.com",
				Phone_number: new("+7123456789"),

				Password: "password",
			}, nil},
		{
			"User error in full_name",
			&User{
				ID:           2,
				Full_name:    "",
				Email:        "example@gmail.com",
				Phone_number: new("+7123456789"),
				Password:     "password",
			}, fmt.Errorf("invalid user full_name")},
		{
			"User error in email",
			&User{
				ID:           3,
				Full_name:    "John Doe",
				Email:        "example",
				Phone_number: new("+7123456789"),
				Password:     "password",
			}, fmt.Errorf("invalid user email")},
		{
			"User error in password",
			&User{
				ID:           4,
				Full_name:    "John Doe",
				Email:        "example@gmail.com",
				Phone_number: new("+7123456789"),
				Password:     "",
			}, fmt.Errorf("invalid user password for user")},
	}

	for i, user := range testCase {
		ans := user.item.Validate()
		assert.Equal(t, user.expected, ans,
			fmt.Sprintf("Test validation user error number: %d",
				i,
			),
		)
	}
}

func TestTransactionValidate(t *testing.T) {
	testCase := []DataTest{
		{
			"Test sum transaction = zero",
			&Transaction{
				ID:          1,
				Sum:         0,
				Type:        "Income",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      1,
				Comments:    "",
				TimeCreated: time.Now(),
			},
			fmt.Errorf("invalid transaction sum: %d", 0),
		},
		{
			"Test categoryID transaction Unincelized",
			&Transaction{
				ID:          2,
				Sum:         123,
				Type:        "Income",
				Date:        time.Now(),
				CategoryID:  -1,
				UserID:      1,
				Comments:    "",
				TimeCreated: time.Now(),
			},
			fmt.Errorf("invalid transaction category id: %d", -1),
		},
		{
			"Test type transaction error",
			&Transaction{
				ID:          3,
				Sum:         100,
				Type:        "Somthing else",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      1,
				Comments:    "",
				TimeCreated: time.Now(),
			},
			fmt.Errorf("invalid transaction type: %v", "Somthing else"),
		},
		{
			"Test comments transaction error",
			&Transaction{
				ID:          1,
				Sum:         12,
				Type:        "Income",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      1,
				Comments:    strings.Repeat("a", 1001),
				TimeCreated: time.Now(),
			},
			fmt.Errorf("invalid transaction comments: %d", 1001),
		},
	}

	for i, transact := range testCase {
		ans := transact.item.Validate()
		assert.Equal(t, transact.expected, ans,
			fmt.Sprintf("Test validation transaction error number: %d, %s",
				i, transact.name,
			),
		)
	}
}

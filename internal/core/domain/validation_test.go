package core_domain

import (
	"errors"
	"testing"
	"time"

	core_errors "github.com/Hodorev-Evgeny/ExpensesTracker/internal/core/errors"
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
			"Test valid user",
			&User{
				ID:        1,
				Full_name: "Test Name",
				Email:     "test@gmail.com",
				Password:  "password123",
				Time_add:  time.Now(),
			},
			nil,
		},
		{
			"Test incorrect user id",
			&User{
				ID:        -1,
				Full_name: "Test Name",
				Email:     "test@gmail.com",
				Password:  "password123",
				Time_add:  time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test empty full name in user",
			&User{
				ID:        1,
				Full_name: "",
				Email:     "test@gmail.com",
				Password:  "password123",
				Time_add:  time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test empty email in user",
			&User{
				ID:        1,
				Full_name: "Test Name",
				Email:     "",
				Password:  "password123",
				Time_add:  time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test incorrect email in user",
			&User{
				ID:        1,
				Full_name: "Test Name",
				Email:     "invalid-email",
				Password:  "password123",
				Time_add:  time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test empty password in user",
			&User{
				ID:        1,
				Full_name: "Test Name",
				Email:     "test@gmail.com",
				Password:  "",
				Time_add:  time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test zero time in user",
			&User{
				ID:        1,
				Full_name: "Test Name",
				Email:     "test@gmail.com",
				Password:  "password123",
				Time_add:  time.Time{},
			},
			core_errors.ErrorBadRequest,
		},
	}

	for i, transact := range testCase {
		err := transact.item.Validate()
		if !errors.Is(err, transact.expected) {
			t.Log("Name: ", transact.name)
			t.Errorf("Test validation failed for user %d want: %d expected: %d", i, transact.expected, err)
		}
	}
}

func TestTransactionValidate(t *testing.T) {
	testCase := []DataTest{
		{
			"Test valid transaction",
			&Transaction{
				ID:          1,
				Sum:         5000,
				Type:        "Income",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      1,
				Comments:    "Salary",
				TimeCreated: time.Now(),
			},
			nil,
		},
		{
			"Test incorrect transaction id",
			&Transaction{
				ID:          -1,
				Sum:         5000,
				Type:        "income",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      1,
				Comments:    "Salary",
				TimeCreated: time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test negative sum in transaction",
			&Transaction{
				ID:          1,
				Sum:         -5000,
				Type:        "income",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      1,
				Comments:    "Salary",
				TimeCreated: time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test zero sum in transaction",
			&Transaction{
				ID:          1,
				Sum:         0,
				Type:        "income",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      1,
				Comments:    "Salary",
				TimeCreated: time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test empty type in transaction",
			&Transaction{
				ID:          1,
				Sum:         5000,
				Type:        "",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      1,
				Comments:    "Salary",
				TimeCreated: time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test incorrect type in transaction",
			&Transaction{
				ID:          1,
				Sum:         5000,
				Type:        "unknown",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      1,
				Comments:    "Salary",
				TimeCreated: time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test zero date in transaction",
			&Transaction{
				ID:          1,
				Sum:         5000,
				Type:        "income",
				Date:        time.Time{},
				CategoryID:  2,
				UserID:      1,
				Comments:    "Salary",
				TimeCreated: time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test incorrect category id in transaction",
			&Transaction{
				ID:          1,
				Sum:         5000,
				Type:        "income",
				Date:        time.Now(),
				CategoryID:  -2,
				UserID:      1,
				Comments:    "Salary",
				TimeCreated: time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test incorrect user id in transaction",
			&Transaction{
				ID:          1,
				Sum:         5000,
				Type:        "income",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      -1,
				Comments:    "Salary",
				TimeCreated: time.Now(),
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test empty comments in transaction",
			&Transaction{
				ID:          1,
				Sum:         5000,
				Type:        "Income",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      1,
				Comments:    "",
				TimeCreated: time.Now(),
			},
			nil,
		},
		{
			"Test zero creation time in transaction",
			&Transaction{
				ID:          1,
				Sum:         5000,
				Type:        "income",
				Date:        time.Now(),
				CategoryID:  2,
				UserID:      1,
				Comments:    "Salary",
				TimeCreated: time.Time{},
			},
			core_errors.ErrorBadRequest,
		},
	}

	for i, transact := range testCase {
		err := transact.item.Validate()
		if !errors.Is(err, transact.expected) {
			t.Log("Name: ", transact.name)
			t.Errorf("Test validation failed for transaction %d want: %d expected: %d", i, transact.expected, err)
		}
	}
}

func TestCategoryValidate(t *testing.T) {
	testCase := []DataTest{
		{
			"Test valid category",
			&Category{
				ID:      2,
				Name:    "Food",
				User_ID: 1,
			},
			nil,
		},
		{
			"Test incorrect category id",
			&Category{
				ID:      -2,
				Name:    "Food",
				User_ID: 1,
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test empty category name",
			&Category{
				ID:      2,
				Name:    "",
				User_ID: 1,
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test incorrect user id in category",
			&Category{
				ID:      2,
				Name:    "Food",
				User_ID: -12,
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test zero user id in category",
			&Category{
				ID:      2,
				Name:    "Food",
				User_ID: 0,
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test category without limit",
			&Category{
				ID:       2,
				Name:     "Food",
				User_ID:  1,
				Limit_id: nil,
			},
			nil,
		},
		{
			"Test incorrect limit id in category",
			&Category{
				ID:       2,
				Name:     "Food",
				User_ID:  1,
				Limit_id: new(int),
			},
			core_errors.ErrorBadRequest,
		},
	}

	for i, transact := range testCase {
		err := transact.item.Validate()
		if !errors.Is(err, transact.expected) {
			t.Log("Name: ", transact.name)
			t.Errorf("Test validation failed for category %d want: %d expected: %d", i, transact.expected, err)
		}
	}
}

func TestStaticValidation(t *testing.T) {
	testCase := []DataTest{
		{
			"Test negative average income in static",
			&Static{
				SumIncome:      10000,
				SumExpenditure: 4000,
				Difference:     6000,
				CountOperation: 4,
				AVGIncome:      -5000,
				AVGExpenditure: 2000,
				CostCategory:   "Food",
				MaxIncome:      7000,
				MaxExpenditure: 3000,
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test negative average expenditure in static",
			&Static{
				SumIncome:      10000,
				SumExpenditure: 4000,
				Difference:     6000,
				CountOperation: 4,
				AVGIncome:      5000,
				AVGExpenditure: -2000,
				CostCategory:   "Food",
				MaxIncome:      7000,
				MaxExpenditure: 3000,
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test negative max income in static",
			&Static{
				SumIncome:      10000,
				SumExpenditure: 4000,
				Difference:     6000,
				CountOperation: 4,
				AVGIncome:      5000,
				AVGExpenditure: 2000,
				CostCategory:   "Food",
				MaxIncome:      -7000,
				MaxExpenditure: 3000,
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test negative max expenditure in static",
			&Static{
				SumIncome:      10000,
				SumExpenditure: 4000,
				Difference:     6000,
				CountOperation: 4,
				AVGIncome:      5000,
				AVGExpenditure: 2000,
				CostCategory:   "Food",
				MaxIncome:      7000,
				MaxExpenditure: -3000,
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test negative share category in static",
			&Static{
				SumIncome:      10000,
				SumExpenditure: 4000,
				Difference:     6000,
				CountOperation: 4,
				AVGIncome:      5000,
				AVGExpenditure: 2000,
				CostCategory:   "Food",
				ShareCategory: map[string]float64{
					"Food": -20,
				},
				MaxIncome:      7000,
				MaxExpenditure: 3000,
			},
			core_errors.ErrorBadRequest,
		},
	}

	for i, transact := range testCase {
		err := transact.item.Validate()
		if !errors.Is(err, transact.expected) {
			t.Log("Name: ", transact.name)
			t.Errorf("Test validation failed for static %d want: %d expected: %d", i, transact.expected, err)
		}
	}
}

func TestLimitValidation(t *testing.T) {
	testCase := []DataTest{
		{
			"Test valid limit",
			&Limit{
				ID:          1,
				Duration:    time.Now().AddDate(0, 1, 0),
				AmountLimit: 10000,
			},
			nil,
		},
		{
			"Test incorrect limit id",
			&Limit{
				ID:          -1,
				Duration:    time.Now().AddDate(0, 1, 0),
				AmountLimit: 10000,
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test zero duration in limit",
			&Limit{
				ID:          1,
				Duration:    time.Time{},
				AmountLimit: 10000,
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test negative amount in limit",
			&Limit{
				ID:          1,
				Duration:    time.Now().AddDate(0, 1, 0),
				AmountLimit: -10000,
			},
			core_errors.ErrorBadRequest,
		},
		{
			"Test zero amount in limit",
			&Limit{
				ID:          1,
				Duration:    time.Now().AddDate(0, 1, 0),
				AmountLimit: 0,
			},
			core_errors.ErrorBadRequest,
		},
	}

	for i, transact := range testCase {
		err := transact.item.Validate()
		if !errors.Is(err, transact.expected) {
			t.Log("Name: ", transact.name)
			t.Errorf("Test validation failed for limit %d want: %d expected: %d", i, transact.expected, err)
		}
	}
}

func TestPatchUser(t *testing.T) {

	testCase := []DataTest{
		{
			name: "Update user full name",
			item: &UserPatch{
				Full_name: Nullable[string]{Value: new("Test Name"), Set: true},
			},
			expected: nil,
		},
		{
			name: "Update user email",
			item: &UserPatch{
				Email: Nullable[string]{Value: new("test@gmail.com"), Set: true},
			},
			expected: nil,
		},
		{
			name: "Update user phone",
			item: &UserPatch{
				Phone: Nullable[string]{Value: new("+79991234567"), Set: true},
			},
			expected: nil,
		},
		{
			name: "Update all user fields",
			item: &UserPatch{
				Full_name: Nullable[string]{Value: new("Test Name"), Set: true},
				Email:     Nullable[string]{Value: new("test@gmail.com"), Set: true},
				Phone:     Nullable[string]{Value: new("+79991234567"), Set: true},
			},
			expected: nil,
		},
		{
			name:     "Empty user patch",
			item:     &UserPatch{},
			expected: nil,
		},
	}

	for i, transact := range testCase {
		val := User{
			ID:        1,
			Full_name: "Test Name",
			Email:     "test@gmail.com",
			Password:  "password123",
			Time_add:  time.Now(),
		}
		err := val.ApplyPatch(*transact.item.(*UserPatch))
		if err != nil && !errors.Is(err, transact.expected) {
			t.Log("Name: ", transact.name)
			t.Errorf("Test apply patch failed for user %d want: %d expected: %d", i, transact.expected, err)
		}
	}
}

func TestPatchTransactions(t *testing.T) {
	testCase := []DataTest{
		{
			name: "Update transaction sum",
			item: &PatchTransaction{
				Sum: Nullable[int]{Value: new(1500), Set: true},
			},
			expected: nil,
		},
		{
			name: "Update transaction category",
			item: &PatchTransaction{
				CategoryId: Nullable[int]{Value: new(2), Set: true},
			},
			expected: nil,
		},
		{
			name: "Update transaction type",
			item: &PatchTransaction{
				TypeTransaction: Nullable[string]{Value: new("Expenditure"), Set: true},
			},
			expected: nil,
		},
		{
			name: "Update transaction comments",
			item: &PatchTransaction{
				Comments: Nullable[string]{Value: new("Покупка продуктов"), Set: true},
			},
			expected: nil,
		},
		{
			name: "Update all transaction fields",
			item: &PatchTransaction{
				Sum:             Nullable[int]{Value: new(5000), Set: true},
				CategoryId:      Nullable[int]{Value: new(3), Set: true},
				TypeTransaction: Nullable[string]{Value: new("Income"), Set: true},
				Comments:        Nullable[string]{Value: new("Income"), Set: true},
			},
			expected: nil,
		},
		{
			name:     "Empty transaction patch",
			item:     &PatchTransaction{},
			expected: nil,
		},
	}

	for i, transact := range testCase {
		val := Transaction{
			ID:          1,
			Sum:         5000,
			Type:        "Income",
			Date:        time.Now(),
			CategoryID:  2,
			UserID:      1,
			Comments:    "Salary",
			TimeCreated: time.Now(),
		}
		err := val.ApplyPatch(*transact.item.(*PatchTransaction))
		if !errors.Is(err, transact.expected) {
			t.Log("Name: ", transact.name)
			t.Errorf("Test apply patch failed for transaction %d want: %d expected: %d", i, transact.expected, err)
		}
	}
}

func TestPatchLimits(t *testing.T) {
	testCase := []DataTest{
		{name: "Update limit duration",
			item: &PatchLimit{
				Duraction: Nullable[time.Time]{Value: new(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)), Set: true},
			},
			expected: nil,
		},
		{
			name: "Update limit amount",
			item: &PatchLimit{
				Amount_limit: Nullable[int]{Value: new(10000), Set: true},
			},
			expected: nil,
		},
		{
			name: "Update all limit fields",
			item: &PatchLimit{
				Duraction:    Nullable[time.Time]{Value: new(time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)), Set: true},
				Amount_limit: Nullable[int]{Value: new(20000), Set: true},
			},
			expected: nil,
		},
		{
			name:     "Empty limit patch",
			item:     &PatchLimit{},
			expected: nil,
		},
	}

	for i, transact := range testCase {
		val := Limit{
			ID:          1,
			Duration:    time.Now().AddDate(0, 1, 0),
			AmountLimit: 10000,
		}
		err := val.ApplyPatch(*transact.item.(*PatchLimit))
		if !errors.Is(err, transact.expected) {
			t.Log("Name: ", transact.name)
			t.Errorf("Test apply patch failed for limit %d want: %d expected: %d", i, transact.expected, err)
		}
	}
}

func TestPatchCategory(t *testing.T) {
	testCase := []DataTest{
		{
			name: "Update category name",
			item: &CategoryUpdate{
				Name: Nullable[string]{Value: new("Product"), Set: true},
			},
			expected: nil,
		},
		{
			name: "Update category limit",
			item: &CategoryUpdate{
				Limit_id: Nullable[int]{Value: new(2), Set: true},
			},
			expected: nil,
		},
		{
			name: "Update all category fields",
			item: &CategoryUpdate{
				Name:     Nullable[string]{Value: new("product"), Set: true},
				Limit_id: Nullable[int]{Value: new(3), Set: true},
			},
			expected: nil,
		},
		{
			name:     "Empty category patch",
			item:     &CategoryUpdate{},
			expected: nil,
		}}

	for i, transact := range testCase {
		val := Category{
			ID:      2,
			Name:    "Food",
			User_ID: 1,
		}
		err := val.Update(*transact.item.(*CategoryUpdate))
		if !errors.Is(err, transact.expected) {
			t.Log("Name: ", transact.name)
			t.Errorf("Test apply patch failed for category %d want: %d expected: %d", i, transact.expected, err)
		}
	}
}

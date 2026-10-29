package core_domain

import (
	"time"
)

type Static struct {
	SumIncome      int
	SumExpenditure int
	Difference     float64
	CountOperation int
	AVGIncome      float64
	AVGExpenditure float64
	CostCategory   string
	ShareCategory  map[string]float64
	MaxIncome      int
	MaxExpenditure int
}

type FiltersStatic struct {
	To         *time.Time
	From       *time.Time
	UserID     int
	CategoryID *int
}

func NewFiltersStatic(
	to, from *time.Time,
	categoryID *int,
	userID int,
) FiltersStatic {
	return FiltersStatic{
		To:         to,
		From:       from,
		CategoryID: categoryID,
		UserID:     userID,
	}
}

func NewStatic(
	category []Category,
	transaction []Transaction,
) *Static {
	var sumIncome, sumExpenditure, countOperation, maxIcome, maxExpenditure, countIcome, countExpenditure int
	var avgIncome, avgExpenditure float64
	var costCategory string

	countOperation = len(transaction)
	tmpShare := make(map[int][]int)
	for _, t := range transaction {
		if t.Type == "Income" {
			if t.Sum > maxIcome {
				maxIcome = t.Sum
			}

			sumIncome += t.Sum
			countIcome++
		} else if t.Type == "Expenditure" {
			if t.Sum > maxExpenditure {
				maxExpenditure = t.Sum
			}

			sumExpenditure += t.Sum
			countExpenditure++

			if val, ok := tmpShare[t.CategoryID]; ok {
				tmpShare[t.CategoryID] = append(val, t.Sum)
			} else {
				tmpShare[t.CategoryID] = []int{t.Sum}
			}
		}
	}

	shareCategory := make(map[string]float64)
	tmpCostCategory := 0
	for _, c := range category {
		if val, ok := tmpShare[c.ID]; ok {
			s := sum(val)
			if s > tmpCostCategory {
				tmpCostCategory = s
				costCategory = c.Name
			}
			shareCategory[c.Name] = float64(s) / float64(sumExpenditure) * 100
		}
	}

	difference := float64(sumIncome - sumExpenditure)
	if countIcome == 0 {
		avgIncome = 0
	} else {
		avgIncome = float64(sumIncome) / float64(countIcome)
	}
	if countExpenditure == 0 {
		avgExpenditure = 0
	} else {
		avgExpenditure = float64(sumExpenditure) / float64(countExpenditure)
	}

	return &Static{
		SumIncome:      sumIncome,
		SumExpenditure: sumExpenditure,
		Difference:     difference,
		CountOperation: countOperation,
		AVGIncome:      avgIncome,
		AVGExpenditure: avgExpenditure,
		CostCategory:   costCategory,
		ShareCategory:  shareCategory,
		MaxIncome:      maxIcome,
		MaxExpenditure: maxExpenditure,
	}
}

func sum(list []int) int {
	s := 0
	for _, i := range list {
		s += i
	}
	return s
}

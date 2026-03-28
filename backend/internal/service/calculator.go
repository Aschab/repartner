package service

import (
	"errors"
	"sort"

	"pack-calculator/internal/domain"
)

var (
	ErrInvalidOrderQty    = errors.New("order_quantity must be greater than 0")
	ErrEmptyPackSizes     = errors.New("pack sizes cannot be empty")
	ErrInvalidPackSize    = errors.New("pack sizes must be positive integers")
	ErrDuplicatePackSizes = errors.New("duplicate pack sizes are not allowed")
)

// Calculator defines the interface for pack calculation.
type Calculator interface {
	Calculate(orderQty int, packSizes []int) (domain.CalculationResult, error)
}

type calculator struct{}

// NewCalculator creates a new Calculator instance.
func NewCalculator() Calculator {
	return &calculator{}
}

// Calculate determines the optimal pack combination for the given order quantity.
func (c *calculator) Calculate(orderQty int, packSizes []int) (domain.CalculationResult, error) {
	// Validate inputs
	if orderQty <= 0 {
		return domain.CalculationResult{}, ErrInvalidOrderQty
	}

	if len(packSizes) == 0 {
		return domain.CalculationResult{}, ErrEmptyPackSizes
	}

	// Validate and deduplicate pack sizes
	seen := make(map[int]bool)
	for _, size := range packSizes {
		if size <= 0 {
			return domain.CalculationResult{}, ErrInvalidPackSize
		}
		if seen[size] {
			return domain.CalculationResult{}, ErrDuplicatePackSizes
		}
		seen[size] = true
	}

	// Sort pack sizes ascending for calculation
	sortedPacks := make([]int, len(packSizes))
	copy(sortedPacks, packSizes)
	sort.Ints(sortedPacks)

	maxPack := sortedPacks[len(sortedPacks)-1]

	// Search for optimal solution
	// Search candidate shipped totals from orderQty through orderQty + maxPack - 1
	// Find the smallest shippable total >= orderQty, and for that total the fewest packs.
	limit := orderQty + maxPack
	bestPackCounts, bestTotal, bestPacks, ok := c.solveUpTo(limit, orderQty, sortedPacks)
	if !ok {
		return domain.CalculationResult{}, errors.New("no valid pack combination found")
	}

	// Build result
	result := domain.CalculationResult{
		RequestedQty: orderQty,
		TotalShipped: bestTotal,
		TotalPacks:   bestPacks,
		Packs:        make([]domain.PackSelection, 0),
	}

	// Sort pack sizes descending for output
	packSizesSorted := make([]int, 0, len(bestPackCounts))
	for size := range bestPackCounts {
		packSizesSorted = append(packSizesSorted, size)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(packSizesSorted)))

	for _, size := range packSizesSorted {
		if count := bestPackCounts[size]; count > 0 {
			result.Packs = append(result.Packs, domain.PackSelection{
				PackSize: size,
				Count:    count,
			})
		}
	}

	return result, nil
}

// solveUpTo uses dynamic programming once up to limit, then selects:
// 1) the smallest reachable total >= minTarget
// 2) for that total, the fewest packs
func (c *calculator) solveUpTo(limit, minTarget int, packSizes []int) (map[int]int, int, int, bool) {
	const inf = (1 << 31) - 1

	// dp[i] = minimum packs needed to make total i
	dp := make([]int, limit+1)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0

	// usedPack[i] = pack size chosen last to reach total i
	usedPack := make([]int, limit+1)

	for i := 1; i <= limit; i++ {
		for _, p := range packSizes {
			if p <= i && dp[i-p] != inf {
				candidate := dp[i-p] + 1
				if candidate < dp[i] {
					dp[i] = candidate
					usedPack[i] = p
				}
			}
		}
	}

	// Rule 2 takes precedence: choose the smallest reachable total >= minTarget.
	bestTotal := -1
	for t := minTarget; t <= limit; t++ {
		if dp[t] != inf {
			bestTotal = t
			break
		}
	}

	if bestTotal == -1 {
		return nil, 0, 0, false
	}

	// Reconstruct pack counts for bestTotal.
	packCounts := make(map[int]int)
	remaining := bestTotal
	for remaining > 0 {
		p := usedPack[remaining]
		packCounts[p]++
		remaining -= p
	}

	return packCounts, bestTotal, dp[bestTotal], true
}

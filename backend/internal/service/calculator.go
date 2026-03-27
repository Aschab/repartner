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
	var bestTotal int
	var bestPacks int
	var bestPackCounts map[int]int

	for t := orderQty; t < orderQty+maxPack; t++ {
		packCounts, totalPacks, ok := c.solveExact(t, sortedPacks)
		if !ok {
			continue
		}

		if bestPackCounts == nil || t < bestTotal || (t == bestTotal && totalPacks < bestPacks) {
			bestTotal = t
			bestPacks = totalPacks
			bestPackCounts = packCounts
		}

		// If we found an exact match, we can stop early for this total
		// but we continue to find minimum total first
		if bestPackCounts != nil && t > bestTotal {
			break
		}
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

// solveExact uses dynamic programming to find the minimum number of packs
// to reach exactly the target total.
func (c *calculator) solveExact(target int, packSizes []int) (map[int]int, int, bool) {
	const inf = 1<<31 - 1

	// dp[i] = minimum packs needed to make total i
	dp := make([]int, target+1)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0

	// usedPack[i] = pack size chosen to reach total i
	usedPack := make([]int, target+1)

	// Fill DP table
	for i := 1; i <= target; i++ {
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

	// Check if target is reachable
	if dp[target] == inf {
		return nil, 0, false
	}

	// Reconstruct pack counts
	packCounts := make(map[int]int)
	remaining := target
	for remaining > 0 {
		p := usedPack[remaining]
		packCounts[p]++
		remaining -= p
	}

	return packCounts, dp[target], true
}

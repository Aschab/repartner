package service

import (
	"testing"
)

func TestCalculator_Calculate(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name          string
		orderQty      int
		packSizes     []int
		wantTotal     int
		wantPacks     int
		wantPacksMap  map[int]int
		wantErr       error
	}{
		// Exact business examples from spec
		{
			name:         "order 1 returns 250",
			orderQty:     1,
			packSizes:    []int{250, 500, 1000, 2000, 5000},
			wantTotal:    250,
			wantPacks:    1,
			wantPacksMap: map[int]int{250: 1},
		},
		{
			name:         "order 250 returns 250",
			orderQty:     250,
			packSizes:    []int{250, 500, 1000, 2000, 5000},
			wantTotal:    250,
			wantPacks:    1,
			wantPacksMap: map[int]int{250: 1},
		},
		{
			name:         "order 251 returns 500",
			orderQty:     251,
			packSizes:    []int{250, 500, 1000, 2000, 5000},
			wantTotal:    500,
			wantPacks:    1,
			wantPacksMap: map[int]int{500: 1},
		},
		{
			name:         "order 501 returns 750",
			orderQty:     501,
			packSizes:    []int{250, 500, 1000, 2000, 5000},
			wantTotal:    750,
			wantPacks:    2,
			wantPacksMap: map[int]int{500: 1, 250: 1},
		},
		{
			name:         "order 12001 returns 12250",
			orderQty:     12001,
			packSizes:    []int{250, 500, 1000, 2000, 5000},
			wantTotal:    12250,
			wantPacks:    4,
			wantPacksMap: map[int]int{5000: 2, 2000: 1, 250: 1},
		},

		// Tie-breaker cases - same shipped quantity, fewer packs wins
		{
			name:         "prefer fewer packs: 500 as one pack not two 250s",
			orderQty:     500,
			packSizes:    []int{250, 500},
			wantTotal:    500,
			wantPacks:    1,
			wantPacksMap: map[int]int{500: 1},
		},
		{
			name:         "prefer fewer packs: 1000 as one pack",
			orderQty:     1000,
			packSizes:    []int{250, 500, 1000},
			wantTotal:    1000,
			wantPacks:    1,
			wantPacksMap: map[int]int{1000: 1},
		},

		// Flexibility cases with non-standard pack sizes
		{
			name:         "flexible packs: order 10 with [3,5,9]",
			orderQty:     10,
			packSizes:    []int{3, 5, 9},
			wantTotal:    10,
			wantPacks:    2,
			wantPacksMap: map[int]int{5: 2},
		},
		{
			name:         "flexible packs: order 11 with [3,5,9]",
			orderQty:     11,
			packSizes:    []int{3, 5, 9},
			wantTotal:    11,
			wantPacks:    3,
			wantPacksMap: map[int]int{5: 1, 3: 2},
		},

		// Order independence - different input order should give same result
		{
			name:         "order independence: shuffled pack sizes",
			orderQty:     501,
			packSizes:    []int{5000, 250, 2000, 1000, 500},
			wantTotal:    750,
			wantPacks:    2,
			wantPacksMap: map[int]int{500: 1, 250: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Calculate(tt.orderQty, tt.packSizes)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.TotalShipped != tt.wantTotal {
				t.Errorf("TotalShipped = %d, want %d", result.TotalShipped, tt.wantTotal)
			}

			if result.TotalPacks != tt.wantPacks {
				t.Errorf("TotalPacks = %d, want %d", result.TotalPacks, tt.wantPacks)
			}

			if result.RequestedQty != tt.orderQty {
				t.Errorf("RequestedQty = %d, want %d", result.RequestedQty, tt.orderQty)
			}

			// Verify pack counts
			gotPacksMap := make(map[int]int)
			for _, p := range result.Packs {
				gotPacksMap[p.PackSize] = p.Count
			}

			if len(gotPacksMap) != len(tt.wantPacksMap) {
				t.Errorf("got %d different pack sizes, want %d", len(gotPacksMap), len(tt.wantPacksMap))
			}

			for size, count := range tt.wantPacksMap {
				if gotPacksMap[size] != count {
					t.Errorf("pack size %d: got count %d, want %d", size, gotPacksMap[size], count)
				}
			}

			// Verify packs are sorted descending
			for i := 1; i < len(result.Packs); i++ {
				if result.Packs[i].PackSize > result.Packs[i-1].PackSize {
					t.Errorf("packs not sorted descending: %v", result.Packs)
				}
			}
		})
	}
}

func TestCalculator_ValidationErrors(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name      string
		orderQty  int
		packSizes []int
		wantErr   error
	}{
		{
			name:      "order quantity zero",
			orderQty:  0,
			packSizes: []int{250, 500},
			wantErr:   ErrInvalidOrderQty,
		},
		{
			name:      "order quantity negative",
			orderQty:  -1,
			packSizes: []int{250, 500},
			wantErr:   ErrInvalidOrderQty,
		},
		{
			name:      "empty pack sizes",
			orderQty:  100,
			packSizes: []int{},
			wantErr:   ErrEmptyPackSizes,
		},
		{
			name:      "pack size zero",
			orderQty:  100,
			packSizes: []int{250, 0, 500},
			wantErr:   ErrInvalidPackSize,
		},
		{
			name:      "pack size negative",
			orderQty:  100,
			packSizes: []int{250, -100, 500},
			wantErr:   ErrInvalidPackSize,
		},
		{
			name:      "duplicate pack sizes",
			orderQty:  100,
			packSizes: []int{250, 500, 250},
			wantErr:   ErrDuplicatePackSizes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := calc.Calculate(tt.orderQty, tt.packSizes)
			if err != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestCalculator_OutputConsistency(t *testing.T) {
	calc := NewCalculator()

	// Test that TotalShipped and TotalPacks are consistent with Packs
	result, err := calc.Calculate(12001, []int{250, 500, 1000, 2000, 5000})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify TotalShipped matches sum of pack sizes * counts
	var calculatedTotal int
	var calculatedPacks int
	for _, p := range result.Packs {
		calculatedTotal += p.PackSize * p.Count
		calculatedPacks += p.Count
	}

	if calculatedTotal != result.TotalShipped {
		t.Errorf("TotalShipped mismatch: reported %d, calculated %d", result.TotalShipped, calculatedTotal)
	}

	if calculatedPacks != result.TotalPacks {
		t.Errorf("TotalPacks mismatch: reported %d, calculated %d", result.TotalPacks, calculatedPacks)
	}
}

func TestCalculator_PacksSortedDescending(t *testing.T) {
	calc := NewCalculator()

	result, err := calc.Calculate(12001, []int{250, 500, 1000, 2000, 5000})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify packs are sorted descending by pack size
	for i := 1; i < len(result.Packs); i++ {
		if result.Packs[i].PackSize > result.Packs[i-1].PackSize {
			t.Errorf("packs not sorted descending: %v", result.Packs)
		}
	}
}

// Additional edge case tests
func TestCalculator_EdgeCases(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name         string
		orderQty     int
		packSizes    []int
		wantTotal    int
		wantPacks    int
		wantPacksMap map[int]int
	}{
		{
			name:         "exact match with largest pack",
			orderQty:     5000,
			packSizes:    []int{250, 500, 1000, 2000, 5000},
			wantTotal:    5000,
			wantPacks:    1,
			wantPacksMap: map[int]int{5000: 1},
		},
		{
			name:         "just over largest pack",
			orderQty:     5001,
			packSizes:    []int{250, 500, 1000, 2000, 5000},
			wantTotal:    5250,
			wantPacks:    2,
			wantPacksMap: map[int]int{5000: 1, 250: 1},
		},
		{
			name:         "single pack size available",
			orderQty:     100,
			packSizes:    []int{50},
			wantTotal:    100,
			wantPacks:    2,
			wantPacksMap: map[int]int{50: 2},
		},
		{
			name:         "single pack size not exact fit",
			orderQty:     101,
			packSizes:    []int{50},
			wantTotal:    150,
			wantPacks:    3,
			wantPacksMap: map[int]int{50: 3},
		},
		// Large order stress test
		{
			name:         "large order 500000 with prime-like pack sizes",
			orderQty:     500000,
			packSizes:    []int{23, 31, 53},
			wantTotal:    500000,
			wantPacks:    9438,
			wantPacksMap: map[int]int{53: 9429, 31: 7, 23: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Calculate(tt.orderQty, tt.packSizes)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.TotalShipped != tt.wantTotal {
				t.Errorf("TotalShipped = %d, want %d", result.TotalShipped, tt.wantTotal)
			}

			if result.TotalPacks != tt.wantPacks {
				t.Errorf("TotalPacks = %d, want %d", result.TotalPacks, tt.wantPacks)
			}

			gotPacksMap := make(map[int]int)
			for _, p := range result.Packs {
				gotPacksMap[p.PackSize] = p.Count
			}

			for size, count := range tt.wantPacksMap {
				if gotPacksMap[size] != count {
					t.Errorf("pack size %d: got count %d, want %d", size, gotPacksMap[size], count)
				}
			}
		})
	}
}

package utils

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		// 正常情况：正数相加
		{"正数相加", 1, 2, 3},
		{"正数相加大数", 100, 200, 300},
		// 负数情况
		{"负数相加", -1, -2, -3},
		{"一负一正", -5, 3, -2},
		// 零的情况
		{"零加正数", 0, 5, 5},
		{"零加负数", 0, -3, -3},
		{"零加零", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

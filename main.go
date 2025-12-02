package main

import (
	"exp/test"
	"fmt"
	"testing"
)

func main() {
	// 创建一个测试上下文
	tests := []struct {
		name string
		fn   func(*testing.T)
	}{
		{"TestParse", test.TestParse},
		{"TestNegative", test.TestNegative},
		{"NullTest", test.NullTest},
		{"NotTest", test.NotTest},
		{"EmptyStringTest", test.EmptyStringTest},
		{"TestLogic", test.TestLogic},
		{"TestMutipleFunction", test.TestMutipleFunction},
		{"TestMutipleExpFunction", test.TestMutipleExpFunction},
		{"TestFunctionParams", test.TestFunctionParams},
		{"TestUnEquals", test.TestUnEquals},
		{"TestArithmetic", test.TestArithmetic},
		{"TestString", test.TestString},
		{"TestParams", test.TestParams},
		{"TestDateCompare", test.TestDateCompare},
		{"TestDateMoreThan", test.TestDateMoreThan},
		{"TestDateLessThan", test.TestDateLessThan},
		{"TestEDATE", test.TestEDATE},
		{"TestEODateStart", test.TestEODateStart},
		{"TestEODateEnd", test.TestEODateEnd},
		{"TestNowTime", test.TestNowTime},
		{"TestRound1", test.TestRound1},
		{"TestRound2", test.TestRound2},
		{"TestRound3", test.TestRound3},
		{"TestTimeSpanDays", test.TestTimeSpanDays},
		{"TestTimeSpanHours", test.TestTimeSpanHours},
		{"TestTimeSpanMinutes", test.TestTimeSpanMinutes},
		{"TestTimeSpanSeconds", test.TestTimeSpanSeconds},
		{"TestTimeSpanMillSeconds", test.TestTimeSpanMillSeconds},
		{"TestRoundAndTimeSpan", test.TestRoundAndTimeSpan},
	}

	// 运行所有测试
	for _, tt := range tests {
		t := &testing.T{}
		fmt.Printf("Running %s...\n", tt.name)
		tt.fn(t)

		// 检查测试是否失败
		if t.Failed() {
			fmt.Printf("%s FAILED\n", tt.name)
		} else {
			fmt.Printf("%s PASSED\n", tt.name)
		}
	}
}

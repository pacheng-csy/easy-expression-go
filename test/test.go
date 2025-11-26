package test

import (
	expression "exp/src"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	expStr := " 2 + 3* -3 > -9 || [SUM] (1,2,3) < 4"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(1), expression.InterfaceToFloat64(value))
}

func NullTest(t *testing.T) {
	expStr := "[ISNULL](obj)"
	exp, _ := expression.CreateExpression(expStr)
	dic := map[string]interface{}{
		"obj": nil,
	}
	exp.LoadArgumentWithDictionary(dic)
	value := exp.Execute()
	assert.Equal(t, float64(1), expression.InterfaceToFloat64(value))
}

func EmptyStringTest(t *testing.T) {
	expStr := "a==''"
	exp, _ := expression.CreateExpression(expStr)
	dic := map[string]interface{}{
		"a": "",
	}
	exp.LoadArgumentWithDictionary(dic)
	value := exp.Execute()
	assert.Equal(t, float64(1), expression.InterfaceToFloat64(value))
}

func TestNegative(t *testing.T) {
	expStr := "3 * -2"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(-6.0), expression.InterfaceToFloat64(value))
}

func TestLogic(t *testing.T) {
	expStr := "3 * (1 + 2) <= 5 || !(8 / (4 - 2) > [SUM](1,2,3))"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(1.0), expression.InterfaceToFloat64(value))
}

func TestMutipleFunction(t *testing.T) {
	expStr := "[SUM]([SUM](1,2),[SUM](3,4),[AVG](5,6,7))"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(16.0), expression.InterfaceToFloat64(value))
}

func TestMutipleExpFunction(t *testing.T) {
	expStr := "3 * (1 + 2) + [SUM]([SUM](1,2),6 / 2,[AVG](5,6,7))"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(21.0), expression.InterfaceToFloat64(value))
}

func TestFunctionParams(t *testing.T) {
	expStr := "[EQUALS](12+3,15)"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(1.0), expression.InterfaceToFloat64(value))
}

func TestUnEquals(t *testing.T) {
	expStr := "4 != 4"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(0.0), expression.InterfaceToFloat64(value))
}

func TestArithmetic(t *testing.T) {
	expStr := "3 * (1 + 2) + 5 - (30 / (4 - 2) % [SUM](1,2,3))"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(11.0), expression.InterfaceToFloat64(value))
}

func TestString(t *testing.T) {
	expStr := "a * (b + c) > d & [Contains](srcText,text)"
	dic := map[string]interface{}{
		"a":       "3",
		"b":       "1",
		"c":       "2",
		"d":       "4",
		"srcText": "abc",
		"text":    "bc",
	}
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgumentWithDictionary(dic)
	value := exp.Execute()
	assert.Equal(t, float64(1.0), expression.InterfaceToFloat64(value))
}

func TestParams(t *testing.T) {
	expStr := "a * (b + c) + 5 - (30 / (d - 2) % [SUM](1,2,3))"
	dic := map[string]interface{}{
		"a": "3",
		"b": "1",
		"c": "2",
		"d": "4",
	}
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgumentWithDictionary(dic)
	value := exp.Execute()
	assert.Equal(t, float64(11.0), expression.InterfaceToFloat64(value))
}

func TestDateCompare(t *testing.T) {
	expStr := "'2024-05-27' == a"
	dic := map[string]interface{}{
		"a": "2024-05-27",
	}
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgumentWithDictionary(dic)
	value := exp.Execute()
	assert.Equal(t, float64(1.0), expression.InterfaceToFloat64(value))
}

func TestDateMoreThan(t *testing.T) {
	expStr := "'2024-05-27' > a"
	dic := map[string]interface{}{
		"a": "2024-05-26",
	}
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgumentWithDictionary(dic)
	value := exp.Execute()
	assert.Equal(t, float64(1.0), expression.InterfaceToFloat64(value))
}

func TestDateLessThan(t *testing.T) {
	expStr := "'2024-05-27' < a"
	dic := map[string]interface{}{
		"a": "2024-05-26",
	}
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgumentWithDictionary(dic)
	value := exp.Execute()
	assert.Equal(t, float64(0.0), expression.InterfaceToFloat64(value))
}

func TestEDATE(t *testing.T) {
	expStr := "[TIMETOSTRING]([EDATE]('2024-05-27',2,D),yyyyMMdd)"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, "20240529", value)
}

func TestEODateStart(t *testing.T) {
	expStr := "[TIMETOSTRING]([EODATE]('2024-05-27',2,S),yyyyMMdd)"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, "20240701", value)
}

func TestEODateEnd(t *testing.T) {
	expStr := "[TIMETOSTRING]([EODATE]('2024-05-27',2,E),yyyyMMdd)"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, "20240731", value)
}

func TestNowTime(t *testing.T) {
	expStr := "[TIMETOSTRING]([NOWTIME](),yyyyMMdd)"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, time.Now().Format("20060102"), value)
}

func TestRound1(t *testing.T) {
	expStr := "[ROUND](11.34,1,-1)"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(11.3), expression.InterfaceToFloat64(value))
}

func TestRound2(t *testing.T) {
	expStr := "[ROUND](11.34,1,0)"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(11.3), expression.InterfaceToFloat64(value))
}

func TestRound3(t *testing.T) {
	expStr := "[ROUND](11.34,1,1)"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(11.4), expression.InterfaceToFloat64(value))
}

func TestTimeSpanDays(t *testing.T) {
	expStr := "[DAYS]('2024-10-15'-'2024-10-10')"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(5.0), expression.InterfaceToFloat64(value))
}

func TestTimeSpanHours(t *testing.T) {
	expStr := "[HOURS]('2024-10-15'-'2024-10-10')"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(120.0), expression.InterfaceToFloat64(value))
}

func TestTimeSpanMinutes(t *testing.T) {
	expStr := "[MINUTES]('2024-10-15'-'2024-10-10')"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(7200.0), expression.InterfaceToFloat64(value))
}

func TestTimeSpanSeconds(t *testing.T) {
	expStr := "[SECONDS]('2024-10-15'-'2024-10-10')"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(432000.0), expression.InterfaceToFloat64(value))
}

func TestTimeSpanMillSeconds(t *testing.T) {
	expStr := "[MILLSECONDS]('2024-10-15'-'2024-10-10')"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(432000000.0), expression.InterfaceToFloat64(value))
}

func TestRoundAndTimeSpan(t *testing.T) {
	expStr := "[ROUND]([DAYS]('2024-10-15'-'2024-10-10') / 30,1,0)"
	exp, _ := expression.CreateExpression(expStr)
	exp.LoadArgument()
	value := exp.Execute()
	assert.Equal(t, float64(0.2), expression.InterfaceToFloat64(value))
}

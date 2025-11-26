package easyExpression

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

type match_Scope struct {
	ChildrenExpressionString string
	EndIndex                 int
	Status                   bool
}

type Expression struct {
	//错误消息
	ErrorMessgage string
	//状态(标识表达式是否可以解析)
	Status bool
	//元素类型
	ElementType ElementType
	//包含关键字的完整表达式
	SourceExpressionString string
	//数据值
	DataString interface{}
	//当前层级实际值
	RealityString interface{}
	//运算符
	Operators []*Operator
	//函数类型
	FunctionType FunctionType
	//若表达式为函数,则可以调用此委托来计算函数输出值,计算时根据函数枚举值来确定要转换的函数类型
	Function     func(FormulaAction, ...any) interface{}
	FunctionName string
	//子表达式
	ExpressionChildren []*Expression
}

func CreateExpression(expressionStr string) (*Expression, error) {
	if len(expressionStr) == 0 {
		return nil, fmt.Errorf("表达式不能为空")
	}
	expStr := strings.Replace(expressionStr, "||", "|", -1)
	expStr = strings.Replace(expStr, "\\\\", "\\", -1)
	expStr = strings.Replace(expStr, "&&", "&", -1)
	expStr = strings.Replace(expStr, "==", "=", -1)
	expStr = strings.Trim(expStr, " ")
	exp := Expression{
		SourceExpressionString: expStr,
		DataString:             "",
	}
	defer func() {
		if err := recover(); err != nil {
			exp.ErrorMessgage = err.(error).Error()
		}
	}()
	if tryParse(&exp) {
		return &exp, nil
	}
	return nil, fmt.Errorf(exp.ErrorMessgage)
}

/****************************************load****************************************************/

func (e *Expression) LoadArgumentWithDictionary(keyValues map[string]interface{}) []KeyValuePairElementString {
	result := make([]KeyValuePairElementString, 0)
	e.loadArgumentWithDictionary(keyValues, &result, e.ElementType == ElementFunction)
	return result
}

func (e *Expression) LoadArgument() {
	// 如果是数据节点或函数节点且有数据内容，则设置RealityString
	if e.ElementType == ElementData ||
		(e.ElementType == ElementFunction && e.DataString != "") {
		e.RealityString = e.DataString
	}

	// 递归处理所有子表达式
	for _, childExp := range e.ExpressionChildren {
		childExp.LoadArgument()
	}
}

func (e *Expression) loadArgumentWithDictionary(keyValues map[string]interface{}, result *[]KeyValuePairElementString, zeroInit bool) {
	if e.DataString != "" {
		if e.ElementType == ElementFunction {
			// 处理函数参数替换
			allParams := e.GetAllParams()
			for _, param := range allParams {
				if v, exists := keyValues[param.Key]; exists {

					if v == nil {
						e.DataString = nil
					} else {
						replacement := "0"
						if v != "" {
							replacement = v.(string)
						}
						e.DataString = strings.ReplaceAll(e.DataString.(string), param.Key, replacement)
					}

				}
			}
			e.RealityString = e.DataString

			// 集合函数空值默认初始化为0
			if e.FunctionType == FunctionAvg || e.FunctionType == FunctionSum {
				zeroInit = true
			}
		} else {
			// 处理普通数据节点
			if v, exists := keyValues[e.DataString.(string)]; exists {
				if v != nil && v == "" && zeroInit {
					e.RealityString = "0"
				} else if v == nil {
					e.RealityString = nil
				} else {
					e.RealityString = v.(string)
				}
				*result = append(*result, KeyValuePairElementString{Key: e.DataString.(string), Value: e.RealityString})
			} else {
				e.RealityString = e.DataString
			}
		}
	} else {
		e.RealityString = e.DataString
	}

	// 递归处理子表达式
	for _, childExp := range e.ExpressionChildren {
		childExp.loadArgumentWithDictionary(keyValues, result, zeroInit)
	}
}

func (e *Expression) GetAllParams() []KeyValuePairElementType {
	var results []KeyValuePairElementType

	if e.ElementType == ElementData && len(e.ExpressionChildren) == 0 {
		// 处理基础数据节点
		results = append(results, KeyValuePairElementType{
			Key:   strings.Replace(e.DataString.(string), "\\", "", -1),
			Value: e.ElementType,
		})
	} else {
		// 递归获取子节点参数
		results = append(results, e.getChildrenAllParams(e)...)
	}

	return results
}

func (e *Expression) getChildrenAllParams(parent *Expression) []KeyValuePairElementType {
	var childrenResults []KeyValuePairElementType

	for _, childExp := range e.ExpressionChildren {
		switch childExp.ElementType {
		case ElementExpression:
			// 递归处理表达式类型子节点
			childrenResults = append(childrenResults, childExp.getChildrenAllParams(childExp)...)

		case ElementFunction:
			if len(childExp.ExpressionChildren) > 0 {
				// 有子表达式则递归处理
				childrenResults = append(childrenResults, childExp.getChildrenAllParams(childExp)...)
			} else {
				// 处理函数参数
				paramList := strings.Split(childExp.DataString.(string), ",")
				for _, param := range paramList {
					childrenResults = append(childrenResults, KeyValuePairElementType{
						Key:   strings.Replace(param, "\\", "", -1),
						Value: ElementFunction,
					})
				}
			}

		case ElementData, ElementReference:
			// 处理数据和引用类型
			paramType := ElementData
			if parent != nil && parent.ElementType != ElementExpression {
				paramType = parent.ElementType
			}

			childrenResults = append(childrenResults, KeyValuePairElementType{
				Key:   strings.Replace(childExp.DataString.(string), "\\", "", -1),
				Value: paramType,
			})
		}
	}

	return childrenResults
}

/****************************************load****************************************************/

/****************************************parse****************************************************/

func tryParse(exp *Expression) bool {
	parse(exp)
	exp.RebuildExpression()
	return true
}

func IsOver(expressionString string) bool {
	if len(expressionString) == 0 {
		return true
	} else {
		return !(Contains(expressionString, '(') || Contains(expressionString, '[') || Contains(expressionString, '&') || Contains(expressionString, '|') || Contains(expressionString, '!') || Contains(expressionString, '>') || Contains(expressionString, '<') || Contains(expressionString, '=') || Contains(expressionString, '+') || Contains(expressionString, '-') || Contains(expressionString, '*') || Contains(expressionString, '/') || Contains(expressionString, '%'))
	}
}

func Contains(text string, contains byte) bool {
	var lastChar byte
	byteArray := []byte(text)
	for i := 0; i < len(byteArray); i++ {
		if text[i] == contains {
			if lastChar != '\\' {
				return true
			}
		}
		lastChar = byteArray[i]
	}
	return false
}

func SetMatchMode(currentChar byte, lastMode MatchMode) (matchMode MatchMode, endTag byte) {
	//go里没有nullable类型，常量又不能取地址，所以此处用空格字符代替nil
	switch currentChar {
	case '(':
		return MatchModeScope, ')'
	case '"':
		return MatchModeScope, '"'
	case '\'':
		return MatchModeScope, '\''
	case '[':
		return MatchModeFunction, ']'
	case '&':
		return MatchModeLogicSymbol, ' '
	case '|':
		return MatchModeLogicSymbol, ' '
	case '!':
		return MatchModeLogicSymbol, ' '
	case '+':
		return MatchModeArithmeticSymbol, ' '
	case '-':
		//有可能是负号，也有可能是减号;上一个block是符号或者none，这此处应该当作负号处理
		if lastMode == MatchModeUnknown || lastMode == MatchModeArithmeticSymbol || lastMode == MatchModeLogicSymbol || lastMode == MatchModeRelationSymbol {
			return MatchModeData, ' '
		}
		return MatchModeArithmeticSymbol, ' '
	case '*':
		return MatchModeArithmeticSymbol, ' '
	case '/':
		return MatchModeArithmeticSymbol, ' '
	case '%':
		return MatchModeArithmeticSymbol, ' '
	case '<':
		return MatchModeRelationSymbol, ' '
	case '>':
		return MatchModeRelationSymbol, ' '
	case '=':
		/*=继承上一个相邻符号的类型，比如<=,>=，此时=号为关系运算符；上一个为逻辑运算符的话，此处=为逻辑运算符，比如 !=；如果上一个block不为符号，那么此时=为等于（关系运算符）
		因此，只有上一个block为逻辑运算符时，才返回logicSymbol，其他情况返回relationSymbol
		*/
		if lastMode == MatchModeLogicSymbol {
			return MatchModeLogicSymbol, ' '
		}
		return MatchModeRelationSymbol, ' '
	case '\\':
		return MatchModeEscapeCharacter, ' '
	default:
		return MatchModeData, ' '
	}
}

func parse(exp *Expression) {
	lastBlock := MatchModeUnknown
	for index := 0; index < len(exp.SourceExpressionString); index++ {
		var matchScope match_Scope
		currentChar := exp.SourceExpressionString[index]
		mode, endTag := SetMatchMode(currentChar, lastBlock)
		switch mode {
		case MatchModeScope:
			if currentChar == endTag {
				//'' 或者 "" 实际上应该认作数据类型
				matchScope = findDataEnd(currentChar, exp.SourceExpressionString, index)
				tempStr := fmt.Sprintf("%c%s%c", currentChar, matchScope.ChildrenExpressionString, endTag)
				dataExp := &Expression{
					ElementType:            ElementData,
					SourceExpressionString: tempStr,
					DataString:             matchScope.ChildrenExpressionString,
				}
				exp.ExpressionChildren = append(exp.ExpressionChildren, dataExp)
				lastBlock = MatchModeData
				index = matchScope.EndIndex
				continue
			} else {
				matchScope = findEnd(currentChar, endTag, exp.SourceExpressionString, index)
			}
			exp.Status = matchScope.Status
			break
		case MatchModeRelationSymbol:
			var relationSymbolStr = getFullSymbol(exp.SourceExpressionString, index, mode)
			//去除可能存在的空字符
			var relationSymbol = convertOperator(strings.Replace(relationSymbolStr, " ", "", -1))
			exp.Operators = append(exp.Operators, &relationSymbol)
			exp.ElementType = ElementExpression
			//如果关系运算符为单字符，则索引+0，如果为多字符（<和=中间有空格，需要忽略掉），则跳过这段。eg: <；<=；<  =；
			index += len(relationSymbolStr) - 1
			lastBlock = mode
			continue
		case MatchModeLogicSymbol:
			var logicSymbolStr = getFullSymbol(exp.SourceExpressionString, index, mode)
			var logicSymbol = convertOperator(strings.Replace(logicSymbolStr, " ", "", -1))
			//因为! 既可以单独修饰一个数据，当作逻辑非，也可以与=联合修饰两个数据，当作不等于，所以此处需要进行二次判定。如果是!=，则此符号为关系运算符
			exp.Operators = append(exp.Operators, &logicSymbol)
			exp.ElementType = ElementExpression
			index += len(logicSymbolStr) - 1
			lastBlock = mode
			continue
		case MatchModeArithmeticSymbol:
			var operatorSymbol = convertOperator(fmt.Sprintf("%c", currentChar))
			exp.Operators = append(exp.Operators, &operatorSymbol)
			exp.ElementType = ElementExpression
			lastBlock = mode
			continue
		case MatchModeFunction:
			matchScope = findEnd('[', endTag, exp.SourceExpressionString, index)
			//确定函数类型
			var executeType, function = GetFunctionType(matchScope.ChildrenExpressionString)
			functionStr := "[" + matchScope.ChildrenExpressionString + "]"
			//如果是函数，则匹配函数内的表达式,eg: [sum](****)
			matchScope = findEnd('(', ')', exp.SourceExpressionString, matchScope.EndIndex+1)
			functionStr += "(" + matchScope.ChildrenExpressionString + ")"
			functionExp := &Expression{
				ElementType:            ElementFunction,
				FunctionType:           executeType,
				Function:               function,
				FunctionName:           executeType.String(),
				SourceExpressionString: functionStr,
				DataString:             matchScope.ChildrenExpressionString,
			}

			exp.ExpressionChildren = append(exp.ExpressionChildren, functionExp)
			var paramList = splitParamObject(matchScope.ChildrenExpressionString)
			for _, v := range paramList {
				paramExp, _ := CreateExpression(v)
				functionExp.ExpressionChildren = append(functionExp.ExpressionChildren, paramExp)
			}
			//函数解析完毕后直接从函数后面位置继续
			index = matchScope.EndIndex
			lastBlock = mode
			continue
		case MatchModeData:
			if currentChar == ' ' {
				continue
			}
			lastBlock = mode
			str, dataMtachMode := GetFullData(exp.SourceExpressionString, index, lastBlock)
			if len(str) != 0 {
				//todo 排除转义符长度
				if str == exp.SourceExpressionString {
					exp.ElementType = ElementData
					exp.DataString = str
					return
				}
				dataExp, _ := CreateExpression(str)
				if dataMtachMode == MatchModeScope && currentChar == '-' {
					//如果在Data分支下获取完整数据包含范围描述符号，即小括号，则认为这个负号修饰的是表达式，增加一个负号运算符
					symbol := Negative
					exp.Operators = append(exp.Operators, &symbol)
					continue
				}
				exp.ExpressionChildren = append(exp.ExpressionChildren, dataExp)
			}
			index += len(str) - 1
			continue
		case MatchModeEscapeCharacter:
			//跳过转义符号
			index++
			lastBlock = mode
			continue
		default:
			break
		}
		if !exp.Status {
			break
		}
		// 递归解析子表达式
		var isOver = exp.ElementType == ElementData || IsOver(matchScope.ChildrenExpressionString)
		if !isOver {
			expressionChildren, _ := CreateExpression(matchScope.ChildrenExpressionString)
			exp.ExpressionChildren = append(exp.ExpressionChildren, expressionChildren)
		}
		// 跳过已解析的块
		index = matchScope.EndIndex
		lastBlock = mode
	}
}

func splitParamObject(srcString string) []string {
	var result []string
	paramString := ""
	areaLevel := 0
	for i := 0; i < len(srcString); i++ {
		var currentChar = srcString[i]
		//()或[]封闭空间内的参数分隔符 , 需要忽略，因为它属于子表达式范围，不用在本层级分析，只把它当作普通字符即可
		switch currentChar {
		case ',':
			if len(paramString) != 0 && areaLevel == 0 {
				result = append(result, paramString)
				paramString = ""
				continue
			}
			break
		case '(', '[':
			//封闭空间开始,提升层级
			areaLevel++
			break
		case ')', ']':
			//封闭空间结束,降低层级
			areaLevel--
			break
		default:
			break
		}
		paramString = fmt.Sprintf("%s%c", paramString, currentChar)
	}
	if len(paramString) != 0 {
		result = append(result, paramString)
	}
	return result
}

func findEnd(startTag byte, endTag byte, exp string, index int) match_Scope {
	result := match_Scope{
		Status:                   true,
		EndIndex:                 -1,
		ChildrenExpressionString: "",
	}
	currentLevel := 0
	expArray := []byte(exp)
	for ; index < len(exp); index++ {
		var currentChar = expArray[index]
		//跳过转义符及后面一个字符
		if currentChar == '\\' {
			result.ChildrenExpressionString = fmt.Sprintf("%s%c", result.ChildrenExpressionString, expArray[index])
			index++
			result.ChildrenExpressionString = fmt.Sprintf("%s%c", result.ChildrenExpressionString, expArray[index])
			continue
		}
		// 第一次匹配到startTag不加层级，因为它的层级就是0
		if currentChar == startTag {
			currentLevel++
			if currentLevel == 1 {
				continue
			}
		} else if currentChar == endTag {
			currentLevel--
		}
		// 层级相同且与结束标志一致，则返回结束标志索引
		if currentLevel == 0 && currentChar == endTag {
			result.EndIndex = index
			break
		}
		result.ChildrenExpressionString = fmt.Sprintf("%s%c", result.ChildrenExpressionString, currentChar)
	}
	if result.EndIndex == -1 {
		result.Status = false
	}
	return result
}

func findDataEnd(tag byte, exp string, index int) match_Scope {
	result := match_Scope{
		Status:                   true,
		EndIndex:                 -1,
		ChildrenExpressionString: "",
	}
	expArray := []byte(exp)
	for i := index + 1; i < len(exp); i++ {
		if expArray[i] == tag {
			result.EndIndex = i
			break
		}
		result.ChildrenExpressionString = fmt.Sprintf("%s%c", result.ChildrenExpressionString, expArray[i])
	}
	return result
}

func convertOperator(currentChar string) Operator {
	switch currentChar {
	case "&":
		return And
	case "|":
		return Or
	case "!":
		return Not
	case "+":
		return Plus
	case "-":
		//负号特殊,此处算作减号
		return Subtract
	case "*":
		return Multiply
	case "/":
		return Divide
	case "%":
		return Mod
	case ">":
		return GreaterThan
	case "<":
		return LessThan
	case "=":
		return Equals
	case "!=":
		return UnEquals
	case "<=", "=<":
		return LessThanOrEquals
	case ">=", "=>":
		return GreaterThanOrEquals
	}
	return None
}

func getFullSymbol(exp string, startIndex int, matchMode MatchMode) string {
	expArray := []byte(exp)
	if startIndex == len(exp) {
		return fmt.Sprintf("%c", expArray[len(exp)-1])
	}
	result := fmt.Sprintf("%c", exp[startIndex])
	for i := startIndex + 1; i < len(exp); i++ {
		if exp[i] == ' ' && i-startIndex == len(result) {
			result = fmt.Sprintf("%s%c", result, expArray[i])
			continue
		}
		mode, _ := SetMatchMode(exp[i], matchMode)
		if mode == MatchModeRelationSymbol && matchMode == MatchModeRelationSymbol {
			result = fmt.Sprintf("%s%c", result, expArray[i])
			break
		} else if mode == MatchModeLogicSymbol && expArray[startIndex] == '!' && matchMode == MatchModeLogicSymbol {
			result = fmt.Sprintf("%s%c", result, expArray[i])
			break
		}
		if mode == MatchModeData {
			break
		}
		matchMode = mode
	}
	return result
}

func GetFunctionType(key string) (executeType FunctionType, function func(FormulaAction, ...any) interface{}) {
	key = strings.ToLower(key)
	switch key {
	case "sum":
		return FunctionSum, FormulaAction.Sum
	case "avg":
		return FunctionAvg, FormulaAction.Avg
	case "contains":
		return FunctionContains, FormulaAction.Contains
	case "excluding":
		return FunctionContainsExcept, FormulaAction.Excluding
	case "equals":
		return FunctionEquals, FormulaAction.Equals
	case "startwith":
		return FunctionStartWith, FormulaAction.StartWith
	case "endwith":
		return FunctionEndWith, FormulaAction.EndWith
	case "different":
		return FunctionDifferent, FormulaAction.Different
	case "round":
		return FunctionRound, FormulaAction.Round
	case "edate":
		return FunctionEDate, FormulaAction.EDate
	case "eodate":
		return FunctionEoDate, FormulaAction.EODate
	case "nowtime":
		return FunctionNowTime, FormulaAction.NowTime
	case "timetostring":
		return FunctionTimeToString, FormulaAction.TimeToString
	case "days":
		return FunctionDays, FormulaAction.Days
	case "hours":
		return FunctionHours, FormulaAction.Hours
	case "minutes":
		return FunctionMinutes, FormulaAction.Minutes
	case "seconds":
		return FunctionSeconds, FormulaAction.Seconds
	case "millseconds":
		return FunctionMillSeconds, FormulaAction.MillSeconds
	case "isnull":
		return FunctionIsNull, FormulaAction.IsNull
	}
	panic(key + " 函数未定义")
}

func GetFullData(exp string, startIndex int, matchMode MatchMode) (value string, mode MatchMode) {
	expArray := []byte(exp)
	if startIndex == len(exp) {
		return fmt.Sprintf("%c", expArray[len(expArray)-1]), MatchModeData
	}
	result := fmt.Sprintf("%c", expArray[startIndex])
	for i := startIndex + 1; i < len(exp); i++ {
		mode, _ := SetMatchMode(exp[i], matchMode)
		switch mode {
		case MatchModeData:
			result = fmt.Sprintf("%s%c", result, expArray[i])
			matchMode = mode
			continue
		case MatchModeLogicSymbol:
			return result, MatchModeLogicSymbol
		case MatchModeArithmeticSymbol:
			return result, MatchModeArithmeticSymbol
		case MatchModeRelationSymbol:
			return result, MatchModeRelationSymbol
		case MatchModeScope:
			var matchScope = findEnd('(', ')', exp, i)
			return matchScope.ChildrenExpressionString, MatchModeScope
		case MatchModeEscapeCharacter:
			//跳过转义符及后面一个字符
			result = fmt.Sprintf("%s%c", result, expArray[i])
			result = fmt.Sprintf("%s%c", result, expArray[i+1])
			i++
			matchMode = mode
			continue
		default:
			return result, MatchModeData
		}
	}
	return result, MatchModeData
}

/****************************************parse****************************************************/

/****************************************execute****************************************************/

// Execute  * 运算优先级从高到低为：
// * 小括号：()
// * 非：!
// * 乘除：* /
// * 加减：+ -
// * 关系运算符：< > =
// * 逻辑运算符：& ||
// *
// * 如果是逻辑表达式，则返回值只有0或1，分别代表false和true/*
func (e *Expression) Execute() interface{} {
	var result = executeChildren(e)
	return result[0]
}

func executeChildren(exp *Expression) []interface{} {
	var childrenResults []interface{}
	if len(exp.ExpressionChildren) == 0 {
		v, err := exp.executeNode(exp)
		if err != nil {
			panic(err)
		}
		childrenResults = append(childrenResults, v)
		return childrenResults
	}
	for _, childExp := range exp.ExpressionChildren {
		v, err := exp.executeNode(childExp)
		if err != nil {
			panic(err)
		}
		childrenResults = append(childrenResults, v)
	}

	/*
	 * 优先级
	 * 1. 算术运算
	 * 2. 关系运算
	 * 3. 逻辑运算
	 *
	 * 【注】：因为针对优先级进行了表达式树的重构，所以每一层级的所有运算符都是同一优先级，因此，这里按照顺序执行即可
	 */
	if len(exp.Operators) == 0 {
		return childrenResults
	}
	var result = childrenResults[0]
	//计算逻辑与和逻辑或,顺序执行
	for i, _ := range exp.Operators {
		//非运算和负数特殊，它只需要一个操作数就可完成计算，其他运算符至少需要两个
		value := childrenResults[i]
		if *exp.Operators[i] != Not && *exp.Operators[i] != Negative {
			value = childrenResults[i+1]
		}
		switch *exp.Operators[i] {
		case None:
			break
		case And:
			if InterfaceToFloat64(result) != 0 && InterfaceToFloat64(value) != 0 {
				result = 1
			} else {
				result = 0
			}
			break
		case Or:
			if InterfaceToFloat64(result) != 0 || InterfaceToFloat64(value) != 0 {
				result = 1
			} else {
				result = 0
			}
			break
		case Not:
			if InterfaceToFloat64(value) == 0 {
				result = 1
			} else {
				result = 0
			}
			break
		case Plus:
			result = InterfaceToFloat64(result) + InterfaceToFloat64(value)
			break
		case Subtract:
			if isTime(result) && isTime(value) {
				result = result.(time.Time).Sub(value.(time.Time))
			} else {
				result = InterfaceToFloat64(result) - InterfaceToFloat64(value)
			}
			break
		case Multiply:
			result = InterfaceToFloat64(result) * InterfaceToFloat64(value)
			break
		case Divide:
			result = InterfaceToFloat64(result) / InterfaceToFloat64(value)
			break
		case Mod:
			result = math.Mod(InterfaceToFloat64(result), InterfaceToFloat64(value))
			break
		case GreaterThan:
			//当前数据是否为日期，如果为日期则按日期比较方式
			if !isTime(result) && !isTime(value) {
				if InterfaceToFloat64(result) > InterfaceToFloat64(value) {
					result = 1
				} else {
					result = 0
				}
			} else {
				if result.(time.Time).After(value.(time.Time)) {
					result = 1
				} else {
					result = 0
				}
			}
			break
		case LessThan:
			//当前数据是否为日期，如果为日期则按日期比较方式
			if !isTime(result) && !isTime(value) {
				if InterfaceToFloat64(result) < InterfaceToFloat64(value) {
					result = 1
				} else {
					result = 0
				}
			} else {
				if result.(time.Time).Before(value.(time.Time)) {
					result = 1
				} else {
					result = 0
				}
			}

			break
		case Equals:
			//当前数据是否为日期，如果为日期则按日期比较方式
			if !isTime(result) && !isTime(value) {
				if result == value {
					result = 1
				} else {
					result = 0
				}
			} else {
				if result.(time.Time) == value.(time.Time) {
					result = 1
				} else {
					result = 0
				}
			}
			break
		case UnEquals:
			if !isTime(result) && !isTime(value) {
				if InterfaceToFloat64(result) != InterfaceToFloat64(value) {
					result = 1
				} else {
					result = 0
				}
			} else {
				if result.(time.Time) == value.(time.Time) {
					result = 1
				} else {
					result = 0
				}
			}
			break
		case GreaterThanOrEquals:
			if !isTime(result) && !isTime(value) {
				if InterfaceToFloat64(result) >= InterfaceToFloat64(value) {
					result = 1
				} else {
					result = 0
				}
			} else {
				if result.(time.Time).After(value.(time.Time)) || result == value {
					result = 1
				} else {
					result = 0
				}
			}
			break
		case LessThanOrEquals:
			if !isTime(result) && !isTime(value) {
				if InterfaceToFloat64(result) <= InterfaceToFloat64(value) {
					result = 1
				} else {
					result = 0
				}
			} else {
				if result.(time.Time).Before(value.(time.Time)) || result == value {
					result = 1
				} else {
					result = 0
				}
			}
			break
		case Negative:
			result = InterfaceToFloat64(value) * -1
			break
		default:
			break
		}
	}
	childrenResults = []interface{}{}
	childrenResults = append(childrenResults, result)
	return childrenResults
}

func (e *Expression) executeNode(childExp *Expression) (interface{}, error) {
	switch childExp.ElementType {
	case ElementExpression:
		return childExp.Execute(), nil

	case ElementData:
		return convert2ObjectValue(childExp.RealityString, childExp.SourceExpressionString)

	case ElementFunction:
		if childExp.Function == nil {
			return nil, fmt.Errorf("at %s: 不存在函数实例 %v", e.SourceExpressionString, childExp.FunctionType)
		}

		var result interface{} = 0.0
		switch childExp.FunctionType {
		case FunctionNone:
			result = childExp.Function(FormulaAction{})

		// 集合函数(参数不固定)
		case FunctionSum, FunctionAvg:
			var paramList []interface{}
			if e.countNonDataChildren(childExp) != 0 {
				for _, child := range childExp.ExpressionChildren {
					childrenResults := executeChildren(child)
					paramList = append(paramList, childrenResults...)
				}
			} else {
				if childExp.RealityString == "" {
					return nil, fmt.Errorf("at %s: 函数 %v 形参 %s 映射到实参 %s 错误",
						e.SourceExpressionString, childExp.FunctionType, childExp.DataString, childExp.RealityString)
				}
				dataArray := strings.Split(childExp.RealityString.(string), ",")
				for _, item := range dataArray {
					paramList = append(paramList, item)
				}
			}
			return childExp.Function(FormulaAction{}, paramList...), nil

		// 固定参数函数
		case FunctionCustomer, FunctionEDate, FunctionEoDate,
			FunctionNowTime, FunctionTimeToString, FunctionRound,
			FunctionContains, FunctionContainsExcept, FunctionEquals,
			FunctionStartWith, FunctionEndWith, FunctionDifferent,
			FunctionDays, FunctionHours, FunctionMinutes,
			FunctionSeconds, FunctionMillSeconds, FunctionIsNull:

			var paramsList []interface{}
			if e.countNonDataChildren(childExp) != 0 {
				for _, child := range childExp.ExpressionChildren {
					childrenResults := executeChildren(child)
					paramsList = append(paramsList, childrenResults...)
				}
			} else {
				if childExp.RealityString == nil {
					paramsList = nil
				} else if childExp.RealityString != "" {
					items := strings.Split(childExp.RealityString.(string), ",")
					for _, item := range items {
						paramsList = append(paramsList, item)
					}
				}
			}

			res := childExp.Function(FormulaAction{}, paramsList)
			return res, nil

		default:
			return nil, fmt.Errorf("at %s: 未知函数类型 %v", e.SourceExpressionString, childExp.FunctionType)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("at %s: 未知表达式节点", e.SourceExpressionString)
	}
}

// 辅助方法：计算非Data类型的子节点数量
func (e *Expression) countNonDataChildren(exp *Expression) int {
	count := 0
	for _, child := range exp.ExpressionChildren {
		if child.ElementType != ElementData {
			count++
		}
	}
	return count
}

func convert2ObjectValue(tag interface{}, sourceExpressionString string) (interface{}, error) {
	if tag == nil {
		return nil, nil
	}
	switch tag {
	case "true":
		return 1.0, nil
	case "false":
		return 0.0, nil
	case "1":
		return 1.0, nil
	case "0":
		return 0.0, nil
	default:
		if strings.HasSuffix(tag.(string), "%") {
			percentStr := strings.TrimSuffix(tag.(string), "%")
			percentResult, err := strconv.ParseFloat(percentStr, 64)
			if err != nil {
				return nil, fmt.Errorf("at %s: %s 不是数值类型", sourceExpressionString, tag)
			}
			return percentResult * 0.01, nil
		}

		if result, err := strconv.ParseFloat(tag.(string), 64); err == nil {
			return result, nil
		}

		if t, err := time.Parse(time.RFC3339, tag.(string)); err == nil {
			return t, nil
		}

		// Try other common date formats if RFC3339 fails
		if t, err := time.Parse("2006-01-02", tag.(string)); err == nil {
			return t, nil
		}
		if t, err := time.Parse("2006-01-02 15:04:05", tag.(string)); err == nil {
			return t, nil
		}

		return tag, nil
	}
}

/****************************************execute****************************************************/

/****************************************build****************************************************/

// RebuildExpression 根据运算优先级重组表达式树
func (e *Expression) RebuildExpression() {
	if len(e.Operators) == 0 {
		return
	}

	for {
		// 获取所有不同的优先级级别
		levels := make(map[int]bool)
		for _, op := range e.Operators {
			levels[op.GetOperatorInfo().Level] = true
		}

		// 如果只有一种优先级，直接返回
		if len(levels) == 1 {
			return
		}

		// 找出最高优先级
		maxLevel := 0
		for level := range levels {
			if level > maxLevel {
				maxLevel = level
			}
		}

		// 获取需要合并的操作符索引组
		operatorGroups := GetTargetLevelOperators(e.Operators, maxLevel)

		for _, group := range operatorGroups {
			// 因为是倒序，所以起始位置是反的
			startIndex := group[len(group)-1]
			endIndex := group[0]

			childCount := endIndex - startIndex + 2
			children := e.getNewChildren(startIndex, childCount)

			// 获取对应的操作符
			childrenOperators := make([]*Operator, len(group))
			for i, idx := range group {
				childrenOperators[i] = e.Operators[idx]
			}
			// 合并为新表达式
			newExp := buildChildren(children, childrenOperators)
			// 在原集合中删除合并的部分并插入新表达式
			e.ExpressionChildren = append(
				e.ExpressionChildren[:startIndex],
				append([]*Expression{&newExp}, e.ExpressionChildren[startIndex:]...)...,
			)
			// 从e.ExpressionChildren中删除children
			removeElementsInCollection(&e.ExpressionChildren, children)
			// 删除已合并的操作符（需要从后往前删除以避免索引变化）
			sort.Sort(sort.Reverse(sort.IntSlice(group)))
			for _, idx := range group {
				e.Operators = append(e.Operators[:idx], e.Operators[idx+1:]...)
			}
		}
	}
}

func removeElementsInCollection(srcExps *[]*Expression, removeList []*Expression) {
	//遍历srcExps,移除存在于removeList中的元素
	for i := len(*srcExps) - 1; i >= 0; i-- {
		for j := len(removeList) - 1; j >= 0; j-- {
			if (*srcExps)[i] == removeList[j] {
				removeList = append(removeList[:j], removeList[j+1:]...)
				*srcExps = append((*srcExps)[:i], (*srcExps)[i+1:]...)
				break
			}
		}
	}
}

func (e *Expression) getNewChildren(startIndex int, count int) []*Expression {
	var newChildren []*Expression
	for i := startIndex; i < len(e.ExpressionChildren); i++ {
		if count <= i-startIndex {
			return newChildren
		}
		newChildren = append(newChildren, e.ExpressionChildren[i])
	}
	return newChildren
}

func buildChildren(expressions []*Expression, operators []*Operator) Expression {
	var dataString = expressions[0].DataString.(string)
	if expressions[0].ElementType == ElementFunction {
		dataString = expressions[0].SourceExpressionString
	}

	for i := 1; i < len(expressions); i++ {
		childStr := expressions[i].DataString.(string)
		if expressions[i].ElementType == ElementExpression {
			childStr = expressions[i].SourceExpressionString
		}
		dataString += operators[i-1].GetOperatorInfo().Value + childStr
	}
	exp := Expression{
		ExpressionChildren:     expressions,
		Operators:              operators,
		DataString:             dataString,
		ElementType:            ElementExpression,
		SourceExpressionString: dataString,
		Status:                 true,
		FunctionType:           FunctionNone,
	}
	return exp
}

// GetTargetLevelOperators 获取相同级别且连续的子表达式运算符
func GetTargetLevelOperators(oldOperators []*Operator, level int) [][]int {
	/*
	 * eg:
	 * 序列为{2,3,1,3,3,2,3,3}, 输入为3，最终输出为各元素的索引集合，{1},{3,4},{6,7}
	 */
	var result [][]int
	var operators []int

	// 此处倒序循环是为了方便后续做删除操作，否则删除后索引的变化会导致数组越界
	for i := len(oldOperators) - 1; i >= 0; i-- {
		if oldOperators[i].GetOperatorInfo().Level == level {
			operators = append(operators, i)
		} else {
			if len(operators) > 0 {
				// Go中没有直接的DeepCopy，需要手动复制slice
				tmp := make([]int, len(operators))
				copy(tmp, operators)
				result = append(result, tmp)
				operators = operators[:0] // 清空slice
			}
		}
	}
	if len(operators) != 0 {
		result = append(result, operators)
	}
	return result
}

/****************************************build****************************************************/

/****************************************check****************************************************/

func (e *Expression) Check() error {
	err := checkExpression(e)
	if err != nil {
		return err
	}
	return nil
}

func checkExpression(expression *Expression) error {
	/*
	 * 1. 除了非运算只需要一个数据，其他的运算符至少需要2个数据
	 */
	if len(expression.Operators) > 0 {
		notOperatorCount := 0
		for _, op := range expression.Operators {
			if *op == Not {
				notOperatorCount++
			}
		}

		expectedChildren := len(expression.Operators) - notOperatorCount + 1
		if len(expression.ExpressionChildren) != expectedChildren {
			return fmt.Errorf("expression check error: data not match operator")
		}
	}

	for _, child := range expression.ExpressionChildren {
		if err := checkExpression(child); err != nil {
			return err
		}
	}

	return nil
}

/****************************************check****************************************************/

/****************************************private****************************************************/

func isTime(v interface{}) bool {
	_, ok := v.(time.Time)
	return ok
}

/****************************************private****************************************************/

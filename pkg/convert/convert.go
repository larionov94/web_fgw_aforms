package convert

import (
	"fgw_web_aforms/pkg/common"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const SkipNumOfStackFrame = 3

func pathToStrCode() string {
	var funcName, fileName, lineNumber, filePath = common.FileWithFuncAndLineNum(SkipNumOfStackFrame)

	return fmt.Sprintf("%s -> %s -> %d -> %s", funcName, fileName, lineNumber, filePath)
}

// ConvStrToInt конвертировать строку в число.
func ConvStrToInt(str string) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- значение: [%v]", err.Error(), pathToStrCode(), value)

		return 0
	}

	return value
}

// ConvStrToFloat конвертировать строку в вещественное число.
func ConvStrToFloat(str string) float64 {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- значение: [%v]", err.Error(), pathToStrCode(), value)

		return 0
	}

	return value
}

// ParseFormFieldInt преобразует поле в целое число, полученное из HTTP запроса.
func ParseFormFieldInt(r *http.Request, fieldName string) int {

	formValue := r.FormValue(fieldName)
	if formValue == "" {
		formValue = "0"

		return 0
	}
	value, err := strconv.Atoi(formValue)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- поле: [%s] --- значение: [%v]", err.Error(), pathToStrCode(), fieldName, value)

		return 0
	}

	return value
}

// ParseFormFieldFloat преобразует поле в вещественное число, полученное из HTTP запроса.
func ParseFormFieldFloat(r *http.Request, fieldName string) float64 {
	formValue := r.FormValue(fieldName)
	if formValue == "" {
		formValue = "0"

		return 0
	}
	value, err := strconv.ParseFloat(formValue, 64)
	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- поле: [%s] --- значение: [%v]", err.Error(), pathToStrCode(), fieldName, value)

		return 0
	}

	return value
}

// ParseFormFieldBool преобразует поле в булево значение, полученное из HTTP запроса.
func ParseFormFieldBool(r *http.Request, fieldName string) bool {
	formValue := r.FormValue(fieldName)
	if formValue == "" {
		formValue = "false"
	} else {
		formValue = "true"
	}

	value, err := strconv.ParseBool(formValue)

	if err != nil {
		log.Printf("Ошибка: [%s] --- ссылка на код: [ %s ] --- поле: [%s] --- значение: [%v]", err.Error(), pathToStrCode(), fieldName, value)

		return false
	}

	return value
}

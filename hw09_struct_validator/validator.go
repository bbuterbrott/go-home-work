package hw09_struct_validator //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	validateTag = "validate"
)

type ValidationError struct {
	Field string
	Err   error
}

type LenValidationError struct {
	len int
}

func (e LenValidationError) Error() string {
	return fmt.Sprintf("length must be %v", e.len)
}

type RegexpValidationError struct {
	regexp string
}

func (e RegexpValidationError) Error() string {
	return fmt.Sprintf("must match regex '%v'", e.regexp)
}

type InValidationError struct {
	values []string
}

func (e InValidationError) Error() string {
	return fmt.Sprintf("value must be one of '%v'", e.values)
}

func (e InValidationError) Is(target error) bool {
	return errors.As(e, &target)
}

type MinValidationError struct {
	min int
}

func (e MinValidationError) Error() string {
	return fmt.Sprintf("value must be greater than '%v'", e.min)
}

type MaxValidationError struct {
	max int
}

func (e MaxValidationError) Error() string {
	return fmt.Sprintf("value must be less than '%v'", e.max)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var b strings.Builder
	b.WriteString("Validation errors: ")
	for i, err := range v {
		b.WriteString(fmt.Sprintf("field '%v': %v", err.Field, err.Err.Error()))
		if i != len(v)-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}

func Validate(v interface{}) error {
	var errs ValidationErrors

	vValue := reflect.ValueOf(v)
	vType := vValue.Type()
	for i := 0; i < vType.NumField(); i++ {
		structField := vType.Field(i)
		value := vValue.Field(i)
		tag := structField.Tag.Get(validateTag)
		if len(tag) == 0 {
			continue
		}

		vals := strings.Split(tag, "|")
		for _, val := range vals {
			splitVal := strings.Split(val, ":")
			if len(splitVal) != 2 {
				log.Printf("incorrect tag '%v' for field '%v'", tag, structField.Name)
				continue
			}

			validationName := splitVal[0]
			validationValue := splitVal[1]
			fieldErrs := validateField(validationName, validationValue, structField.Name, value)
			if fieldErrs != nil {
				errs = append(errs, fieldErrs...)
			}
		}
	}

	return errs
}

func validateField(validationName string, validationValue string, fieldName string, fieldValue reflect.Value) []ValidationError {
	fmt.Printf("validating field '%v' with tag name: '%v', tag value: '%v' and field value '%v'\n", fieldName, validationName, validationValue, fieldValue)
	switch validationName {
	case "len":
		return validateLen(validationValue, fieldName, fieldValue)

	case "regexp":
		return validateRegexp(validationValue, fieldName, fieldValue)

	case "in":
		return validateIn(validationValue, fieldName, fieldValue)

	case "min":
		return validateMin(validationValue, fieldName, fieldValue)

	case "max":
		return validateMax(validationValue, fieldName, fieldValue)

	default:
		fmt.Printf("unknown validation '%v'\n", validationName)
	}

	return nil
}

//nolint:dupl Конечно, эти методы очень похожи, но я считаю, что лучше логику оставить раздельной, чтобы её можно было исправлять по отдельности
func validateLen(validationValue string, fieldName string, fieldValue reflect.Value) []ValidationError {
	var errs []ValidationError
	expLen, err := strconv.Atoi(validationValue)
	if err != nil {
		fmt.Println("len validation value must be int")
		return nil
	}

	k := fieldValue.Kind()
	switch k {
	case reflect.String:
		s := fieldValue.String()
		if len(s) != expLen {
			return append(errs, ValidationError{Field: fieldName, Err: LenValidationError{expLen}})
		}
		return nil

	case reflect.Slice:
		for i := 0; i < fieldValue.Len(); i++ {
			sliceElem := fieldValue.Index(i)
			if valErrs := validateLen(validationValue, fieldName+"["+strconv.Itoa(i)+"]", sliceElem); valErrs != nil {
				errs = append(errs, valErrs...)
			}
		}
		return errs

	default:
		fmt.Println("invalid field type for len validation")
		return nil
	}
}

func validateRegexp(validationValue string, fieldName string, fieldValue reflect.Value) []ValidationError {
	var errs []ValidationError
	r, err := regexp.Compile(validationValue)
	if err != nil {
		fmt.Printf("incorrect regexp '%v'\n", validationValue)
		return nil
	}

	k := fieldValue.Kind()
	switch k {
	case reflect.String:
		s := fieldValue.String()
		if len(s) == 0 {
			return nil
		}
		if !r.MatchString(s) {
			return append(errs, ValidationError{Field: fieldName, Err: RegexpValidationError{validationValue}})
		}
		return nil

	case reflect.Slice:
		for i := 0; i < fieldValue.Len(); i++ {
			sliceElem := fieldValue.Index(i)
			if valErrs := validateRegexp(validationValue, fieldName+"["+strconv.Itoa(i)+"]", sliceElem); valErrs != nil {
				errs = append(errs, valErrs...)
			}
		}
		return errs

	default:
		fmt.Println("invalid field type for regex validation")
		return nil
	}
}

func validateIn(validationValue string, fieldName string, fieldValue reflect.Value) []ValidationError {
	var errs []ValidationError
	expValues := strings.Split(validationValue, ",")

	k := fieldValue.Kind()
	switch k {
	case reflect.String:
		s := fieldValue.String()
		for _, expValue := range expValues {
			if expValue == s {
				return nil
			}
		}
		return append(errs, ValidationError{Field: fieldName, Err: InValidationError{expValues}})

	case reflect.Int:
		v := fieldValue.Int()

		for _, expStringValue := range expValues {
			expValue, err := strconv.Atoi(expStringValue)
			if err != nil {
				fmt.Println("all in validation value must be int")
				return nil
			}

			if int(v) == expValue {
				return nil
			}
		}
		return append(errs, ValidationError{Field: fieldName, Err: InValidationError{expValues}})

	case reflect.Slice:
		for i := 0; i < fieldValue.Len(); i++ {
			sliceElem := fieldValue.Index(i)
			if valErrs := validateIn(validationValue, fieldName+"["+strconv.Itoa(i)+"]", sliceElem); valErrs != nil {
				errs = append(errs, valErrs...)
			}
		}
		return errs

	default:
		fmt.Println("invalid field type for in validation")
		return nil
	}
}

//nolint:dupl
func validateMin(validationValue string, fieldName string, fieldValue reflect.Value) []ValidationError {
	var errs []ValidationError
	min, err := strconv.Atoi(validationValue)
	if err != nil {
		fmt.Println("min validation value must be int")
		return nil
	}

	k := fieldValue.Kind()
	switch k {
	case reflect.Int:
		v := fieldValue.Int()
		if int(v) < min {
			return append(errs, ValidationError{Field: fieldName, Err: MinValidationError{min}})
		}
		return nil

	case reflect.Slice:
		for i := 0; i < fieldValue.Len(); i++ {
			sliceElem := fieldValue.Index(i)
			if valErrs := validateMin(validationValue, fieldName+"["+strconv.Itoa(i)+"]", sliceElem); valErrs != nil {
				errs = append(errs, valErrs...)
			}
		}
		return errs

	default:
		fmt.Println("invalid field type for min validation")
		return nil
	}
}

func validateMax(validationValue string, fieldName string, fieldValue reflect.Value) []ValidationError { //nolint:dupl
	var errs []ValidationError
	max, err := strconv.Atoi(validationValue)
	if err != nil {
		fmt.Println("max validation value must be int")
		return nil
	}

	k := fieldValue.Kind()
	switch k {
	case reflect.Int:
		v := fieldValue.Int()
		if int(v) > max {
			return append(errs, ValidationError{Field: fieldName, Err: MaxValidationError{max}})
		}

	case reflect.Slice:
		for i := 0; i < fieldValue.Len(); i++ {
			sliceElem := fieldValue.Index(i)
			if valErrs := validateMax(validationValue, fieldName+"["+strconv.Itoa(i)+"]", sliceElem); valErrs != nil {
				errs = append(errs, valErrs...)
			}
		}
		return errs

	default:
		fmt.Println("invalid field type for max validation")
	}

	return nil
}

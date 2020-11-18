package hw09_struct_validator //nolint:golint,stylecheck

import (
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

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var b strings.Builder
	b.WriteString("Validation errors: ")
	for i, err := range v {
		b.WriteString("field '")
		b.WriteString(err.Field)
		b.WriteString("': ")
		b.WriteString(err.Err.Error())
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
			err := validateField(validationName, validationValue, structField.Name, value)
			if err != nil {
				errs = append(errs, *err)
			}
		}
	}

	return errs
}

func validateField(validationName string, validationValue string, fieldName string, fieldValue reflect.Value) *ValidationError {
	fmt.Printf("validating field '%v' with tag name: '%v', tag value: '%v' and field value '%v'\n", fieldName, validationName, validationValue, fieldValue)
	switch validationName {
	case "len":
		return validateLen(validationValue, fieldName, fieldValue)

	case "regexp":
		return validateRegex(validationValue, fieldName, fieldValue)

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

func validateLen(validationValue string, fieldName string, fieldValue reflect.Value) *ValidationError {
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
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("length must be %v", expLen)}
		}

	case reflect.Slice:
		sliceErrs := make([]ValidationError, 0, fieldValue.Len())
		for i := 0; i < fieldValue.Len(); i++ {
			sliceElem := fieldValue.Index(i)
			if valErr := validateLen(validationValue, fieldName, sliceElem); valErr != nil {
				sliceErrs = append(sliceErrs, *valErr)
			}
		}
		if len(sliceErrs) != 0 {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("all elements in slice must have length %v", expLen)}
		}

	default:
		fmt.Println("invalid field type for len validation")
		return nil
	}

	return nil
}

func validateRegex(validationValue string, fieldName string, fieldValue reflect.Value) *ValidationError {
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
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("must match regex '%v'", validationValue)}
		}

	case reflect.Slice:
		sliceErrs := make([]ValidationError, 0)
		for i := 0; i < fieldValue.Len(); i++ {
			sliceElem := fieldValue.Index(i)
			if valErr := validateRegex(validationValue, fieldName, sliceElem); valErr != nil {
				sliceErrs = append(sliceErrs, *valErr)
			}
		}
		if len(sliceErrs) != 0 {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("all elements in slice must match regex '%v'", validationValue)}
		}

	default:
		fmt.Println("invalid field type for regex validation")
		return nil
	}

	return nil
}

func validateIn(validationValue string, fieldName string, fieldValue reflect.Value) *ValidationError {
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
		return &ValidationError{Field: fieldName, Err: fmt.Errorf("value must be one of '%v'", expValues)}

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
		return &ValidationError{Field: fieldName, Err: fmt.Errorf("value must be one of '%v'", expValues)}

	case reflect.Slice:
		sliceErrs := make([]ValidationError, 0)
		for i := 0; i < fieldValue.Len(); i++ {
			sliceElem := fieldValue.Index(i)
			if valErr := validateIn(validationValue, fieldName, sliceElem); valErr != nil {
				sliceErrs = append(sliceErrs, *valErr)
			}
		}
		if len(sliceErrs) != 0 {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("all elements in slice must have one value of '%v'", expValues)}
		}

	default:
		fmt.Println("invalid field type for in validation")
		return nil
	}

	return nil
}

//nolint:dupl Конечно, эти методы очень похожи, но я считаю, что лучше логику оставить раздельной, чтобы её можно было исправлять по отдельности
func validateMin(validationValue string, fieldName string, fieldValue reflect.Value) *ValidationError {
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
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("value must be greater than '%v'", min)}
		}
		return nil

	case reflect.Slice:
		sliceErrs := make([]ValidationError, 0)
		for i := 0; i < fieldValue.Len(); i++ {
			sliceElem := fieldValue.Index(i)
			if valErr := validateMin(validationValue, fieldName, sliceElem); valErr != nil {
				sliceErrs = append(sliceErrs, *valErr)
			}
		}
		if len(sliceErrs) != 0 {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("all elements in slice must have value greater than '%v'", min)}
		}

	default:
		fmt.Println("invalid field type for min validation")
		return nil
	}

	return nil
}

func validateMax(validationValue string, fieldName string, fieldValue reflect.Value) *ValidationError { //nolint:dupl
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
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("value must be less than '%v'", max)}
		}
		return nil

	case reflect.Slice:
		sliceErrs := make([]ValidationError, 0)
		for i := 0; i < fieldValue.Len(); i++ {
			sliceElem := fieldValue.Index(i)
			if valErr := validateMax(validationValue, fieldName, sliceElem); valErr != nil {
				sliceErrs = append(sliceErrs, *valErr)
			}
		}
		if len(sliceErrs) != 0 {
			return &ValidationError{Field: fieldName, Err: fmt.Errorf("all elements in slice must have value less than '%v'", max)}
		}

	default:
		fmt.Println("invalid field type for max validation")
		return nil
	}

	return nil
}

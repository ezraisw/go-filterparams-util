package fputil

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cbrand/go-filterparams/definition"
)

type DataType interface {
	IsFilterAllowed(filter *definition.Filter) bool
	Parse(value interface{}) (interface{}, bool)
}

var (
	trueStrings = map[string]bool{
		"true": true,
		"yes":  true,
		"y":    true,
	}

	falseStrings = map[string]bool{
		"false": true,
		"no":    true,
		"n":     true,
	}
)

type boolDataType struct {
}

func BoolDataType() DataType {
	return &boolDataType{}
}

func (dt boolDataType) IsFilterAllowed(filter *definition.Filter) bool {
	return filter.Identification == definition.FilterEq.Identification ||
		filter.Identification == definition.FilterIn.Identification
}

func (dt boolDataType) Parse(value interface{}) (interface{}, bool) {
	switch v := value.(type) {
	case string:
		v = strings.ToLower(v)
		if trueStrings[v] {
			return true, true
		}
		if falseStrings[v] {
			return false, true
		}
	case bool:
		return v, true
	}

	return nil, false
}

type uintDataType struct {
	bits int
}

func UintDataType(bits int) DataType {
	return &uintDataType{
		bits: bits,
	}
}

func (dt uintDataType) IsFilterAllowed(filter *definition.Filter) bool {
	return filter.Identification == definition.FilterEq.Identification ||
		filter.Identification == definition.FilterGt.Identification ||
		filter.Identification == definition.FilterGte.Identification ||
		filter.Identification == definition.FilterLt.Identification ||
		filter.Identification == definition.FilterLte.Identification ||
		filter.Identification == definition.FilterIn.Identification
}

func (dt uintDataType) Parse(value interface{}) (interface{}, bool) {
	switch v := value.(type) {
	case string:
		if n, err := strconv.ParseUint(v, 10, dt.bits); err == nil {
			return n, true
		}
	case uint64:
		return v, true
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	}

	return nil, false
}

type intDataType struct {
	bits int
}

func IntDataType(bits int) DataType {
	return &intDataType{
		bits: bits,
	}
}

func (dt intDataType) IsFilterAllowed(filter *definition.Filter) bool {
	return filter.Identification == definition.FilterEq.Identification ||
		filter.Identification == definition.FilterGt.Identification ||
		filter.Identification == definition.FilterGte.Identification ||
		filter.Identification == definition.FilterLt.Identification ||
		filter.Identification == definition.FilterLte.Identification ||
		filter.Identification == definition.FilterIn.Identification
}

func (dt intDataType) Parse(value interface{}) (interface{}, bool) {
	switch v := value.(type) {
	case string:
		if n, err := strconv.ParseInt(v, 10, dt.bits); err == nil {
			return n, true
		}
	case int64:
		return v, true
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	}

	return nil, false
}

type floatDataType struct {
	bits int
}

func FloatDataType(bits int) DataType {
	return &floatDataType{
		bits: bits,
	}
}

func (dt floatDataType) IsFilterAllowed(filter *definition.Filter) bool {
	return filter.Identification == definition.FilterEq.Identification ||
		filter.Identification == definition.FilterGt.Identification ||
		filter.Identification == definition.FilterGte.Identification ||
		filter.Identification == definition.FilterLt.Identification ||
		filter.Identification == definition.FilterLte.Identification ||
		filter.Identification == definition.FilterIn.Identification
}

func (dt floatDataType) Parse(value interface{}) (interface{}, bool) {
	switch v := value.(type) {
	case string:
		if n, err := strconv.ParseFloat(v, dt.bits); err == nil {
			return n, true
		}
	case float32:
		return float64(v), true
	case float64:
		return v, true
	}

	return nil, false
}

type stringDataType struct {
}

func StringDataType() DataType {
	return &stringDataType{}
}

func (dt stringDataType) IsFilterAllowed(filter *definition.Filter) bool {
	return filter.Identification == definition.FilterEq.Identification ||
		filter.Identification == definition.FilterLike.Identification ||
		filter.Identification == definition.FilterILike.Identification ||
		filter.Identification == definition.FilterIn.Identification
}

func (dt stringDataType) Parse(value interface{}) (interface{}, bool) {
	switch v := value.(type) {
	case string:
		return v, true
	}

	return fmt.Sprintf("%v", value), false
}

type timeDataType struct {
	inFormat  string
	outFormat string
}

func TimeDataType(inFormat string, outFormat string) DataType {
	return &timeDataType{
		inFormat:  inFormat,
		outFormat: outFormat,
	}
}

func (dt timeDataType) IsFilterAllowed(filter *definition.Filter) bool {
	return filter.Identification == definition.FilterEq.Identification ||
		filter.Identification == definition.FilterGt.Identification ||
		filter.Identification == definition.FilterGte.Identification ||
		filter.Identification == definition.FilterLt.Identification ||
		filter.Identification == definition.FilterLte.Identification ||
		filter.Identification == definition.FilterIn.Identification
}

func (dt timeDataType) Parse(value interface{}) (interface{}, bool) {
	switch v := value.(type) {
	case string:
		if t, err := time.Parse(dt.inFormat, v); err == nil {
			if dt.outFormat != "" {
				return t.Format(dt.outFormat), true
			}

			return t, true
		}
	}

	return nil, false
}

type nilableDataType struct {
	dataType DataType
}

func NilableDataType(dataType DataType) DataType {
	return nilableDataType{
		dataType: dataType,
	}
}

func (dt nilableDataType) IsFilterAllowed(filter *definition.Filter) bool {
	return dt.dataType.IsFilterAllowed(filter)
}

func (dt nilableDataType) Parse(value interface{}) (interface{}, bool) {
	switch value {
	case "nil", nil:
		return nil, true
	}

	return dt.dataType.Parse(value)
}

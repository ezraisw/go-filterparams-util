package fputil

import (
	"reflect"
	"time"

	"github.com/cbrand/go-filterparams"
	"github.com/cbrand/go-filterparams/definition"
)

type CleanSpec struct {
	Name     string
	DataType DataType
	Rules    []Rule
}

func GetCleanSpecs(v interface{}) []CleanSpec {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		panic("not a struct")
	}

	specs := make([]CleanSpec, 0)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		name := getName(field)
		if name == "" {
			continue
		}

		dataType := getDataType(field)
		if dataType == nil {
			continue
		}

		specs = append(specs, CleanSpec{
			Name:     name,
			DataType: dataType,

			// TODO: Possibly from tags?
			Rules: []Rule{},
		})
	}

	return specs
}

func getDataType(field reflect.StructField) DataType {
	nilable := false

	ft := field.Type
	if ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
		nilable = true
	}

	var dataType DataType
	switch ft.Kind() {
	case reflect.Bool:
		dataType = BoolDataType()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dataType = IntDataType(ft.Bits())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dataType = UintDataType(ft.Bits())
	case reflect.Float32, reflect.Float64:
		dataType = FloatDataType(ft.Bits())
	case reflect.String:
		dataType = StringDataType()
	case reflect.Struct:
		if ft.PkgPath() != "time" || ft.Name() != "Time" {
			return nil
		}

		dataType = TimeDataType(time.RFC3339, "")
	default:
		return nil
	}

	if nilable {
		dataType = NilableDataType(dataType)
	}

	return dataType
}

func CleanQuery(queryData *filterparams.QueryData, specs ...CleanSpec) *filterparams.QueryData {
	filter := CleanFilter(queryData.GetFilter(), specs...)
	orders := CleanOrders(queryData.GetOrders(), specs...)
	return filterparams.NewQueryData(filter, orders)
}

func CleanFilter(filter interface{}, specs ...CleanSpec) interface{} {
	mSpecs := make(map[string]CleanSpec)
	for _, spec := range specs {
		mSpecs[spec.Name] = spec
	}

	c := filterCleaning{
		mSpecs: mSpecs,
	}

	newFilter := c.clean(filter)

	if newFilter == true {
		return nil
	}

	return newFilter
}

type filterCleaning struct {
	mSpecs map[string]CleanSpec
}

// Clean the given filter.
//
// Can return either: true, false, or any derivative of *definition.Param
func (c filterCleaning) clean(filter interface{}) interface{} {
	switch f := filter.(type) {
	case *definition.And:
		return c.cleanAnd(f)
	case *definition.Or:
		return c.cleanOr(f)
	case *definition.Negate:
		return c.cleanNegate(f)
	case *definition.Parameter:
		return c.cleanParameter(f)
	}

	return false
}

func (c filterCleaning) cleanAnd(a *definition.And) interface{} {
	left := c.clean(a.Left)
	if left == false {
		return false
	}

	right := c.clean(a.Right)
	if right == false {
		return false
	}

	if left == true {
		// If both are true, it should return true through this path.
		return right
	} else if right == true {
		return left
	}

	newAnd := *a
	newAnd.Left = left
	newAnd.Right = right

	return &newAnd
}

func (c filterCleaning) cleanOr(o *definition.Or) interface{} {
	left := c.clean(o.Left)
	if left == true {
		return true
	}

	right := c.clean(o.Right)
	if right == true {
		return true
	}

	if left == false {
		// If both are false, it should return false through this path.
		return right
	} else if right == false {
		return left
	}

	newOr := *o
	newOr.Left = left
	newOr.Right = right

	return &newOr
}

func (c filterCleaning) cleanNegate(n *definition.Negate) interface{} {
	negated := c.clean(n.Negated)
	if negated == false {
		return true
	} else if negated == true {
		return false
	}

	newNegate := *n
	newNegate.Negated = negated

	return &newNegate
}

func (c filterCleaning) cleanParameter(pm *definition.Parameter) interface{} {
	spec, ok := c.mSpecs[pm.Name]
	if !ok {
		return false
	}

	if !spec.DataType.IsFilterAllowed(pm.Filter) {
		return false
	}

	newValue, ok := spec.DataType.Parse(pm.Value)
	if !ok {
		return false
	}

	for _, rule := range spec.Rules {
		if !rule.Validate(newValue) {
			return false
		}
	}

	newPm := *pm
	newPm.Value = newValue

	return &newPm
}

func CleanOrders(orders []*definition.Order, specs ...CleanSpec) []*definition.Order {
	mNames := make(map[string]bool)
	for _, spec := range specs {
		mNames[spec.Name] = true
	}

	newOrders := make([]*definition.Order, 0)
	for _, order := range orders {
		if mNames[order.GetOrderBy()] {
			newOrders = append(newOrders, order)
		}
	}

	return newOrders
}

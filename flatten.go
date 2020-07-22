package transform

import (
	"fmt"
	"reflect"
)

type Flattener struct {
	r        map[string]interface{}
	sep      string
	sliceSep string
}

// NewFlattener returns a new Flattener
func NewFlattener() *Flattener {
	return &Flattener{
		sep:      "/",
		sliceSep: ".",
		r:        make(map[string]interface{}),
	}
}

// SetSeparator //
func (f *Flattener) SetSeparator(sep string) { f.sep = sep }

// SetSliceSeparator //
func (f *Flattener) SetSliceSeparator(sep string) { f.sliceSep = sep }

// Flatten takes in a map[string]interface{} and returns a flattened map[string]interface{},
// nested fields are separated using f.sep and slice elements are indexed using f.sliceSep
func (f *Flattener) Flatten(in map[string]interface{}) (map[string]interface{}, error) {
	var err error
	for prefix, value := range in {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Interface {
			v = v.Elem()
		}
		err = f.flatten(prefix, v)
		if err != nil {
			return nil, err
		}
	}
	return f.r, nil
}
func (f *Flattener) flattenMap(prefix string, v reflect.Value) error {
	for _, k := range v.MapKeys() {
		if k.Kind() == reflect.Interface {
			k = k.Elem()
		}
		var newprefix string
		switch k.Kind() {
		case reflect.String:
			newprefix = prefix + f.sep + k.String()
		case reflect.Int:
			newprefix = prefix + f.sep + fmt.Sprintf("%d", k.Int())
		case reflect.Bool:
			newprefix = prefix + f.sep + fmt.Sprintf("%t", k.Bool())
		case reflect.Float64:
			newprefix = prefix + f.sep + fmt.Sprintf("%v", k.Float())
		default:
			return fmt.Errorf("unsupported map key '%v' type: %v", k, k.Kind())
		}
		f.flatten(newprefix, v.MapIndex(k))
	}
	return nil
}
func (f *Flattener) flattenSlice(prefix string, v reflect.Value) error {
	var err error
	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		if val.Kind() == reflect.Interface {
			val = val.Elem()
		}
		newprefix := prefix + f.sliceSep + fmt.Sprintf("%d", i)
		err = f.flatten(newprefix, val)
		if err != nil {
			return err
		}
	}
	return nil
}
func (f *Flattener) flatten(prefix string, v reflect.Value) error {
	var err error
	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Bool:
		f.r[prefix] = v.Bool()
	case reflect.Int:
		f.r[prefix] = v.Int()
	case reflect.Float64:
		f.r[prefix] = v.Float()
	case reflect.Map:
		err = f.flattenMap(prefix, v)
	case reflect.Slice:
		err = f.flattenSlice(prefix, v)
	case reflect.String:
		f.r[prefix] = v.String()
	default:
		return fmt.Errorf("unexpected type: %v, %v", v.Kind(), v)
	}
	return err
}

package utils

import "reflect"

func ConvertBoolToValue(b bool) float64 {
	if b {
		return 1
	} else {
		return 0
	}
}

func CopyStruct(src, dst interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcValue := srcVal.Field(i)
		name := srcVal.Type().Field(i).Name

		dstValue := dstVal.FieldByName(name)
		if !dstValue.IsValid() {
			continue
		}
		if dstValue.Type() != srcValue.Type() {
			continue
		}

		dstValue.Set(srcValue)
	}
}

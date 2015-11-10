package underscore

import (
	"errors"
	"reflect"
)

var EMPTY_ARRAY = make([]interface{}, 0)

func Map(source, selector interface{}) (interface{}, error) {
	selectorRT := reflect.TypeOf(selector)
	if !(selectorRT.Kind() == reflect.Func && selectorRT.NumIn() == 2 && selectorRT.NumOut() == 2) {
		return nil, errors.New("underscore: Map's selector is nil")
	}

	selectorRV := reflect.ValueOf(selector)
	resultRT := reflect.SliceOf(
		selectorRV.Type().Out(0),
	)
	if source == nil {
		return makeArray(resultRT), nil
	}

	resultRV := makeArray(resultRT)
	sourceRV := reflect.ValueOf(source)
	switch sourceRV.Kind() {
		case reflect.Array:
		case reflect.Slice:
			if sourceRV.Len() == 0 {
				return makeArray(resultRT), nil
			}

			for i := 0; i < sourceRV.Len(); i++ {
				values := selectorRV.Call(
					[]reflect.Value{
						sourceRV.Index(i),
						reflect.ValueOf(i),
					},
				)
				if !values[1].IsNil() {
					return nil, values[1].Interface().(error)
				}

				resultRV = reflect.Append(resultRV, values[0])
			}
		case reflect.Map:
			keyRVs := sourceRV.MapKeys()
			if len(keyRVs) == 0 {
				return makeArray(resultRT), nil
			}

			for i := 0; i < len(keyRVs); i++ {
				values := selectorRV.Call(
					[]reflect.Value{
						sourceRV.MapIndex(keyRVs[i]),
						keyRVs[i],
					},
				)
				if !values[1].IsNil() {
					return nil, values[1].Interface().(error)
				}

				resultRV = reflect.Append(resultRV, values[0])
			}
	}
	
	return resultRV.Interface(), nil
}

//chain
func (this *Query) Map(selector interface{}) Queryer {
	if this.err == nil {
		this.source, this.err = Map(this.source, selector)
	}
	return this
}
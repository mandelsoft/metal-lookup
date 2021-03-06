/*
 * Copyright 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package convert

import (
	"reflect"
	"strings"
)

var boolType = reflect.TypeOf((*bool)(nil)).Elem()

func BoolType() reflect.Type {
	return boolType
}

func AsBool(s interface{}) (bool, error) {
	if s == nil {
		return false, nil
	}
	switch v := s.(type) {
	case bool:
		return v, nil
	case *bool:
		return *v, nil
	default:
		i, err := ConvertTo(s, boolType)
		if err == nil {
			return i.(bool), nil
		}
		return false, err
	}
}

func Bool(s interface{}) bool {
	if v, err := AsBool(s); err == nil {
		return v
	}
	return false
}

func BestEffortBool(s interface{}) bool {
	v, err := AsBool(s)
	if err == nil {
		return v
	}

	if str, err := AsString(s); err == nil {
		switch strings.ToLower(str) {
		case "false", "off":
			return false
		case "true", "on":
			return true
		}
		return len(str) > 0
	}

	if i, err := AsInt(s); err == nil {
		return i != 0
	}

	value := reflect.ValueOf(s)
	switch value.Kind() {
	case reflect.Ptr:
		return !value.IsNil()
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Chan:
		return value.Len() > 0
	default:
		return false
	}
}

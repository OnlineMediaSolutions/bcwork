package structs

import (
	"errors"
	"reflect"
)

var types = make(map[string]reflect.Type)

//Register new type for dynamic instance
func Regsiter(j interface{}) {
	t := reflect.TypeOf(j)
	types[t.PkgPath()+"."+t.Name()] = t
}

//Register new type for dynamic instance with custom name
func RegsiterName(name string, j interface{}) {
	t := reflect.TypeOf(j)
	types[name] = t
}

func TypeNames() []string {
	names := make([]string, 0, len(types))
	for k := range types {
		names = append(names, k)
	}

	return names
}

func NewInstance(name string) (interface{}, error) {
	if v, ok := types[name]; ok {
		jrunner := reflect.New(v).Interface()
		return jrunner, nil
	}

	return nil, errors.New("Class '" + name + "' not found in class registry (make sure to use init() to register the class)")
}

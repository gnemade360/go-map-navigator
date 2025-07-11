package mapnavigator

import (
	"fmt"
	"reflect"
	"strconv"

	models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
)

// MapNavigator provides functionality for navigating and modifying nested map and slice data structures.
// It supports various navigation patterns including wildcard operations, conditional modifications,
// and property creation.
type MapNavigator struct {
	// NodeModifier is applied to nodes during navigation when the "-" key is encountered
	NodeModifier   models.MapNodeModifier
	// keys stores the current navigation path
	keys           []string
	// ReadOnly when true, prevents modifications to the data structure
	ReadOnly       bool
	// CreateProperty when true, creates missing properties during navigation
	CreateProperty bool
}

// VisitMapNode navigates a map with interface{} keys using the provided key path.
// It supports wildcard operations (*), modifier application (-), and nested navigation.
// Returns the value at the specified path or an error if navigation fails.
func (m *MapNavigator) VisitMapNode(mp map[interface{}]interface{}, ks ...string) (interface{}, error) {
	if len(m.keys) == 0 {
		m.keys = ks
	}
	k := ks[0]

	if mp == nil || len(k) == 0 {
		return nil, fmt.Errorf("map is nil")
	}
	if k == "*" {
		toReturn := make(map[interface{}]interface{}, 0)
		var err error
		for itmKey, itm := range mp {
			if obj, er := m.VisitNode(itm, ks[1:]...); er != nil {
				err = er
				continue
			} else {

				kind := reflect.TypeOf(itm).Kind()
				if (kind != reflect.Map && kind != reflect.Slice) || len(ks) < 2 {
					if !m.ReadOnly {
						mp[itmKey] = obj
					}
				}
				if !m.ReadOnly {
					toReturn[itmKey] = obj
				}

			}
		}
		return toReturn, err
	} else if k == "-" {
		if m.NodeModifier != nil {
			nmp := m.NodeModifier.ModifyNode(mp)

			return nmp, nil
		}
	} else if obj, ok := mp[k]; ok && obj != nil {
		//do something here
		if ttr, err := m.VisitNode(obj, ks[1:]...); err == nil {
			kind := reflect.TypeOf(obj).Kind()
			if (kind != reflect.Map && kind != reflect.Slice) || len(ks) < 2 {
				if !m.ReadOnly {
					mp[k] = ttr
				}
			}
			return ttr, nil
		} else {
			return nil, err
		}
	}

	if m.CreateProperty && len(ks) == 1 {
		mp[ks[0]] = ""
		return m.VisitMapNode(mp, ks...)
	}
	return nil, fmt.Errorf("key \"%v\" not found in map error in mapnavigator", k)
}

// VisitMapStringNode navigates a map with string keys using the provided key path.
// It supports wildcard operations (*), modifier application (-), and nested navigation.
// Returns the value at the specified path or an error if navigation fails.
func (m *MapNavigator) VisitMapStringNode(mp map[string]interface{}, ks ...string) (interface{}, error) {
	//if len(m.keys) == 0 {
	m.keys = ks
	//}
	k := ks[0]

	if mp == nil || len(k) == 0 {
		return nil, fmt.Errorf("map is nil")
	}
	if k == "*" {
		//toReturn := make(map[string]interface{}, 0)
		toReturn := make([]interface{}, 0)
		var err error
		for itmKey, itm := range mp {
			if obj, er := m.VisitNode(itm, ks[1:]...); er != nil {
				err = er
				if itm != nil {
					kind := reflect.TypeOf(itm).Kind()
					if (kind != reflect.Map && kind != reflect.Slice) || len(ks) < 2 {
						if !m.ReadOnly {
							if o, deleted := obj.(*MapNavigatorDeleted); deleted && o != nil {
								delete(mp, itmKey)
							} else {
								mp[itmKey] = obj
							}
						}
					}
					//toReturn[itmKey] = obj
					if !m.ReadOnly {
						if o, deleted := obj.(*MapNavigatorDeleted); deleted && o != nil {
						} else {
							toReturn = appendToToreturn(kind, toReturn, obj)
						}
					}
				}
				//continue

			} else {
				if itm != nil {

					kind := reflect.TypeOf(itm).Kind()
					if (kind != reflect.Map && kind != reflect.Slice) || len(ks) < 2 {
						if !m.ReadOnly {
							if o, deleted := obj.(*MapNavigatorDeleted); deleted && o != nil {
								delete(mp, itmKey)
							} else {
								mp[itmKey] = obj

							}
						}
					}
				}
				//toReturn[itmKey] = obj
				if !m.ReadOnly {
					if o, deleted := obj.(*MapNavigatorDeleted); deleted && o != nil {
					} else {
						toReturn = appendToToreturn(itm, toReturn, obj)
					}
				}
			}
		}
		return toReturn, err
	} else if k == "-" {
		if m.NodeModifier != nil {
			nmp := m.NodeModifier.ModifyNode(mp)
			return nmp, nil
		}
	} else if obj, ok := mp[k]; ok {
		//do something here
		if ttr, err := m.VisitNode(obj, ks[1:]...); err == nil && obj != nil {
			kind := reflect.TypeOf(obj).Kind()
			if (kind != reflect.Map && kind != reflect.Slice) || GetLength(ks) < 2 {
				if !m.ReadOnly {
					if o, deleted := ttr.(*MapNavigatorDeleted); deleted && o != nil {
						delete(mp, k)
					} else {
						mp[k] = ttr
					}
				}
			}
			return ttr, nil
		} else if obj == nil && ttr != nil {
			if !m.ReadOnly {
				if o, deleted := ttr.(*MapNavigatorDeleted); deleted && o != nil {
					delete(mp, k)
				} else {
					mp[k] = ttr
				}
			}
		} else {
			return nil, err
		}
	}
	if m.CreateProperty && len(ks) == 1 {
		mp[ks[0]] = ""
		return m.VisitMapStringNode(mp, ks...)
	}
	return nil, fmt.Errorf("key \"%v\" not found in map error in mapnavigator", k)
}

// GetLength returns the effective length of a key slice, accounting for wildcard suffixes.
// It subtracts 1 from the length if the last key is a wildcard (*).
func GetLength(ks []string) int {
	sub := 0
	length := len(ks)
	if length > 0 && ks[length-1] == "*" {
		sub = sub + 1
	}
	return length - sub
}

// VisitSliceNode navigates a slice using the provided key path.
// It supports numeric indices, wildcard operations (*), and modifier application (-).
// Returns the value at the specified path or an error if navigation fails.
func (m *MapNavigator) VisitSliceNode(arr []interface{}, ks ...string) (interface{}, error) {
	//if len(m.keys) == 0 {
	m.keys = ks
	//}

	k := ks[0]
	if k == "-" {
		if m.NodeModifier != nil {
			nmp := m.NodeModifier.ModifyNode(arr)

			return nmp, nil
		}
	}
	index := -1
	if i, er := strconv.Atoi(k); er == nil {
		index = i
	}

	//do something here
	if childKind := reflect.TypeOf(arr); childKind.Kind() == reflect.Slice {
		if index == -1 {
			toReturn := make([]interface{}, 0)
			var err error

			for i, itm := range arr {
				if obj, er := m.VisitNode(itm, ks[1:]...); er != nil {
					err = er
					if itm != nil {
						kind := reflect.TypeOf(itm).Kind()
						if (kind != reflect.Map && kind != reflect.Slice) || len(ks) < 2 {
							if !m.ReadOnly {
								if _, deleted := obj.(*MapNavigatorDeleted); !deleted {
									arr[i] = obj
								}
							}
						}
						if _, deleted := obj.(*MapNavigatorDeleted); !deleted {
							toReturn = appendToToreturn(kind, toReturn, obj)
						}
					}
					continue
				} else {
					kind := reflect.TypeOf(itm).Kind()
					if (kind != reflect.Map && kind != reflect.Slice) || len(ks) < 2 {
						if !m.ReadOnly {
							if _, deleted := obj.(*MapNavigatorDeleted); !deleted {
								arr[i] = obj
							}
						}
					}

					if _, deleted := obj.(*MapNavigatorDeleted); !deleted {
						toReturn = appendToToreturn(kind, toReturn, obj)
					}
				}

			}

			return toReturn, err
		} else {
			if len(arr) > index {
				result, err := m.VisitNode(arr[index], ks[1:]...)
				if err == nil && !m.ReadOnly && len(ks) == 1 {
					// Only update the array element if we're at the last key
					// (i.e., we're modifying the element itself, not a nested property)
					arr[index] = result
				}
				return result, err
			}
		}
	} else {
		return nil, fmt.Errorf("according to key, array is expected")
	}

	return nil, nil
}

func appendToToreturn(itm interface{}, toReturn []interface{}, obj interface{}) []interface{} {
	var kind reflect.Kind
	if itm != nil {
		kind = reflect.TypeOf(itm).Kind()
	} else if obj != nil {
		kind = reflect.TypeOf(obj).Kind()
	} else {
		return toReturn
	}
	if kind == reflect.Slice {
		toReturn = append(toReturn, obj.([]interface{})...)
	} else {
		toReturn = append(toReturn, obj)
	}
	return toReturn
}

// VisitValueNode handles terminal values during navigation.
// It applies modifiers if present and validates that no additional keys remain.
// Returns the (potentially modified) value or an error if keys remain.
func (m *MapNavigator) VisitValueNode(mp interface{}, ks ...string) (interface{}, error) {
	//if len(m.keys) == 0 {
	m.keys = ks
	//}
	if len(ks) > 0 {
		if m.NodeModifier != nil {
			mp = m.NodeModifier.ModifyNode(mp)
		}
		return mp, fmt.Errorf("invalid number of keys sent")
	}
	return mp, nil
}

// VisitNode is the main entry point for navigating any data structure.
// It automatically detects the type (map, slice, or value) and delegates to the appropriate visitor.
// Returns the value at the specified path or an error if navigation fails.
func (m *MapNavigator) VisitNode(mp interface{}, ks ...string) (interface{}, error) {
	//if len(m.keys) == 0 {
	m.setKeys(ks)
	//}
	if len(ks) == 0 {
		if m.NodeModifier != nil && !m.ReadOnly {
			mp = m.NodeModifier.ModifyNode(mp)
		}
		return mp, nil
	}

	if mp == nil {
		return nil, fmt.Errorf("can not use this function with nil")
	}
	//kind := reflect.TypeOf(mp).Kind()
	kind := reflect.TypeOf(mp).Kind()
	switch kind {
	case reflect.Map:
		if mps, ok := mp.(map[string]interface{}); ok {
			var tr interface{}
			var err error
			if tr, err = m.VisitMapStringNode(mps, ks...); err == nil && len(m.keys) >= 2 {
				//tr, err = m.VisitNode(tr, ks[1:]...)
			}
			return tr, err

		}
		if mpi, ok := mp.(map[interface{}]interface{}); ok {
			var tr interface{}
			var err error
			if tr, err = m.VisitMapNode(mpi, ks...); err == nil && len(m.keys) >= 2 {
				//tr, err = m.VisitNode(tr, ks[1:]...)
			}
			return tr, err
		}
	case reflect.Slice:
		toReturn, er := m.VisitSliceNode(mp.([]interface{}), ks...)

		return toReturn, er
	default:

		return m.VisitValueNode(mp, ks...)
	}
	return mp, nil
}

func (m *MapNavigator) setKeys(ks []string) {
	//fmt.Println("\n\n====---> Old======map-navigator======>(˚¥˚)<================Old==>")
	//fmt.Printf("\nNew: %v\n", ks)
	//fmt.Printf("\nOld: %v\n", m.keys)
	//fmt.Println("================>(ΩΩΩ)<==============<==END Old==")
	m.keys = ks
}

//func (m *MapNavigator) ExecuteModifier(mp interface{}, ks []string) (interface{}, error) {
//    if len(ks) == 0 {
//        if m.NodeModifier != nil {
//            mp = m.NodeModifier.ModifyNode(mp)
//        }
//        return mp, nil
//    }
//    return mp, nil
//}

// KeyConfig represents configuration for a navigation key, including array index information.
// This type is used internally for key parsing and navigation logic.
type KeyConfig struct {
	// OriginalKey is the unprocessed key string
	OriginalKey string
	// Key is the processed key string
	Key         string
	// IsArray indicates if this key represents an array access
	IsArray     bool
	// Index is the array index if IsArray is true
	Index       int
}

//func extractKey(k string) *KeyConfig {
//    kc := &KeyConfig{OriginalKey: k}
//    if newKeys := strings.Split(k, "#"); len(newKeys) > 0 {
//        kc.IsArray = true
//        k = newKeys[0]
//        if indx, er := strconv.Atoi(newKeys[1]); er == nil {
//            kc.Index = indx
//        } else {
//            kc.Index = -1
//        }
//    } else if len(newKeys) == 1 {
//        kc.Key = newKeys[0]
//        return kc
//    } else {
//        return nil
//    }
//
//
//}

// MapNavigator options
// MapNavigatorOption defines a function type for configuring MapNavigator instances.
// This allows for flexible configuration using the functional options pattern.
type MapNavigatorOption func(*MapNavigator)

// NewMapNavigator accepts a slice of option functions as the rest arguments
// NewMapNavigator creates a new MapNavigator instance with the specified node modifier.
// The modifier will be applied to nodes when the "-" key is encountered during navigation.
func NewMapNavigator(fn models.MapNodeModifier) *MapNavigator {
	obj := &MapNavigator{}
	obj.NodeModifier = fn
	return obj
}

// NewMapNavigatorFunc creates a new MapNavigator instance with a function-based modifier.
// This is a convenience function for cases where the modifier is a simple function.
func NewMapNavigatorFunc(fn models.MapNodeModifierFunc) *MapNavigator {
	obj := &MapNavigator{}
	obj.NodeModifier = models.MapNodeModifier(fn)
	return obj
}

// MapNavigatorDeleted is a marker type used to indicate that a value should be deleted.
// When a modifier returns this type, the navigator will remove the corresponding key/value.
type MapNavigatorDeleted struct {
	// Original holds the original value before deletion
	Original interface{}
}

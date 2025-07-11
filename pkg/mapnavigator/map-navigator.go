package mapnavigator

import (
	"fmt"
	"reflect"
	"strconv"

	models "github.com/passionintellectual/go-map-navigator/pkg/mapnavigator/map-nav-models"
)

type MapNavigator struct {
	NodeModifier   models.MapNodeModifier
	keys           []string
	ReadOnly       bool
	CreateProperty bool
}

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

func GetLength(ks []string) int {
	sub := 0
	length := len(ks)
	if length > 0 && ks[length-1] == "*" {
		sub = sub + 1
	}
	return length - sub
}

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
				return m.VisitNode(arr[index], ks[1:]...)
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

type KeyConfig struct {
	OriginalKey string
	Key         string
	IsArray     bool
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
type MapNavigatorOption func(*MapNavigator)

// NewMapNavigator accepts a slice of option functions as the rest arguments
func NewMapNavigator(fn models.MapNodeModifier) *MapNavigator {
	obj := &MapNavigator{}
	obj.NodeModifier = fn
	return obj
}

func NewMapNavigatorFunc(fn models.MapNodeModifierFunc) *MapNavigator {
	obj := &MapNavigator{}
	obj.NodeModifier = models.MapNodeModifier(fn)
	return obj
}

type MapNavigatorDeleted struct {
	Original interface{}
}

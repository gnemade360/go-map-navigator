package models

type MapNodeModifierFunc func(interface{}) interface{}

func (m MapNodeModifierFunc) ModifyNode(i interface{}) interface{} {
	return m(i)
}

type MapNodeModifier interface {
	ModifyNode(interface{}) interface{}
}

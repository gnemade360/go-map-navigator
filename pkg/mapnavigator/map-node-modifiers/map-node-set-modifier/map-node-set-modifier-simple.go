package setmodifier

// SimpleMapNodeSetModifier is a simplified version without external dependencies
type SimpleMapNodeSetModifier struct {
	Value interface{}
}

func (m *SimpleMapNodeSetModifier) ModifyNode(node interface{}) interface{} {
	return m.Value
}
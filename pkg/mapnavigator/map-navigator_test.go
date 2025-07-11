package mapnavigator_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/gnemade360/go-map-navigator/pkg/mapnavigator"
	models "github.com/gnemade360/go-map-navigator/pkg/mapnavigator/map-nav-models"
	"gopkg.in/yaml.v3"
)

func TestGetFromMapNavigator_VisitNode_ShouldGetIDFromMap(t *testing.T) {

	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}
	mp := GetMap()

	got, err := m.VisitNode(mp, "id")
	assert.Nil(t, err)
	assert.Equal(t, got, mp["id"])

}

func TestGetFromMapNavigator_VisitNode_ShouldGetNameFromMap(t *testing.T) {
	mm := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}

	mp := GetMap()
	got, err := mm.VisitNode(mp, "dashboardMetadata", "name")
	expected := mp["dashboardMetadata"].(map[string]interface{})["name"]
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestGetFromMapNavigator_VisitNode_ShouldGetArray(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}
	mp := GetMap()

	got, err := m.VisitNode(mp, "tiles", "*", "name")
	expected := []interface{}{}
	for _, itm := range mp["tiles"].([]interface{}) {
		name := itm.(map[string]interface{})["name"]
		expected = append(expected, name)
	}
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestGetFromMapNavigator_VisitNode_ShouldGetObjectArray(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}
	mp := GetMap()

	got, err := m.VisitNode(mp, "tiles", "*", "bounds")
	expected := []interface{}{}
	for _, itm := range mp["tiles"].([]interface{}) {
		name := itm.(map[string]interface{})["bounds"]
		expected = append(expected, name)
	}
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestGetFromMapNavigator_VisitNode_ShouldGetIndexBasedItemOfArray(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}
	mp := GetMap()

	for i, itm := range mp["tiles"].([]interface{}) {
		name := itm.(map[string]interface{})["name"]
		got, err := m.VisitNode(mp["tiles"], strconv.Itoa(i), "name")

		assert.Nil(t, err)
		assert.Equal(t, got, name)
	}

}
func TestGetFromMapNavigator_VisitNode_ShouldGetAllItemsOfArray(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(o interface{}) interface{} {
			return o
		}),
		ReadOnly: true,
	}
	mp := GetMap()
	got, err := m.VisitNode(mp["tiles"], "*", "name")
	expected := []interface{}{}
	for _, itm := range mp["tiles"].([]interface{}) {
		name := itm.(map[string]interface{})["name"]
		expected = append(expected, name)
	}
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestModifyMapNavigator_VisitNode_ShouldModifyMapValue(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			v = v.(string) + "----END"
			return v
		}),
	}
	mp := GetMap()
	id := mp["id"]
	got, err := m.VisitNode(mp, "id")
	assert.Equal(t, got, fmt.Sprintf("%v----END", id))
	assert.Equal(t, mp["id"], got)
	assert.Nil(t, err)
}

func TestModifyMapNavigator_VisitNode_ShouldModifyNestedValue(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			if v == "Markdown" {
				v = v.(string) + "---123"
			}
			return v
		}),
	}
	mp := GetMap()
	expected := []string{}
	for _, itm := range mp["tiles"].([]interface{}) {
		name := itm.(map[string]interface{})["name"]
		expected = append(expected, name.(string))
	}
	_, err := m.VisitNode(mp, "tiles", "*", "name")

	for i, itm := range mp["tiles"].([]interface{}) {
		name := itm.(map[string]interface{})["name"]
		assert.Equal(t, name, m.NodeModifier.ModifyNode(expected[i]))
	}

	assert.Nil(t, err)
}

func TestModifyMapNavigator_VisitNode_ShouldModifyMap(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			v.(map[string]interface{})["top"] = "9000"
			return v
		}),
		ReadOnly: true,
	}
	mp := GetMap()
	expected := []map[string]interface{}{}
	for _, itm := range mp["tiles"].([]interface{}) {
		bounds := itm.(map[string]interface{})["bounds"]
		expected = append(expected, bounds.(map[string]interface{}))
	}
	_, err := m.VisitNode(mp, "tiles", "*", "bounds")

	for i, itm := range mp["tiles"].([]interface{}) {
		bounds := itm.(map[string]interface{})["bounds"]
		//fmt.Println(bounds.(map[string]interface{})["top"])
		assert.Equal(t, bounds, m.NodeModifier.ModifyNode(expected[i]))
	}

	assert.Nil(t, err)
}

func TestModifyMapNavigator_VisitNode_ShouldNotModifyIndexBasedSliceItemsIfReadonly(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			v.(map[string]interface{})["top"] = 9000
			return v
		}),
		ReadOnly: true,
	}
	mp := GetMap()

	received, err := m.VisitNode(mp, "tiles", "2", "bounds")
	assert.Equal(t, received.(map[string]interface{})["top"], float64(646))

	assert.Equal(t, mp["tiles"].([]interface{})[2].(map[string]interface{})["bounds"].(map[string]interface{})["top"], float64(646))

	assert.Nil(t, err)
}

func TestModifyMapNavigator_VisitNode_ShouldModifyIndexBasedSliceItems(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			if vMap, ok := v.(map[string]interface{}); ok {
				// Only modify if it's the bounds object (has "top" field)
				if _, hasTop := vMap["top"]; hasTop {
					vMap["top"] = 9000
					return vMap
				}
			}
			return v
		}),
	}
	mp := GetMap()

	_, err := m.VisitNode(mp, "tiles", "2", "bounds")

	assert.Nil(t, err)
	assert.Equal(t, mp["tiles"].([]interface{})[2].(map[string]interface{})["bounds"].(map[string]interface{})["top"], 9000)
	assert.NotEqual(t, mp["tiles"].([]interface{})[1].(map[string]interface{})["bounds"].(map[string]interface{})["top"], 9000)
}

func TestModifyMapNavigator_VisitNode_ShouldModifyMapString(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			v = fmt.Sprintf("%v---->Modified", v)
			return v
		}),
	}
	mp := GetMap()
	owner := mp["dashboardMetadata"].(map[string]interface{})["owner"]
	_, err := m.VisitNode(mp, "dashboardMetadata", "owner")

	assert.Equal(t, mp["dashboardMetadata"].(map[string]interface{})["owner"], m.NodeModifier.ModifyNode(owner))

	assert.Nil(t, err)
}

func TestModifyMapNavigator_VisitNode_ShouldModifyMapSlice(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			tags := v.([]interface{})
			newTags := make([]interface{}, 0)
			for _, tag := range tags {

				newTags = append(newTags, fmt.Sprintf("%v--->Modified", tag))
			}
			return newTags
		}),
	}
	mp := GetMap()
	tags := mp["dashboardMetadata"].(map[string]interface{})["tags"]
	_, err := m.VisitNode(mp, "dashboardMetadata", "tags")

	assert.Equal(t, mp["dashboardMetadata"].(map[string]interface{})["tags"], m.NodeModifier.ModifyNode(tags))

	assert.Nil(t, err)
}
func TestModifyMapNavigator_VisitNode_ShouldModifySliceInSlice(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return fmt.Sprintf("%v--->Modified", v)
		}),
	}
	mp := GetMap()

	_, err := m.VisitNode(mp, "testArr", "*", "*")

	assert.Equal(t, mp["testArr"].([]interface{})[0].([]interface{})[0], m.NodeModifier.ModifyNode("t"))
	assert.Equal(t, mp["testArr"].([]interface{})[0].([]interface{})[1], m.NodeModifier.ModifyNode("u"))

	assert.Equal(t, mp["testArr"].([]interface{})[1].([]interface{})[0], m.NodeModifier.ModifyNode("a"))
	assert.Equal(t, mp["testArr"].([]interface{})[1].([]interface{})[1], m.NodeModifier.ModifyNode("b"))

	assert.Nil(t, err)
}

func TestModifyMapNavigator_VisitNode_ShouldModifySliceOfObjInSlice(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return fmt.Sprintf("%v--->Modified", v)
		}),
	}
	mp := GetMap()

	_, err := m.VisitNode(mp, "testObjArray", "*", "*", "name")

	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[0].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("0-0"))
	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[1].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("0-1"))

	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[0].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("1-0"))
	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[1].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("1-1"))

	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[0].(map[string]interface{})["city"], "city-0-0")
	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[1].(map[string]interface{})["city"], "city-0-1")

	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[0].(map[string]interface{})["city"], "city-1-0")
	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[1].(map[string]interface{})["city"], "city-1-1")

	assert.Nil(t, err)
}

func TestModifyMapNavigator_VisitNode_ShouldModifySliceOfObjInSliceForAll(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return fmt.Sprintf("%v--->Modified", v)
		}),
	}
	mp := GetMap()

	_, err := m.VisitNode(mp, "testObjArray", "*", "*", "*")

	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[0].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("0-0"))
	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[1].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("0-1"))

	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[0].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("1-0"))
	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[1].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("1-1"))

	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[0].(map[string]interface{})["city"], m.NodeModifier.ModifyNode("city-0-0"))
	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[1].(map[string]interface{})["city"], m.NodeModifier.ModifyNode("city-0-1"))

	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[0].(map[string]interface{})["city"], m.NodeModifier.ModifyNode("city-1-0"))
	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[1].(map[string]interface{})["city"], m.NodeModifier.ModifyNode("city-1-1"))

	assert.Nil(t, err)
}

func TestShouldReturnArray(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: nil,
	}
	mp := GetMap("test-data/summarytest.json")
	//itms := mp["dashboards"].([]interface{})
	//itms = itms[:1]
	//mp["dashboards"] = itms
	arr, err := m.VisitNode(mp, "*", "*")
	assert.Nil(t, err)
	assert.Equal(t, len(arr.([]interface{})), 2)
}

func TestShouldReturnError(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return fmt.Sprintf("%v--->Modified", v)
		}),
	}
	mp := GetMap("test-data/summarytest.json")

	_, err := m.VisitNode(mp, "*", "*", "*", "*", "*", "*", "*")
	//arr, err := m.VisitNode(mp, "*")
	assert.NotNil(t, err)

}
func TestShouldReturnTags(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return fmt.Sprintf("%v--->Modified", v)
			//return v.(map[string]interface{})["id"]
		}),
	}
	mp := GetMap("test-data/summarytest.json")

	a, err := m.VisitNode(mp, "*", "*", "tags", "*")
	//arr, err := m.VisitNode(mp, "*")
	assert.Nil(t, err)
	assert.Equal(t, len(a.([]interface{})), 2)
	assert.Equal(t, len(a.([]interface{})[0].([]interface{})), 3)
	assert.Equal(t, len(a.([]interface{})[1].([]interface{})), 3)
}

func TestShouldReturnArrayAndVisitAll(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return fmt.Sprintf("%v--->Modified", v)
		}),
	}
	mp := GetMap("test-data/summarytest.json")
	itms := mp["dashboards"].([]interface{})
	itms = itms[:1]
	mp["dashboards"] = itms
	arr, err := m.VisitNode(mp, "*", "*", "id")
	//arr, err := m.VisitNode(mp, "*")
	assert.Nil(t, err)
	assert.Equal(t, len(arr.([]interface{})), len(itms))
	assert.Equal(t, arr.([]interface{})[0], "1ffb752a-0818-4508-8c9f-b7de10f2f3c7--->Modified")
}

func TestModifyMapNavigator_VisitNodeYamlMap_ShouldModifySliceOfObjInSliceForAll(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return fmt.Sprintf("%v--->Modified", v)
		}),
	}
	mp := GetYamlMap()

	_, err := m.VisitNode(mp, "testObjArray", "*", "*", "*")

	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[0].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("0-0"))
	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[1].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("0-1"))

	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[0].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("1-0"))
	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[1].(map[string]interface{})["name"], m.NodeModifier.ModifyNode("1-1"))

	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[0].(map[string]interface{})["city"], m.NodeModifier.ModifyNode("city-0-0"))
	assert.Equal(t, mp["testObjArray"].([]interface{})[0].([]interface{})[1].(map[string]interface{})["city"], m.NodeModifier.ModifyNode("city-0-1"))

	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[0].(map[string]interface{})["city"], m.NodeModifier.ModifyNode("city-1-0"))
	assert.Equal(t, mp["testObjArray"].([]interface{})[1].([]interface{})[1].(map[string]interface{})["city"], m.NodeModifier.ModifyNode("city-1-1"))

	assert.Nil(t, err)
}

func TestModifyMapNavigator_VisitNodeYamlMap_ShouldModifyWhenParentObjectKey(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			if vmap, isMap := v.(map[string]interface{}); isMap && vmap != nil {
				if url, ok := vmap["url"]; ok {
					//do something here
					vmap["url"] = fmt.Sprintf("modified--%v", url)
				}
			}
			return v
		}),
	}
	mp := GetAppDetectionRuleSampleJSON()

	_, err := m.VisitNode(mp, "-")
	assert.Equal(t, mp["url"], "modified--<nil>")

	assert.Nil(t, err)
}
func TestModifyMapNavigator_VisitNode_ShouldDeleteSliceOfObjInSliceForAll(t *testing.T) {

	testFn := func(t *testing.T, expectedCity string) {
		m := &MapNavigator{
			NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if vmap, isMap := v.(map[string]interface{}); isMap && vmap != nil {
					if vmap["name"] == expectedCity {
						return &MapNavigatorDeleted{Original: v}
					}
				}
				return v
			}),
		}
		mp := GetMap()
		assert.Equal(t, len(mp["shouldDeleteItemArray"].([]interface{})), 3)

		_, err := m.VisitNode(mp, "shouldDeleteItemArray", "*")
		//shouldDeleteItemArray
		assert.Equal(t, len(mp["shouldDeleteItemArray"].([]interface{})), 2)
		for _, itm := range mp["shouldDeleteItemArray"].([]interface{}) {
			assert.NotEqual(t, itm.(map[string]interface{})["city"], expectedCity)
		}
		assert.Nil(t, err)
	}
	testFn(t, "0-2")
	testFn(t, "0-1")
	testFn(t, "0-0")

}

func TestModifyMapNavigator_VisitNode_ShouldDeleteSliceOfObjInObjectForAll(t *testing.T) {

	testFn := func(t *testing.T, fn func(interface{}) bool) {
		m := &MapNavigator{
			NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if vmap, isMap := v.(map[string]interface{}); isMap && vmap != nil {
					if fn(vmap) {
						return &MapNavigatorDeleted{Original: v}
					}
				}
				return v
			}),
		}
		mp := GetMap()
		originalLength := len(mp["tiles"].([]interface{}))
		originalCount := 0
		for _, itm := range mp["tiles"].([]interface{}) {
			if fn(itm) {
				originalCount++
			}
		}

		_, err := m.VisitNode(mp, "tiles", "*")
		assert.Equal(t, len(mp["tiles"].([]interface{})), originalLength-originalCount)

		for _, itm := range mp["tiles"].([]interface{}) {
			result := fn(itm)
			assert.False(t, result)
		}
		assert.Nil(t, err)
	}
	testFn(t, func(i interface{}) bool {
		if imap, isMap := i.(map[string]interface{}); isMap && imap != nil {
			return imap["name"] == "Markdown"
		}
		return false
	})

}

func TestModifyMapNavigator_VisitNode_ShouldDeleteWithinObject(t *testing.T) {

	testFn := func(t *testing.T, fn func(interface{}) bool) {
		m := &MapNavigator{
			NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
				if fn(v) {
					return &MapNavigatorDeleted{Original: v}
				}
				return v
			}),
		}
		mp := GetMap()

		_, err := m.VisitNode(mp, "tiles", "1", "*")
		_, ok := mp["tiles"].([]interface{})[1].(map[string]interface{})["name"]
		assert.False(t, ok)
		assert.Nil(t, err)
	}
	testFn(t, func(i interface{}) bool {
		if imap, isMap := i.(string); isMap && len(imap) > 0 {
			return imap == "Markdown"
		}
		return false
	})

}

func TestModifyMapNavigator_VisitNode_ShouldDeleteWithinObject1(t *testing.T) {

	testFn := func(t *testing.T) {
		m := &MapNavigator{
			NodeModifier: models.MapNodeModifierFunc(func(i interface{}) interface{} {
				if imap, isMap := i.(map[string]interface{}); isMap && len(imap) > 0 {
					if _, ok := imap["bounds"]; ok {
						//do something here
						//imap["bounds"] = &MapNavigatorDeleted{Original: imap["bounds"]}
						delete(imap, "bounds")
						return imap
					}
				}
				return i

			}),
		}
		mp := GetMap()

		_, err := m.VisitNode(mp["tiles"].([]interface{}), "*")
		_, ok := mp["tiles"].([]interface{})[0].(map[string]interface{})["bounds"]
		assert.False(t, ok)
		assert.Nil(t, err)
	}
	testFn(t)

}
func GetMap(f ...string) map[string]interface{} {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	s := "test-data/test.json"
	if len(f) > 0 && len(f[0]) > 0 {
		s = f[0]
	}
	path = filepath.Join(path, s)
	data, err := ioutil.ReadFile(path)
	mp := make(map[string]interface{})
	var _ = json.Unmarshal(data, &mp)

	return mp

}
func GetAppDetectionRuleSampleJSON() map[string]interface{} {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path = filepath.Join(path, "test-data/appdetectionrule.json")
	data, err := ioutil.ReadFile(path)
	mp := make(map[string]interface{})
	var _ = json.Unmarshal(data, &mp)
	return mp
}
func GetYamlMap() map[interface{}]interface{} {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	path = filepath.Join(path, "test-data/test.yaml")
	data, err := ioutil.ReadFile(path)
	mp := make(map[interface{}]interface{})
	var _ = yaml.Unmarshal(data, &mp)

	return mp
}

// Additional tests for edge cases and untested scenarios

func TestMapNavigator_CreatePropertyIfAbsent(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return "modified-" + v.(string)
		}),
		CreateProperty: true,
	}
	
	// Test creating property in map[string]interface{}
	t.Run("create property in string map", func(t *testing.T) {
		mp := map[string]interface{}{
			"existing": "value",
		}
		
		result, err := m.VisitNode(mp, "newField")
		assert.Nil(t, err)
		assert.Equal(t, "modified-", result)
		assert.Equal(t, "modified-", mp["newField"])
	})
	
	// Test creating property in map[interface{}]interface{}
	t.Run("create property in interface map", func(t *testing.T) {
		mp := map[interface{}]interface{}{
			"existing": "value",
		}
		
		result, err := m.VisitNode(mp, "newField")
		assert.Nil(t, err)
		assert.Equal(t, "modified-", result)
		assert.Equal(t, "modified-", mp["newField"])
	})
	
	// Test creating nested property disabled
	t.Run("create nested property disabled", func(t *testing.T) {
		// Create a new MapNavigator with CreateProperty disabled
		mDisabled := &MapNavigator{
			NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
				return "modified-" + v.(string)
			}),
			CreateProperty: false,
		}
		
		mp := map[string]interface{}{
			"level1": map[string]interface{}{},
		}
		
		_, err := mDisabled.VisitNode(mp, "level1", "newField")
		assert.NotNil(t, err)
		if err != nil {
			assert.Contains(t, err.Error(), "not found")
		}
	})
}

func TestMapNavigator_NilHandling(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return v
		}),
		ReadOnly: true,
	}
	
	t.Run("nil map navigation", func(t *testing.T) {
		_, err := m.VisitNode(nil, "key")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "can not use this function with nil")
	})
	
	t.Run("nil map string node", func(t *testing.T) {
		var mp map[string]interface{} = nil
		_, err := m.VisitMapStringNode(mp, "key")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "map is nil")
	})
	
	t.Run("nil map interface node", func(t *testing.T) {
		var mp map[interface{}]interface{} = nil
		_, err := m.VisitMapNode(mp, "key")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "map is nil")
	})
}

func TestMapNavigator_EmptyKeyHandling(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return v
		}),
	}
	
	t.Run("empty key in map string node", func(t *testing.T) {
		mp := map[string]interface{}{"field": "value"}
		_, err := m.VisitMapStringNode(mp, "")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "map is nil")
	})
	
	t.Run("empty key in map interface node", func(t *testing.T) {
		mp := map[interface{}]interface{}{"field": "value"}
		_, err := m.VisitMapNode(mp, "")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "map is nil")
	})
}

func TestMapNavigator_SpecialCharacterKeys(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return v.(string) + "-modified"
		}),
	}
	
	specialKeys := []string{
		"key.with.dots",
		"key-with-dashes",
		"key_with_underscores",
		"key with spaces",
		"key@with#special$chars",
		"123numeric",
	}
	
	for _, key := range specialKeys {
		t.Run(fmt.Sprintf("special key: %s", key), func(t *testing.T) {
			mp := map[string]interface{}{
				key: "value",
			}
			
			result, err := m.VisitNode(mp, key)
			assert.Nil(t, err)
			assert.Equal(t, "value-modified", result)
			assert.Equal(t, "value-modified", mp[key])
		})
	}
}

func TestMapNavigator_LargeNestedStructure(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return "found"
		}),
		ReadOnly: true,
	}
	
	// Create deeply nested structure
	nested := map[string]interface{}{"value": "target"}
	for i := 0; i < 10; i++ {
		nested = map[string]interface{}{"level": nested}
	}
	
	// Navigate to deepest level
	keys := make([]string, 11)
	for i := 0; i < 10; i++ {
		keys[i] = "level"
	}
	keys[10] = "value"
	
	result, err := m.VisitNode(nested, keys...)
	assert.Nil(t, err)
	assert.Equal(t, "target", result) // ReadOnly is true, so modifier should not be applied
}

func TestMapNavigator_NewMapNavigatorFunc(t *testing.T) {
	modifierFunc := models.MapNodeModifierFunc(func(v interface{}) interface{} {
		return "modified"
	})
	
	navigator := NewMapNavigatorFunc(modifierFunc)
	assert.NotNil(t, navigator)
	assert.NotNil(t, navigator.NodeModifier)
	
	// Test that the modifier works
	result := navigator.NodeModifier.ModifyNode("test")
	assert.Equal(t, "modified", result)
}

func TestMapNavigator_GetLength(t *testing.T) {
	tests := []struct {
		name     string
		keys     []string
		expected int
	}{
		{
			name:     "empty keys",
			keys:     []string{},
			expected: 0,
		},
		{
			name:     "single key",
			keys:     []string{"key"},
			expected: 1,
		},
		{
			name:     "multiple keys",
			keys:     []string{"key1", "key2", "key3"},
			expected: 3,
		},
		{
			name:     "keys ending with wildcard",
			keys:     []string{"key1", "key2", "*"},
			expected: 2,
		},
		{
			name:     "only wildcard",
			keys:     []string{"*"},
			expected: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLength(tt.keys)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapNavigator_WildcardOnNonSlice(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return v
		}),
		ReadOnly: true,
	}
	
	// Test wildcard on a string value (not a slice)
	mp := map[string]interface{}{
		"field": "string-value",
	}
	
	_, err := m.VisitNode(mp, "field", "*")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "invalid number of keys sent")
}

func TestMapNavigator_ModifySliceWithNonNumericIndex(t *testing.T) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return v
		}),
	}
	
	arr := []interface{}{"item1", "item2", "item3"}
	
	// Test with non-numeric index that's not "*"
	_, err := m.VisitSliceNode(arr, "invalid")
	assert.Nil(t, err) // Should handle gracefully by treating as wildcard
}

func TestMapNavigator_DashModifier(t *testing.T) {
	modifiedValue := "dash-modified"
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return modifiedValue
		}),
	}
	
	t.Run("dash modifier on map", func(t *testing.T) {
		mp := map[string]interface{}{
			"field": "value",
		}
		
		result, err := m.VisitMapStringNode(mp, "-")
		assert.Nil(t, err)
		assert.Equal(t, modifiedValue, result)
	})
	
	t.Run("dash modifier on interface map", func(t *testing.T) {
		mp := map[interface{}]interface{}{
			"field": "value",
		}
		
		result, err := m.VisitMapNode(mp, "-")
		assert.Nil(t, err)
		assert.Equal(t, modifiedValue, result)
	})
	
	t.Run("dash modifier on slice", func(t *testing.T) {
		arr := []interface{}{"item1", "item2"}
		
		result, err := m.VisitSliceNode(arr, "-")
		assert.Nil(t, err)
		assert.Equal(t, modifiedValue, result)
	})
}

func BenchmarkMapNavigator_VisitNode(b *testing.B) {
	m := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			return v
		}),
		ReadOnly: true,
	}
	
	mp := GetMap()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.VisitNode(mp, "dashboardMetadata", "tags", "*")
	}
}

func ExampleMapNavigator_VisitNode() {
	// Create a navigator with a custom modifier
	navigator := &MapNavigator{
		NodeModifier: models.MapNodeModifierFunc(func(v interface{}) interface{} {
			if s, ok := v.(string); ok {
				return "[MODIFIED] " + s
			}
			return v
		}),
	}
	
	// Create sample data
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"tags": []interface{}{"admin", "user"},
		},
	}
	
	// Navigate and modify
	result, _ := navigator.VisitNode(data, "user", "name")
	fmt.Println(result) // Output: [MODIFIED] John
}
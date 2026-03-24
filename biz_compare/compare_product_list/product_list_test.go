package compare_product_list

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// DiffResult represents a difference result with product ID
type DiffResult struct {
	ProductID string
	Field     string
	Diff      string
}

// ResponseJSON represents the JSON structure as generic map
type ResponseJSON map[string]interface{}

// sortProductListByID sorts product_list by product_id numerically
func sortProductListByID(data map[string]interface{}) {
	if data == nil {
		return
	}

	productList, ok := data["product_list"].([]interface{})
	if !ok {
		return
	}

	sort.Slice(productList, func(i, j int) bool {
		itemI, okI := productList[i].(map[string]interface{})
		itemJ, okJ := productList[j].(map[string]interface{})
		if !okI || !okJ {
			return false
		}

		productInfoI, okI := itemI["product_info"].(map[string]interface{})
		productInfoJ, okJ := itemJ["product_info"].(map[string]interface{})
		if !okI || !okJ {
			return false
		}

		idStrI, _ := productInfoI["product_id"].(string)
		idStrJ, _ := productInfoJ["product_id"].(string)

		idI, _ := strconv.ParseInt(idStrI, 10, 64)
		idJ, _ := strconv.ParseInt(idStrJ, 10, 64)

		return idI < idJ
	})

	data["product_list"] = productList
}

// shouldIgnoreField checks if the field should be ignored
func shouldIgnoreField(path string) bool {
	// Fields to ignore during comparison
	ignoredFields := []string{
		"request_id",
		"server_time",
		"latency_info",
	}

	for _, field := range ignoredFields {
		if path == field || strings.HasSuffix(path, "."+field) {
			return true
		}
	}
	return false
}

// getProductIDFromPath extracts product_id from the path
func getProductIDFromPath(path string) string {
	// Look for product_info.product_id in the path
	parts := strings.Split(path, ".")
	for i, part := range parts {
		if part == "product_info" && i+1 < len(parts) {
			// Check if next part is product_id
			if i+2 < len(parts) && parts[i+1] == "product_id" {
				// Try to get value from path - we'll handle this differently
			}
		}
	}
	return ""
}

// deepCompare compares two values deeply and returns differences
func deepCompare(path, currentProductID string, v1, v2 reflect.Value) []string {
	var diffs []string

	// Skip ignored fields
	if shouldIgnoreField(path) {
		return diffs
	}

	// Handle nil values
	if !v1.IsValid() && !v2.IsValid() {
		return diffs
	}
	if !v1.IsValid() {
		diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s, 差异: 只在第一个文件中存在", currentProductID, path))
		return diffs
	}
	if !v2.IsValid() {
		diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s, 差异: 只在第二个文件中存在", currentProductID, path))
		return diffs
	}

	// If types are different
	if v1.Type() != v2.Type() {
		diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s, 差异: 类型不一致 - 类型1: %s, 类型2: %s", currentProductID, path, v1.Type(), v2.Type()))
		return diffs
	}

	// If v1 is interface, get the underlying value
	if v1.Kind() == reflect.Interface && !v1.IsNil() {
		v1 = v1.Elem()
	}
	if v2.Kind() == reflect.Interface && !v2.IsNil() {
		v2 = v2.Elem()
	}

	// Recheck type after unwrapping interface
	if v1.Type() != v2.Type() {
		diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s, 差异: 类型不一致 - 类型1: %s, 类型2: %s", currentProductID, path, v1.Type(), v2.Type()))
		return diffs
	}

	switch v1.Kind() {
	case reflect.Map:
		keys1 := v1.MapKeys()
		keys2 := v2.MapKeys()
		keySet1 := make(map[string]bool)
		keySet2 := make(map[string]bool)

		for _, key := range keys1 {
			keySet1[key.String()] = true
		}
		for _, key := range keys2 {
			keySet2[key.String()] = true
		}

		// Check keys only in v1
		for _, key := range keys1 {
			if !keySet2[key.String()] {
				diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s.%s, 差异: 只在第一个文件中存在", currentProductID, path, key.String()))
			}
		}
		// Check keys only in v2
		for _, key := range keys2 {
			if !keySet1[key.String()] {
				diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s.%s, 差异: 只在第二个文件中存在", currentProductID, path, key.String()))
			}
		}

		// Compare common keys
		for _, key := range keys1 {
			if keySet2[key.String()] {
				newPath := path
				if newPath == "" {
					newPath = key.String()
				} else {
					newPath = path + "." + key.String()
				}

				// Track product_id
				newProductID := currentProductID

				// Get the value
				v1Val := v1.MapIndex(key)

				// Unwrap interface if needed
				if v1Val.Kind() == reflect.Interface && !v1Val.IsNil() {
					v1Val = v1Val.Elem()
				}

				// Try to get product_id from various levels
				if v1Val.IsValid() && v1Val.Kind() == reflect.Map {
					idKey := reflect.ValueOf("product_id")
					if idVal := v1Val.MapIndex(idKey); idVal.IsValid() {
						// Unwrap interface
						if idVal.Kind() == reflect.Interface && !idVal.IsNil() {
							idVal = idVal.Elem()
						}
						if idStr, ok := idVal.Interface().(string); ok {
							newProductID = idStr
						}
					}
				}

				diffs = append(diffs, deepCompare(newPath, newProductID, v1.MapIndex(key), v2.MapIndex(key))...)
			}
		}

	case reflect.Slice:
		if v1.Len() != v2.Len() {
			diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s, 差异: 数组长度不一致 - 长度1: %d, 长度2: %d", currentProductID, path, v1.Len(), v2.Len()))
			return diffs
		}
		for i := 0; i < v1.Len(); i++ {
			newPath := fmt.Sprintf("%s[%d]", path, i)

			// Try to get product_id for this slice element
			newProductID := currentProductID
			if strings.HasSuffix(path, ".product_list") {
				elem := v1.Index(i)
				// Unwrap interface if needed
				if elem.Kind() == reflect.Interface && !elem.IsNil() {
					elem = elem.Elem()
				}
				if elem.IsValid() && elem.Kind() == reflect.Map {
					// First check if there's product_info key
					productInfoKey := reflect.ValueOf("product_info")
					if piVal := elem.MapIndex(productInfoKey); piVal.IsValid() {
						if piVal.Kind() == reflect.Interface && !piVal.IsNil() {
							piVal = piVal.Elem()
						}
						if piVal.IsValid() && piVal.Kind() == reflect.Map {
							idKey := reflect.ValueOf("product_id")
							if idVal := piVal.MapIndex(idKey); idVal.IsValid() {
								if idVal.Kind() == reflect.Interface && !idVal.IsNil() {
									idVal = idVal.Elem()
								}
								if idStr, ok := idVal.Interface().(string); ok {
									newProductID = idStr
								}
							}
						}
					}
				}
			}

			diffs = append(diffs, deepCompare(newPath, newProductID, v1.Index(i), v2.Index(i))...)
		}

	case reflect.Struct:
		for i := 0; i < v1.NumField(); i++ {
			field := v1.Type().Field(i)
			newPath := path
			if newPath == "" {
				newPath = field.Name
			} else {
				newPath = path + "." + field.Name
			}
			diffs = append(diffs, deepCompare(newPath, currentProductID, v1.Field(i), v2.Field(i))...)
		}

	case reflect.String:
		if v1.String() != v2.String() {
			diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s, 差异: 值不一致 - 值1: \"%s\", 值2: \"%s\"", currentProductID, path, v1.String(), v2.String()))
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v1.Int() != v2.Int() {
			diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s, 差异: 值不一致 - 值1: %d, 值2: %d", currentProductID, path, v1.Int(), v2.Int()))
		}

	case reflect.Float32, reflect.Float64:
		if v1.Float() != v2.Float() {
			diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s, 差异: 值不一致 - 值1: %f, 值2: %f", currentProductID, path, v1.Float(), v2.Float()))
		}

	case reflect.Bool:
		if v1.Bool() != v2.Bool() {
			diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s, 差异: 值不一致 - 值1: %v, 值2: %v", currentProductID, path, v1.Bool(), v2.Bool()))
		}

	default:
		if !reflect.DeepEqual(v1.Interface(), v2.Interface()) {
			diffs = append(diffs, fmt.Sprintf("商品ID: %s, 字段: %s, 差异: 值不一致 - 值1: %v, 值2: %v", currentProductID, path, v1.Interface(), v2.Interface()))
		}
	}

	return diffs
}

// compareJSONFiles compares two JSON files and returns whether they are the same, differences, and error
func compareJSONFiles(file1Path, file2Path string) (bool, []string, error) {
	// Read files
	data1, err := os.ReadFile(file1Path)
	if err != nil {
		return false, nil, fmt.Errorf("读取文件1失败: %v", err)
	}
	data2, err := os.ReadFile(file2Path)
	if err != nil {
		return false, nil, fmt.Errorf("读取文件2失败: %v", err)
	}

	// Parse JSON as generic map to handle all fields
	var json1, json2 map[string]interface{}
	if err := json.Unmarshal(data1, &json1); err != nil {
		return false, nil, fmt.Errorf("解析文件1 JSON失败: %v", err)
	}
	if err := json.Unmarshal(data2, &json2); err != nil {
		return false, nil, fmt.Errorf("解析文件2 JSON失败: %v", err)
	}

	// Sort product_list by product_id numerically
	sortProductListByID(json1["data"].(map[string]interface{}))
	sortProductListByID(json2["data"].(map[string]interface{}))

	// Deep compare using reflection
	v1 := reflect.ValueOf(json1)
	v2 := reflect.ValueOf(json2)

	diffs := deepCompare("", "", v1, v2)

	if len(diffs) == 0 {
		return true, nil, nil
	}
	return false, diffs, nil
}

// TestCompareJSONFiles is the unit test for comparing JSON files
func TestCompareJSONFiles(t *testing.T) {
	// Get current directory
	dir := "./"

	// Compare optimize.json and origin.json
	same, diffs, err := compareJSONFiles(dir+"optimize.json", dir+"origin.json")
	if err != nil {
		t.Fatalf("比较文件失败: %v", err)
	}

	if same {
		fmt.Println("文件内容相同")
	} else {
		fmt.Println("文件内容存在差异:")
		for _, diff := range diffs {
			fmt.Println(diff)
		}
	}
}

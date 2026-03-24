package compare_product_meta

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"
)

// ProductInfo 商品信息结构
type ProductInfo struct {
	ProductID      string      `json:"product_id"`
	Name           []NameItem  `json:"name"`
	CategoryID     string      `json:"category_id"`
	LimitNum       int         `json:"limit_num"`
	LimitType      int         `json:"limit_type"`
	PromotionID    int         `json:"promotion_id"`
	PropertyList   []Property  `json:"property_list"`
	RechargePoints int        `json:"recharge_points"`
	Price          Price       `json:"price"`
	MultiOption    []Price     `json:"multi_option"`
	PayType        int         `json:"pay_type"`
}

// NameItem 名称项
type NameItem struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

// Property 属性
type Property struct {
	ID         int       `json:"id"`
	Num        int       `json:"num"`
	Name       []NameItem `json:"name"`
	CategoryID int       `json:"category_id"`
	TypeKey    string   `json:"type_key"`
	Gift       *Gift    `json:"gift"`
}

// Gift 赠品
type Gift struct {
	Type int `json:"type"`
}

// Price 价格
type Price struct {
	PlatChannelProductIDMap map[string]string `json:"plat_channel_product_id_map"`
	AreaCurrencyMap         map[string]Currency `json:"area_currency_map"`
	PayLevelCode            string              `json:"pay_level_code"`
	PayLevelID              int                 `json:"pay_level_id"`
	Num                     int                 `json:"num"`
}

// Currency 货币
type Currency struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

// Data 数据结构
type Data struct {
	ProductInfo []ProductInfo `json:"product_info"`
}

// Response 响应结构
type Response struct {
	Code      int    `json:"code"`
	Info      string `json:"info"`
	RequestID string `json:"request_id"`
	Data      Data   `json:"data"`
}

// DiffItem 差异项
type DiffItem struct {
	ProductID string
	FieldName string
	Diff      string
}

// 读取并解析JSON文件
func readJSONFile(filepath string) (*Response, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// 按product_id数值排序
func sortByProductID(products []ProductInfo) {
	sort.Slice(products, func(i, j int) bool {
		idI := 0
		idJ := 0
		fmt.Sscanf(products[i].ProductID, "%d", &idI)
		fmt.Sscanf(products[j].ProductID, "%d", &idJ)
		return idI < idJ
	})
}

// 比较两个NameItem数组
func compareNameItems(fieldName string, p1, p2 []NameItem, diffs []DiffItem, productID string) []DiffItem {
	if len(p1) != len(p2) {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName, Diff: fmt.Sprintf("数组长度不同: %d vs %d", len(p1), len(p2))})
		return diffs
	}
	for i := range p1 {
		if p1[i].Lang != p2[i].Lang {
			diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fmt.Sprintf("%s[%d].lang", fieldName, i), Diff: fmt.Sprintf("%s vs %s", p1[i].Lang, p2[i].Lang)})
		}
		if p1[i].Text != p2[i].Text {
			diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fmt.Sprintf("%s[%d].text", fieldName, i), Diff: fmt.Sprintf("%s vs %s", p1[i].Text, p2[i].Text)})
		}
	}
	return diffs
}

// 比较两个Gift
func compareGifts(fieldName string, g1, g2 *Gift, diffs []DiffItem, productID string) []DiffItem {
	if g1 == nil && g2 == nil {
		return diffs
	}
	if g1 == nil || g2 == nil {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName, Diff: fmt.Sprintf("gift: nil vs %v", g2)})
		return diffs
	}
	if g1.Type != g2.Type {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".type", Diff: fmt.Sprintf("%d vs %d", g1.Type, g2.Type)})
	}
	return diffs
}

// 比较两个Property
func compareProperties(fieldName string, p1, p2 Property, diffs []DiffItem, productID string) []DiffItem {
	if p1.ID != p2.ID {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".id", Diff: fmt.Sprintf("%d vs %d", p1.ID, p2.ID)})
	}
	if p1.Num != p2.Num {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".num", Diff: fmt.Sprintf("%d vs %d", p1.Num, p2.Num)})
	}
	if p1.CategoryID != p2.CategoryID {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".category_id", Diff: fmt.Sprintf("%d vs %d", p1.CategoryID, p2.CategoryID)})
	}
	if p1.TypeKey != p2.TypeKey {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".type_key", Diff: fmt.Sprintf("%s vs %s", p1.TypeKey, p2.TypeKey)})
	}
	diffs = compareNameItems(fieldName+".name", p1.Name, p2.Name, diffs, productID)
	diffs = compareGifts(fieldName+".gift", p1.Gift, p2.Gift, diffs, productID)
	return diffs
}

// 比较两个Price
func comparePrices(fieldName string, p1, p2 Price, diffs []DiffItem, productID string) []DiffItem {
	if p1.PayLevelCode != p2.PayLevelCode {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".pay_level_code", Diff: fmt.Sprintf("%s vs %s", p1.PayLevelCode, p2.PayLevelCode)})
	}
	if p1.PayLevelID != p2.PayLevelID {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".pay_level_id", Diff: fmt.Sprintf("%d vs %d", p1.PayLevelID, p2.PayLevelID)})
	}
	if p1.Num != p2.Num {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".num", Diff: fmt.Sprintf("%d vs %d", p1.Num, p2.Num)})
	}

	// 比较 plat_channel_product_id_map
	if len(p1.PlatChannelProductIDMap) != len(p2.PlatChannelProductIDMap) {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".plat_channel_product_id_map", Diff: fmt.Sprintf("map长度不同: %d vs %d", len(p1.PlatChannelProductIDMap), len(p2.PlatChannelProductIDMap))})
	} else {
		for k, v1 := range p1.PlatChannelProductIDMap {
			if v2, ok := p2.PlatChannelProductIDMap[k]; !ok {
				diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".plat_channel_product_id_map." + k, Diff: "键不存在"})
			} else if v1 != v2 {
				diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".plat_channel_product_id_map." + k, Diff: fmt.Sprintf("%s vs %s", v1, v2)})
			}
		}
		for k := range p2.PlatChannelProductIDMap {
			if _, ok := p1.PlatChannelProductIDMap[k]; !ok {
				diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".plat_channel_product_id_map." + k, Diff: "键不存在"})
			}
		}
	}

	// 比较 area_currency_map
	if len(p1.AreaCurrencyMap) != len(p2.AreaCurrencyMap) {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".area_currency_map", Diff: fmt.Sprintf("map长度不同: %d vs %d", len(p1.AreaCurrencyMap), len(p2.AreaCurrencyMap))})
	} else {
		for k, v1 := range p1.AreaCurrencyMap {
			if v2, ok := p2.AreaCurrencyMap[k]; !ok {
				diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".area_currency_map." + k, Diff: "键不存在"})
			} else {
				if v1.Currency != v2.Currency {
					diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".area_currency_map." + k + ".currency", Diff: fmt.Sprintf("%s vs %s", v1.Currency, v2.Currency)})
				}
				if v1.Amount != v2.Amount {
					diffs = append(diffs, DiffItem{ProductID: productID, FieldName: fieldName + ".area_currency_map." + k + ".amount", Diff: fmt.Sprintf("%s vs %s", v1.Amount, v2.Amount)})
				}
			}
		}
	}
	return diffs
}

// 比较两个ProductInfo
func compareProductInfo(p1, p2 ProductInfo, diffs []DiffItem) []DiffItem {
	productID := p1.ProductID

	if p1.ProductID != p2.ProductID {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: "product_id", Diff: fmt.Sprintf("%s vs %s", p1.ProductID, p2.ProductID)})
	}
	diffs = compareNameItems("name", p1.Name, p2.Name, diffs, productID)
	if p1.CategoryID != p2.CategoryID {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: "category_id", Diff: fmt.Sprintf("%s vs %s", p1.CategoryID, p2.CategoryID)})
	}
	if p1.LimitNum != p2.LimitNum {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: "limit_num", Diff: fmt.Sprintf("%d vs %d", p1.LimitNum, p2.LimitNum)})
	}
	if p1.LimitType != p2.LimitType {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: "limit_type", Diff: fmt.Sprintf("%d vs %d", p1.LimitType, p2.LimitType)})
	}
	if p1.PromotionID != p2.PromotionID {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: "promotion_id", Diff: fmt.Sprintf("%d vs %d", p1.PromotionID, p2.PromotionID)})
	}
	if p1.RechargePoints != p2.RechargePoints {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: "recharge_points", Diff: fmt.Sprintf("%d vs %d", p1.RechargePoints, p2.RechargePoints)})
	}
	if p1.PayType != p2.PayType {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: "pay_type", Diff: fmt.Sprintf("%d vs %d", p1.PayType, p2.PayType)})
	}

	// 比较 property_list
	if len(p1.PropertyList) != len(p2.PropertyList) {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: "property_list", Diff: fmt.Sprintf("数组长度不同: %d vs %d", len(p1.PropertyList), len(p2.PropertyList))})
	} else {
		for i := range p1.PropertyList {
			diffs = compareProperties(fmt.Sprintf("property_list[%d]", i), p1.PropertyList[i], p2.PropertyList[i], diffs, productID)
		}
	}

	// 比较 price
	diffs = comparePrices("price", p1.Price, p2.Price, diffs, productID)

	// 比较 multi_option
	if len(p1.MultiOption) != len(p2.MultiOption) {
		diffs = append(diffs, DiffItem{ProductID: productID, FieldName: "multi_option", Diff: fmt.Sprintf("数组长度不同: %d vs %d", len(p1.MultiOption), len(p2.MultiOption))})
	} else {
		for i := range p1.MultiOption {
			diffs = comparePrices(fmt.Sprintf("multi_option[%d]", i), p1.MultiOption[i], p2.MultiOption[i], diffs, productID)
		}
	}

	return diffs
}

// TestCompareProductMeta 深度比较两个JSON文件
func TestCompareProductMeta(t *testing.T) {
	originPath := "origin.json"
	optimizePath := "optimize.json"

	// 读取文件
	originResp, err := readJSONFile(originPath)
	if err != nil {
		t.Fatalf("读取origin.json失败: %v", err)
	}

	optimizeResp, err := readJSONFile(optimizePath)
	if err != nil {
		t.Fatalf("读取optimize.json失败: %v", err)
	}

	// 创建用于比较的副本，忽略request_id
	originData := originResp.Data
	optimizeData := optimizeResp.Data

	// 按product_id排序
	sortByProductID(originData.ProductInfo)
	sortByProductID(optimizeData.ProductInfo)

	// 深度比较
	var diffs []DiffItem

	// 比较商品数量
	if len(originData.ProductInfo) != len(optimizeData.ProductInfo) {
		fmt.Printf("商品数量不同: %d vs %d\n", len(originData.ProductInfo), len(optimizeData.ProductInfo))
	}

	// 比较每个商品
	minLen := len(originData.ProductInfo)
	if len(optimizeData.ProductInfo) < minLen {
		minLen = len(optimizeData.ProductInfo)
	}

	for i := 0; i < minLen; i++ {
		diffs = compareProductInfo(originData.ProductInfo[i], optimizeData.ProductInfo[i], diffs)
	}

	// 处理商品数量不一致的情况
	if len(originData.ProductInfo) > len(optimizeData.ProductInfo) {
		for i := len(optimizeData.ProductInfo); i < len(originData.ProductInfo); i++ {
			productID := originData.ProductInfo[i].ProductID
			diffs = append(diffs, DiffItem{ProductID: productID, FieldName: "", Diff: "仅在origin.json中存在"})
		}
	}
	if len(optimizeData.ProductInfo) > len(originData.ProductInfo) {
		for i := len(originData.ProductInfo); i < len(optimizeData.ProductInfo); i++ {
			productID := optimizeData.ProductInfo[i].ProductID
			diffs = append(diffs, DiffItem{ProductID: productID, FieldName: "", Diff: "仅在optimize.json中存在"})
		}
	}

	if len(diffs) == 0 {
		fmt.Println("文件内容相同")
	} else {
		// 按商品ID分组差异
		productDiffs := make(map[string][]DiffItem)
		for _, diff := range diffs {
			productDiffs[diff.ProductID] = append(productDiffs[diff.ProductID], diff)
		}

		// 按商品ID排序输出
		var sortedIDs []int
		for idStr := range productDiffs {
			var id int
			fmt.Sscanf(idStr, "%d", &id)
			sortedIDs = append(sortedIDs, id)
		}
		sort.Ints(sortedIDs)

		for _, id := range sortedIDs {
			idStr := fmt.Sprintf("%d", id)
			fmt.Printf("商品ID: %s\n", idStr)
			for _, diff := range productDiffs[idStr] {
				if diff.FieldName == "" {
					fmt.Printf("  %s\n", diff.Diff)
				} else {
					fmt.Printf("  字段: %s, 差异: %s\n", diff.FieldName, diff.Diff)
				}
			}
		}

		if len(diffs) > 20 {
			fmt.Printf("\n... 共发现 %d 处差异\n", len(diffs))
		}
	}
}

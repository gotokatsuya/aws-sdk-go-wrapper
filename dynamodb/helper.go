// DynamoDB utility

package dynamodb

import (
	AWS "github.com/awslabs/aws-sdk-go/aws"
	SDK "github.com/awslabs/aws-sdk-go/gen/dynamodb"

	"fmt"
	"strconv"
)

type Any interface{}

// Create new AttributeValue from the type of value
func createAttributeValue(v Any) SDK.AttributeValue {
	switch t := v.(type) {
	case string:
		return SDK.AttributeValue{
			S: AWS.String(t),
		}
	case int, int32, int64, uint, uint32, uint64, float32, float64:
		return SDK.AttributeValue{
			N: AWS.String(fmt.Sprint(t)),
		}
	default:
		return SDK.AttributeValue{}
	}
}

// Retrieve value from DynamoDB type
func getItemValue(val SDK.AttributeValue) Any {
	switch {
	case val.N != nil:
		data, _ := strconv.Atoi(*val.N)
		return data
	case val.S != nil:
		return *val.S
	case val.BOOL != nil:
		return *val.BOOL
	case len(val.B) > 0:
		return val.B
	case len(val.M) > 0:
		return Unmarshal(val.M)
	case len(val.NS) > 0:
		var data []int
		for _, vString := range val.NS {
			vInt, _ := strconv.Atoi(vString)
			data = append(data, vInt)
		}
		return data
	case len(val.SS) > 0:
		var data []string
		for _, vString := range val.SS {
			data = append(data, vString)
		}
		return data
	case len(val.BS) > 0:
		var data [][]byte
		for _, vBytes := range val.BS {
			data = append(data, vBytes)
		}
		return data
	case len(val.L) > 0:
		var data []interface{}
		for _, vAny := range val.L {
			data = append(data, getItemValue(vAny))
		}
		return data
	}
	return nil
}

// Convert DynamoDB Item to map data
func Unmarshal(item map[string]SDK.AttributeValue) map[string]interface{} {
	data := make(map[string]interface{})
	for key, val := range item {
		data[key] = getItemValue(val)
	}
	return data
}

func NewProvisionedThroughput(read, write int64) *SDK.ProvisionedThroughput {
	return &SDK.ProvisionedThroughput{
		ReadCapacityUnits:  AWS.Long(read),
		WriteCapacityUnits: AWS.Long(write),
	}
}

//=======================
//  KeySchema
//=======================

// Create new KeySchema slice
func NewKeySchema(elements ...*SDK.KeySchemaElement) []SDK.KeySchemaElement {
	if len(elements) > 1 {
		schema := make([]SDK.KeySchemaElement, 2, 2)
		schema[0] = *elements[0]
		schema[1] = *elements[1]
		return schema
	} else {
		schema := make([]SDK.KeySchemaElement, 1, 1)
		schema[0] = *elements[0]
		return schema
	}
}

// Create new single KeySchema
func NewKeyElement(keyName, keyType string) *SDK.KeySchemaElement {
	return &SDK.KeySchemaElement{
		AttributeName: AWS.String(keyName),
		KeyType:       AWS.String(keyType),
	}
}

// Create new single KeySchema for HashKey
func NewHashKeyElement(keyName string) *SDK.KeySchemaElement {
	return NewKeyElement(keyName, SDK.KeyTypeHash)
}

// Create new single KeySchema for RangeKey
func NewRangeKeyElement(keyName string) *SDK.KeySchemaElement {
	return NewKeyElement(keyName, SDK.KeyTypeRange)
}

//=======================
//  AttributeDefinition
//=======================

// Convert multiple definition to single slice
func NewAttributeDefinitions(attr ...SDK.AttributeDefinition) []SDK.AttributeDefinition {
	return attr
}

// Create new definition of table
func NewAttributeDefinition(attrName, attrType string) SDK.AttributeDefinition {
	newAttr := SDK.AttributeDefinition{}
	var typ *string
	switch attrType {
	case "S", "N", "B", "BOOL", "L", "M", "SS", "NS", "BS":
		typ = AWS.String(attrType)
	default:
		return newAttr
	}
	newAttr.AttributeName = AWS.String(attrName)
	newAttr.AttributeType = typ
	return newAttr
}

// NewStringAttribute returns a table AttributeDefinition for string
func NewStringAttribute(attrName string) SDK.AttributeDefinition {
	return NewAttributeDefinition(attrName, "S")
}

// NewNumberAttribute returns a table AttributeDefinition for number
func NewNumberAttribute(attrName string) SDK.AttributeDefinition {
	return NewAttributeDefinition(attrName, "N")
}

// NewByteAttribute returns a table AttributeDefinition for byte
func NewByteAttribute(attrName string) SDK.AttributeDefinition {
	return NewAttributeDefinition(attrName, "B")
}

// NewBoolAttribute returns a table AttributeDefinition for boolean
func NewBoolAttribute(attrName string) SDK.AttributeDefinition {
	return NewAttributeDefinition(attrName, "BOOL")
}

package evaluator

import (
    "fmt"
    "reflect"
    "strings"
    "github.com/Ayush/rule-engine/internal/ast"
)

func EvaluateRule(node *ast.Node, data map[string]interface{}) (bool, error) {
    if node == nil {
        return false, fmt.Errorf("invalid node")
    }

    switch node.Type {
    case "operator":
        results := make([]bool, len(node.Children))
        for i, child := range node.Children {
            result, err := EvaluateRule(child, data)
            if err != nil {
                return false, err
            }
            results[i] = result
        }

        switch node.Operator {
        case "AND":
            for _, result := range results {
                if !result {
                    return false, nil
                }
            }
            return true, nil
        case "OR":
            for _, result := range results {
                if result {
                    return true, nil
                }
            }
            return false, nil
        default:
            return false, fmt.Errorf("unknown operator: %s", node.Operator)
        }

    case "operand":
        value, exists := data[node.Field]
        if !exists {
            return false, fmt.Errorf("field %s not found in data", node.Field)
        }
        return compareValues(value, node.Value, node.Operator)
    }

    return false, fmt.Errorf("invalid node type: %s", node.Type)
}

func compareValues(a, b interface{}, operator string) (bool, error) {
    // Convert to same type if possible
    aValue := reflect.ValueOf(a)
    bValue := reflect.ValueOf(b)

    // Handle string comparison
    if aValue.Kind() == reflect.String || bValue.Kind() == reflect.String {
        aStr := strings.ToLower(fmt.Sprintf("%v", a))
        bStr := strings.ToLower(fmt.Sprintf("%v", b))
        switch operator {
        case "=":
            return aStr == bStr, nil
        case "!=":
            return aStr != bStr, nil
        default:
            return false, fmt.Errorf("invalid operator %s for strings", operator)
        }
    }

    // Handle numeric comparison
    aFloat, aErr := toFloat64(a)
    bFloat, bErr := toFloat64(b)
    if aErr != nil || bErr != nil {
        return false, fmt.Errorf("cannot compare values: %v and %v", a, b)
    }

    switch operator {
    case ">":
        return aFloat > bFloat, nil
    case "<":
        return aFloat < bFloat, nil
    case "=":
        return aFloat == bFloat, nil
    case ">=":
        return aFloat >= bFloat, nil
    case "<=":
        return aFloat <= bFloat, nil
    case "!=":
        return aFloat != bFloat, nil
    default:
        return false, fmt.Errorf("unknown operator: %s", operator)
    }
}

func toFloat64(v interface{}) (float64, error) {
    switch value := v.(type) {
    case int:
        return float64(value), nil
    case int32:
        return float64(value), nil
    case int64:
        return float64(value), nil
    case float32:
        return float64(value), nil
    case float64:
        return value, nil
    default:
        return 0, fmt.Errorf("cannot convert %v to float64", v)
    }
}

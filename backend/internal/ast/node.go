package ast

type Node struct {
    Type      string        `json:"type"`      // "operator" or "operand"
    Operator  string        `json:"operator"`  // "AND", "OR", ">", "<", "=", etc.
    Children  []*Node       `json:"children"`  // For multi-way tree
    Field     string        `json:"field"`     // For operands
    Value     interface{}   `json:"value"`     // For operands
}
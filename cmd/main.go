package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Ayush/rule-engine/internal/evaluator"
	"github.com/Ayush/rule-engine/internal/parser"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

// Rule represents a single rule with its ID and expression
type Rule struct {
    ID         string `json:"id"`
    Expression string `json:"expression"`
}

// RuleStore manages rule storage
type RuleStore struct {
    rules map[string]string  // map[ruleID]expression
    mu    sync.RWMutex
}

// CombineRequest represents the request body for combining rules
type CombineRequest struct {
    RuleIDs []string `json:"rule_ids"`
}

// EvaluateRequest represents the request body for evaluating rules
type EvaluateRequest struct {
    RuleIDs  []string                `json:"rule_ids"`
    QueryData map[string]interface{} `json:"query_data"`
}

// Response represents the API response structure
type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

var store = &RuleStore{
    rules: make(map[string]string),
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Use `mongo.Connect` directly to connect to MongoDB
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }

    // Check the connection
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Could not connect to MongoDB:", err)
    }

    // Set the database
    db = client.Database("rule-engine")
	
    router := mux.NewRouter()

    // API routes
    router.HandleFunc("/api/rules", createRule).Methods("POST")
    router.HandleFunc("/api/rules/combine", combineRules).Methods("POST")
    router.HandleFunc("/api/rules/evaluate", evaluateRules).Methods("POST")
    router.HandleFunc("/api/rule", getAllRules).Methods("GET")


    // CORS middleware
    corsMiddleware := handlers.CORS(
        handlers.AllowedOrigins([]string{"*"}),     // Allow all origins
        handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
        handlers.AllowedHeaders([]string{"Content-Type"}),
    )

    // Start server with CORS middleware
    fmt.Println("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", corsMiddleware(router)))
}

func createRule(w http.ResponseWriter, r *http.Request) {
    var rule Rule
    if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
        sendResponse(w, false, nil, "Invalid request body")
        return
    }

    // Trim spaces from rule ID
    rule.ID = strings.TrimSpace(rule.ID)

    // Validate rule syntax
    if _, err := parser.ParseRule(rule.Expression); err != nil {
        sendResponse(w, false, nil, fmt.Sprintf("Invalid rule syntax: %v", err))
        return
    }

    // Remove existing rules with the same ID
    collection := db.Collection("rules")
    _, err := collection.DeleteMany(context.TODO(), bson.M{"id": rule.ID})
    if err != nil {
        sendResponse(w, false, nil, fmt.Sprintf("Failed to delete existing rules: %v", err))
        return
    }

    // Insert the new rule into MongoDB
    _, err = collection.InsertOne(context.TODO(), rule)
    if err != nil {
        sendResponse(w, false, nil, fmt.Sprintf("Failed to store rule in DB: %v", err))
        return
    }

    sendResponse(w, true, rule, "")
}

func getAllRules(w http.ResponseWriter, r *http.Request) {
    // Query all rules from the MongoDB collection
    collection := db.Collection("rules")
    cursor, err := collection.Find(context.TODO(), bson.M{})
    if err != nil {
        sendResponse(w, false, nil, "Failed to get rules from DB")
        return
    }
    defer cursor.Close(context.TODO())

    var rules []Rule
    if err := cursor.All(context.TODO(), &rules); err != nil {
        sendResponse(w, false, nil, "Failed to decode rules")
        return
    }

    sendResponse(w, true, rules, "")
}



func combineRules(w http.ResponseWriter, r *http.Request) {
    var req CombineRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendResponse(w, false, nil, "Invalid request body")
        return
    }

    // Get rules from MongoDB
    collection := db.Collection("rules")
    filter := bson.M{"id": bson.M{"$in": req.RuleIDs}}
    cursor, err := collection.Find(context.TODO(), filter)
    if err != nil {
        sendResponse(w, false, nil, "Failed to get rules from DB")
        return
    }

    var rules []Rule
    if err := cursor.All(context.TODO(), &rules); err != nil {
        sendResponse(w, false, nil, "Failed to decode rules")
        return
    }

    if len(rules) == 0 {
        sendResponse(w, false, nil, "No valid rules found")
        return
    }

    // Combine rules
    var expressions []string
    for _, rule := range rules {
        expressions = append(expressions, rule.Expression)
    }
    combinedRule := combineRuleExpressions(expressions)

    // Respond with the combined rule (no DB storage)
    response := map[string]interface{}{
        "combined_expression": combinedRule,
    }

    sendResponse(w, true, response, "")
}



func evaluateRules(w http.ResponseWriter, r *http.Request) {
    var req EvaluateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendResponse(w, false, nil, "Invalid request body")
        return
    }

    // Get rules from MongoDB
    collection := db.Collection("rules")
    filter := bson.M{"id": bson.M{"$in": req.RuleIDs}}
    cursor, err := collection.Find(context.TODO(), filter)
    if err != nil {
        sendResponse(w, false, nil, "Failed to get rules from DB")
        return
    }

    var rules []Rule
    if err := cursor.All(context.TODO(), &rules); err != nil {
        sendResponse(w, false, nil, "Failed to decode rules")
        return
    }

    if len(rules) == 0 {
        sendResponse(w, false, nil, "No valid rules found")
        return
    }

    // Combine rules
    var expressions []string
    for _, rule := range rules {
        expressions = append(expressions, rule.Expression)
    }
    combinedRule := combineRuleExpressions(expressions)

    // Parse combined rule
    ast, err := parser.ParseRule(combinedRule)
    if err != nil {
        sendResponse(w, false, nil, fmt.Sprintf("Error parsing combined rules: %v", err))
        return
    }

    // Evaluate rule
    result, err := evaluator.EvaluateRule(ast, req.QueryData)
    if err != nil {
        sendResponse(w, false, nil, fmt.Sprintf("Error evaluating rules: %v", err))
        return
    }

    // Store evaluation result in MongoDB
    evaluation := bson.M{"rule_ids": req.RuleIDs, "query_data": req.QueryData, "result": result}
    evalCollection := db.Collection("evaluations")
    _, err = evalCollection.InsertOne(context.TODO(), evaluation)
    if err != nil {
        sendResponse(w, false, nil, "Failed to store evaluation in DB")
        return
    }

    response := map[string]interface{}{
        "combined_expression": combinedRule,
        "result":             result,
    }

    sendResponse(w, true, response, "")
}


func combineRuleExpressions(rules []string) string {
    if len(rules) == 0 {
        return ""
    }
    if len(rules) == 1 {
        return rules[0]
    }

    // Wrap each rule in parentheses and join with AND
    wrappedRules := make([]string, len(rules))
    for i, rule := range rules {
        if strings.Contains(rule, "OR") || strings.Contains(rule, " or ") {
            wrappedRules[i] = "(" + rule + ")"
        } else {
            wrappedRules[i] = rule
        }
    }
    return strings.Join(wrappedRules, " AND ")
}

func sendResponse(w http.ResponseWriter, success bool, data interface{}, errMsg string) {
    w.Header().Set("Content-Type", "application/json")
    
    response := Response{
        Success: success,
        Data:    data,
        Error:   errMsg,
    }

    if !success {
        w.WriteHeader(http.StatusBadRequest)
    }

    json.NewEncoder(w).Encode(response)
}
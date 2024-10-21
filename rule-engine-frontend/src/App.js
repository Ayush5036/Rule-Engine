import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
    const [ruleID, setRuleID] = useState('');
    const [ruleExpression, setRuleExpression] = useState('');
    const [evaluationData, setEvaluationData] = useState('');
    const [combinedExpression, setCombinedExpression] = useState('');
    const [evaluateResult, setEvaluateResult] = useState('');
    const [allRules, setAllRules] = useState([]);

    // Fetch all rules on page load
    useEffect(() => {
        fetchRules();
    }, []);

    const fetchRules = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/rule');
            const data = await response.json();
            if (data.success) {
                setAllRules(data.data);
            } else {
                alert(`Error fetching rules: ${data.error}`);
            }
        } catch (error) {
            alert('Error connecting to backend: ' + error.message);
        }
    };

    const createRule = async () => {
        const rule = { id: ruleID, expression: ruleExpression };
        try {
            const response = await fetch('http://localhost:8080/api/rules', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(rule),
            });

            const data = await response.json();
            if (data.success) {
                alert('Rule created successfully');
                fetchRules();  // Refresh the list of rules after creation
            } else {
                alert(`Error creating rule: ${data.error}`);
            }
        } catch (error) {
            alert('Error connecting to backend: ' + error.message);
        }
    };

    const combineRules = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/rules/combine', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ rule_ids: [ruleID] }), // Adjust according to the rules you want to combine
            });

            const data = await response.json();
            if (data.success) {
                setCombinedExpression(data.data.combined_expression);
            } else {
                alert(`Error combining rules: ${data.error}`);
            }
        } catch (error) {
            alert('Error connecting to backend: ' + error.message);
        }
    };

    const evaluateRules = async () => {
        let parsedData;
        try {
            parsedData = JSON.parse(evaluationData); // Validate the input
        } catch (e) {
            alert('Invalid query format. Please use correct JSON format.');
            return;
        }

        const evaluationRequest = {
            rule_ids: [ruleID], 
            query_data: parsedData, 
        };

        try {
            const response = await fetch('http://localhost:8080/api/rules/evaluate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(evaluationRequest),
            });

            const data = await response.json();
            if (data.success) {
                setEvaluateResult(JSON.stringify(data.data.result));
            } else {
                alert(`Error evaluating rule: ${data.error}`);
            }
        } catch (error) {
            alert('Error connecting to backend: ' + error.message);
        }
    };

    return (
        <div className="app-container">
            {/* Sidebar for displaying all rules */}
            <div className="sidebar">
                <h3>All Rules</h3>
                <ul>
                    {allRules.map(rule => (
                        <li key={rule.id}>{rule.id}: {rule.expression}</li>
                    ))}
                </ul>
            </div>

            {/* Main content for create, combine, evaluate */}
            <div className="main-content">
                <div className="box">
                    <h3>Create Rule</h3>
                    <input
                        type="text"
                        placeholder="Rule ID"
                        value={ruleID}
                        onChange={(e) => setRuleID(e.target.value)}
                    />
                    <input
                        type="text"
                        placeholder="Rule Expression"
                        value={ruleExpression}
                        onChange={(e) => setRuleExpression(e.target.value)}
                    />
                    <button className="button evaluate" onClick={createRule}>
                        Create Rule
                    </button>
                </div>

                <div className="box">
                    <h3>Combine Rules</h3>
                    <button className="button" onClick={combineRules}>
                        Combine Rules
                    </button>
                    <p>Combined Expression: {combinedExpression}</p>
                </div>

                <div className="box">
                    <h3>Evaluate Rules</h3>
                    <textarea
                        placeholder="Enter Query Data (JSON format)"
                        value={evaluationData}
                        onChange={(e) => setEvaluationData(e.target.value)}
                    />
                    <button className="button evaluate" onClick={evaluateRules}>
                        Evaluate Rules
                    </button>
                    <p>Evaluation Result: {evaluateResult}</p>
                </div>
            </div>
        </div>
    );
}

export default App;

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"core"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// StartMCPServer initializes and starts the standard I/O MCP server.
func StartMCPServer() error {
	// 1. 서버 생성
	s := server.NewMCPServer(
		"fehu-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// 도구(Tool): DB 열기
	openDbTool := mcp.NewTool("open_db",
		mcp.WithDescription("Open or create a database file"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Database name without .db extension. Use absolute path (e.g. /Users/john/budget) or relative to current working directory.")),
	)
	s.AddTool(openDbTool, handleOpenDB)

	// 도구(Tool): 계정 목록 조회
	getAccTool := mcp.NewTool("get_accounts",
		mcp.WithDescription("Get a list of accounts in the Fehu database with balances. Can optionally filter by name or description."),
		mcp.WithString("name", mcp.Description("Optional exact account name to fetch")),
		mcp.WithString("desc", mcp.Description("Optional description keyword to filter by")),
	)
	s.AddTool(getAccTool, handleGetAccounts)

	// 도구(Tool): 새 계정 생성
	createAccTool := mcp.NewTool("create_account",
		mcp.WithDescription("Create a new account in the database"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Account name (e.g. expense:food)")),
		mcp.WithString("desc", mcp.Description("Optional description for the account")),
	)
	s.AddTool(createAccTool, handleCreateAccount)

	// 도구(Tool): 여러 계정 일괄 생성
	batchCreateAccTool := mcp.NewTool("batch_create_accounts",
		mcp.WithDescription("Create multiple accounts at once using a JSON array"),
		mcp.WithString("accounts_json", mcp.Required(), mcp.Description(`JSON array of accounts. Format: [{"name":"expense:food", "desc":"optional"}]`)),
	)
	s.AddTool(batchCreateAccTool, handleBatchCreateAccounts)

	// 도구(Tool): 계정 수정
	updateAccTool := mcp.NewTool("update_account",
		mcp.WithDescription("Update an existing account's name or description"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Current account name")),
		mcp.WithString("new_name", mcp.Description("New account name (if renaming)")),
		mcp.WithString("desc", mcp.Description("New description for the account")),
	)
	s.AddTool(updateAccTool, handleUpdateAccount)

	// 도구(Tool): 계정 삭제
	deleteAccTool := mcp.NewTool("delete_account",
		mcp.WithDescription("Delete an existing account"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Account name to delete")),
	)
	s.AddTool(deleteAccTool, handleDeleteAccount)

	// 도구(Tool): 트랜잭션 목록 조회
	getTxnTool := mcp.NewTool("get_transactions",
		mcp.WithDescription("Get a list of transactions in the Fehu database"),
		mcp.WithNumber("id", mcp.Description("Optional exact transaction ID to fetch")),
		mcp.WithString("desc", mcp.Description("Optional description keyword to filter by")),
		mcp.WithString("time_range", mcp.Description("Optional time range filter (e.g. '2024-01-01;00:00:00~2024-12-31;23:59:59')")),
	)
	s.AddTool(getTxnTool, handleGetTransactions)

	// 도구(Tool): 새 트랜잭션 생성
	createTxnTool := mcp.NewTool("create_transaction",
		mcp.WithDescription("Create a new transaction"),
		mcp.WithString("record", mcp.Required(), mcp.Description("Record format (e.g. income:salary<50000;asset:bank>50000)")),
		mcp.WithString("desc", mcp.Description("Optional description for the transaction (can include #tags)")),
		mcp.WithString("time", mcp.Description("Optional time (YYYY-MM-DD;HH:MM:SS format)")),
	)
	s.AddTool(createTxnTool, handleCreateTransaction)

	// 도구(Tool): 여러 트랜잭션 일괄 생성
	batchCreateTxnTool := mcp.NewTool("batch_create_transactions",
		mcp.WithDescription("Create multiple transactions at once using a JSON array"),
		mcp.WithString("transactions_json", mcp.Required(), mcp.Description(`JSON array of transactions. Format: [{"record":"acc1<100;acc2>100", "desc":"optional", "time":"YYYY-MM-DD;HH:MM:SS"}]`)),
	)
	s.AddTool(batchCreateTxnTool, handleBatchCreateTransactions)

	// 도구(Tool): 트랜잭션 수정
	updateTxnTool := mcp.NewTool("update_transaction",
		mcp.WithDescription("Update an existing transaction's description or time"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Transaction ID to update")),
		mcp.WithString("desc", mcp.Description("New description for the transaction (can include #tags)")),
		mcp.WithString("time", mcp.Description("New time (YYYY-MM-DD;HH:MM:SS format)")),
	)
	s.AddTool(updateTxnTool, handleUpdateTransaction)

	// 도구(Tool): 트랜잭션 레코드(금액/계정) 수정
	updateTxnRecordTool := mcp.NewTool("update_transaction_record",
		mcp.WithDescription("Replace all records (account flows and amounts) of an existing transaction"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Transaction ID to update")),
		mcp.WithString("record", mcp.Required(), mcp.Description("New record format (e.g. income:salary<50000;asset:bank>50000)")),
	)
	s.AddTool(updateTxnRecordTool, handleUpdateTransactionRecord)

	// 도구(Tool): 여러 트랜잭션 일괄 수정
	batchUpdateTxnTool := mcp.NewTool("batch_update_transactions",
		mcp.WithDescription("Update multiple transactions at once using a JSON array"),
		mcp.WithString("updates_json", mcp.Required(), mcp.Description(`JSON array of updates. Format: [{"id": 1, "desc": "new desc (optional)", "time": "2024-01-01;12:00:00 (optional)"}]`)),
	)
	s.AddTool(batchUpdateTxnTool, handleBatchUpdateTransactions)

	// 도구(Tool): 트랜잭션 삭제
	deleteTxnTool := mcp.NewTool("delete_transaction",
		mcp.WithDescription("Delete an existing transaction"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Transaction ID to delete")),
	)
	s.AddTool(deleteTxnTool, handleDeleteTransaction)

	// 도구(Tool): 태그 목록 조회
	getTagTool := mcp.NewTool("get_tags",
		mcp.WithDescription("Get a list of all tags in the Fehu database"),
	)
	s.AddTool(getTagTool, handleGetTags)

	// 도구(Tool): 태그 수정
	updateTagTool := mcp.NewTool("update_tag",
		mcp.WithDescription("Update an existing tag's name or description"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Current tag name")),
		mcp.WithString("new_name", mcp.Description("New tag name (if renaming)")),
		mcp.WithString("desc", mcp.Description("New description for the tag")),
	)
	s.AddTool(updateTagTool, handleUpdateTag)

	// 도구(Tool): 태그 삭제
	deleteTagTool := mcp.NewTool("delete_tag",
		mcp.WithDescription("Delete an existing tag"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Tag name to delete")),
	)
	s.AddTool(deleteTagTool, handleDeleteTag)

	// 도구(Tool): 계산기(calc)
	calcTool := mcp.NewTool("calc",
		mcp.WithDescription("Evaluate a calc expression in Fehu"),
		mcp.WithString("expression", mcp.Required(), mcp.Description("The expression to evaluate (e.g. sum(acc(__all__, 'expense*'))) ")),
	)
	s.AddTool(calcTool, handleCalc)

	// 도구(Tool): 계산기 변수 저장(def)
	defCalcTool := mcp.NewTool("def_calc",
		mcp.WithDescription("Evaluate a calc expression and save the result as a named variable (def command)"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Variable name to save as")),
		mcp.WithString("expression", mcp.Required(), mcp.Description("The expression to evaluate (e.g. acc(__all__, 'expense*'))")),
	)
	s.AddTool(defCalcTool, handleDefCalc)

	// 도구(Tool): 재무 요약(Summary)
	getSummaryTool := mcp.NewTool("get_summary",
		mcp.WithDescription("Get a quick financial summary including total assets, liabilities, equity, income, expenses, net worth, and check the accounting equation"),
		mcp.WithString("time_range", mcp.Description("Optional time range filter (e.g. '2024-01-01;00:00:00~2024-12-31;23:59:59')")),
	)
	s.AddTool(getSummaryTool, handleGetSummary)

	// 도구(Tool): DB 닫기
	closeDbTool := mcp.NewTool("close_db",
		mcp.WithDescription("Close the current database"),
	)
	s.AddTool(closeDbTool, handleCloseDB)

	// 도구(Tool): README 문서 제공
	readmeTool := mcp.NewTool("get_readme",
		mcp.WithDescription("Get the README documentation for Fehu to understand its features, commands, and architecture"),
	)
	s.AddTool(readmeTool, handleGetReadme)

	return server.ServeStdio(s)
}

func handleOpenDB(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Error: %v", r))
		}
	}()
	args := request.GetArguments()
	name := args["name"].(string)

	core.Open(name + ".db")
	return mcp.NewToolResultText(fmt.Sprintf("Database %s.db opened successfully", name)), nil
}

func handleGetAccounts(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Error: %v", r))
		}
	}()
	args := request.GetArguments()
	name, _ := args["name"].(string)
	desc, _ := args["desc"].(string)

	var ret []int64
	if name != "" {
		ret = []int64{core.GetAccByName(name)}
	} else if desc != "" {
		ret = core.GetAccByDesc(desc)
	} else {
		ret = core.GetAcc()
	}
	return mcp.NewToolResultText(core.PrintAccs(ret)), nil
}

func handleCreateAccount(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Error: %v", r))
		}
	}()
	args := request.GetArguments()
	name := core.SureName(args["name"].(string))
	desc, _ := args["desc"].(string)

	id := core.NewAcc(name, desc)
	return mcp.NewToolResultText(fmt.Sprintf("account #%d created", id)), nil
}

func handleBatchCreateAccounts(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Batch account creation failed.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()
	jsonStr := args["accounts_json"].(string)

	var accs []struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &accs); err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Failed to parse accounts JSON.\nError: %v\n\nPlease ensure your input matches this format:\n[{\"name\":\"expense:food\", \"desc\":\"optional\"}]\n\nProvided input:\n%s", err, jsonStr)), nil
	}

	var createdIDs []string
	for _, a := range accs {
		id := core.NewAcc(core.SureName(a.Name), a.Desc)
		createdIDs = append(createdIDs, fmt.Sprintf("%d", id))
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully created %d accounts. IDs: %s", len(createdIDs), strings.Join(createdIDs, ", "))), nil
}

func handleUpdateAccount(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Account update failed.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()
	name := args["name"].(string)

	if descArg, ok := args["desc"]; ok {
		if desc, ok := descArg.(string); ok {
			core.SureID(core.AltAcc(name, &desc))
		}
	}

	if newNameArg, ok := args["new_name"]; ok {
		if newName, ok := newNameArg.(string); ok && newName != "" {
			core.SureID(core.AltRenameAcc(name, newName))
			name = newName
		}
	}
	return mcp.NewToolResultText(fmt.Sprintf("Account '%s' updated successfully", name)), nil
}

func handleDeleteAccount(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Account deletion failed.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()
	name := args["name"].(string)

	core.SureID(core.DelAcc(name))
	return mcp.NewToolResultText(fmt.Sprintf("Account '%s' deleted successfully", name)), nil
}

func handleGetTransactions(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Error: %v", r))
		}
	}()
	args := request.GetArguments()
	desc, _ := args["desc"].(string)
	timeRange, _ := args["time_range"].(string)

	var ret []int64
	if idArg, ok := args["id"]; ok {
		var id int64
		switch v := idArg.(type) {
		case float64:
			id = int64(v)
		case int:
			id = int64(v)
		case string:
			fmt.Sscanf(v, "%d", &id)
		}
		ret = []int64{id}
	} else if desc != "" {
		ret = core.GetTxnByDesc(desc)
	} else if timeRange != "" {
		tokens := strings.Split(timeRange, "~")
		var A, B *time.Time
		if len(tokens) > 0 && tokens[0] != "" {
			t := core.ParseTime(tokens[0])
			A = &t
		}
		if len(tokens) > 1 && tokens[1] != "" {
			t := core.ParseTime(tokens[1])
			B = &t
		}
		ret = core.GetTxnByTime(A, B)
	} else {
		ret = core.GetTxn()
	}
	return mcp.NewToolResultText(core.PrintTxns(ret)), nil
}

func handleCreateTransaction(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Transaction creation failed. Please check if the accounts exist and the format is correct.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()
	recordStr := args["record"].(string)

	desc, _ := args["desc"].(string)
	timeStr, _ := args["time"].(string)

	pats := ParseTxnPattern(recordStr)
	t := time.Now()
	if timeStr != "" {
		t = core.ParseTime(timeStr)
	}

	tid := core.NewTxn(desc, t)
	for _, p := range pats {
		aid := core.SureID(core.GetAccByName(p.Name))
		core.NewRecord(tid, aid, p.Amount)
	}
	return mcp.NewToolResultText(fmt.Sprintf("txn #%d created", tid)), nil
}

func handleBatchCreateTransactions(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Batch transaction creation failed.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()
	jsonStr := args["transactions_json"].(string)

	var txns []struct {
		Record string `json:"record"`
		Desc   string `json:"desc"`
		Time   string `json:"time"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &txns); err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Failed to parse transactions JSON.\nError: %v\n\nPlease ensure your input matches this format:\n[{\"record\":\"acc1<100;acc2>100\", \"desc\":\"optional\", \"time\":\"YYYY-MM-DD;HH:MM:SS\"}]\n\nProvided input:\n%s", err, jsonStr)), nil
	}

	var createdIDs []string
	for _, t := range txns {
		pats := ParseTxnPattern(t.Record)
		ts := time.Now()
		if t.Time != "" {
			ts = core.ParseTime(t.Time)
		}

		tid := core.NewTxn(t.Desc, ts)
		for _, p := range pats {
			aid := core.SureID(core.GetAccByName(p.Name))
			core.NewRecord(tid, aid, p.Amount)
		}
		createdIDs = append(createdIDs, fmt.Sprintf("%d", tid))
	}

	return mcp.NewToolResultText(fmt.Sprintf("Successfully created %d transactions. IDs: %s", len(createdIDs), strings.Join(createdIDs, ", "))), nil
}

func handleUpdateTransaction(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Transaction update failed.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()

	var id int64
	switch v := args["id"].(type) {
	case float64:
		id = int64(v)
	case int:
		id = int64(v)
	case string:
		fmt.Sscanf(v, "%d", &id)
	}

	var descPtr *string
	if descArg, ok := args["desc"]; ok {
		if desc, ok := descArg.(string); ok {
			descPtr = &desc
		}
	}

	var timePtr *time.Time
	if timeArg, ok := args["time"]; ok {
		if timeStr, ok := timeArg.(string); ok && timeStr != "" {
			t := core.ParseTime(timeStr)
			timePtr = &t
		}
	}

	core.SureID(core.AltTxn(id, descPtr, timePtr))
	return mcp.NewToolResultText(fmt.Sprintf("txn #%d updated successfully", id)), nil
}

func handleUpdateTransactionRecord(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Transaction record update failed.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()

	var id int64
	switch v := args["id"].(type) {
	case float64:
		id = int64(v)
	case int:
		id = int64(v)
	case string:
		fmt.Sscanf(v, "%d", &id)
	}

	recordStr := args["record"].(string)
	pats := ParseTxnPattern(recordStr)

	var records []core.Record
	for _, p := range pats {
		aid := core.SureID(core.GetAccByName(p.Name))
		records = append(records, core.Record{TID: id, AID: aid, Amount: p.Amount})
	}

	core.SureID(core.AltTxnRecord(id, records))
	return mcp.NewToolResultText(fmt.Sprintf("txn #%d records updated successfully", id)), nil
}

func handleBatchUpdateTransactions(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Batch transaction update failed.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()
	jsonStr := args["updates_json"].(string)

	var updates []struct {
		ID   int64   `json:"id"`
		Desc *string `json:"desc"`
		Time *string `json:"time"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &updates); err != nil {
		return mcp.NewToolResultText(fmt.Sprintf("Failed to parse transaction updates JSON.\nError: %v\n\nPlease ensure your input matches this format:\n[{\"id\": 1, \"desc\": \"new desc\", \"time\": \"YYYY-MM-DD;HH:MM:SS\"}]\n\nProvided input:\n%s", err, jsonStr)), nil
	}

	var updatedIDs []string
	for _, u := range updates {
		var timePtr *time.Time
		if u.Time != nil && *u.Time != "" {
			t := core.ParseTime(*u.Time)
			timePtr = &t
		}
		core.SureID(core.AltTxn(u.ID, u.Desc, timePtr))
		updatedIDs = append(updatedIDs, fmt.Sprintf("%d", u.ID))
	}
	return mcp.NewToolResultText(fmt.Sprintf("Successfully updated %d transactions. IDs: %s", len(updatedIDs), strings.Join(updatedIDs, ", "))), nil
}

func handleDeleteTransaction(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Transaction deletion failed.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()

	var id int64
	switch v := args["id"].(type) {
	case float64:
		id = int64(v)
	case int:
		id = int64(v)
	case string:
		fmt.Sscanf(v, "%d", &id)
	}

	core.DelTxn(id)
	return mcp.NewToolResultText(fmt.Sprintf("txn #%d deleted successfully", id)), nil
}

func handleGetTags(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Error: %v", r))
		}
	}()
	ret := core.GetTag()
	return mcp.NewToolResultText(core.PrintTags(ret)), nil
}

func handleUpdateTag(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Tag update failed.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()
	name := args["name"].(string)

	if descArg, ok := args["desc"]; ok {
		if desc, ok := descArg.(string); ok {
			core.SureID(core.AltTag(name, &desc))
		}
	}

	if newNameArg, ok := args["new_name"]; ok {
		if newName, ok := newNameArg.(string); ok && newName != "" {
			core.SureName(newName)
			core.AltRenameTag(name, newName)
			name = newName
		}
	}
	return mcp.NewToolResultText(fmt.Sprintf("Tag '%s' updated successfully", name)), nil
}

func handleDeleteTag(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Tag deletion failed.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()
	name := args["name"].(string)

	core.SureID(core.DelTag(name))
	return mcp.NewToolResultText(fmt.Sprintf("Tag '%s' deleted successfully", name)), nil
}

func handleCalc(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Error: %v", r))
		}
	}()
	args := request.GetArguments()
	expr, ok := args["expression"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid expression")
	}

	resStmt := core.CalcStmt(expr)
	return mcp.NewToolResultText(resStmt.String()), nil
}

func handleDefCalc(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Error: %v", r))
		}
	}()
	args := request.GetArguments()
	name := args["name"].(string)
	expr := args["expression"].(string)

	core.SureName(name)
	core.DefStmt(name, expr)
	return mcp.NewToolResultText(fmt.Sprintf("Variable '%s' defined successfully", name)), nil
}

func handleGetSummary(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Failed to get summary.\nInternal Error: %v", r))
		}
	}()
	args := request.GetArguments()
	timeRange, _ := args["time_range"].(string)

	targetTable := "__all__"
	periodStr := "All Time"
	if timeRange != "" {
		tokens := strings.Split(timeRange, "~")
		start, end := "", ""
		if len(tokens) > 0 {
			start = tokens[0]
		}
		if len(tokens) > 1 {
			end = tokens[1]
		}
		targetTable = fmt.Sprintf("between(__all__, '%s', '%s')", start, end)
		periodStr = timeRange
	}

	evalStr := func(expr string) string {
		return core.CalcStmt(expr).String()
	}

	formatMoney := func(valStr string) string {
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return valStr
		}
		s := fmt.Sprintf("%.2f", val)
		s = strings.TrimSuffix(s, ".00")
		parts := strings.Split(s, ".")
		intPart := parts[0]
		decPart := ""
		if len(parts) > 1 {
			decPart = "." + parts[1]
		}
		isNeg := false
		if strings.HasPrefix(intPart, "-") {
			isNeg = true
			intPart = intPart[1:]
		}
		res := ""
		for i, c := range intPart {
			if i > 0 && (len(intPart)-i)%3 == 0 {
				res += ","
			}
			res += string(c)
		}
		if isNeg {
			res = "-" + res
		}
		return res + decPart
	}

	assets := formatMoney(evalStr(fmt.Sprintf("sum(acc(%s, 'asset*'))", targetTable)))
	liabs := formatMoney(evalStr(fmt.Sprintf("sum(acc(%s, 'liability*'))", targetTable)))
	equity := formatMoney(evalStr(fmt.Sprintf("sum(acc(%s, 'equity*'))", targetTable)))
	income := formatMoney(evalStr(fmt.Sprintf("sum(acc(%s, 'income*'))", targetTable)))
	expenses := formatMoney(evalStr(fmt.Sprintf("sum(acc(%s, 'expense*'))", targetTable)))

	netWorth := formatMoney(evalStr(fmt.Sprintf("sum(acc(%s, 'asset*')) - sum(acc(%s, 'liability*'))", targetTable, targetTable)))
	netIncome := formatMoney(evalStr(fmt.Sprintf("sum(acc(%s, 'income*')) - sum(acc(%s, 'expense*'))", targetTable, targetTable)))

	// 회계 등식 검증: 복식부기에서 모든 계정의 합은 0이어야 합니다.
	totalStr := evalStr(fmt.Sprintf("sum(acc(%s, '*'))", targetTable))
	total, _ := strconv.ParseFloat(totalStr, 64)

	eqCheck := "[FAIL] Imbalanced (Sum of all accounts != 0)"
	if total > -0.001 && total < 0.001 {
		eqCheck = "[PASS] Balanced (Assets = Liabilities + Equity)"
	}

	summary := fmt.Sprintf(`Financial Summary (%s):
Total Assets:      %s
Total Liabilities: %s
Total Equity:      %s
Total Income:      %s
Total Expenses:    %s
-----------------------------
Net Worth:         %s
Net Income:        %s
Accounting Eq:     %s`, periodStr, assets, liabs, equity, income, expenses, netWorth, netIncome, eqCheck)

	return mcp.NewToolResultText(summary), nil
}

func handleCloseDB(ctx context.Context, request mcp.CallToolRequest) (res *mcp.CallToolResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = mcp.NewToolResultText(fmt.Sprintf("Error: %v", r))
		}
	}()
	core.Close()
	return mcp.NewToolResultText("Database closed successfully"), nil
}

func handleGetReadme(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 1. 현재 작업 경로에서 확인
	content, err := os.ReadFile("README.md")
	if err == nil {
		return mcp.NewToolResultText(string(content)), nil
	}

	// 2. 상위 폴더에서 확인 (개발 환경 고려)
	content, err = os.ReadFile("../README.md")
	if err == nil {
		return mcp.NewToolResultText(string(content)), nil
	}

	// 3. 실행 파일이 위치한 폴더에서 확인 (배포 환경 고려)
	if exe, err := os.Executable(); err == nil {
		if content, err := os.ReadFile(filepath.Join(filepath.Dir(exe), "README.md")); err == nil {
			return mcp.NewToolResultText(string(content)), nil
		}
	}

	return mcp.NewToolResultText("Error: Could not find README.md file. Please ensure it is in the same directory as the executable."), nil
}

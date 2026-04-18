package core_test

import (
	"core"
	"os"
	"path/filepath"
	"testing"

	"github.com/Rhymond/go-money"
)

// setup opens a temporary SQLite DB and registers cleanup.
func setup(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.db")
	if err := core.Open(path); err != nil {
		t.Fatalf("Open: %v", err)
	}
	core.Code = "USD"
	if money.GetCurrency(core.Code) == nil {
		money.AddCurrency(core.Code, core.Code+" ", "1 $", ".", ",", 8)
	}
	t.Cleanup(core.Close)
}

// ── Open ─────────────────────────────────────────────────────────────────────

func TestOpen_InvalidPath(t *testing.T) {
	// 존재하지 않는 디렉터리 안의 경로 → SQLite가 디렉터리를 만들 수 없어 에러
	err := core.Open("/nonexistent_dir_fehu_test/sub/test.db")
	if err == nil {
		core.Close()
		t.Error("expected error for non-creatable path, got nil")
	}
}

// ── Account ───────────────────────────────────────────────────────────────────

func TestNewAcc_And_GetAcc(t *testing.T) {
	setup(t)

	id, err := core.NewAcc("checking", "my account")
	if err != nil {
		t.Fatalf("NewAcc: %v", err)
	}
	if id <= 0 {
		t.Errorf("expected positive ID, got %d", id)
	}

	ids, err := core.GetAcc()
	if err != nil {
		t.Fatalf("GetAcc: %v", err)
	}
	found := false
	for _, v := range ids {
		if v == id {
			found = true
		}
	}
	if !found {
		t.Errorf("created account %d not found in GetAcc result", id)
	}
}

func TestGetAccByName_NotFound(t *testing.T) {
	setup(t)

	_, err := core.GetAccByName("nonexistent")
	if err == nil {
		t.Error("expected error for missing account, got nil")
	}
}

func TestAltRenameAcc_CascadesChildren(t *testing.T) {
	setup(t)

	if _, err := core.NewAcc("parent", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := core.NewAcc("parent:child", ""); err != nil {
		t.Fatal(err)
	}

	if _, err := core.AltRenameAcc("parent", "newparent"); err != nil {
		t.Fatalf("AltRenameAcc: %v", err)
	}

	// 부모 이름 변경 확인
	if _, err := core.GetAccByName("newparent"); err != nil {
		t.Error("renamed parent not found")
	}
	// 자식도 cascade 변경 확인
	if _, err := core.GetAccByName("newparent:child"); err != nil {
		t.Error("cascaded child not found")
	}
	// 기존 이름은 없어야 함
	if _, err := core.GetAccByName("parent"); err == nil {
		t.Error("old parent name should not exist")
	}
}

func TestAltRenameAcc_DuplicateName(t *testing.T) {
	setup(t)

	core.NewAcc("acc1", "")
	core.NewAcc("acc2", "")

	_, err := core.AltRenameAcc("acc1", "acc2")
	if err == nil {
		t.Error("expected error when renaming to existing name, got nil")
	}
}

func TestDelAcc_WithRecords_Blocked(t *testing.T) {
	setup(t)

	aid, _ := core.NewAcc("wallet", "")
	aid2, _ := core.NewAcc("expense", "")
	tid, _ := core.NewTxn("test", 0)
	core.NewRecord(tid, aid, money.New(1000, "USD"))
	core.NewRecord(tid, aid2, money.New(-1000, "USD"))

	_, err := core.DelAcc("wallet")
	if err == nil {
		t.Error("expected error deleting account with records, got nil")
	}
}

// ── Transaction ───────────────────────────────────────────────────────────────

func TestNewTxn_And_GetTxn(t *testing.T) {
	setup(t)

	id, err := core.NewTxn("lunch", 1000000)
	if err != nil {
		t.Fatalf("NewTxn: %v", err)
	}

	ids := core.GetTxn()
	found := false
	for _, v := range ids {
		if v == id {
			found = true
		}
	}
	if !found {
		t.Errorf("created txn %d not found in GetTxn", id)
	}
}

func TestAltTxnRecord_DoesNotDeleteTransaction(t *testing.T) {
	// 버그 수정 검증: AltTxnRecord가 txn 자체를 삭제하면 안 됨
	setup(t)

	aid1, _ := core.NewAcc("assets", "")
	aid2, _ := core.NewAcc("expenses", "")
	tid, _ := core.NewTxn("original", 1000)

	core.NewRecord(tid, aid1, money.New(5000, "USD"))
	core.NewRecord(tid, aid2, money.New(-5000, "USD"))

	newRecords := []core.Record{
		{TID: tid, AID: aid1, Amount: money.New(9000, "USD")},
		{TID: tid, AID: aid2, Amount: money.New(-9000, "USD")},
	}

	if _, err := core.AltTxnRecord(tid, newRecords); err != nil {
		t.Fatalf("AltTxnRecord: %v", err)
	}

	// 트랜잭션이 여전히 존재해야 함
	txn, err := core.GetTxnByID(tid)
	if err != nil {
		t.Fatalf("transaction was deleted after AltTxnRecord: %v", err)
	}
	if txn.ID != tid {
		t.Errorf("expected txn ID %d, got %d", tid, txn.ID)
	}

	// 새 레코드가 반영됐는지 확인
	if len(txn.Record) != 2 {
		t.Errorf("expected 2 records, got %d", len(txn.Record))
	}
}

func TestDelTxn(t *testing.T) {
	setup(t)

	tid, _ := core.NewTxn("to-delete", 0)
	if _, err := core.DelTxn(tid); err != nil {
		t.Fatalf("DelTxn: %v", err)
	}

	_, err := core.GetTxnByID(tid)
	if err == nil {
		t.Error("expected error after deleting transaction, got nil")
	}
}

// ── GetAccAmount error propagation ───────────────────────────────────────────

func TestGetAccAmount_EmptyAccount(t *testing.T) {
	setup(t)

	aid, _ := core.NewAcc("empty", "")
	amount, err := core.GetAccAmount(aid)
	if err != nil {
		t.Fatalf("GetAccAmount: %v", err)
	}
	if !amount.IsZero() {
		t.Errorf("expected zero amount, got %s", amount.Display())
	}
}

// ── Tag ───────────────────────────────────────────────────────────────────────

func TestNewTag_And_GetTag(t *testing.T) {
	setup(t)

	id, err := core.NewTag("food", "food expenses")
	if err != nil {
		t.Fatalf("NewTag: %v", err)
	}

	ids, err := core.GetTag()
	if err != nil {
		t.Fatalf("GetTag: %v", err)
	}
	found := false
	for _, v := range ids {
		if v == id {
			found = true
		}
	}
	if !found {
		t.Errorf("created tag %d not found in GetTag", id)
	}
}

func TestDelTag(t *testing.T) {
	setup(t)

	id, _ := core.NewTag("temp", "")
	if _, err := core.DelTag("temp"); err != nil {
		t.Fatalf("DelTag: %v", err)
	}

	_, err := core.GetTagByID(id)
	if err == nil {
		t.Error("expected error after deleting tag, got nil")
	}
}

// ── LoadAllDefsFromDB ─────────────────────────────────────────────────────────

func TestLoadAllDefsFromDB_OnOpen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "defs.db")

	// 첫 번째 오픈: var 저장
	if err := core.Open(path); err != nil {
		t.Fatal(err)
	}
	if err := core.DefStmt("myvar_load_test", "1+1"); err != nil {
		core.Close()
		t.Fatalf("DefStmt: %v", err)
	}
	core.Close()

	// Vars에서 수동 제거 (같은 프로세스에서 재오픈 시뮬레이션)
	delete(core.Vars, "myvar_load_test")

	// 두 번째 오픈: var가 DB에서 자동으로 로드돼야 함
	if err := core.Open(path); err != nil {
		t.Fatalf("second Open: %v", err)
	}
	defer core.Close()

	if _, ok := core.Vars["myvar_load_test"]; !ok {
		t.Error("var 'myvar_load_test' was not loaded from DB on Open")
	}

	// 정리
	delete(core.Vars, "myvar_load_test")
}

// ── GetTxnByDesc ──────────────────────────────────────────────────────────────

func TestGetTxnByDesc(t *testing.T) {
	setup(t)

	core.NewTxn("coffee shop", 1000)
	core.NewTxn("grocery store", 2000)
	core.NewTxn("coffee machine", 3000)

	ids, err := core.GetTxnByDesc("coffee")
	if err != nil {
		t.Fatalf("GetTxnByDesc: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 results for 'coffee', got %d", len(ids))
	}
}

// ── SureName ─────────────────────────────────────────────────────────────────

func TestSureName_ValidName(t *testing.T) {
	setup(t)

	name, err := core.SureName("validName")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "validName" {
		t.Errorf("expected 'validName', got %q", name)
	}
}

func TestSureName_InvalidChars(t *testing.T) {
	setup(t)

	_, err := core.SureName("invalid name!")
	if err == nil {
		t.Error("expected error for invalid characters, got nil")
	}
}

// helper: skip test if CGo binary can't run (Windows Defender etc.)
func init() {
	if os.Getenv("FEHU_SKIP_INTEGRATION") == "1" {
		return
	}
}

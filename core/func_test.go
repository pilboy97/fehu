package core_test

import (
	"core"
	"testing"
)

func TestParsePeriod_ValidTimes(t *testing.T) {
	p, err := core.ParsePeriod("2024-01-01;00:00:00", "2024-12-31;23:59:59")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.St == nil || p.Ed == nil {
		t.Fatal("expected both St and Ed to be non-nil")
	}
	if *p.St >= *p.Ed {
		t.Errorf("expected St < Ed, got St=%d Ed=%d", *p.St, *p.Ed)
	}
}

func TestParsePeriod_InvalidStartTime(t *testing.T) {
	// 버그 수정 전: 에러를 무시하고 ts=0(1970년)을 사용했음
	_, err := core.ParsePeriod("not-a-date", "2024-12-31;23:59:59")
	if err == nil {
		t.Fatal("expected error for invalid start time, got nil")
	}
}

func TestParsePeriod_InvalidEndTime(t *testing.T) {
	_, err := core.ParsePeriod("2024-01-01;00:00:00", "not-a-date")
	if err == nil {
		t.Fatal("expected error for invalid end time, got nil")
	}
}

func TestParsePeriod_EmptyStrings(t *testing.T) {
	// 빈 문자열은 nil 포인터 (필터 없음)이어야 함
	p, err := core.ParsePeriod("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.St != nil || p.Ed != nil {
		t.Errorf("expected St=nil Ed=nil for empty strings, got St=%v Ed=%v", p.St, p.Ed)
	}
}

func TestParsePeriod_OnlyStart(t *testing.T) {
	p, err := core.ParsePeriod("2024-06-01;00:00:00", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.St == nil {
		t.Error("expected St to be non-nil")
	}
	if p.Ed != nil {
		t.Error("expected Ed to be nil")
	}
}

func TestParseTime_ValidFormat(t *testing.T) {
	ts, err := core.ParseTime("2024-01-15;12:30:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts <= 0 {
		t.Errorf("expected positive unix timestamp, got %d", ts)
	}
}

func TestParseTime_InvalidFormat(t *testing.T) {
	_, err := core.ParseTime("2024/01/15")
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

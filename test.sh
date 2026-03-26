#!/bin/bash

# 이전 테스트 DB 파일이 있다면 삭제하여 깨끗한 상태에서 시작
rm -f test_db.db

echo "--- Fehu 프로그램 검증 스크립트 시작 ---"
echo ""

# 1. DB 열기 및 통화 코드 설정
echo ">>> 1. DB 열기 및 통화 코드 설정 (KRW)"
echo ""

# 2. 계정 생성
echo ">>> 2. 계정 생성: asset:bank, income:salary, expense:food"
fehu.exe -d test_db -c "new acc asset:bank"
fehu.exe -d test_db -c "new acc income:salary"
fehu.exe -d test_db -c "new acc expense:food"
echo ""

# 3. 거래 기록 (급여 수령, 식비 지출)
echo ">>> 3. 거래 기록: income:salary<100000;expense:food>30000;asset:bank>70000"
fehu.exe -d test_db -c "new txn income:salary<100000;expense:food>30000;asset:bank>70000 -d='월급 및 식비'"
echo ""

# 4. 계정 목록 조회 및 잔액 확인
echo ">>> 4. 계정 목록 조회"
fehu.exe -d test_db -c "get acc"
echo ""

# 5. 거래 목록 조회
echo ">>> 5. 거래 목록 조회"
fehu.exe -d test_db -c "get txn"
echo ""

# 6. 계산기 기능 사용 (총 지출 계산)
echo ">>> 6. 계산기: 총 지출 계산 (sum(acc(__all__, 'expense*')))"
fehu.exe -d test_db -c "calc sum(acc(__all__, 'expense*'))"
echo ""

# 7. 변수 정의 및 사용
echo ">>> 7. 변수 정의 및 사용: def my_expense sum(acc(__all__, 'expense*'))"
fehu.exe -d test_db -c "def my_expense sum(acc(__all__, 'expense*'))"
fehu.exe -d test_db -c "calc my_expense"
echo ""

# 8. DB 닫기
echo ">>> 8. DB 닫기"
fehu.exe -d test_db -c "close"
echo ""

echo "--- Fehu 프로그램 검증 스크립트 완료 ---"

# 테스트 DB 파일 삭제
rm -f test_db.db

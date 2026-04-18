package main

import (
	"cli"
	"core"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
)

var env *Env
var BenchMode bool
var initDB string
var cmd string
var mcpMode bool

func init() {
	env = NewEnv()
}

func main() {
	//플래그 정의
	flag.BoolVar(&mcpMode, "mcp", false, "start mcp server")
	flag.StringVar(&initDB, "d", "", "start with opening db")
	flag.BoolVar(&BenchMode, "b", false, "print elapsed time")
	flag.StringVar(&cmd, "c", "", "execute command")
	flag.StringVar(&core.Code, "CODE", "USD", "set currency code")
	flag.Parse()

	// 입력된 통화 코드가 기본 목록에 없다면, 암호화폐나 커스텀 화폐로 간주하고 소수점 8자리로 자동 등록합니다.
	if money.GetCurrency(core.Code) == nil {
		money.AddCurrency(core.Code, core.Code+" ", "1 $", ".", ",", 8)
	}

	// DB 자동 오픈 플래그가 있으면 MCP 서버 시작 전에도 DB를 엽니다.
	if len(initDB) != 0 {
		if err := core.Open(initDB + ".db"); err != nil {
			log.Fatalf("failed to open database: %v", err)
		}
	}

	// "mcp" 인자가 전달된 경우 일반 CLI/REPL 로직을 무시하고 MCP 서버만 단독 실행
	if mcpMode {
		if err := StartMCPServer(); err != nil {
			log.Fatalf("MCP Server error: %v", err)
		}
		return
	}

	// CLI 정의
	var CLI = cli.NewCLI(func(cmd string) error {
		//명령어가 입력될때

		//시작 시점 기록
		st := time.Now()

		//명령어 해석
		res, err := parser.Parse(cmd)
		if err != nil {
			return err
		}

		//명령어 실행
		err = Proc(res)

		if err != nil {
			return err
		}

		//종료 시점 기록
		ed := time.Now()
		//만약 벤치 모드라면
		if BenchMode {
			e := ed.Sub(st)

			// 실행시간 출력
			fmt.Printf("%d (ms) elapsed\n", e.Milliseconds())
		}

		return nil
	})
	CLI.Prefix = func() string {
		return DBName() + " > "
	}

	if len(cmd) != 0 {
		if err := CLI.Exec(cmd); err != nil && err != cli.ErrShutdownSystem {
			log.Print(err)
		}
	}

	// 루프 실행
	println("Fehu started")
	var isAlive = true
	for isAlive {
		func() {
			defer func() {
				// 예상치 못한 panic 처리
				if r := recover(); r != nil {
					log.Print(r)
				}
			}()

			if err := CLI.Run(os.Stdin); err == cli.ErrShutdownSystem {
				isAlive = false
			} else if err != nil {
				log.Print(err)
			}
		}()
	}

	//프로그램 종료라면
	if len(env.Path()) > 0 {
		//DB 닫기
		core.DB.Close()
	}
}

// DBName returns the base name of the currently opened database file without the extension.
func DBName() string {
	path := filepath.Base(env.Path())
	return strings.SplitN(path, ".", 2)[0]
}

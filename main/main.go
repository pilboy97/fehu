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
)

var env *Env
var BenchMode bool
var initDB string
var cmd string

func init() {
	env = NewEnv()
}

func main() {
	//플래그 정의
	flag.StringVar(&initDB, "d", "", "start with opening db")
	flag.BoolVar(&BenchMode, "b", false, "print elapsed time")
	flag.StringVar(&cmd, "c", "", "execute command")
	flag.StringVar(&core.Code, "CODE", "KRW", "set currency code")
	flag.Parse()

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

	if len(initDB) != 0 {
		core.Open(initDB + ".db")
	}

	if len(cmd) != 0 {
		defer func() {
			if r := recover(); r != nil {
				log.Print(r)
			}
		}()

		CLI.Exec(cmd)
	}

	// 루프 실행
	println("Fehu started")
	var isAlive = true
	// 프로그램 활성화 상태
	for isAlive {
		func() bool {
			// 오류가 발생했다면 false를 리턴
			defer func() {
				//예외 처리
				if r := recover(); r != nil {
					if e, ok := r.(error); ok {
						//만약 프로그램 종료 예외라면
						if e == cli.ErrShutdownSystem {
							isAlive = false
							// 프로그램 실행상태 변경 후 종료
							return
						}
					}
					log.Print(r)
				}
			}()

			//표준 입력으로 입력받은 명령어 실행
			CLI.Run(os.Stdin)
			//정상 종료
			return true
		}()
	}

	//프로그램 종료라면
	if len(env.Path()) > 0 {
		//DB 닫기
		core.DB.Close()
	}
}
func DBName() string {
	path := filepath.Base(env.Path())
	return strings.SplitN(path, ".", 2)[0]
}

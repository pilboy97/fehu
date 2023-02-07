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
	flag.StringVar(&initDB, "d", "", "start with opening db")
	flag.BoolVar(&BenchMode, "b", false, "print elapsed time")
	flag.StringVar(&cmd, "c", "", "execute command")
	flag.StringVar(&core.Code, "CODE", "KRW", "set currency code")
	flag.Parse()

	var CLI = cli.NewCLI(func(cmd string) error {
		st := time.Now()
		res, err := parser.Parse(cmd)
		if err != nil {
			return err
		}

		err = Proc(res)

		if err != nil {
			return err
		}
		ed := time.Now()
		if BenchMode {
			e := ed.Sub(st)

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

		err := CLI.OnCmd(cmd)
		if err != nil {
			if err == cli.ErrShutdownSystem {
				return
			}
			panic(err)
		}
		return
	}

	println("Fehu started")
	for {
		ok := func() bool {
			defer func() {
				if r := recover(); r != nil {
					log.Print(r)
				}
			}()

			err := CLI.Run(os.Stdin)
			if err != nil {
				if err == cli.ErrShutdownSystem {
					return true
				}
				panic(err)
			}
			return true
		}()

		if ok {
			break
		}
	}

	if len(env.Path()) > 0 {
		core.DB.Close()
	}
}
func DBName() string {
	path := filepath.Base(env.Path())
	return strings.SplitN(path, ".", 2)[0]
}

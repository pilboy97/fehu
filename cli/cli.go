package cli

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

var ErrShutdownSystem = errors.New("system is shutdowned")

type CLI struct {
	ch chan struct{}

	Prefix func() string
	OnCmd  func(cmd string) error
}

func NewCLI(f func(cmd string) error) *CLI {
	return &CLI{
		OnCmd: f,
	}
}
func (c *CLI) IsAlive() bool {
	select {
	case <-c.ch:
		return false
	default:
		return true
	}
}
func (c *CLI) Done() {
	if !c.IsAlive() {
		return
	}

	close(c.ch)
}

func (cli *CLI) Run(r io.Reader) {
	cli.ch = make(chan struct{})
	defer cli.Done()

	stdin := bufio.NewScanner(r)

	fmt.Print(cli.Prefix())

	for cli.IsAlive() {
		for stdin.Scan() {
			cmd := stdin.Text()
			cli.Exec(cmd)

			fmt.Print(cli.Prefix())
		}
	}
}
func (cli *CLI) Exec(str string) {
	err := cli.OnCmd(str)
	if err != nil {
		if err == ErrShutdownSystem {
			cli.Done()
		}
		panic(err)
	}
}

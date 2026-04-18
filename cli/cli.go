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

func (cli *CLI) Run(r io.Reader) error {
	cli.ch = make(chan struct{})
	defer cli.Done()

	stdin := bufio.NewScanner(r)

	fmt.Print(cli.Prefix())

	for cli.IsAlive() {
		for stdin.Scan() {
			cmd := stdin.Text()
			if err := cli.Exec(cmd); err != nil {
				return err
			}
			fmt.Print(cli.Prefix())
		}
	}
	return nil
}

func (cli *CLI) Exec(str string) error {
	err := cli.OnCmd(str)
	if err != nil {
		if err == ErrShutdownSystem {
			cli.Done()
		}
		return err
	}
	return nil
}

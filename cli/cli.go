package cli

import (
	"bufio"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

var ErrShutdownSystem = errors.New("system is shutdowned")

type CLI struct {
	isAlive bool

	Prefix func() string
	OnCmd  func(cli *CLI, cmd string) error
}

func NewCLI(f func(cli *CLI, cmd string) error) *CLI {
	return &CLI{
		OnCmd: f,
	}
}
func (c *CLI) IsAlive() bool {
	return c.isAlive
}
func (c *CLI) Done() {
	c.isAlive = false
}

func (cli *CLI) Run(r io.Reader) error {
	var err error

	cli.isAlive = true
	defer cli.Done()

	stdin := bufio.NewScanner(r)

	fmt.Print(cli.Prefix())

	for cli.IsAlive() {
		for stdin.Scan() {
			cmd := stdin.Text()
			if err = cli.OnCmd(cli, cmd); err != nil {
				if err == ErrShutdownSystem {
					cli.Done()
				}
				return err
			}

			fmt.Print(cli.Prefix())
		}
	}
	return nil
}

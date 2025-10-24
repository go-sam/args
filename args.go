package args

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-sam/utils"
)

type Parser struct {
	args      []string
	validArgs map[string]bool
}

func New() *Parser {
	return &Parser{
		args:      os.Args[1:],
		validArgs: make(map[string]bool),
	}
}

func (p *Parser) registerArgs(short, long string) {
	p.validArgs["-"+short] = true
	p.validArgs["--"+long] = true
}

func (p *Parser) matchFlag(index int, short, long string) bool {
	if index >= len(p.args) {
		return false
	}
	arg := p.args[index]
	return arg == "-"+short || arg == "--"+long
}

func (p *Parser) String(short, long string, target *string) bool {
	p.registerArgs(short, long)

	for i := range p.args {
		if p.matchFlag(i, short, long) {
			if i+1 < len(p.args) {
				*target = p.args[i+1]
				return true
			}
		}
	}
	return false
}

func (p *Parser) Bool(short, long string, target *bool) {
	p.registerArgs(short, long)

	for i := range p.args {
		if p.matchFlag(i, short, long) {
			*target = true
			return
		}
	}
}

func (p *Parser) Integer(short, long string, target *int) bool {
	p.registerArgs(short, long)

	for i := range p.args {
		if p.matchFlag(i, short, long) {
			if i+1 < len(p.args) {
				if value, err := utils.TryParseInt(p.args[i+1]); err {
					*target = value
					return true
				}
			}
		}
	}
	return false
}

func (p *Parser) GetStringValue(short, long string) (string, bool) {
	p.registerArgs(short, long)

	for i := range p.args {
		if p.matchFlag(i, short, long) {
			if i+1 < len(p.args) {
				return p.args[i+1], true
			}
		}
	}
	return "", false
}

func (p *Parser) HasFlag(short, long string) bool {
	p.registerArgs(short, long)

	for i := range p.args {
		if p.matchFlag(i, short, long) {
			return true
		}
	}
	return false
}

func (p *Parser) Help() bool {
	return p.HasFlag("h", "help")
}

func (p *Parser) ValidateArgs() error {
	for _, arg := range p.args {
		if !strings.HasPrefix(arg, "-") {
			continue
		}

		if !p.validArgs[arg] {
			return fmt.Errorf("unknown argument: %s", arg)
		}
	}

	return nil
}

package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/twpayne/chezmoi/lib/chezmoi"
	"github.com/twpayne/go-vfs"
)

var lastpassCommand = &cobra.Command{
	Use:   "lastpass",
	Short: "Execute the LastPass CLI",
	RunE:  makeRunE(config.runLastPassCommand),
}

type LastPassCommandConfig struct {
	Lpass string
}

func init() {
	rootCommand.AddCommand(lastpassCommand)
	config.LastPass.Lpass = "lpass"
	config.addFunc("lastpass", config.lastpassFunc)
}

func (c *Config) runLastPassCommand(fs vfs.FS, cmd *cobra.Command, args []string) error {
	return c.exec(append([]string{c.LastPass.Lpass}, args...))
}

func (c *Config) lastpassFunc(id string) interface{} {
	name := c.LastPass.Lpass
	args := []string{"show", "-j", id}
	if c.Verbose {
		fmt.Printf("%s %s\n", name, strings.Join(args, " "))
	}
	output, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		chezmoi.ReturnTemplateFuncError(fmt.Errorf("lastpass %q: %s show -j %q: %v\n%s", id, name, id, err, output))
	}
	var data []map[string]interface{}
	if err := json.Unmarshal(output, &data); err != nil {
		chezmoi.ReturnTemplateFuncError(fmt.Errorf("lastpass %q: %s show -j %q: %v", id, name, id, err))
	}
	return data
}

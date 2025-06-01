package configurator

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/google/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/zerops-dev/di/app"
	"github.com/zerops-dev/di/di/vflags"
	"github.com/zerops-dev/di/env"
)

type Configurator interface {
	Register(string, interface{})
}

var InitSet = wire.NewSet(
	wire.Bind(new(Configurator), new(Handler)),
	New,
)

type setup struct {
	RegisterFlags bool
}

func WithRegisterFlags(in bool) func(*setup) {
	return func(s *setup) {
		s.RegisterFlags = in
	}
}

func New(command *cobra.Command, option ...func(*setup)) *Handler {
	s := setup{
		RegisterFlags: true,
	}
	for _, option := range option {
		option(&s)
	}

	h := &Handler{
		command: command,
		viper:   viper.New(),
		setup:   s,
	}
	h.command.Flags().StringVarP(&h.config, "config", "c", "", `configuration file`)
	_ = h.command.Flags().MarkHidden("config")

	return h
}

type Handler struct {
	flags      []vflags.Value
	command    *cobra.Command
	viper      *viper.Viper
	config     string
	usedConfig string
	envFiles   []string
	setup      setup
}

type Flag struct {
	Name  string
	Value *string
}

type Test struct {
	Value []string
}

func (h *Handler) Write(filename string) error {
	return h.viper.WriteConfigAs(filename)
}

func (h *Handler) AddEnvFile(filename string) {
	h.envFiles = append(h.envFiles, filename)
}

func (h *Handler) Read(app *app.ApplicationSetup, cmd *cobra.Command) error {
	env.New()
	env.Load(h.envFiles...)
	h.viper.AutomaticEnv()
	if err := h.viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	h.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	wd, err := os.Getwd()
	if err == nil {
		h.viper.AddConfigPath(wd)
	}
	if app.ConfigPath == "" {
		h.viper.AddConfigPath(path.Join("/etc", app.Service))
	} else {
		h.viper.AddConfigPath(app.ConfigPath)
	}
	h.viper.SetConfigName("config")
	if h.viper.IsSet("config") {
		h.viper.SetConfigFile(h.viper.GetString("config"))
	}
	h.viper.SetTypeByDefaultValue(true)
	h.viper.ReadInConfig()
	h.usedConfig = h.viper.ConfigFileUsed()
	h.readConfigDir()
	for _, f := range h.flags {
		if err := f.SetValue(h.viper.Get(f.Name)); err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) readConfigDir() {
	if h.usedConfig == "" {
		return
	}
	configDir := h.usedConfig + ".d"
	dir, err := os.Stat(configDir)
	if err != nil {
		return
	}
	if !dir.IsDir() {
		return
	}
	rDir, err := os.ReadDir(configDir)
	if err != nil {
		return
	}
	sort.Slice(rDir, func(i, j int) bool {
		return rDir[i].Name() < rDir[j].Name()
	})
	for _, file := range rDir {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".yml") {
			continue
		}
		fmt.Fprintf(os.Stderr, "used config: %s\n", path.Join(configDir, file.Name()))
		h.viper.SetConfigFile(path.Join(configDir, file.Name()))
		h.viper.MergeInConfig()
	}

}

func (h *Handler) UsedConfigFile() string {
	return h.usedConfig
}

func (h *Handler) Register(prefix string, cfg interface{}) {
	flags := h.command.Flags()
	for _, f := range vflags.Parse(prefix, cfg) {
		if prefix == "" || h.setup.RegisterFlags {
			f.AppendToFlags(flags)
		}
		h.flags = append(h.flags, f)
		h.viper.SetDefault(f.Name, f.Valuer.Value())
	}
}

func (h *Handler) ConfigCommand(app *app.ApplicationSetup) *cobra.Command {
	file := "-"
	var actual bool
	cmd := &cobra.Command{
		Use:   "config",
		Short: "print configuration",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			if actual {
				if err := h.Read(app, cmd.Parent()); err != nil {
					return err
				}
			}
			out := os.Stdout
			if file != "-" {
				var err error
				out, err = os.Create(file)
				if err != nil {
					return err
				}
			}
			data := make(map[string]interface{})
			for _, v := range h.flags {
				if !v.Configurable {
					continue
				}
				if v.Virtual {
					continue
				}
				SetToYaml(strings.Split(v.Name, "."), data, v.Valuer.Value())
			}
			delete(data, "file")
			delete(data, "actual")
			delete(data, "help")
			delete(data, "config")
			return yaml.NewEncoder(out).Encode(data)
		},
	}
	cmd.Hidden = true
	cmd.Flags().StringVarP(&h.config, "config", "c", "", `configuration file`)
	cmd.Flags().StringVarP(&file, "file", "f", file, `output to file, use "-" for stdout`)
	cmd.Flags().BoolVarP(&actual, "actual", "a", actual, `output actual configuration`)
	return cmd
}

func SetToYaml(parts []string, storage map[string]interface{}, value interface{}) {
	if len(parts) <= 1 {
		storage[parts[0]] = value
		return
	}

	if _, e := storage[parts[0]]; !e {
		storage[parts[0]] = make(map[string]interface{})
	}
	SetToYaml(parts[1:], storage[parts[0]].(map[string]interface{}), value)
}

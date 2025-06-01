package app

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/google/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Version string
type Name string
type Service string
type Description string
type DescriptionLong string
type Exec string
type Cancel context.CancelFunc

func (c Cancel) CancelFunc() context.CancelFunc {
	return context.CancelFunc(c)
}

type Context context.Context
type BuildTime time.Time
type StartTime time.Time
type Commit string

const EnvPanicFork = "VSH_PANIC_FORK"

type ApplicationSetup struct {
	Name            string
	Version         string
	Commit          string
	BuildTime       time.Time
	StartTime       time.Time
	Cancel          context.CancelFunc
	Context         context.Context
	Exec            string
	Description     string
	DescriptionLong string
	Service         string
	ConfigPath      string
	rootCmd         *cobra.Command
	viper           *viper.Viper
}

func New(ctx context.Context, name string) *ApplicationSetup {
	cancelCtx, cancelFunc := context.WithCancel(ctx)
	return &ApplicationSetup{
		Name:      name,
		StartTime: time.Now(),
		BuildTime: time.Now(),
		Cancel:    cancelFunc,
		Context:   cancelCtx,
		viper:     viper.New(),
	}
}

func (a *ApplicationSetup) SetBuildTime(in string) {
	i, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		i = 0
	}
	a.BuildTime = time.Unix(i, 0)
}

func (a *ApplicationSetup) Run() {
	if err := a.rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func (a *ApplicationSetup) RootCommand() *cobra.Command {
	if a.rootCmd == nil {
		a.rootCmd = &cobra.Command{
			Use:   a.Name,
			Short: a.Description,
			Long:  a.DescriptionLong,
		}
	}
	return a.rootCmd
}

func (a *ApplicationSetup) ForkForPanicLogging() bool {
	_, exists := os.LookupEnv(EnvPanicFork)
	return !exists
}

func NewBuildTime(app *ApplicationSetup) BuildTime {
	return BuildTime(app.BuildTime)
}

func NewStartTime(app *ApplicationSetup) StartTime {
	return StartTime(app.StartTime)
}

func NewVersion(app *ApplicationSetup) Version {
	return Version(app.Version)
}

func NewName(app *ApplicationSetup) Name {
	return Name(app.Name)
}

func NewCancel(app *ApplicationSetup) Cancel {
	return Cancel(app.Cancel)
}

func NewContext(app *ApplicationSetup) Context {
	return app.Context
}

func NewDescriptionLong(app *ApplicationSetup) DescriptionLong {
	return DescriptionLong(app.DescriptionLong)
}

func NewDescription(app *ApplicationSetup) Description {
	return Description(app.Description)
}

func NewService(app *ApplicationSetup) Service {
	return Service(app.Service)
}

func NewExec(app *ApplicationSetup) Exec {
	return Exec(app.Exec)
}

func NewCommit(app *ApplicationSetup) Commit {
	return Commit(app.Commit)
}

var Set = wire.NewSet(
	NewService,
	NewDescription,
	NewDescriptionLong,
	NewExec,
	NewBuildTime,
	NewStartTime,
	NewVersion,
	NewContext,
	NewName,
	NewCancel,
	NewCommit,
)

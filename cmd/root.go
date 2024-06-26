package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/licat233/genzero/config"
	"github.com/licat233/genzero/global"
	"github.com/licat233/genzero/modules/api"
	"github.com/licat233/genzero/modules/logic"
	"github.com/licat233/genzero/modules/model"
	"github.com/licat233/genzero/modules/pb"
	"github.com/licat233/genzero/tools"
	"github.com/licat233/genzero/update"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version number of " + config.ProjectName,
	Run: func(cmd *cobra.Command, args []string) {
		tools.Success("current version: " + config.CurrentVersion)
	},
}

var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"up", "u"},
	Short:   "Upgrade " + config.ProjectName + " to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		update.Update()
	},
}

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Create the default " + config.ProjectName + " configuration file in the current directory",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s requires at least one argument", cmd.CommandPath())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var initConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"c"},
	Short:   "Create the default " + config.ProjectName + " configuration file in the current directory, or specified directory",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.C.CreateYaml(); err != nil {
			tools.Warning(err.Error())
			os.Exit(1)
		}
		tools.Success("Done.")
	},
}

var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate model code",
	Run: func(cmd *cobra.Command, args []string) {
		config.C.Model.Status = true
		if err := Initialize(); err != nil {
			tools.Error("Failed to initialize: " + err.Error())
			os.Exit(1)
		}
		if err := model.New().Run(); err != nil {
			tools.Warning(err.Error())
			os.Exit(1)
		}
		tools.Success("Done.")
	},
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Generate .api files",
	Run: func(cmd *cobra.Command, args []string) {
		config.C.Api.Status = true
		if err := Initialize(); err != nil {
			tools.Error("Failed to initialize: " + err.Error())
			os.Exit(1)
		}
		if err := api.New().Run(); err != nil {
			tools.Warning(err.Error())
			os.Exit(1)
		}
		tools.Success("Done.")
	},
}

var pbCmd = &cobra.Command{
	Use:     "pb",
	Aliases: []string{"proto"},
	Short:   "Generate .proto files",
	Run: func(cmd *cobra.Command, args []string) {
		config.C.Pb.Status = true
		if err := Initialize(); err != nil {
			tools.Error("Failed to initialize: " + err.Error())
			os.Exit(1)
		}
		if err := pb.New().Run(); err != nil {
			tools.Warning(err.Error())
			os.Exit(1)
		}
		tools.Success("Done.")
	},
}

var logicCmd = &cobra.Command{
	Use:   "logic",
	Short: "Modify logic files, this feature has not been developed yet",
	Run: func(cmd *cobra.Command, args []string) {
		config.C.Logic.Status = true
		if err := Initialize(); err != nil {
			tools.Error("Failed to initialize: " + err.Error())
			os.Exit(1)
		}
		if err := logic.New().Run(); err != nil {
			tools.Warning(err.Error())
			os.Exit(1)
		}
		tools.Success("Done.")
	},
}

var apilogicCmd = &cobra.Command{
	Use:   "api",
	Short: "Modify api logic files",
	Run: func(cmd *cobra.Command, args []string) {
		config.C.Logic.Api.Status = true
		if err := Initialize(); err != nil {
			tools.Error("Failed to initialize: " + err.Error())
			os.Exit(1)
		}
		if err := logic.New().Run(); err != nil {
			tools.Warning(err.Error())
			os.Exit(1)
		}
		tools.Success("Done.")
	},
}

var rpclogicCmd = &cobra.Command{
	Use:   "rpc",
	Short: "Modify rpc logic files",
	Run: func(cmd *cobra.Command, args []string) {
		config.C.Logic.Rpc.Status = true
		if err := Initialize(); err != nil {
			tools.Error("Failed to initialize: " + err.Error())
			os.Exit(1)
		}
		if err := logic.New().Run(); err != nil {
			tools.Warning(err.Error())
			os.Exit(1)
		}
		tools.Success("Done.")
	},
}

var startCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"run"},
	Short:   "Use yaml file configuration to start " + config.ProjectName,
	Run: func(cmd *cobra.Command, args []string) {
		config.UseConf = true
		runByConf()
	},
}

var rootCmd = &cobra.Command{
	Use:        config.ProjectName,
	Aliases:    []string{},
	SuggestFor: []string{},
	Short:      "This is a tool to generate gozero service based on mysql",
	GroupID:    "",
	Long:       fmt.Sprintf("This is a tool to generate gozero service based on mysql.\nThe goctl tool must be installed before use.\ncurrent version: %s\nGithub: https://github.com/licat233/genzero", config.CurrentVersion),
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("%s requires at least one argument", cmd.CommandPath())
		}
		return nil
	},
	Version: config.CurrentVersion,
	Run: func(cmd *cobra.Command, args []string) {
		// run()
	},
}

func runByConf() {
	if err := Initialize(); err != nil {
		tools.Error("Failed to initialize: " + err.Error())
		os.Exit(1)
	}
	tasks := make([]tools.TaskFunc, 0)
	if config.C.Model.Status {
		tasks = append(tasks, func() error {
			return model.New().Run()
		})
	}

	if config.C.Pb.Status {
		tasks = append(tasks, func() error {
			return pb.New().Run()
		})
	}

	if config.C.Api.Status {
		tasks = append(tasks, func() error {
			return api.New().Run()
		})
	}
	if config.C.Logic.Status {
		tasks = append(tasks, func() error {
			return logic.New().Run()
		})
	}

	err := tools.RunConcurrentTasks(tasks)
	if err != nil {
		tools.Error(err.Error())
		os.Exit(1)
	}

	tools.Success("Done.")
}

var IsDev bool

func init() {
	config.C = config.New()

	initConfigCmd.PersistentFlags().StringVar(&config.InitConfSrc, "dir", config.DefaultConfigFileName, "file location for yaml configuration")

	rootCmd.PersistentFlags().BoolVar(&IsDev, "dev", false, "dev mode, print error message")
	rootCmd.PersistentFlags().StringVar(&config.ConfSrc, "conf", config.DefaultConfigFileName, "file location for yaml configuration")

	rootCmd.PersistentFlags().StringVar(&config.C.DB.DSN, "dsn", "", "data source name (DSN) to use when connecting to the database")
	rootCmd.PersistentFlags().StringVar(&config.C.DB.Src, "src", "", "sql file to use when connecting to the database")
	rootCmd.PersistentFlags().StringSliceVar(&config.C.DB.Tables, "tables", []string{}, "need to generate tables, default is all tables，split multiple value by ','")
	rootCmd.PersistentFlags().StringSliceVar(&config.C.DB.IgnoreTables, "ignore_tables", []string{}, "ignore table string, default is none，split multiple value by ','")
	rootCmd.PersistentFlags().StringSliceVar(&config.C.DB.IgnoreColumns, "ignore_columns", []string{}, "ignore column string, default is none，split multiple value by ','")

	pbCmd.PersistentFlags().StringVar(&config.C.Pb.FileStyle, "file_style", config.LowerCamelCase, "proto file naming style: "+config.StyleList)
	pbCmd.PersistentFlags().StringVar(&config.C.Pb.Package, "pkg", "", "proto package")
	pbCmd.PersistentFlags().StringVar(&config.C.Pb.GoPackage, "gopkg", "", "proto go package")
	pbCmd.PersistentFlags().BoolVar(&config.C.Pb.Multiple, "multiple", false, "proto multiple")
	pbCmd.PersistentFlags().StringVar(&config.C.Pb.Dir, "dir", "", "proto output directory")
	pbCmd.PersistentFlags().StringVar(&config.C.Pb.ServiceName, "service_name", "", "service name, default value is database name")
	pbCmd.PersistentFlags().StringSliceVar(&config.C.Pb.Tables, "tables", []string{}, "need to generate tables, default is all tables，split multiple value by ','")
	pbCmd.PersistentFlags().StringSliceVar(&config.C.Pb.IgnoreTables, "ignore_tables", []string{}, "ignore table string, default is none，split multiple value by ','")
	pbCmd.PersistentFlags().StringSliceVar(&config.C.Pb.IgnoreColumns, "ignore_columns", []string{}, "ignore column string, default is none，split multiple value by ','")

	apiCmd.PersistentFlags().StringVar(&config.C.Api.JsonStyle, "json_style", config.SnakeCase, "JSON field naming style: "+config.StyleList)
	apiCmd.PersistentFlags().StringVar(&config.C.Api.Jwt, "jwt", "", "api jwt")
	apiCmd.PersistentFlags().StringSliceVar(&config.C.Api.Middleware, "middleware", []string{}, "api middleware")
	apiCmd.PersistentFlags().StringVar(&config.C.Api.Prefix, "prefix", "", "api prefix")
	apiCmd.PersistentFlags().BoolVar(&config.C.Api.Multiple, "multiple", false, "api multiple")
	apiCmd.PersistentFlags().StringVar(&config.C.Api.Dir, "dir", "", "api output directory")
	apiCmd.PersistentFlags().StringVar(&config.C.Api.ServiceName, "service_name", "", "service name, default value is database name")
	apiCmd.PersistentFlags().StringSliceVar(&config.C.Api.Tables, "tables", []string{}, "need to generate tables, default is all tables，split multiple value by ','")
	apiCmd.PersistentFlags().StringSliceVar(&config.C.Api.IgnoreTables, "ignore_tables", []string{}, "ignore table string, default is none，split multiple value by ','")
	apiCmd.PersistentFlags().StringSliceVar(&config.C.Api.IgnoreColumns, "ignore_columns", []string{}, "ignore column string, default is none，split multiple value by ','")

	modelCmd.PersistentFlags().StringVar(&config.C.Model.Dir, "dir", "", "model output directory")
	modelCmd.PersistentFlags().StringVar(&config.C.Api.ServiceName, "service_name", "", "service name, default value is database name")
	modelCmd.PersistentFlags().StringSliceVar(&config.C.Api.Tables, "tables", []string{}, "need to generate tables, default is all tables，split multiple value by ','")
	modelCmd.PersistentFlags().StringSliceVar(&config.C.Api.IgnoreTables, "ignore_tables", []string{}, "ignore table string, default is none，split multiple value by ','")
	modelCmd.PersistentFlags().StringSliceVar(&config.C.Api.IgnoreColumns, "ignore_columns", []string{}, "ignore column string, default is none，split multiple value by ','")

	apilogicCmd.PersistentFlags().BoolVar(&config.C.Logic.Api.UseRpc, "use_rpc", false, "use rpc for api")
	// apilogicCmd.PersistentFlags().BoolVar(&config.C.Logic.Api.RpcMultiple, "rpc_multiple", false, "is multiple rpc ?")
	apilogicCmd.PersistentFlags().StringVar(&config.C.Logic.Api.FileStyle, "file_style", config.LowerCamelCase, "file naming style: "+config.StyleList)
	apilogicCmd.PersistentFlags().StringVar(&config.C.Logic.Api.Dir, "dir", "", "api logic directory")
	apilogicCmd.PersistentFlags().StringSliceVar(&config.C.Logic.Api.Tables, "tables", []string{}, "need to generate tables, default is all tables，split multiple value by ','")
	apilogicCmd.PersistentFlags().StringSliceVar(&config.C.Logic.Api.IgnoreTables, "ignore_tables", []string{}, "ignore table string, default is none，split multiple value by ','")
	// apilogicCmd.PersistentFlags().StringSliceVar(&config.C.Logic.Api.IgnoreColumns, "ignore_columns", []string{}, "ignore column string, default is none，split multiple value by ','")

	rpclogicCmd.PersistentFlags().BoolVar(&config.C.Logic.Rpc.Multiple, "multiple", false, "is multiple ?")
	rpclogicCmd.PersistentFlags().StringVar(&config.C.Logic.Rpc.FileStyle, "file_style", config.LowerCamelCase, "file naming style: "+config.StyleList)
	rpclogicCmd.PersistentFlags().StringVar(&config.C.Logic.Rpc.Dir, "dir", "", "rpc logic directory")
	rpclogicCmd.PersistentFlags().StringSliceVar(&config.C.Logic.Rpc.Tables, "tables", []string{}, "need to generate tables, default is all tables，split multiple value by ','")
	// rpclogicCmd.PersistentFlags().StringSliceVar(&config.C.Logic.Rpc.IgnoreColumns, "ignore_columns", []string{}, "ignore column string, default is none，split multiple value by ','")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(pbCmd)
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(modelCmd)
	rootCmd.AddCommand(startCmd)

	logicCmd.AddCommand(apilogicCmd)
	logicCmd.AddCommand(rpclogicCmd)

	rootCmd.AddCommand(logicCmd)

	initCmd.AddCommand(initConfigCmd)
}

func Initialize() error {
	if config.ConfSrc == "." {
		config.ConfSrc = config.DefaultConfigFileName
		config.UseConf = true
	}
	if !config.UseConf {
		config.UseConf = strings.HasSuffix(config.ConfSrc, ".yaml") || strings.HasSuffix(config.ConfSrc, ".yam")
	}
	if config.UseConf {
		if err := config.C.ConfigureByYaml(); err != nil {
			return fmt.Errorf("read config faild: %v", err)
		}
	}

	if err := config.C.Validate(); err != nil {
		return err
	}

	if err := global.InitSchema(); err != nil {
		return err
	}
	return nil
}

func Execute() {
	defer func() {
		if !IsDev {
			if err := recover(); err != nil {
				tools.Warning(fmt.Sprintf("%v", err))
			}
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		tools.Warning(err.Error())
		os.Exit(1)
	}
}

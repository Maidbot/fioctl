package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var targetsTagCmd = &cobra.Command{
	Use:   "tag <target> [<target>...]",
	Short: "Apply a comman separated list of tags to one or more targets.",
	Run:   doTargetsTag,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	targetsCmd.AddCommand(targetsTagCmd)
	targetsTagCmd.PersistentFlags().StringP("tags", "T", "", "comma,separate,list")
	targetsTagCmd.PersistentFlags().BoolP("no-tail", "", false, "Don't tail output of CI Job")

	targetsCmd.MarkPersistentFlagRequired("tags")

	if err := viper.BindPFlags(targetsTagCmd.PersistentFlags()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func doTargetsTag(cmd *cobra.Command, args []string) {
	factory := viper.GetString("factory")
	tags := strings.Split(viper.GetString("tags"), ",")

	targets, err := api.TargetsList(factory)
	if err != nil {
		fmt.Print("ERROR: ")
		fmt.Println(err)
		os.Exit(1)
	}

	for idx := range args {
		if target, ok := targets.Signed.Targets[args[idx]]; ok {
			custom, err := api.TargetCustom(target)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("Changing tags of %s from %s -> %s\n", args[idx], custom.Tags, tags)
		} else {
			fmt.Printf("Target(%s) not found in targets.json\n", args[idx])
			os.Exit(1)
		}
	}

	url, err := api.TargetUpdateTags(factory, args, tags)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("CI URL: %s\n", url)
	if !viper.GetBool("no-tail") {
		api.JobservTail(url)
	}
}
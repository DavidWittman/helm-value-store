package cmd

import (
	"errors"
	"fmt"

	"github.com/skuid/helm-value-store/dynamo"
	"github.com/skuid/helm-value-store/store"
	"github.com/spf13/cobra"
)

type getCmdArgs struct {
	table  string
	labels selectorSet
	name   string
	uuid   string
}

var getArgs = getCmdArgs{}

var getCmd = &cobra.Command{
	Use:   "get-values",
	Short: "get the values of a release",
	Run:   get,
}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&getArgs.table, "table", "helm-charts", "Name of table")
	getCmd.Flags().StringVar(&getArgs.uuid, "uuid", "", "The UUID to get.")
	getCmd.Flags().VarP(&getArgs.labels, "label", "l", `The labels to filter by. Each label should have the format "k=v".
    	Can be specified multiple times, or a comma-separated get.`)
	getCmd.Flags().StringVar(&getArgs.name, "name", "", "The name of the release")

}

func hasReleases(releases store.Releases, message string) {
	if len(releases) == 0 {
		exitOnErr(errors.New(message))
	}
}

func get(cmd *cobra.Command, args []string) {
	rs, err := dynamo.NewReleaseStore(getArgs.table)
	exitOnErr(err)

	releases := store.Releases{}

	if len(getArgs.uuid) > 0 {
		release, err := rs.Get(getArgs.uuid)
		exitOnErr(err)
		releases = append(releases, *release)

	} else if len(getArgs.name) > 0 || len(getArgs.labels) > 0 {
		releases, err = rs.List(getArgs.labels.ToMap())
		exitOnErr(err)

		hasReleases(releases, "No releases match those labels!")

		if len(getArgs.name) > 0 {
			releases = filterByName(releases, getArgs.name)
		}
		hasReleases(releases, "No releases match that name and those labels!")

	} else {
		exitOnErr(errors.New("Must supply a UUID, release name, or labels!"))
	}

	for i, release := range releases {
		if i > 0 && i <= len(releases)-1 {
			fmt.Println("---")
		}
		fmt.Printf("# %s: %s, %s\n", release.Name, release.UniqueID, release.Labels)
		fmt.Print(release.Values)
	}

}

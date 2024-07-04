package main

import "docute/cmd"

func main() {
	root := cmd.RootCMD()
	root.AddCommand(cmd.GenerateCMD(), cmd.WatchCMD(), cmd.HostCMD())

	cmd.Execute(root)
}

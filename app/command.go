package app

func (a *ApplicationSetup) RegisterCommands() {
	a.RootCommand().AddCommand(
		CreateBashCmd(a),
		CreateBuildTimeCmd(a),
		CreateCommitCmd(a),
		CreateDocCmd(),
		CreateManCmd(a),
		CreateNameCmd(a),
		CreateVersionCmd(a),
	)
}

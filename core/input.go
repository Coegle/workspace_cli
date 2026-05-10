package core

type BaseInput struct {
	FeatBranchName string
	Cfg            *Config
}

func (b *BaseInput) GetWorkspaceName() string {
	return b.Cfg.GetSafeDirName(b.FeatBranchName)
}

type AddInput struct {
	BaseInput
	ServiceRepoNames []string
	BaseBranchName   string
}

type RemoveInput struct {
	BaseInput
	ServiceRepoNames []string
}

type DestroyInput struct {
	BaseInput
	Force bool
}

type OpenInput struct {
	BaseInput
}

type SyncInput struct {
	ServiceRepoNames []string
	BaseBranchName   string
	Cfg              *Config
}

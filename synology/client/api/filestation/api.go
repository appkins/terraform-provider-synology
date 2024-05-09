package filestation

type Api interface {
	CreateFolder(folderPath string, name string, forceParent bool) (*CreateFolderResponse, error)
	ListShares() (*ListShareResponse, error)
}

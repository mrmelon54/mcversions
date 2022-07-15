package structure

// PistonMetaPackage is used to parse the data from the api.
type PistonMetaPackage struct {
	Arguments              *PistonMetaPackageArguments   `json:"arguments"`
	AssetIndex             *PistonMetaPackageAssetIndex  `json:"assetIndex"`
	Assets                 string                        `json:"assets"`
	ComplianceLevel        int                           `json:"complianceLevel"`
	Downloads              *PistonMetaPackageDownloads   `json:"downloads"`
	ID                     string                        `json:"id"`
	JavaVersion            *PistonMetaPackageJavaVersion `json:"javaVersion"`
	Libraries              []*PistonMetaPackageLibrary   `json:"libraries"`
	MainClass              string                        `json:"mainClass"`
	MinimumLauncherVersion int                           `json:"minimumLauncherVersion"`
	ReleaseTime            string                        `json:"releaseTime"`
	Time                   string                        `json:"time"`
	Type                   string                        `json:"type"`
}

// PistonMetaPackageArguments is used to store the game and jvm arguments
type PistonMetaPackageArguments struct {
	Game []any `json:"game"`
	Jvm  []any `json:"jvm"`
}

// PistonMetaPackageAssetIndex is used to store the asset hashes
type PistonMetaPackageAssetIndex struct {
	Id        string `json:"id"`
	Sha1      string `json:"sha1"`
	Size      int    `json:"size"`
	TotalSize int    `json:"totalSize"`
	Url       string `json:"url"`
}

// PistonMetaPackageDownloads is used to store the client and server download information
type PistonMetaPackageDownloads struct {
	Client         *PistonMetaPackageDownloadsData `json:"client"`
	ClientMappings *PistonMetaPackageDownloadsData `json:"client_mappings"`
	Server         *PistonMetaPackageDownloadsData `json:"server"`
	ServerMappings *PistonMetaPackageDownloadsData `json:"server_mappings"`
}

// PistonMetaPackageDownloadsData is used to store download objects.
type PistonMetaPackageDownloadsData struct {
	Sha1 string `json:"sha1"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
}

// PistonMetaPackageJavaVersion is used to store the valid Java version
type PistonMetaPackageJavaVersion struct {
	Component    string `json:"component"`
	MajorVersion int    `json:"majorVersion"`
}

// PistonMetaPackageLibrary is used to store the library data
type PistonMetaPackageLibrary struct {
	Downloads *PistonMetaPackageLibraryDownload `json:"downloads"`
	Name      string                            `json:"name"`
}

// PistonMetaPackageLibraryDownload is used to store the library download
type PistonMetaPackageLibraryDownload struct {
	Artifact *PistonMetaPackageLibraryDownloadArtifact `json:"artifact"`
}

// PistonMetaPackageLibraryDownloadArtifact is used to store the library download artifact
type PistonMetaPackageLibraryDownloadArtifact struct {
	Path string `json:"path"`
	Sha1 string `json:"sha1"`
	Size int    `json:"size"`
	Url  string `json:"url"`
}

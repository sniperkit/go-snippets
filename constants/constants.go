package constants

//TrueString string version of true boolean
const TrueString = "true"

// Subcommands supported by CLI
const BatchCommand = "batch"
const BuildCommand = "build"
const InitCommand = "init"
const ListCommand = "list"
const PublishCommand = "publish"
const PullCommand = "pull"
const RunCommand = "run"
const SearchCommand = "search"
const ValidateCommand = "validate"
const VersionCommand = "version"

//CacheFromFlag defines the docker cache-from option to utilize a previous built image
const CacheFromFlag = "cache-from"

//ShortCacheFromFlag defines the shorthand version of the cache-from option
const ShortCacheFromFlag = "c"

//JobDirectoryFlag defines the location of the seed spec and Dockerfile
const JobDirectoryFlag = "directory"

//ShortJobDirectoryFlag defines the shorthand location of the seed spec and Dockerfile
const ShortJobDirectoryFlag = "d"

//DockerfileFlag defines the dockerfile to use to build the image
const DockerfileFlag = "dockerfile"

//ShortDockerfileFlag defines the shorthand dockerfile to use to build the image
const ShortDockerfileFlag = "D"

//ManifestFlag defines the seed manifest file
const ManifestFlag = "manifest"

//ShortManifestFlag defines the shorthand manifest file
const ShortManifestFlag = "m"

//SettingFlag defines the SettingFlag
const SettingFlag = "setting"

//ShortSettingFlag defines the shorthand SettingFlag
const ShortSettingFlag = "e"

//MountFlag defines the MountFlag
const MountFlag = "mount"

//ShortMountFlag defines the shorthand MountFlag
const ShortMountFlag = "m"

//InputsFlag defines the InputFlag
const InputsFlag = "inputs"

//ShortInputsFlag defines the shorthand input flag
const ShortInputsFlag = "i"

//InputsFlag defines the InputFlag
const JsonFlag = "json"

//ShortInputsFlag defines the shorthand input flag
const ShortJsonFlag = "j"

//JobOutputDirFlag defines the job output directory
const JobOutputDirFlag = "outDir"

//ShortJobOutputDirFlag defines the shorthand output directory
const ShortJobOutputDirFlag = "o"

//ShortImgNameFlag defines image name to run
const ShortImgNameFlag = "in"

//ImgNameFlag defines image name to run
const ImgNameFlag = "imageName"

//RmFlag defines if the docker image should be removed after docker run is executed
const RmFlag = "rm"

//QuietFlag defines if output from the docker image being run should be suppressed
const QuietFlag = "quiet"

//QuietFlag shorthand flag that defines if output from the docker image being run should be suppressed
const ShortQuietFlag = "q"

//SchemaFlag defines a schema file to validate seed against
const SchemaFlag = "schema"

//ShortSchemaFlag shorthand flag that defines schema file to validate seed against
const ShortSchemaFlag = "s"

//RegistryFlag defines registry
const RegistryFlag = "registry"

//ShortRegistryFlag shorthand flag that defines registry
const ShortRegistryFlag = "r"

//OrgFlag defines organization
const OrgFlag = "org"

//ShortOrgFlag shorthand flag that defines organization
const ShortOrgFlag = "o"

//FilterFlag defines filter
const FilterFlag = "filter"

//ShortFilterFlag shorthand flag that defines filter
const ShortFilterFlag = "f"

//UserFlag defines user
const UserFlag = "user"

//ShortUserFlag shorthand flag that defines user
const ShortUserFlag = "u"

//PassFlag defines password
const PassFlag = "password"

//ShortPassFlag shorthand flag that defines password
const ShortPassFlag = "p"

//ForcePublishFlag forces a publish - don't try to deconflict
const ForcePublishFlag = "f"

//PkgVersionMinor specifies to bump package minor version
const PkgVersionMinor = "pm"

//PkgVersionMajor specifies to bump package major version
const PkgVersionMajor = "P"

//PkgVersionPatch specifies to bump package patch version
const PkgVersionPatch = "pp"

//JobVersionMinor specifies to bump job minor version
const JobVersionMinor = "jm"

//JobVersionMajor specifies to bump job major version
const JobVersionMajor = "J"

//JobVersionPatch specifies to bump job patch version
const JobVersionPatch = "jp"

//BatchFlag defines whether to run in batch mode
const BatchFlag = "batch"

//ShortBatchFlag - shorthand flag for batch
const ShortBatchFlag = "b"

//RepeatFlag defines how many times to run a docker image
const RepeatFlag = "repetitions"

//ShortRepeatFlag - shorthand flag for repetitions
const ShortRepeatFlag = "rep"

//VersionFlag defines version of seed spec to use
const VersionFlag = "version"

//ShortVersionFlag - shorthand flag for version
const ShortVersionFlag = "v"

//ResultsFileManifestName defines the filename for the results_manifest file
const ResultsFileManifestName = "seed.outputs.json"

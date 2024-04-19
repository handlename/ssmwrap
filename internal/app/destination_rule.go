package app

type DestinationType string

const (
	DestinationTypeEnv  DestinationType = "env"
	DestinationTypeFile DestinationType = "file"
)

type DestinationRule struct {
	Type DestinationType

	// To is address of destination.
	To string

	TypeEnvOptions  *DestinationTypeEnvOptions
	TypeFileOptions *DestinationTypeFileOptions
}

type DestinationTypeEnvOptions struct {
	// Prefix is a prefix for environment variable.
	// For example, if Prefix is PREFIX, then the environment variable name will be PREFIX_NAME.
	Prefix string

	// EntirePath is a flag to export entire path as environment variable name.
	// For example, if EntirePath is true and the path is /a/b/c, then the environment variable name will be A_B_C.
	// If EntirePath is false, then the environment variable name will be C.
	EntirePath bool
}

type DestinationTypeFileOptions struct {
	// Mode is a file mode of exported file.
	// If Mode is 0, then the default file mode is used defined in FileExporter.
	Mode int

	// Uid is a user id of exported file.
	// If Uid is 0, then the default user id is used defined in FileExporter.
	Uid int

	// Gid is a group id of exported file.
	// If Gid is 0, then the default group id is used defined in FileExporter.
	Gid int
}

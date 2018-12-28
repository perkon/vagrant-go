package vagrant_go

type Config struct {
	// BinaryName is the name of the vagrant executable that's going to be used. It must be present in $PATH.
	BinaryName string
}

func DefaultConfig() *Config {
	return &Config{
		BinaryName: "vagrant",
	}
}

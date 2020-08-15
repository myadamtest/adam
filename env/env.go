package env

type Env struct {
	OutDir       string
	ExcludeTable map[string]*struct{}
}

func GetEnv() Env {
	return Env{}
}

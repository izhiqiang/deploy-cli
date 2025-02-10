package hosts

type Host struct {
	Host       string `yml:"host"`
	Port       int    `yml:"port"`
	User       string `yml:"user"`
	Password   string `yml:"password"`
	PrivateKey string `yml:"private_key,omitempty"`
	Timeout    int    `yml:"timeout"`
}

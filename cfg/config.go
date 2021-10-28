package cfg

var Config Configure

type Configure struct {
	Debug bool
	Mysql string
	Proxy string
}

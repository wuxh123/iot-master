package types

type Project interface {
	//执行脚本
	Run(name string) error
}

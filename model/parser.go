package model


type _link struct {
	Name     string
	Protocol string
}

type _variable struct {
	Link     string
	Type     string
	Addr     string
	Path     string
	Default  string
	Writable bool //可写，用于输出（如开关）
}

type _batch struct {
	Name string //唯一
	Link string
	Type string
	Addr string
	Size int

	Results []struct {
		Offset int
		Path   string //Variable
	}
}

package install

import "io"

type Installer interface {
	SetOperator(insertFunc Operator)
	ExecReder(r io.Reader) error
	ExecFile(fileName string) error
	SetDrop(b bool)
}

type Operator interface {
	Insert(data string) error
	Update(data string) error
	Create(data string) error
	Drop(data string) error
	Select(data string) error
}

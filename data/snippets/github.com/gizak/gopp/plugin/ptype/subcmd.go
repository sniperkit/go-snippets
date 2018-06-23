package ptype

type SubcmdRunner interface {
	RunSubcmd(args []string) error
	Subcmd() string
	Usage() string
	Descrip() string
}

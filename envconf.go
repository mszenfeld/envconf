package envconf

import "strings"

type Loader struct {
	prefix string
}

func NewLoader() *Loader {
	return &Loader{}
}

func (l *Loader) Load(s interface{}) error {
	return nil
}

func (l *Loader) SetPrefix(p string) {
	l.prefix = strings.ToUpper(p)
}

func (l *Loader) Prefix() string {
	return l.prefix
}

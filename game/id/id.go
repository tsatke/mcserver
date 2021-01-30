package id

import "strings"

type ID [2]string

const defaultNamespace = "minecraft"

func ParseID(name string) ID {
	frags := strings.SplitN(name, ":", 2)
	if len(frags) == 1 {
		if frags[0] == "" {
			return ID{}
		}
		return ID{defaultNamespace, frags[0]}
	}
	return ID{frags[0], frags[1]}
}

func (id ID) String() string    { return id.Namespace() + ":" + id.Name() }
func (id ID) Namespace() string { return id[0] }
func (id ID) Name() string      { return id[1] }

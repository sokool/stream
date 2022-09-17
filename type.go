package stream

import (
	"fmt"
	"github.com/google/uuid"
	"path"
	"reflect"
	"regexp"
	"strings"
)

type Type string

func NewType(v any) (Type, error) {
	var s string

	t := reflect.TypeOf(v)
	s = t.Name()
	switch t.Kind() {
	case reflect.Pointer:
		s = t.Elem().Name()
	case reflect.String:
		s = v.(string)
	default:

	}

	if s = strings.ReplaceAll(strings.TrimSpace(s), " ", ""); len(s) == 0 {
		return "", Err("name can not be empty")
	}

	return Type(strings.Title(s)), nil
}

func (t Type) Hash() string {
	return uid(t.String()).String()
}

func (t Type) String() string {
	return string(t)
}

func (t Type) IsZero() bool {
	return t == ""
}

func (t Type) CutPrefix(of Type) Type {
	a, b := t.String(), of.String()
	if strings.Index(a, b) == -1 {
		return t
	}

	return Type(strings.Replace(a, b, "", 1))
}

type info struct {
	id, path, uuid, hash, name string
	pkg, typ                   string
	reflect                    struct {
		typ      reflect.Type
		value    reflect.Value
		instance any
	}
}

func infoOf(v any) *info {
	rv := reflect.ValueOf(v)
	rt := rv.Type()
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	var s = rt.String()
	var t = strings.Title(rt.Name())
	//if n, ok := v.(Namer); ok {
	//	s = n.Name()
	//}
	//
	var id string
	//if n, ok := v.(Identifier); ok {
	//	id = n.Id()
	//}
	//
	s = regexp.MustCompile(`[. _]`).ReplaceAllString(s, " ")
	s = strings.Title(s)
	s = strings.Replace(s, " ", "", -1)

	_, p := path.Split(rt.PkgPath())
	return &info{
		id:   id,
		path: rt.PkgPath() + "/" + rt.Name(),
		uuid: uuid.NewSHA1(uuid.NameSpaceDNS, []byte(rt.PkgPath()+"/"+rt.Name())).String(),
		hash: uuid.NewSHA1(uuid.NameSpaceDNS, []byte(fmt.Sprintf("%p", v))).String(),
		name: s,
		pkg:  strings.Title(p),
		typ:  t,
		reflect: struct {
			typ      reflect.Type
			value    reflect.Value
			instance any
		}{rt, rv, v},
	}
}

func (o info) String() string {
	return fmt.Sprintf(`
  id: %s
name: %s
path: %s
uuid: %s
hash: %s
`, o.id, o.name, o.path, o.uuid, o.hash)
}

func uid(s string) ID {
	return ID(uuid.NewSHA1(uuid.NameSpaceDNS, []byte(s)).String())
}

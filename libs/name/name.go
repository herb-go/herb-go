package name

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Name struct {
	Raw                         string
	Parents                     string
	ParentsList                 []string
	Title                       string
	Lower                       string
	Camel                       string
	Pascal                      string
	LowerWithParentDotSeparated string
	LowerWithParentPath         string
	PascalWithParents           string
}

func (n *Name) LowerPath(filename ...string) string {
	fmt.Println(path.Join(n.ParentsList...))
	return path.Join(path.Join(n.ParentsList...), n.Lower, path.Join(filename...))
}

func fieldsep(r rune) bool {
	return r == ' ' || r == '_' || r == '-'
}

var regWithParent, _ = regexp.Compile("^[a-zA-Z][/0-9a-zA-Z\\s\\_\\-]*$")
var regWithoutParent, _ = regexp.Compile("^[a-zA-Z][0-9a-zA-Z\\s\\_\\-]*$")

func listToPascal(values ...string) string {
	var result string
	for _, v := range values {
		if commonInitialisms[strings.ToUpper(v)] {
			result = result + strings.ToUpper(v)
		} else {
			result = result + strings.ToUpper(v[0:1]) + v[1:]
		}
	}
	return result
}
func New(withparents bool, s ...string) (*Name, error) {
	all := strings.Join(s, " ")
	if all == "" {
		return &Name{}, nil
	}
	var match bool
	if withparents {
		match = regWithParent.MatchString(all)
	} else {
		match = regWithoutParent.MatchString(all)
	}
	if !match {
		return nil, fmt.Errorf("name \"%s\" is not available.\nOnly alphanumeric character (0-9,a-z,A-Z) \"-\"_\"and space are allowed in name", all)
	}

	list := strings.Split(all, "/")
	plist := list[0 : len(list)-1]
	parentsList := []string{}
	parentsPascalList := []string{}
	for _, v := range plist {
		if v != "" {
			parentsList = append(parentsList, v)
			parentsPascalList = append(parentsPascalList, listToPascal(strings.FieldsFunc(v, fieldsep)...))
		}
	}
	r := list[len(list)-1]
	s = strings.FieldsFunc(r, fieldsep)

	n := &Name{
		Raw:         r,
		ParentsList: parentsList,
		Parents:     strings.Join(parentsList, "/"),
	}

	n.Title = strings.ToUpper(n.Raw[0:1]) + n.Raw[1:]
	if len(s) == 0 {
		return n, nil
	}
	n.PascalWithParents = strings.Join(parentsPascalList, "")
	n.Pascal = listToPascal(s...)
	n.PascalWithParents = n.PascalWithParents + n.Pascal
	n.Camel = s[0][0:1] + n.Pascal[1:]
	n.Lower = strings.ToLower(n.Camel)
	n.LowerWithParentDotSeparated = n.Lower
	n.LowerWithParentPath = n.Lower
	if len(n.ParentsList) > 0 {
		n.LowerWithParentDotSeparated = strings.ToLower(strings.Join(parentsPascalList, ".")) + "." + n.LowerWithParentDotSeparated
		n.LowerWithParentPath = filepath.Join(n.Parents, n.LowerWithParentPath)
	}
	return n, nil
}

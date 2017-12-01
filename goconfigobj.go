/*
Package goconfigobj is Golang Version for configobj in Python.
Python Version: https://github.com/DiffSK/configobj

Just Parse UTF-8 no BOM text, not parse list type, be careful.
*/
package goconfigobj

import "io"
import "bufio"
import "regexp"
import "strings"

//all regexp from https://github.com/DiffSK/configobj/blob/master/src/configobj/__init__.py
var (
	keyWordRegexp          = regexp.MustCompile(`^(\s*)((?:".*?")|(?:'.*?')|(?:[^'"=].*?))\s*=\s*(.*)$`)
	sectionMarkerRegexp    = regexp.MustCompile(`^(\s*)((?:\[\s*)+)((?:"\s*\S.*?\s*")|(?:'\s*\S.*?\s*')|(?:[^'"\s].*?))((?:\s*\])+)(\s*(?:\#.*)?)?$`)
	valueRegexp            = regexp.MustCompile(`^(?:(?:((?:(?:(?:".*?")|(?:'.*?')|(?:[^'",\#][^,\#]*?))\s*,\s*)*)((?:".*?")|(?:'.*?')|(?:[^'",\#\s][^,]*?))?)|(,))(\s*(?:\#.*)?)?$`)
	tripleTrailer          = `(\s*(?:#.*)?)?$`
	singleLineSingleRegexp = regexp.MustCompile(`r"^'''(.*?)'''` + tripleTrailer)
	singleLineDoubleRegexp = regexp.MustCompile(`^"""(.*?)"""` + tripleTrailer)
	multiLineSingleRegexp  = regexp.MustCompile(`^(.*?)'''` + tripleTrailer)
	multiLineDoubleRegexp  = regexp.MustCompile(`^(.*?)"""` + tripleTrailer)
)

//Section is the base struct to save datas
type Section struct {
	parent *Section
	name   string
	sects  map[string]*Section
	depth  int
	data   map[string]string
}

//ConfigObj have base section to save others
//the baseSection parent is nil and depth set to zero
type ConfigObj struct {
	baseSection *Section
}

//NewConfigObj create a new ConfigObj
func NewConfigObj(r io.Reader) *ConfigObj {
	co := &ConfigObj{
		baseSection: newSection("default", 0, nil),
	}

	co.Parse(r)

	return co
}

//Parse parse data from reader
func (co *ConfigObj) Parse(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	lines := make([]string, 0, 1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) == 0 {
		return nil
	}

	parentSect := co.baseSection.Parent()
	curSect := co.baseSection
	maxLine := len(lines)
	curIndex := 0
	curDepth := curSect.depth

	for ; curIndex < maxLine; curIndex++ {
		line := strings.TrimSpace(lines[curIndex])
		if line == "" || strings.IndexAny(line, "#;") == 0 {
			continue
		}

		mat := sectionMarkerRegexp.FindStringSubmatch(line)
		if len(mat) != 0 {
			//its a section
			sectOpen := mat[2]
			sectClose := mat[4]
			curDepth = len(sectOpen)
			if curDepth != len(sectClose) || curDepth < 1 {
				//XXX ignore this section
				curDepth = curSect.depth
				continue
			}

			if curDepth < curSect.depth {
				parentSect = co.matchParentSectDepth(curSect, curDepth)
			} else if curDepth == curSect.Depth() {
				parentSect = curSect.Parent()
			} else if curDepth == curSect.Depth()+1 {
				parentSect = curSect
			} else {
				//default use baseSection
				parentSect = co.baseSection
			}

			sectName := strings.Trim(mat[3], `"'`)
			sect := newSection(sectName, curDepth, parentSect)
			parentSect.AddSection(sectName, sect)
			curSect = sect
			continue
		}

		//should be a valid `key=value` line
		mat = keyWordRegexp.FindStringSubmatch(line)
		if len(mat) != 0 {
			value := mat[3]
			if strings.HasPrefix(value, `'''`) || strings.HasPrefix(value, `"""`) {
				//multiLine value
				value, curIndex = co.multiLineValue(value, lines, curIndex, maxLine)
			}

			key := strings.Trim(mat[2], `":'`)
			value = strings.Trim(value, `"'`)
			curSect.SetValue(key, value)
		}
	}

	return nil
}

func (co *ConfigObj) multiLineValue(value string, lines []string, curIndex, maxLine int) (string, int) {
	singleLineRegexp := singleLineSingleRegexp
	multiLineRegexp := multiLineSingleRegexp
	quot := value[0:3]
	newValue := value[3:len(value)]
	if quot == `"""` {
		singleLineRegexp = singleLineDoubleRegexp
		multiLineRegexp = multiLineDoubleRegexp
	}

	mat := singleLineRegexp.FindStringSubmatch(value)
	if len(mat) != 0 {
		return mat[1], curIndex
	} else if strings.Index(newValue, quot) != -1 {
		//have some quot in string?
		return newValue, curIndex
	}

	//read multi line
	curIndex++
	for ; curIndex < maxLine; curIndex++ {
		newValue += "\n"
		line := lines[curIndex]
		if strings.Index(line, quot) == -1 {
			newValue += line
			continue
		} else {
			mat := multiLineRegexp.FindStringSubmatch(line)
			if len(mat) != 0 {
				newValue += mat[1]
			}
			break
		}
	}

	return newValue, curIndex
}

func (co *ConfigObj) matchParentSectDepth(sect *Section, depth int) *Section {
	if sect == nil || sect == sect.Parent() {
		return nil
	}

	section := sect.Parent()
	if section.Depth() == depth {
		return section.Parent()
	}

	return co.matchParentSectDepth(section, depth)
}

//Section get section, if not have section, return nil
func (co *ConfigObj) Section(name string) *Section {
	return co.baseSection.Section(name)
}

//Value get string value by key. Please change the value type by yourself
func (co *ConfigObj) Value(key string) string {
	return co.baseSection.Value(key)
}

//AllSections return all sections
func (co *ConfigObj) AllSections() map[string]*Section {
	return co.baseSection.AllSections()
}

//AllDatas return all key-value data
func (co *ConfigObj) AllDatas() map[string]string {
	return co.baseSection.AllDatas()
}

//Section get section, if not have section, return nil
func (s *Section) Section(name string) *Section {
	sect, ok := s.sects[name]
	if ok {
		return sect
	}
	return nil
}

//Value get string value by key. Please change the value type by yourself
func (s *Section) Value(key string) string {
	value, ok := s.data[key]
	if ok {
		return value
	}

	return ""
}

//SetValue set section value
func (s *Section) SetValue(key, value string) {
	s.data[key] = value
}

//AllSections return all sections
func (s *Section) AllSections() map[string]*Section {
	return s.sects
}

//AllDatas return all key-value data
func (s *Section) AllDatas() map[string]string {
	return s.data
}

func newSection(name string, depth int, parent *Section) *Section {
	return &Section{
		depth:  depth,
		parent: parent,
		sects:  make(map[string]*Section),
		name:   name,
		data:   make(map[string]string),
	}
}

//Parent get parent section, return nil if empty
func (s *Section) Parent() *Section {
	return s.parent
}

//Depth return Section depth
func (s *Section) Depth() int {
	return s.depth
}

//AddSection add child section
func (s *Section) AddSection(name string, sect *Section) {
	s.sects[name] = sect
}

type parser struct {
	lines []string
}

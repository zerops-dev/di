package s

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strings"
)

type treeItem struct {
	id    int
	item  *XItem
	out   reflect.Type
	in    []*treeItem
	used  []*treeItem
	index int
}

type Dep struct {
	Id        DepId
	Transform interface{}
	Args      []interface{}
}

func Args(in ...interface{}) []interface{} {
	return in
}

type Dependency struct {
	Name        string
	Type        string
	TypePackage string
	IsPtr       bool
	TypeStar    string

	InitIsPtr       bool
	InitType        string
	InitTypeStar    string
	InitReflectType reflect.Type

	TransformationStart string
	TransformationEnd   string

	Index    int
	InitName string
}

type SetterCall struct {
	Call      string
	Interface reflect.Type
	Type      string
	ReturnErr bool
}

type SetterTarget struct {
	Source      *XItem
	Target      *XItem
	Setter      SetterCall
	IsInterface bool
}

type XItem struct {
	parent *XScope

	Index    int
	IsConfig bool

	Init   interface{}
	Output interface{}

	Name                  DepId
	Deps                  []Dep
	AppInject             string
	Setters               []interface{}
	SettersImplementation []reflect.Type
	SettersCall           []SetterCall
	AppInjectIndexed      string
	RunOnce               bool

	InitFunc        string
	InitFuncPackage string
	InitFuncFile    string
	InitFuncLine    int

	InitFuncWithError bool
	InitName          string
	InitNameStar      string

	InitType        string
	InitReflectType reflect.Type
	InitTypePackage string
	InitTypePtr     bool
	InitTypeStar    string

	OutputType        string
	OutputTypePackage string
	OutputTypePtr     bool
	OutputTypeStar    string

	Dependencies []*Dependency

	ConfigPath string
}

func (s *XItem) getItems() (r []*XItem) {
	return []*XItem{s}
}

func (i *XItem) getIndex(index int) *XItem {
	if i.Index == index {
		return i
	}
	return nil
}

func (i *XItem) getByName(name DepId) *XItem {
	if i.Name == name {
		return i
	}
	return nil
}

func (i *XItem) GetDependencyByType(t string) *Dependency {
	for _, i := range i.Dependencies {
		if i.InitType == t {
			return i
		}
	}
	return nil
}

type DepId int

type XScope struct {
	Items  []dependecyItem
	parent *XScope
}

var _ dependecyItem = &XScope{}

func Scope(i ...interface{}) *XScope {
	s := &XScope{}
	s.Add(i...)
	return s
}

func (s *XScope) getSiblings() (r []*XItem) {
	for _, i := range s.Items {
		if ii, isx := i.(*XItem); isx {
			r = append(r, ii)
		}
	}
	return
}

func (s *XScope) getChildren() (r []*XItem) {
	for _, i := range s.Items {
		if ii, isx := i.(*XScope); isx {
			r = append(r, ii.getItems()...)
		}
	}
	return
}

func (s *XScope) getItems() (r []*XItem) {
	for _, i := range s.Items {
		r = append(r, i.getItems()...)
	}
	return
}

func (s *XScope) getIndex(index int) *XItem {
	for _, i := range s.Items {
		if m := i.getIndex(index); m != nil {
			return m
		}
	}
	return nil
}

func (s *XScope) Add(is ...interface{}) *XScope {
	for _, item := range is {
		switch i := item.(type) {
		case []*XItem:
			for _, j := range i {
				s.Add(j)
			}
		case []*XScope:
			for _, j := range i {
				s.Add(j)
			}
		case *XItem:
			s.Items = append(s.Items, i)
		case *XScope:
			s.Items = append(s.Items, i)
		case []interface{}:
			s.Add(i...)
		default:
			panic(fmt.Sprintf("unknown type %#v", item))
		}
	}

	return s
}

func (s *XScope) getByName(name DepId) *XItem {
	for _, i := range s.Items {
		if m := i.getByName(name); m != nil {
			return m
		}
	}
	return nil
}

func (s *XScope) getDependency(dep *Dependency) error {
	var match []string
	for _, depCandidate := range s.getSiblings() {
		if depCandidate.OutputType == dep.Type {
			match = append(match, depCandidate.InitName)
			dep.TypeStar = depCandidate.InitName
			if depCandidate.InitTypePtr {
				dep.TypeStar = "*" + depCandidate.InitName
			}
		}
	}
	if len(match) > 1 {
		return fmt.Errorf("ambious dependency for %s sibling candidates: %s", dep.InitName, strings.Join(match, ", "))
	}
	if len(match) == 1 {
		return nil
	}

	match = []string{}
	for _, depCandidate := range s.getChildren() {
		if depCandidate.OutputType == dep.Type {
			match = append(match, depCandidate.InitName)
			dep.TypeStar = depCandidate.InitName
			if depCandidate.InitTypePtr {
				dep.TypeStar = "*" + depCandidate.InitName
			}
		}
	}
	if len(match) > 1 {
		return fmt.Errorf("ambious dependency for %s children candidates: %s", dep.InitName, strings.Join(match, ", "))
	}
	if len(match) == 1 {
		return nil
	}

	if s.parent != nil {
		return s.parent.getDependency(dep)
	}
	return nil
}

func (s *XScope) init(d *Di, parent *XScope, index int) (outIndex int, err error) {
	outIndex = index
	s.parent = parent
	for _, i := range s.Items {
		if outIndex, err = i.init(d, s, outIndex); err != nil {
			return
		}
	}
	return
}

func (s *XScope) deps(d *Di) (err error) {
	for _, i := range s.Items {
		if err = i.deps(d); err != nil {
			return
		}
	}
	return
}

func NamedConfig(name DepId, init interface{}, path string) *XItem {
	c := Config(init, path)
	c.Name = name
	return c
}

func Config(init interface{}, path string) *XItem {
	return &XItem{
		IsConfig:   true,
		Init:       init,
		ConfigPath: path,
	}
}

func NamedService(name DepId, init interface{}, options ...Option) *XItem {
	s := Service(init, options...)
	s.Name = name
	return s
}

func WithOutput(output interface{}) Option {
	return func(i *XItem) {
		i.Output = output
	}
}

func WithDeps(deps ...Dep) Option {
	return func(i *XItem) {
		i.Deps = append(i.Deps, deps...)
	}
}

func WithAppInject(name string) Option {
	return func(i *XItem) {
		i.AppInject = name
	}
}

func WithSetter(setter interface{}) Option {
	return func(i *XItem) {
		i.Setters = append(i.Setters, setter)
	}
}

func WithSetterImplementation[T any]() Option {
	return func(i *XItem) {
		i.SettersImplementation = append(i.SettersImplementation, reflect.TypeOf((*T)(nil)).Elem())
	}
}

func WithRunHandler() Option {
	return func(i *XItem) {
		i.RunOnce = true
		i.AppInject = "__runCommand"
	}
}

type Option func(*XItem)

//go:noinline
func Service(init interface{}, options ...Option) *XItem {
	i := &XItem{
		AppInject: "_",
		Init:      init,
	}
	for _, o := range options {
		o(i)
	}
	return i
}

type DiOption func(*Di)

func SkipRunCommand() DiOption {
	return func(d *Di) {
		d.SkipRunCommand = true
	}
}

func RegisterFlags(in bool) DiOption {
	return func(d *Di) {
		d.RegisterFlags = in
	}
}

func RunCommandOnce() DiOption {
	return func(d *Di) {
		d.RunCommandOnce = true
	}
}

func WithCommand(command, short, long string) DiOption {
	return func(d *Di) {
		d.Command = command
		d.CommandShort = short
		d.CommandLong = long
	}
}

func WithRunCommand() DiOption {
	return func(d *Di) {
		d.Command = "run"
		d.CommandShort = "run command"
		d.CommandLong = "Run command."
	}
}

func WithPanicLog() DiOption {
	return func(d *Di) {
		d.LogPanic = true
	}
}

func WithOmitConstructErr() DiOption {
	return func(d *Di) {
		d.OmitConstructErr = true
	}
}

func NewDi(configPrefix string, options ...DiOption) *Di {
	di := &Di{
		ConfigPrefix:  configPrefix,
		Scope:         Scope(),
		RegisterFlags: true,
		Imports: map[string]string{
			"configurator": "github.com/zerops-dev/di/di",
			"app":          "github.com/zerops-dev/di/app",
			"cobra":        "github.com/spf13/cobra",
		},
	}
	for _, o := range options {
		o(di)
	}
	return di
}

type dependecyItem interface {
	init(*Di, *XScope, int) (int, error)
	deps(*Di) (err error)
	getIndex(index int) *XItem
	getByName(DepId) *XItem
	getItems() []*XItem
}

//go:generate gomodrun template2struct
type Di struct {
	tree             map[reflect.Type]*treeItem
	ConfigPrefix     string
	SkipRunCommand   bool
	RegisterFlags    bool
	RunCommandOnce   bool
	Command          string
	CommandShort     string
	CommandLong      string
	LogPanic         bool
	OmitConstructErr bool

	Scope   *XScope
	Package string
	Imports map[string]string
}

func (d *Di) SetCommand(in string) *Di {
	d.Command = in
	return d
}

func (d *Di) SetCommandShort(in string) *Di {
	d.CommandShort = in
	return d
}

func (d *Di) SetCommandLong(in string) *Di {
	d.CommandLong = in
	return d
}

func (d *Di) Init(t interface{}) error {
	if i, is := t.(interface {
		GetGFilePath() string
	}); is {
		d.Package = i.GetGFilePath()
	}

	if _, err := d.Scope.init(d, nil, 0); err != nil {
		return err
	}
	if err := d.Scope.deps(d); err != nil {
		return err
	}
	if err := d.initTree(); err != nil {
		return err
	}
	return nil
}

func (d *Di) CallConfigInit(index int) interface{} {
	return reflect.ValueOf(d.Scope.getIndex(index).Init).Call(nil)[0].Interface()
}

func (d *Di) CallServiceInit(index int, args ...interface{}) interface{} {
	return reflect.ValueOf(d.Scope.getIndex(index).Init).Call(nil)[0].Interface()
}

func (d *Di) CallServiceInitWithError(index int) (interface{}, error) {
	result := reflect.ValueOf(d.Scope.getIndex(index).Init).Call(nil)
	return result[0].Interface(), (result[1].Interface()).(error)
}

func (d *Di) GetImports() (r []string) {
	for name, pkg := range d.Imports {
		if pkg == d.Package {
			continue
		}

		if path.Base(pkg) == name {
			r = append(r, fmt.Sprintf("%q", pkg))
		} else {
			r = append(r, fmt.Sprintf("%s %q", name, pkg))
		}
	}
	return r
}

func (d *Di) GetImport(pkg string) (string, string) {
	if pkg == "." {
		return "", ""
	}
	typeName := ""
	if lastIndex := strings.LastIndex(pkg, "."); lastIndex > -1 {
		typeName = pkg[lastIndex:]
		pkg = pkg[0:lastIndex]
	}

	pkgName := path.Base(pkg)
	for name, storedPkg := range d.Imports {
		if storedPkg == pkg {
			return name, name + typeName
		}
	}
	if pkg == d.Package {
		return "", typeName[1:]
	}
	importedPackageName := d.getImport(pkgName, pkg, 0)
	return importedPackageName, importedPackageName + typeName
}

func (d *Di) getImport(pkgName string, pkg string, index int) string {
	i := pkgName
	if index > 0 {
		i = fmt.Sprintf("%s%d", pkgName, index)
	}
	if _, exists := d.Imports[i]; exists {
		return d.getImport(pkgName, pkg, index+1)
	}
	d.Imports[i] = pkg
	return i
}

func (d *Di) initTree() error {
	d.tree = make(map[reflect.Type]*treeItem)

	var id int
	var roots []*treeItem
	for _, c := range d.Items() {
		px := c.InitReflectType
		if existingItem, hasDep := d.tree[px]; hasDep {
			existingItem.item = c
		} else {
			d.tree[px] = &treeItem{
				id:   id,
				item: c,
				out:  px,
				in:   nil,
			}
			id++
		}
		item := d.tree[px]
		if len(c.Dependencies) == 0 {
			roots = append(roots, item)
		}
		for _, dep := range c.Dependencies {
			cp := dep.InitReflectType
			if dep.InitReflectType.Kind() == reflect.Ptr {
				cp = dep.InitReflectType.Elem()
			}

			if _, hasDep := d.tree[cp]; !hasDep {
				d.tree[cp] = &treeItem{
					id:   id,
					item: nil,
					out:  cp,
					in:   nil,
				}
				id++
			}
			item.in = append(item.in, d.tree[cp])
		}
	}
	sortTree := make([]*treeItem, 0, len(d.tree))
	for _, item := range d.tree {
		sortTree = append(sortTree, item)
	}
	sort.Slice(sortTree, func(i, j int) bool { return sortTree[i].id > sortTree[j].id })
	for _, item := range sortTree {
		if item.item == nil {
			roots = append(roots, item)
		}
		for _, inItem := range item.in {
			inItem.used = append(inItem.used, item)
		}
	}

	setIndex := 1
	stack := roots
	for {
		var newStack []*treeItem
		for _, item := range stack {
			item.index = setIndex
			setIndex++

			var newStackAdd []*treeItem
			for _, i := range item.used {
				if func() bool {
					if i.index > 0 {
						return false
					}
					for _, dep := range i.in {
						if dep.index == 0 {
							return false
						}
					}
					return true
				}() {
					newStack = append(newStack, newStackAdd...)
				}
			}
			newStack = append(newStack, newStackAdd...)
		}
		if len(newStack) == 0 {
			break
		}
		stack = newStack
	}

	sort.Slice(sortTree, func(i, j int) bool { return sortTree[i].index < sortTree[j].index })
	return nil
}

func (i *XItem) deps(d *Di) (err error) {
	for _, dep := range i.Deps {
		depI := d.Scope.getByName(dep.Id)
		if depI == nil {
			return fmt.Errorf("name %v not defined", dep.Id)
		}
		for _, inD := range i.Dependencies {
			if inD.Type == depI.InitType {
				inD.Type = depI.InitName
				inD.TypePackage = ""
				inD.IsPtr = depI.InitTypePtr
				inD.IsPtr = depI.InitTypePtr
				inD.TypeStar = depI.InitName
				if inD.IsPtr {
					inD.TypeStar = "*" + inD.TypeStar
				}
			}
		}

		if dep.Transform != nil {
			tT := reflect.TypeOf(dep.Transform)
			if tT.Kind() == reflect.Func {
				tP := reflect.ValueOf(dep.Transform).Pointer()
				tF := runtime.FuncForPC(tP)
				_, tType := d.GetImport(tF.Name())

				if tT.NumOut() != 1 {
					return fmt.Errorf("unsupported transform function output parameters: %d", tT.NumOut())
				}
				outT := tT.Out(0)
				_, outTType := d.GetImport(outT.PkgPath() + "." + outT.Name())
				inD := i.GetDependencyByType(outTType)
				if inD == nil {
					return fmt.Errorf("transformation type %s not found on service %s", outTType, i.InitName)
				}
				inD.TransformationStart = fmt.Sprintf("%s(", tType)
				inD.TransformationEnd = ")"

				for _, a := range dep.Args {
					inD.TransformationStart += fmt.Sprintf("%#v,", a)
				}
			}
		}
	}

	for _, dep := range i.Dependencies {
		if err := i.parent.getDependency(dep); err != nil {
			return fmt.Errorf("service %s: %s", i.InitName, err.Error())
		}
	}

	return nil
}

func (i *XItem) Implements(t reflect.Type) bool {
	return reflect.PointerTo(i.InitReflectType).Implements(t)
}

func (i *XItem) init(d *Di, parent *XScope, index int) (outIndex int, err error) {
	i.parent = parent
	i.Index = index
	outIndex = i.Index + 1
	i.AppInjectIndexed = "_"
	if i.AppInject != "_" {
		i.AppInjectIndexed = fmt.Sprintf("%s%d", i.AppInject, i.Index)
	}
	defer func() {
		if err != nil {
			err = fmt.Errorf("[%s] %s", i.InitName, err)
		}
	}()
	{
		iT := reflect.TypeOf(i.Init)
		if iT.Kind() != reflect.Func {
			err = fmt.Errorf("unsupported type kind: %s value: %v", iT.Kind().String(), iT.String())
			return
		}
		switch {
		case iT.NumOut() == 2:
			i.InitFuncWithError = true
		case iT.NumOut() == 1:
		default:
			err = fmt.Errorf("unsupported return argument count: %d", iT.NumOut())
			return
		}

		for inIndex := 0; inIndex < iT.NumIn(); inIndex++ {
			inT := iT.In(inIndex)
			dep := &Dependency{
				Index:           inIndex,
				InitName:        fmt.Sprintf("arg%d", inIndex),
				InitReflectType: inT,
			}
			if inT.Kind() == reflect.Ptr {
				dep.IsPtr = true
				inT = inT.Elem()
			}
			dep.TypePackage, dep.Type = d.GetImport(inT.PkgPath() + "." + inT.Name())
			dep.TypeStar = dep.Type
			if dep.IsPtr {
				dep.TypeStar = "*" + dep.Type
			}
			dep.InitType = dep.Type
			dep.InitTypeStar = dep.TypeStar
			dep.InitIsPtr = dep.IsPtr
			i.Dependencies = append(i.Dependencies, dep)
		}

		outT := iT.Out(0)
		if outT.Kind() == reflect.Ptr {
			outT = outT.Elem()
			i.InitTypePtr = true
		}

		i.InitTypePackage, i.InitType = d.GetImport(outT.PkgPath() + "." + outT.Name())
		i.InitReflectType = outT
		i.InitTypeStar = i.InitType
		if i.InitTypePtr {
			i.InitTypeStar = "*" + i.InitType
		}

		initP := reflect.ValueOf(i.Init).Pointer()
		f := runtime.FuncForPC(initP)
		i.InitFuncFile, i.InitFuncLine = f.FileLine(initP)
		i.InitFuncPackage, i.InitFunc = d.GetImport(f.Name())

		for _, setter := range i.Setters {
			setterT := reflect.TypeOf(setter)
			if setterT.Kind() != reflect.Func {
				err = fmt.Errorf("unsupported type kind: %s value: %v", setterT.Kind().String(), setterT.String())
				return
			}

			var returnErr bool
			switch setterT.NumOut() {
			case 0:
			case 1:
				returnErr = true
				if !setterT.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
					err = fmt.Errorf("unsupported output type %s", setterT.Out(0).String())
					return
				}
			default:
				err = fmt.Errorf("unsupported setter output count: %d", setterT.NumOut())
				return
			}

			if setterT.NumIn() != 2 {
				err = fmt.Errorf("unsupported setter argument count: %d", setterT.NumIn())
				return
			}
			arg1 := setterT.In(1)

			setterP := reflect.ValueOf(setter).Pointer()
			f := runtime.FuncForPC(setterP)
			_, setterCall := d.GetImport(f.Name())

			_, t := d.GetImport(arg1.PkgPath() + "." + arg1.Name())

			i.SettersCall = append(i.SettersCall, SetterCall{
				Call:      setterCall,
				Interface: arg1,
				Type:      t,
				ReturnErr: returnErr,
			})
		}

	}

	{
		if i.Output == nil {
			i.OutputType = i.InitType
			i.OutputTypePtr = i.InitTypePtr
			i.OutputTypePackage = i.InitTypePackage
			i.OutputTypeStar = i.InitTypeStar
		} else {
			outT := reflect.TypeOf(i.Output)
			if outT.Kind() == reflect.Ptr {
				outT = outT.Elem()
				i.OutputTypePtr = true
			}

			i.OutputTypePackage, i.OutputType = d.GetImport(outT.PkgPath() + "." + outT.Name())
			i.OutputTypeStar = i.OutputType
			if i.OutputTypePtr {
				i.OutputTypeStar = "*" + i.OutputType
			}
		}
	}

	i.InitName = fmt.Sprintf("Name%s%d", toName(i.OutputType), i.Index)
	i.InitNameStar = i.InitName

	if i.OutputTypePtr {
		i.InitNameStar = "*" + i.InitName
	}

	return
}

func (d *Di) Items() (r []*XItem) {
	return d.Scope.getItems()
}

func (d *Di) Add(is ...interface{}) *Di {
	d.Scope.Add(is...)
	return d
}

func (d *Di) Configs() (r []*XItem) {
	for _, c := range d.Items() {
		if !c.IsConfig {
			continue
		}
		r = append(r, c)
	}
	return
}

func (d *Di) Services() (r []*XItem) {
	for _, c := range d.Items() {
		if c.IsConfig {
			continue
		}
		r = append(r, c)
	}

	sort.Slice(r, func(i, j int) bool {
		sourceIType, sourceITypeExists := d.tree[r[i].InitReflectType]
		if !sourceITypeExists {
			return false
		}
		sourceJType, sourceJTypeExists := d.tree[r[j].InitReflectType]
		if !sourceJTypeExists {
			return true
		}
		return sourceIType.index < sourceJType.index
	})

	return
}

func (d *Di) Setters() (r []SetterTarget) {
	for _, c := range d.Items() {
		if c.IsConfig {
			continue
		}
		if len(c.SettersCall) == 0 {
			continue
		}
		for _, setter := range c.SettersCall {
			for _, target := range d.Items() {
				var add bool
				var isInterface bool
				switch target.InitReflectType.Kind() {
				case reflect.Interface:
					add = target.InitReflectType.Implements(setter.Interface)
					if !add {
						for _, setterImpl := range target.SettersImplementation {
							if setterImpl == setter.Interface {
								add = true
								isInterface = true
								break
							}
						}
					}
				default:
					add = reflect.PointerTo(target.InitReflectType).Implements(setter.Interface)
				}
				if add {
					r = append(r, SetterTarget{
						Source:      c,
						Target:      target,
						Setter:      setter,
						IsInterface: isInterface,
					})
				}
			}
		}
	}

	// sort setters by deps
	sort.Slice(r, func(i, j int) bool {
		sourceIType := d.tree[r[i].Source.InitReflectType]
		sourceJType := d.tree[r[j].Source.InitReflectType]
		targetIType := d.tree[r[i].Target.InitReflectType]
		targetJType := d.tree[r[j].Target.InitReflectType]
		if sourceIType.index == sourceJType.index {
			return targetIType.index < targetJType.index
		}
		return sourceIType.index < sourceJType.index
	})
	return
}

func (d *Di) ServicesFilter(filtered ...string) (r []*XItem) {

	for _, c := range d.Items() {
		if c.IsConfig {
			continue
		}
		var skip bool
		for _, f := range filtered {
			if f == c.OutputType {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		r = append(r, c)
	}
	return
}

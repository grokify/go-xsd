package xsd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	ufs "github.com/grokify/go-util/fs"
	unet "github.com/grokify/go-util/net"
	ustr "github.com/grokify/go-util/str"
)

const (
	goPkgPrefix     = ""
	goPkgSuffix     = "_go"
	protSep         = "://"
	xsdNamespaceUri = "http://www.w3.org/2001/XMLSchema"
)

var (
	loadedSchemas = map[string]*Schema{}
)

type Schema struct {
	elemBase

	XMLName            xml.Name          `xml:"schema"`
	XMLNamespacePrefix string            `xml:"-"`
	XMLNamespaces      map[string]string `xml:"-"`
	XMLIncludedSchemas []*Schema         `xml:"-"`
	XSDNamespacePrefix string            `xml:"-"`
	XSDParentSchema    *Schema           `xml:"-"`

	hasAttrAttributeFormDefault
	hasAttrBlockDefault
	hasAttrElementFormDefault
	hasAttrFinalDefault
	hasAttrLang
	hasAttrId
	hasAttrSchemaLocation
	hasAttrTargetNamespace
	hasAttrVersion
	hasElemAnnotation
	hasElemsAttribute
	hasElemsAttributeGroup
	hasElemsComplexType
	hasElemsElement
	hasElemsGroup
	hasElemsInclude
	hasElemsImport
	hasElemsNotation
	hasElemsRedefine
	hasElemsSimpleType

	loadLocalPath, loadURI string
}

func (me *Schema) allSchemas(loadedSchemas map[string]bool) (schemas []*Schema) {
	schemas = append(schemas, me)
	loadedSchemas[me.loadURI] = true
	for _, ss := range me.XMLIncludedSchemas {
		if v, ok := loadedSchemas[ss.loadURI]; ok && v {
			continue
		}
		schemas = append(schemas, ss.allSchemas(loadedSchemas)...)
	}
	return
}

func (me *Schema) collectGlobals(bag *PkgBag, loadedSchemas map[string]bool) {
	loadedSchemas[me.loadURI] = true
	bag.allAtts = append(bag.allAtts, me.Attributes...)

	bag.allAttGroups = append(bag.allAttGroups, me.AttributeGroups...)

	bag.allElems = append(bag.allElems, me.Elements...)

	bag.allElemGroups = append(bag.allElemGroups, me.Groups...)

	bag.allNotations = append(bag.allNotations, me.Notations...)

	for _, ss := range me.XMLIncludedSchemas {
		if v, ok := loadedSchemas[ss.loadURI]; ok && v {
			continue
		}
		ss.collectGlobals(bag, loadedSchemas)
	}
}

func (me *Schema) globalComplexType(bag *PkgBag, name string, loadedSchemas map[string]bool) (ct *ComplexType) {
	var imp string
	for _, ct = range me.ComplexTypes {
		if bag.resolveQnameRef(ustr.PrefixWithSep(me.XMLNamespacePrefix, ":", ct.Name.String()), "T", &imp) == name {
			return
		}
	}
	loadedSchemas[me.loadURI] = true
	for _, ss := range me.XMLIncludedSchemas {
		if v, ok := loadedSchemas[ss.loadURI]; ok && v {
			//fmt.Printf("Ignoring processed schema: %s\n", ss.loadUri)
			continue
		}
		if ct = ss.globalComplexType(bag, name, loadedSchemas); ct != nil {
			return
		}
	}
	ct = nil
	return
}

func (me *Schema) globalElement(bag *PkgBag, name string) (el *Element) {
	var imp string
	if len(name) > 0 {
		var rname = bag.resolveQnameRef(name, "", &imp)
		for _, el = range me.Elements {
			if bag.resolveQnameRef(ustr.PrefixWithSep(me.XMLNamespacePrefix, ":", el.Name.String()), "", &imp) == rname {
				return
			}
		}
		for _, ss := range me.XMLIncludedSchemas {
			if el = ss.globalElement(bag, name); el != nil {
				return
			}
		}
	}
	el = nil
	return
}

func (me *Schema) globalSubstitutionElems(el *Element, loadedSchemas map[string]bool) (els []*Element) {
	var elName = el.Ref.String()
	if len(elName) == 0 {
		elName = el.Name.String()
	}
	for _, tle := range me.Elements {
		if (tle != el) && (len(tle.SubstitutionGroup) > 0) {
			if tle.SubstitutionGroup.String()[(strings.Index(tle.SubstitutionGroup.String(), ":")+1):] == elName {
				els = append(els, tle)
			}
		}
	}
	loadedSchemas[me.loadURI] = true
	for _, inc := range me.XMLIncludedSchemas {
		if v, ok := loadedSchemas[inc.loadURI]; ok && v {
			//fmt.Printf("Ignoring processed schema: %s\n", inc.loadUri)
			continue
		}
		els = append(els, inc.globalSubstitutionElems(el, loadedSchemas)...)
	}
	return
}

func (me *Schema) MakeGoPkgSrcFile() (goOutFilePath string, err error) {
	var goOutDirPath = filepath.Join(filepath.Dir(me.loadLocalPath), goPkgPrefix+filepath.Base(me.loadLocalPath)+goPkgSuffix)
	goOutFilePath = filepath.Join(goOutDirPath, path.Base(me.loadURI)+".go")
	var bag = newPkgBag(me)
	loadedSchemas := make(map[string]bool)
	for _, inc := range me.allSchemas(loadedSchemas) {
		bag.Schema = inc
		inc.makePkg(bag)
	}
	bag.Schema = me
	me.hasElemAnnotation.makePkg(bag)
	bag.appendFmt(true, "")
	me.makePkg(bag)
	if err = ufs.EnsureDirExists(filepath.Dir(goOutFilePath)); err == nil {
		err = ufs.WriteTextFile(goOutFilePath, bag.assembleSource())
	}
	return
}

func (me *Schema) onLoad(rootAtts []xml.Attr, loadURI, localPath string) (err error) {
	var tmpURL string
	var sd *Schema
	loadedSchemas[loadURI] = me
	me.loadLocalPath, me.loadURI = localPath, loadURI
	me.XMLNamespaces = map[string]string{}
	for _, att := range rootAtts {
		if att.Name.Space == "xmlns" {
			me.XMLNamespaces[att.Name.Local] = att.Value
		} else if len(att.Name.Space) > 0 {

		} else if att.Name.Local == "xmlns" {
			me.XMLNamespaces[""] = att.Value
		}
	}
	for k, v := range me.XMLNamespaces {
		if v == xsdNamespaceUri {
			me.XSDNamespacePrefix = k
		} else if v == me.TargetNamespace.String() {
			me.XMLNamespacePrefix = k
		}
	}
	if len(me.XMLNamespaces["xml"]) == 0 {
		me.XMLNamespaces["xml"] = "http://www.w3.org/XML/1998/namespace"
	}
	me.XMLIncludedSchemas = []*Schema{}
	for _, inc := range me.Includes {
		if tmpURL = inc.SchemaLocation.String(); !strings.Contains(tmpURL, protSep) {
			tmpURL = path.Join(path.Dir(loadURI), tmpURL)
		}
		var ok bool
		var toLoadURI string
		if pos := strings.Index(tmpURL, protSep); pos >= 0 {
			toLoadURI = tmpURL[pos+len(protSep):]
		} else {
			toLoadURI = tmpURL
		}
		if sd, ok = loadedSchemas[toLoadURI]; !ok {
			if sd, err = LoadSchema(tmpURL, len(localPath) > 0); err != nil {
				return
			}
		}
		sd.XSDParentSchema = me
		me.XMLIncludedSchemas = append(me.XMLIncludedSchemas, sd)
	}
	me.initElement(nil)
	return
}

func (me *Schema) RootSchema(pathSchemas []string) *Schema {
	if me.XSDParentSchema != nil {
		for _, sch := range pathSchemas {
			if me.XSDParentSchema.loadURI == sch {
				fmt.Printf("schema loop detected %+v - > %s!\n", pathSchemas, me.XSDParentSchema.loadURI)
				return me
			}
		}
		pathSchemas = append(pathSchemas, me.loadURI)
		return me.XSDParentSchema.RootSchema(pathSchemas)
	}
	return me
}

func ClearLoadedSchemasCache() {
	loadedSchemas = map[string]*Schema{}
}

func loadSchema(r io.Reader, loadUri, localPath string) (sd *Schema, err error) {
	var data []byte
	var rootAtts []xml.Attr
	if data, err = io.ReadAll(r); err == nil {
		var t xml.Token
		sd = new(Schema)
		for xd := xml.NewDecoder(bytes.NewReader(data)); err == nil; {
			if t, err = xd.Token(); err == nil {
				if startEl, ok := t.(xml.StartElement); ok {
					rootAtts = startEl.Attr
					break
				}
			}
		}
		if err = xml.Unmarshal(data, sd); err == nil {
			err = sd.onLoad(rootAtts, loadUri, localPath)
		}
		if err != nil {
			sd = nil
		}
	}
	return
}

func loadSchemaFile(filename string, loadURI string) (sd *Schema, err error) {
	var file *os.File
	if file, err = os.Open(filename); err == nil {
		defer file.Close()
		sd, err = loadSchema(file, loadURI, filename)
	}
	return
}

func LoadSchema(uri string, localCopy bool) (sd *Schema, err error) {
	var protocol, localPath string
	var rc io.ReadCloser

	if pos := strings.Index(uri, protSep); pos < 0 {
		protocol = "http" + protSep
	} else {
		protocol = uri[:pos+len(protSep)]
		uri = uri[pos+len(protSep):]
	}
	if localCopy {
		if localPath = filepath.Join(PkgGen.BaseCodePath, uri); !ufs.FileExists(localPath) {
			if err = ufs.EnsureDirExists(filepath.Dir(localPath)); err == nil {
				err = unet.DownloadFile(protocol+uri, localPath)
			}
		}
		if err == nil {
			if sd, err = loadSchemaFile(localPath, uri); sd != nil {
				sd.loadLocalPath = localPath
			}
		}
	} else if rc, err = unet.OpenRemoteFile(protocol + uri); err == nil {
		defer rc.Close()
		sd, err = loadSchema(rc, uri, "")
	}
	return
}

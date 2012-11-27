package file

import (
	"fmt"
	"github.com/opesun/jsonp"
	"github.com/opesun/nocrud/frame/composables/basics"
	iface "github.com/opesun/nocrud/frame/interfaces"
	"github.com/opesun/nocrud/frame/misc/convert"
	"github.com/opesun/sanitize"
)

type C struct {
	basics.Basics
	fileSys iface.FileSys
	fileBiz map[string][]iface.ReadableFile
	opt     map[string]interface{}
}

func (c *C) Init(ctx iface.Context) {
	c.Basics.Hooks = ctx.Conducting().Hooks()
	c.Basics.Db = ctx.Db()
	c.fileSys = ctx.FileSys()
	c.opt = ctx.Options().Document().All().(map[string]interface{})
	c.fileBiz = map[string][]iface.ReadableFile{}
}

func (c *C) SanitizerMangler(san *sanitize.Extractor) {
	san.AddFuncs(sanitize.FuncMap{
		"file": func(dat interface{}, s sanitize.Scheme) (interface{}, error) {
			temps := c.fileSys.Temporaries()
			if !temps.Exists(s.Key) {
				return nil, fmt.Errorf("Can't find key amongst files.")
			}
			temp := temps.Select(s.Key)
			if len(temp) > 0 {
				c.fileBiz[s.Key] = temp
			}
			return nil, nil
		},
	})
}

func (c *C) moveFiles(subject, id string) error {
	uploads, err := c.fileSys.SelectPlace("uploads")
	if err != nil {
		return err
	}
	target_dir := uploads.Directory(subject, id)
	for folder, files := range c.fileBiz {
		for _, file := range files {
			data, err := file.Read()
			if err != nil {
				return err
			}
			target_file := target_dir.Directory(folder).File(file.Name())
			err = target_file.Create()
			if err != nil {
				return err
			}
			err = target_file.Write(data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Collects the filenames so we can save them in the document.
func fileNames(files map[string][]iface.ReadableFile) map[string][]string {
	ret := map[string][]string{}
	for folder, files := range files {
		if _, has := ret[folder]; !has {
			ret[folder] = []string{}
		}
		for _, file := range files {
			ret[folder] = append(ret[folder], file.Name())
		}
	}
	return ret
}

func merge(a map[string]interface{}, b map[string][]string) map[string]interface{} {
	for i, v := range b {
		a[i] = v
	}
	return a
}

func (c *C) Insert(a iface.Filter, data map[string]interface{}) (iface.Id, error) {
	has_files := len(c.fileBiz) > 0
	if has_files {
		merge(data, fileNames(c.fileBiz))
	}
	id, err := c.Basics.Insert(a, data)
	if err != nil {
		return id, err
	}
	if has_files {
		err := c.moveFiles(a.Subject(), id.String())
		if err != nil {
			return nil, err
		}
	}
	return id, nil
}

// We build the modifier query here, using "$each" where we have multiple files.
// The return value of this can be used with the "$addToSet" modifier.
func createQuery(filenames map[string][]string) map[string]interface{} {
	ret := map[string]interface{}{}
	for folder, slice := range filenames {
		l := len(slice)
		if l == 0 {
			panic("I shouldn't receive empty file slices.")
		} else if l == 1 {
			ret[folder] = slice[0]
		} else {
			each := []interface{}{}
			for _, v := range slice {
				each = append(each, v)
			}
			ret[folder] = map[string]interface{}{
				"$each": each,
			}
		}
	}
	return ret
}

func (c *C) Update(a iface.Filter, data map[string]interface{}) error {
	upd := map[string]interface{}{
		"$set": data,
	}
	has_files := len(c.fileBiz) > 0
	if has_files {
		ids, err := a.Ids()
		if err != nil {
			return err
		}
		err = c.moveFiles(a.Subject(), ids[0].String())
		if err != nil {
			return err
		}
		upd["$addToSet"] = createQuery(fileNames(c.fileBiz))
	}
	err := a.Update(upd)
	if err != nil {
		return err
	}
	if c.Hooks != nil {
		c.Hooks.Select("Updated").Fire(a)
		c.Hooks.Select(a.Subject() + "Updated").Fire(a)
	}
	return nil
}

func (c *C) getScheme(noun, verb string) (map[string]interface{}, error) {
	scheme, ok := jsonp.GetM(c.opt, fmt.Sprintf("nouns.%v.verbs.%v.input", noun, verb))
	if !ok {
		return nil, fmt.Errorf("Can't find scheme for noun %v, verb %v.", noun, verb)
	}
	return scheme, nil
}

func (c *C) DeleteFile(a iface.Filter, data map[string]interface{}) error {
	upd := map[string]interface{}{
		"$pull": map[string]interface{}{
			data["key"].(string): data["file"],
		},
	}
	return a.Update(upd)
}

func (c *C) DeleteAllFiles(a iface.Filter, data map[string]interface{}) error {
	upd := map[string]interface{}{
		"$unset": data["key"].(string),
	}
	return a.Update(upd)
}

func (c *C) New(a iface.Filter) ([]map[string]interface{}, error) {
	scheme, err := c.getScheme(a.Subject(), "Insert")
	if err != nil {
		return nil, err
	}
	return convert.SchemeToFields(scheme, nil)
}

func (c *C) Edit(a iface.Filter) ([]map[string]interface{}, error) {
	doc, err := a.FindOne()
	if err != nil {
		return nil, err
	}
	scheme, err := c.getScheme(a.Subject(), "Update")
	if err != nil {
		return nil, err
	}
	return convert.SchemeToFields(scheme, doc)
}

func (c *C) Install(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$addToSet": map[string]interface{}{
			"Hooks.SanitizerMangler": "file",
			"Hooks.fileTypeHandler": []interface{}{
				"file",
				"FileTypeHandler",
			},
		},
	}
	return o.Update(upd)
}

func (c *C) Uninstall(o iface.Document, resource string) error {
	upd := map[string]interface{}{
		"$pull": map[string]interface{}{
			"Hooks.SanitizerMangler": "file",
			"Hooks.fileTypeHandler": []interface{}{
				"file",
				"FileTypeHandler",
			},
		},
	}
	return o.Update(upd)
}

package temporaries

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
	"mime/multipart"
	"bytes"
)

type Temporaries struct {
	files	map[string][]*multipart.FileHeader
}

func New(f map[string][]*multipart.FileHeader) *Temporaries {
	if f == nil {
		f = map[string][]*multipart.FileHeader{}
	}
	return &Temporaries{
		f,
	}
}

type ReadableFile struct {
	file 	*multipart.FileHeader
}

func (r *ReadableFile) Read() ([]byte, error) {
	buf := new(bytes.Buffer)
	file, err := r.file.Open()
	if err != nil {
		return nil, err
	}
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *ReadableFile) Name() string {
	return r.file.Filename
}

func (t *Temporaries) Select(key string) []iface.ReadableFile {
	val, ok := t.files[key]
	if !ok {
		return nil
	}
	ret := []iface.ReadableFile{}
	for _, v := range val {
		ret = append(ret, &ReadableFile{
			v,
		})
	}
	return ret
}

func (t *Temporaries) Exists(key string) bool {
	_, ok := t.files[key]
	return ok
}

func (t *Temporaries) Keys() []string {
	keys := []string{}
	for i := range t.files {
		keys = append(keys, i)
	}
	return keys
}
package parser

type Zip map[string]string

func (zip Zip) Get(key, value string) string {
	v, ok := zip[key]
	if ok {
		return v
	}
	return value
}

func (zip Zip) Same(key string) string {
	return zip.Get(key, key)
}

func NewZip(l1, l2 []string, errMsg string) (zip Zip) {
	if len(l1) != len(l2) {
		panic(errMsg)
	}
	zip = make(Zip)
	for i, arg := range l1 {
		zip[arg] = l2[i]
	}
	return
}

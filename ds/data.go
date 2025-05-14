package ds

type Data struct {
	k int
	v string
}

func NewData(k int, v string) *Data {
	return &Data{
		k,
		v,
	}
}

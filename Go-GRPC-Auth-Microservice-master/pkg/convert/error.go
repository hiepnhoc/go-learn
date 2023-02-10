package convert

type Error struct {
	from   string
	method string
	in     interface{}
	out    interface{}
	error  error
}

func NewError(from, method string, in, out interface{}, error error) *Error {
	return &Error{from, method, in, out, error}
}

//
//func (e *Error) Error() string {
//	return fmt.Sprintf("[ConvertError] %s - %s [type:%s value:%v] -> [type:%s value:%v] error %v", e.from, e.method, TypeOfObject(e.in), e.in, TypeOfObject(e.out), e.out, e.error)
//}

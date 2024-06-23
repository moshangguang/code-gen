package runtime

func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}
func MustTrue(ok bool, msg interface{}) {
	if ok {
		return
	}
	panic(msg)
}
func MustFalse(ok bool, msg interface{}) {
	if !ok {
		return
	}
	panic(msg)
}

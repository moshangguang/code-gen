package running

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
func True(ok bool, msg interface{}) {
	if ok {
		return
	}
	panic(msg)
}
func False(ok bool, msg interface{}) {
	if !ok {
		return
	}
	panic(msg)
}

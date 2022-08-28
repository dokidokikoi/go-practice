package handlers

func PanicIfUserError(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicIfTaskError(err error) {
	if err != nil {
		panic(err)
	}
}

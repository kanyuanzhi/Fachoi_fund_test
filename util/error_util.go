package util

func CheckError(err error, info string) {
	if err != nil {
		panic(info + ": " + err.Error())
	}
}

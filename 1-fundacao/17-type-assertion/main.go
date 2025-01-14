package main

func main() {
	var minhaVar interface{} = "Wesley"
	println(minhaVar.(string))
	res, ok := minhaVar.(string)
	println(res, ok)
}

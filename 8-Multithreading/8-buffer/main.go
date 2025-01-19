package main

func main() {
	c := make(chan string, 2)
	c <- "Hello"
	c <- "World"

	println(<-c)
	println(<-c)
}

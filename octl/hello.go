package main

func main() {
	// name := os.Args[1]
	// SSH(name)
	// UpLoadFile("/root/Octopoda/octl/go.mod", "./test/")
	SpreadFile("go.mod", "test/", "newfolder/", []string{"pi0"})
}

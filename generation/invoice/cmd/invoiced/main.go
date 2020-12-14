package main

import (
	"fmt"

	server "../../../invoicegen"
	merchant "../../merchant"
)

func main() {
	fmt.Println("docgen: invoice & sales report generation")
	fmt.Println("=========================================")

	fmt.Println("Loading 'hackwave laboratories' document templates")

	fmt.Println("Creating [http server] AND [starting] for document generation...")
	app := server.Init(merchant.Default("hackwave laboratories"))
	app.Start()
}

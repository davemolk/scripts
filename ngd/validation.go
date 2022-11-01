package main

import "log"

func (n *ngd) validateInput(config config) {
	switch {
	case config.query == "" && config.queryExact == "":
		log.Println("no search term supplied...")
	case config.filetype != "":
		n.checkFileType(config.filetype)
	}	
}

func (n *ngd) checkFileType(file string) {
	allowed := []string{"pdf", "doc(x)", "xls(x)", "ppt(x)", "html"}
	for _, v := range allowed {
		if file == v {
			return 
		}
	}
	log.Fatal("filetype unsupported")
}
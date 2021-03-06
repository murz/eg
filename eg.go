package main

import (
	"log"
	"os"
	"bytes"
	"compress/gzip"
	"io"
	"flag"
	"strings"
	"github.com/hoisie/mustache"
	"io/ioutil"
	"regexp"
	"fmt"
	"github.com/murz/eg/proxy"
	"github.com/murz/eg/templates"
)

func main() {
	var path string;
	flag.StringVar(&path, "path", "/", "")
	args := os.Args[1:len(os.Args)] // throw out the first one because it's always eg
	if len(args) > 0 {
		switch args[0] {
		case "new", "n":
			new(args)
		case "help":
			help(args)
		case "build", "b":
			build(args)
		case "run", "r":
			run(args)
		case "remove", "rm", "del", "delete":
			delete(args)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func run(args []string) {
	args = args[1:len(args)] // shave off the 'run' arg
	processFlags(args)
	proxy.Run()

}


func new(args []string) {
	if len(args) < 2 {
		log.Print("ego: Not enough args for `eg new`. Use `eg help` for more info.")
		return
	}
	args = args[1:len(args)] // shave off the 'new' arg
	switch(args[0]) {
	case "app":
		newApp(args)
	case "controller", "ctrlr", "ctrl", "c":
		newController(args)
	case "action", "actn":
		newAction(args)
	}
}

func delete(args []string) {
	if len(args) < 2 {
		log.Print("ego: Not enough args for `eg delete`. Use `eg help` for more info.")
		return
	}
	args = args[1:len(args)] // shave off the 'new' arg
	switch(args[0]) {
	case "action":
		deleteAction(args)
	}
}

func newApp(args []string) {
	if len(args) < 2 {
		log.Print("ego: Not enough args for `eg new app`. Use `eg help` for more info.")
		return
	}
	args = args[1:len(args)] // shave off the 'app' arg
	name := args[0]
	os.Mkdir(name, 0777)
	os.Mkdir(name+"/app", 0777)
	os.Mkdir(name+"/app/controllers", 0777)
	os.Mkdir(name+"/app/helpers", 0777)
	os.Mkdir(name+"/app/models", 0777)
	os.Mkdir(name+"/app/views", 0777)
	os.Mkdir(name+"/app/views/errors", 0777)
	os.Mkdir(name+"/app/assets", 0777)
	os.Mkdir(name+"/app/assets/javascripts", 0777)
	os.Mkdir(name+"/app/assets/stylesheets", 0777)
	os.Mkdir(name+"/app/assets/images", 0777)
	os.Mkdir(name+"/conf", 0777)
	os.Mkdir(name+"/public", 0777)

    routesconf := mustache.Render(string(templates.Routes()), map[string]string{})
    routesconfFile, _ := os.Create(name+"/conf/routes.go")
    routesconfFile.Write([]byte(routesconf))

    dbconf := mustache.Render(string(templates.Databases()), map[string]string{})
    dbconfFile, _ := os.Create(name+"/conf/db.go")
    dbconfFile.Write([]byte(dbconf))

	error404view := mustache.Render(string(error_html_mustache()), map[string]string {
		"Message": "404 Not Found",
	})
	error404viewFile, _ := os.Create(name+"/app/views/errors/404.html")
	error404viewFile.Write([]byte(error404view))

	error501view := mustache.Render(string(error_html_mustache()), map[string]string {
		"Message": "501 Not Implemented",
	})
	error501viewFile, _ := os.Create(name+"/app/views/errors/501.html")
	error501viewFile.Write([]byte(error501view))

	log.Printf("Your new ego application, '%v', was successfully created", args[0])
}

var flags = map[string]string {
	"file": "",
	"method": "GET",
	"path": "/",
	"port": "5000",
}

func newController(args []string) {
	if len(args) < 2 {
		log.Print("ego: Not enough args for `eg new controller`. Use `eg help` for more info.")
		return
	}
	if (!checkDirs([]string{
		"app",
		"app/controllers",
	})) {
		log.Print("ego: You must be in an ego project directory to use `eg new controller`.")
		return
	}
	// defFlag("method", "m", "GET")
	// defFlag("path", "p", "/")
	// flag.Parse()
	name := args[1]
	args = args[2:len(args)] // shave off the 'action name' args
	processFlags(args)

	ctrlFile, err := os.Create("app/controllers/"+strings.ToLower(name)+"_controller.go")
	checkErr(err)

	ctrl := mustache.Render(string(templates.Controller()), map[string]string {
		"Name": strings.Title(name) + "Controller",
		"Embed": "*http.Controller",
	})
	
	ctrlFile.WriteString(ctrl)
	log.Printf("Controller '%v', was successfully created", name)
}

func newAction(args []string) {
	if len(args) < 2 {
		log.Print("ego: Not enough args for `eg new action`. Use `eg help` for more info.")
		return
	}
	if (!checkDirs([]string{
		"app",
		"app/actions",
	})) {
		log.Print("ego: You must be in an ego project directory to use `eg new action`.")
		return
	}
	// defFlag("method", "m", "GET")
	// defFlag("path", "p", "/")
	// flag.Parse()
	name := args[1]
	args = args[2:len(args)] // shave off the 'action name' args
	processFlags(args)

	var actionFile *os.File
	var err error
	if flags["file"] != "" {
		actionFile, err = os.OpenFile("app/actions/"+flags["file"], os.O_WRONLY, 777)
		checkErr(err)
	} else {
		actionFile, err = os.Create("app/actions/"+strings.ToLower(name)+".go")
		checkErr(err)
		cruft := mustache.Render(string(actionfile_go_mustache()), map[string]string {})
		actionFile.WriteString(cruft)
	}

	action := mustache.Render(string(action_go_mustache()), map[string]string {
		"Name": strings.Title(name),
		"Path": flags["path"],
		"Method": flags["method"],
	})
	
	actionFile.Seek(0, 2)
	actionFile.WriteString("\n\n")
	actionFile.WriteString(action)
	log.Printf("Action '%v', was successfully created", name)
}

func deleteAction(args []string) {
	if len(args) < 2 {
		log.Print("ego: Not enough args for `eg rm action`. Use `eg help` for more info.")
		return
	}
	if (!checkDirs([]string{
		"app",
		"app/actions",
	})) {
		log.Print("ego: You must be in an ego project directory to use `eg rm action`.")
		return
	}
	name := args[1]
	args = args[2:len(args)] // shave off the 'action name' args
	processFlags(args)

	if flags["file"] == "" {
		log.Print("ego: Must specify -file with `eg rm action`")
	}

	actionFile, err := os.OpenFile("app/actions/"+flags["file"], os.O_RDWR, 777)
	checkErr(err)

	bytes, err := ioutil.ReadAll(actionFile)
	checkErr(err)

	expString := fmt.Sprintf("(?is)\\n\\nvar %v = actions\\.register.+?}\\)", name)
	exp, err := regexp.Compile(expString)
	checkErr(err)
	pos := exp.FindAllIndex(bytes, -1)

	oldStr := string(bytes)
	

	str := make([]byte, pos[0][1] - pos[0][0])
	actionFile.ReadAt(str, int64(pos[0][0]))

	newStr := strings.Replace(oldStr, string(str), "", -1)
	actionFile.Truncate(0)
	actionFile.Seek(0, 0)

	actionFile.WriteString(newStr)

	log.Printf("pos: %v", pos)
}

func processFlags(args []string) {
	next := ""
	for _, arg := range args {
		if next != "" {
			flags[next] = arg
			next = ""
		} else if arg[0:1] == "-" {
			name := arg[1:]
			if name[0:1] == "-" {
				name = name[1:]
			}
			pieces := strings.Split(name, "=")
			if len(pieces) > 1 {
				flags[pieces[0]] = pieces[1]
			} else {
				next = name
			}
		}
	}
}

func build(args []string) {
	// package views
	// dirlist, err := ioutil.ReadDir("/app/views")
	// if err != nil {
	// 	log.Fatalf("Error reading %s: %s\n", dirname, err)
	// }
	// for _, f := range dirlist {
	// 	filename := path.Join(tm.pkgName, dirname, f.Name())
	// 	if f.IsDir() {
	// 		tm.ParseDir(path.Join(dirname, f.Name()))
	// 	} else {
	// 		tm.Parse(filename)
	// 	}
	// }
}

func help(args []string) {
	log.Print("lol.. todo")
}

func checkDirs(dirs []string) bool {
	// var dirs = []string{
	// 	"app",
	// 	"app/actions",
	// 	"app/models",
	// 	"conf",
	// 	"conf/app.json",
	// 	"public",
	// }
	for _, dir := range dirs {
		if ex, _ := exists(dir); !ex {
			return false
		}
	}
	return true
}

type Flag struct {
	Name string
	Value *string
	Alt *string
}

var  flagMap = make(map[string]*Flag)

func defFlag(name string, alt string, def string) {
	var f = flag.String(name, def, "")
	var altf = flag.String(alt, def, "")
	fl := &Flag{
		Name: name,
		Value: f,
		Alt: altf,
	}
	flagMap[name] = fl
}

func getFlag(name string) *Flag {
	return flagMap[name]
}

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

// Templates
func server_go_mustache() []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte{
0x1f,0x8b,0x08,0x00,0x00,0x09,0x6e,0x88,0x00,0xff,0x34,0xcc,
0x3d,0xca,0x02,0x31,0x10,0x80,0xe1,0x3a,0x73,0x8a,0x61,0xaa,
0xa4,0xd9,0x14,0x5f,0xf7,0x81,0x57,0xd8,0x42,0x0f,0x20,0x63,
0x18,0xd7,0x20,0xf9,0x21,0x3f,0x2b,0x18,0x72,0x77,0x17,0xc1,
0xb7,0x7f,0xde,0xcc,0xee,0xc9,0x9b,0x60,0x60,0x1f,0x01,0x7c,
0xc8,0xa9,0x34,0xd4,0xa0,0x68,0xf3,0xed,0xd1,0x6f,0x8b,0x4b,
0xc1,0x86,0x5e,0xde,0x56,0xb6,0x44,0xa0,0xae,0x48,0x63,0xac,
0x1c,0x64,0x4e,0xcb,0x39,0x5b,0x76,0xcd,0xa7,0x58,0x09,0x0c,
0xc0,0xbd,0x47,0xf7,0xfd,0x68,0x83,0x03,0xd4,0xce,0x05,0x2b,
0x9e,0xf0,0x80,0xcb,0x2a,0xaf,0x8b,0x94,0x5d,0x8a,0xfe,0x69,
0x03,0xaa,0x2e,0xe7,0x1e,0x35,0xfd,0xff,0x1d,0x91,0x81,0xf9,
0x09,0x00,0x00,0xff,0xff,0xb2,0x01,0x5a,0x7e,0x8b,0x00,0x00,
0x00,
	}))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes()
}

func db_json_mustache() []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte{
0x1f,0x8b,0x08,0x00,0x00,0x09,0x6e,0x88,0x00,0xff,0xaa,0xe6,
0xe2,0x54,0x4a,0x29,0xca,0x2c,0x4b,0x2d,0x52,0xb2,0x52,0x50,
0x2a,0xc8,0x2f,0x2e,0x49,0x2f,0x4a,0x2d,0x56,0xd2,0x01,0x0a,
0xe7,0x25,0xe6,0xa6,0x82,0x04,0x4b,0xf2,0x53,0xf2,0xe3,0x53,
0x52,0x73,0xf3,0xc1,0xa2,0xa5,0xc5,0x10,0xa5,0xa8,0xa2,0x05,
0x89,0xc5,0xc5,0xe5,0xf9,0x45,0x29,0x60,0x43,0x60,0x6c,0xae,
0x5a,0x40,0x00,0x00,0x00,0xff,0xff,0xb8,0x0c,0x60,0x76,0x5e,
0x00,0x00,0x00,
	}))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes()
}

func app_json_mustache() []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte{
0x1f,0x8b,0x08,0x00,0x00,0x09,0x6e,0x88,0x00,0xff,0xaa,0xe6,
0xe2,0xe4,0xaa,0x05,0x04,0x00,0x00,0xff,0xff,0x23,0x18,0x17,
0xfc,0x05,0x00,0x00,0x00,
	}))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes()
}

func actionfile_go_mustache() []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte{
0x1f,0x8b,0x08,0x00,0x00,0x09,0x6e,0x88,0x00,0xff,0x2a,0x48,
0x4c,0xce,0x4e,0x4c,0x4f,0x55,0x48,0x4c,0x2e,0xc9,0xcc,0xcf,
0x2b,0xe6,0xe2,0xca,0xcc,0x2d,0xc8,0x2f,0x2a,0x51,0xd0,0xe0,
0xe2,0x54,0x4a,0xcf,0x2c,0xc9,0x28,0x4d,0xd2,0x4b,0xce,0xcf,
0xd5,0xcf,0x2d,0x2d,0xaa,0xd2,0x4f,0x4d,0xcf,0xd7,0x87,0xaa,
0x53,0xc2,0x21,0x9d,0x51,0x52,0x52,0xa0,0xc4,0xa5,0xc9,0x05,
0x08,0x00,0x00,0xff,0xff,0x0b,0xc9,0xe7,0xd0,0x57,0x00,0x00,
0x00,
	}))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes()
}

func action_go_mustache() []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte{
0x1f,0x8b,0x08,0x00,0x00,0x09,0x6e,0x88,0x00,0xff,0x3c,0x8e,
0xc1,0xaa,0x83,0x30,0x10,0x45,0xd7,0xc9,0x57,0x0c,0x2e,0x1e,
0xfa,0x10,0xdd,0x0b,0x5d,0x14,0xda,0x45,0x17,0xd5,0x22,0xfd,
0x81,0xa0,0xd3,0x46,0xd0,0xa8,0xc9,0xa4,0x14,0xc2,0xfc,0x7b,
0xad,0x4a,0x77,0x77,0x0e,0x77,0x0e,0xf7,0xa5,0x2c,0x84,0x50,
0xaa,0x01,0x99,0xe1,0x00,0xaa,0xa1,0x6e,0x34,0x2e,0xab,0xf1,
0xd9,0x39,0x42,0x1b,0xff,0x69,0xa2,0x29,0x3b,0xae,0x38,0x48,
0x71,0x53,0xa4,0x0b,0x88,0x42,0xf8,0x06,0xe6,0x28,0x95,0xe2,
0x8a,0xa4,0xc7,0x76,0x85,0x5b,0xdc,0xf0,0xf9,0x8d,0x8d,0x27,
0x2c,0xe0,0xe1,0x4d,0x03,0xb1,0xc5,0x19,0xfe,0x57,0x57,0x8d,
0xb3,0x47,0x47,0x09,0xec,0x97,0xf3,0x3d,0xc1,0xa2,0x16,0x79,
0x0e,0xf7,0xea,0x54,0x15,0x70,0x19,0xa6,0x1e,0x07,0x34,0x04,
0xa4,0x3b,0xb7,0x6f,0xca,0x96,0x86,0x45,0xf2,0xd6,0x6c,0x8f,
0xe5,0x48,0xbf,0x1e,0xb6,0x52,0x70,0x2a,0x39,0xf9,0x04,0x00,
0x00,0xff,0xff,0x71,0xb1,0xcf,0x2d,0xce,0x00,0x00,0x00,
	}))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes()
}

func error_html_mustache() []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte{
0x1f,0x8b,0x08,0x00,0x00,0x09,0x6e,0x88,0x00,0xff,0x54,0x51,
0xc1,0x6e,0xeb,0x20,0x10,0x3c,0xfb,0x7d,0x05,0xcf,0xbd,0xb4,
0x92,0x1d,0xdb,0x72,0xdd,0x83,0x43,0x23,0xf5,0xd6,0x4b,0xfb,
0x0f,0x1b,0xd8,0x60,0x54,0x0c,0x16,0x90,0x36,0x6e,0x94,0x7f,
0x2f,0x98,0xb4,0x49,0x64,0x09,0x34,0x33,0xec,0x7a,0x76,0x96,
0xfe,0xe7,0x86,0xf9,0x79,0x42,0x32,0xf8,0x51,0x6d,0xfe,0xd1,
0x74,0x65,0x74,0x40,0xe0,0xe1,0xce,0xa8,0x97,0x5e,0xe1,0xe6,
0x78,0x7c,0x43,0xe7,0x40,0xe0,0xe9,0x44,0xab,0x44,0x45,0xd1,
0xf9,0x59,0x21,0x89,0xf5,0xcf,0xb9,0xc7,0x83,0xaf,0x98,0x73,
0x79,0x54,0xb2,0xd8,0xa7,0x20,0x5b,0xc3,0x67,0x72,0x8c,0x38,
0xdb,0x02,0xfb,0x10,0xd6,0xec,0x35,0xef,0xc9,0x1d,0x22,0xae,
0x17,0x76,0x02,0xce,0xa5,0x16,0x3d,0xa9,0x13,0x1e,0xc1,0x0a,
0xa9,0x7f,0xe1,0x29,0x1e,0x57,0x3d,0x76,0x46,0xfb,0x72,0x07,
0xa3,0x54,0x73,0x4f,0xf2,0x57,0x54,0x9f,0xe8,0x25,0x03,0xf2,
0x8e,0x7b,0xcc,0x8b,0x17,0x2b,0x41,0x15,0x0e,0xb4,0x2b,0x1d,
0x5a,0xb9,0x4b,0x1d,0x99,0x51,0xc6,0x86,0x5f,0xb6,0xf0,0xd8,
0x76,0xed,0x2d,0xd7,0x75,0x5d,0x22,0xa2,0xf7,0xd2,0x0d,0xc0,
0xcd,0x57,0x5f,0x93,0x76,0x3a,0x90,0xa6,0x0e,0x87,0x15,0x5b,
0xb8,0xaf,0x8b,0xe5,0x5b,0xb5,0x0f,0x57,0x6f,0x41,0x49,0x11,
0x6c,0x32,0xd4,0x1e,0xed,0xc5,0xeb,0xd0,0x9c,0x9d,0x2a,0xf4,
0x41,0x28,0xdd,0x04,0x6c,0x19,0xaf,0xac,0x57,0x0d,0x8e,0xeb,
0xcb,0x14,0x4e,0x7e,0x63,0xdf,0x3d,0x4d,0x87,0xbf,0x62,0x5a,
0x2d,0x69,0xc6,0xf0,0xab,0x73,0xfa,0x34,0xce,0xbe,0x04,0x3d,
0x34,0xb7,0x2b,0x08,0x38,0xbe,0x4b,0x7a,0x80,0x71,0x6b,0x3f,
0x01,0x00,0x00,0xff,0xff,0x61,0xd4,0x4f,0xc1,0xcc,0x01,0x00,
0x00,
	}))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes()
}

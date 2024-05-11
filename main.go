package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/mermaid"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

const CODE_BLOCK_STYLE = "solarized-dark"
const CSS_STYLE = `
  <style>
    div.container { 
      max-width: 80%;
      margin: auto;
      margin-top: 30px;
      padding: 20px;
      border: 1px solid #ccc;
      box-shadow: 0 0 10px rgba(0,0,0,0.1); 
    }
    /* Styling the navigation bar */
    .navbar {
        overflow: hidden;
        background-color: #333;
        font-family: Arial, sans-serif;
    }

    /* Navigation links style */
    .navbar a {
        float: left;
        display: block;
        color: white;
        text-align: center;
        padding: 14px 20px;
        text-decoration: none;
    }

    /* Change the color of links on hover */
    .navbar a:hover {
        background-color: #ddd;
        color: black;
    }
  </style>
`

func main() {
	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildInput := buildCmd.String("input", "pages/", "Directory of Markdown files")
	buildOutput := buildCmd.String("output", "build/web/", "Output directory for HTML files")

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	serveInput := serveCmd.String("input", "build/web/", "Directory to serve")
	servePort := serveCmd.String("port", "8080", "Port to serve on")

	if len(os.Args) < 2 {
		fmt.Println("expected 'build' or 'serve' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "build":
		buildCmd.Parse(os.Args[2:])
		buildSite(*buildInput, *buildOutput)
	case "serve":
		serveCmd.Parse(os.Args[2:])
		serveSite(*serveInput, *servePort)
	default:
		fmt.Println("expected 'build' or 'serve' subcommands")
		os.Exit(1)
	}
}

func buildSite(inputDir, outputDir string) {
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle(CODE_BLOCK_STYLE),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
			&mermaid.Extender{},
		),
	)

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".md" {
			mdContent, err := ioutil.ReadFile(filepath.Join(inputDir, file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			buffer := bytes.NewBuffer(nil)
			if err := markdown.Convert(mdContent, buffer); err != nil {
				log.Fatal(err)
			}

			htmlContent := fmt.Sprintf(`
        <html>
          <head>%s</head>
          <body>
           <div class='navbar'>
            <a href='/'>Home</a>
           </div>
           <div class='container'>
            %s
           </div>
          </body>
        </html>`,
				CSS_STYLE,
				buffer.String())

			newFilename := filepath.Base(file.Name())
			newFilename = newFilename[:len(newFilename)-len(filepath.Ext(newFilename))] + ".html"
			err = ioutil.WriteFile(filepath.Join(outputDir, newFilename), []byte(htmlContent), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func serveSite(directory, port string) {
	fs := http.FileServer(http.Dir(directory))
	http.Handle("/", fs)

	log.Println("Serving on http://localhost:" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

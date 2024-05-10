### Main structure

```
main.go # Main entrypoint with the commands to build the web, serve it, etc.
internal/
    - markdownProcessor.go
    - markdownProcessor_test.go
    - webBuilder.go
    - webBuilder_test.go
```

### Code ideas

More or less should be like this:

```go
// List all pages in the pages directory
pages := []string{"examplePageA", "examplePageB"}

// Process the markdown files and build the web
// If the page is named as index, it will be the main page
for _, page := range pages {
    htmlContent := markdownProcessor.Process(page + ".md")
    webBuilder.Build(page, htmlContent)
}
```

### Using

I think would be good to use it like: 
```
# To build it
tinycms build --input pages/ --output build/web/

# To serve it
tinycms serve --input build/web/ --port 8080

# The final files should be something like

pages/ # The input raw md files
    - index.md
    - examplePageA.md
    - mediaForPageA.jpg
    - examplePageB.md
    - mediaForPageB.png
build/ # The final converted to html files
    web/ # The final web build for your website
        - index.html
        - examplePageA.html
        - mediaForPageA.jpg
        - examplePageB.html
        - mediaForPageB.png
```

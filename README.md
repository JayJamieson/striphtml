# striphtml

Fork of [jaytaylor/html2text](https://github.com/jaytaylor/html2text) with changes to add more extensive CLI and HTTP server built around it.

## Introduction

Ensure your emails are readable by all!

Turns HTML into raw text, useful for sending fancy HTML emails with an equivalently nicely formatted TXT document as a fallback.

striphtml is a simple golang cli/server for rendering HTML into plaintext.

It requires go 1.21 or newer.

## Download the package

```bash
go get github.com/JayJamieson/striphtml
```

## Example usage

### Library

```go
package main

import (
 "fmt"

 "github.com/JayJamieson/striphtml"
)

func main() {
 inputHTML := `
<html>
  <head>
    <title>My Mega Service</title>
    <link rel=\"stylesheet\" href=\"main.css\">
    <style type=\"text/css\">body { color: #fff; }</style>
  </head>

  <body>
    <div class="logo">
      <a href="http://example.com/"><img src="/logo-image.jpg" alt="Mega Service"/></a>
    </div>

    <h1>Welcome to your new account on my service!</h1>

    <p>
      Here is some more information:

      <ul>
        <li>Link 1: <a href="https://example.com">Example.com</a></li>
        <li>Link 2: <a href="https://example2.com">Example2.com</a></li>
        <li>Something else</li>
      </ul>
    </p>

    <table>
      <thead>
        <tr><th>Header 1</th><th>Header 2</th></tr>
      </thead>
      <tfoot>
        <tr><td>Footer 1</td><td>Footer 2</td></tr>
      </tfoot>
      <tbody>
        <tr><td>Row 1 Col 1</td><td>Row 1 Col 2</td></tr>
        <tr><td>Row 2 Col 1</td><td>Row 2 Col 2</td></tr>
      </tbody>
    </table>
  </body>
</html>`

 text, err := striphtml.FromString(inputHTML, striphtml.Options{PrettyTables: true})
 if err != nil {
  panic(err)
 }
 fmt.Println(text)
}
```

Output:

```txt
Mega Service ( http://example.com/ )

******************************************
Welcome to your new account on my service!
******************************************

Here is some more information:

* Link 1: Example.com ( https://example.com )
* Link 2: Example2.com ( https://example2.com )
* Something else

+-------------+-------------+
|  HEADER 1   |  HEADER 2   |
+-------------+-------------+
| Row 1 Col 1 | Row 1 Col 2 |
| Row 2 Col 1 | Row 2 Col 2 |
+-------------+-------------+
|  FOOTER 1   |  FOOTER 2   |
+-------------+-------------+
```

### Command line

Read HTML from stdin and write plain text to stdout.

```shell
echo '<div>hi</div>' | striphtml
```

As HTTP server.

```shell
striphtml server
```

## Unit-tests

Running the unit-tests is straightforward and standard:

```bash
go test
```

## License

Permissive MIT license.

## Alternatives

- <https://github.com/jaytaylor/html2text> - Original implementation
- <https://github.com/k3a/html2text> - Lightweight

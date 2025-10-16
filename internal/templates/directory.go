package templates

import "html/template"

var DirectoryTemplate = template.Must(template.New("directory").ParseFiles("internal/templates/directory.html"))

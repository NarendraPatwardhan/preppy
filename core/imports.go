package core

import (
	"strings"
)

func GetImports(source []byte) ([]string, error) {
	imports := []string{}
	importQuery := `
		(import_statement (aliased_import (dotted_name) @m))
		(import_statement (dotted_name) @m)
		(import_from_statement module_name: (dotted_name) @m)
	`
	treesitter, error := NewTreeSitter(source, []byte(importQuery))
	if error != nil {
		return imports, error
	}
	cursor := treesitter.Exec()
	matchedNode, _ := cursor.NextMatch()
	for matchedNode != nil {
		captures := matchedNode.Captures
		if len(captures) > 0 {
			for _, capture := range captures {
				start := capture.Node.StartByte()
				end := capture.Node.EndByte()
				required := string(source[start:end])
				required = strings.Split(required, ".")[0]
				imports = append(imports, required)
			}
		}
		matchedNode, _ = cursor.NextMatch()
	}
	return imports, nil
}

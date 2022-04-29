package core

import (
	"context"

	ts "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/python"
)

func parse(source []byte) (*ts.Node, error) {
	ctx := context.Background()
	return ts.ParseCtx(ctx, source, python.GetLanguage())
}

func newQuery(queryText []byte) (*ts.Query, error) {
	return ts.NewQuery(queryText, python.GetLanguage())
}

func newCursor() *ts.QueryCursor {
	return ts.NewQueryCursor()
}

type TreeSitter struct {
	Root  *ts.Node
	Query *ts.Query
}

func NewTreeSitter(source []byte, queryText []byte) (TreeSitter, error) {
	root, err := parse(source)
	if err != nil {
		return TreeSitter{}, err
	}
	query, err := newQuery(queryText)
	if err != nil {
		return TreeSitter{}, err
	}
	return TreeSitter{
		Root:  root,
		Query: query,
	}, nil
}

func (t TreeSitter) Exec() *ts.QueryCursor {
	cursor := newCursor()
	cursor.Exec(t.Query, t.Root)
	return cursor
}

package holddb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert_AddsEdgeToRoot(t *testing.T) {
	// Just inserting a single value
	tree := New()

	tree.Insert([]byte("foobar"), "myval")

	// Checking the final state
	wanted := &Tree{
		root: &node{
			edges: []*edge{
				{
					label: "f"[0],
					node: &node{
						prefix: []byte("foobar"),
						leaf: &leaf{
							key:   "foobar",
							value: "myval",
						},
					},
				},
			},
		},
	}
	assert.Equal(t, wanted, tree)
}

func TestInsert_UpdatesLeaf(t *testing.T) {
	tree := New()

	// Updating an existing value
	tree.Insert([]byte("foobar"), "myval")
	tree.Insert([]byte("foobar"), "myval2")

	// Checking the final state
	wanted := &Tree{
		root: &node{
			edges: []*edge{
				{
					label: "f"[0],
					node: &node{
						prefix: []byte("foobar"),
						leaf: &leaf{
							key:   "foobar",
							value: "myval2",
						},
					},
				},
			},
		},
	}
	assert.Equal(t, wanted, tree)
}

func TestInsert_AddsSortedEdges(t *testing.T) {
	tree := New()

	tree.Insert([]byte("foobar"), "myval")
	tree.Insert([]byte("casbar"), "myval2")
	tree.Insert([]byte("honk"), "myval3")

	// Checking the final state
	wanted := &Tree{
		root: &node{
			edges: []*edge{
				{
					label: "c"[0],
					node: &node{
						prefix: []byte("casbar"),
						leaf: &leaf{
							key:   "casbar",
							value: "myval2",
						},
					},
				},
				{
					label: "f"[0],
					node: &node{
						prefix: []byte("foobar"),
						leaf: &leaf{
							key:   "foobar",
							value: "myval",
						},
					},
				},
				{
					label: "h"[0],
					node: &node{
						prefix: []byte("honk"),
						leaf: &leaf{
							key:   "honk",
							value: "myval3",
						},
					},
				},
			},
		},
	}
	assert.Equal(t, wanted, tree)
}

func TestInsert_SplitsNode(t *testing.T) {
	// Just inserting a single value
	tree := New()

	tree.Insert([]byte("foobar"), "myval")
	tree.Insert([]byte("foobaz"), "myval2")

	// Checking the final state
	wanted := &Tree{
		root: &node{
			edges: []*edge{
				{
					label: "f"[0],
					node: &node{
						prefix: []byte("fooba"),
						edges: []*edge{
							{
								label: "r"[0],
								node: &node{
									prefix: []byte("r"),
									leaf: &leaf{
										key:   "foobar",
										value: "myval",
									},
								},
							},
							{
								label: "z"[0],
								node: &node{
									prefix: []byte("z"),
									leaf: &leaf{
										key:   "foobaz",
										value: "myval2",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, wanted, tree)
}

func TestInsert_SplitsNodeMore(t *testing.T) {
	// Just inserting a single value
	tree := New()

	tree.Insert([]byte("foobar"), "myval")
	tree.Insert([]byte("foobaz"), "myval2")
	tree.Insert([]byte("foobaar"), "myval3")

	// Checking the final state
	wanted := &Tree{
		root: &node{
			edges: []*edge{
				{
					label: "f"[0],
					node: &node{
						prefix: []byte("fooba"),
						edges: []*edge{
							{
								label: "a"[0],
								node: &node{
									prefix: []byte("ar"),
									leaf: &leaf{
										key:   "foobaar",
										value: "myval3",
									},
								},
							},
							{
								label: "r"[0],
								node: &node{
									prefix: []byte("r"),
									leaf: &leaf{
										key:   "foobar",
										value: "myval",
									},
								},
							},
							{
								label: "z"[0],
								node: &node{
									prefix: []byte("z"),
									leaf: &leaf{
										key:   "foobaz",
										value: "myval2",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, wanted, tree)
}

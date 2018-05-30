package graphapi

import (
	"reflect"
)

// ImplementsGraph defines and identifiable node entity
type ImplementsGraph interface {
}

// ImplementNode define an entity which can have dependencies
type ImplementNode interface {
	GetID() string
	GetName() string
	GetType() string
	GetEdges() []ImplementEdge
}

// ImplementEdge defines and identifiable node entity
type ImplementEdge interface {
	GetID() string
	GetFrom() string
	GetTo() string
	GetType() string
	GetProperties() interface{}
}

// Repository defines a normal repository
type Repository interface {
	FindAll() (*Graph, error)
	FindAllActive() (*Graph, error)
	GetType() reflect.Type
	// Map api entity to graph entity
	Map(entity interface{}, others map[string]ImplementNode) ImplementNode
	// Resolve entity dependencies
	Resolve(entity interface{}, mappers map[string]ImplementNode)
}

// Graph define a graph resource
type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Node define a node resource
type Node struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
}

// GetID get source entity
func (n *Node) GetID() string {
	return n.ID
}

// GetName get name
func (n *Node) GetName() string {
	return n.Name
}

// GetType get name
func (n *Node) GetType() string {
	return n.Type
}

// Edge define a edge resource
type Edge struct {
	ID         string      `json:"id"`
	From       string      `json:"from"`
	To         string      `json:"to"`
	Type       string      `json:"type"`
	Properties interface{} `json:"properties"`
}

// GetID get edge entity
func (e *Edge) GetID() string {
	return e.ID
}

// GetFrom get edge entity
func (e *Edge) GetFrom() string {
	return e.From
}

// GetTo get edge entity
func (e *Edge) GetTo() string {
	return e.To
}

// GetType get edge type
func (e *Edge) GetType() string {
	return e.Type
}

// GetProperties get edge entity
func (e *Edge) GetProperties() interface{} {
	return e.Properties
}

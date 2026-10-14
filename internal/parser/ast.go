package parser

// Pos records a position in source code.
type Pos struct {
	Line   int // 0-based
	Column int // 0-based
}

// Node is the interface implemented by all AST nodes.
type Node interface {
	nodePos() Pos
}

// Program is the top-level AST node.
type Program struct {
	Stmts []Node
}

func (p *Program) nodePos() Pos {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].nodePos()
	}
	return Pos{}
}

// FunctionDef represents a def...end block.
type FunctionDef struct {
	Pos          Pos
	Name         string
	Params       []string
	RequiredArgs int // number of params without default values
	Body         []Node
}

func (f *FunctionDef) nodePos() Pos { return f.Pos }

// ThreadDef represents a thread...end block.
type ThreadDef struct {
	Pos  Pos
	Name string
	Body []Node
}

func (t *ThreadDef) nodePos() Pos { return t.Pos }

// IfStmt represents if/elif/else/end.
type IfStmt struct {
	Pos       Pos
	Condition Node
	Body      []Node
	ElseIfs   []ElseIfClause
	Else      []Node
}

func (i *IfStmt) nodePos() Pos { return i.Pos }

// ElseIfClause is an elif branch.
type ElseIfClause struct {
	Pos       Pos
	Condition Node
	Body      []Node
}

// WhileStmt represents while...end.
type WhileStmt struct {
	Pos       Pos
	Condition Node
	Body      []Node
}

func (w *WhileStmt) nodePos() Pos { return w.Pos }

// ReturnStmt represents a return statement.
type ReturnStmt struct {
	Pos   Pos
	Value Node // may be nil
}

func (r *ReturnStmt) nodePos() Pos { return r.Pos }

// Assignment represents var = expr.
type Assignment struct {
	Pos    Pos
	Target Node
	Value  Node
}

func (a *Assignment) nodePos() Pos { return a.Pos }

// BinaryExpr represents a op b.
type BinaryExpr struct {
	Pos   Pos
	Op    string
	Left  Node
	Right Node
}

func (b *BinaryExpr) nodePos() Pos { return b.Pos }

// UnaryExpr represents a prefix unary operator.
type UnaryExpr struct {
	Pos     Pos
	Op      string
	Operand Node
}

func (u *UnaryExpr) nodePos() Pos { return u.Pos }

// CallExpr represents func(args...).
type CallExpr struct {
	Pos  Pos
	Func Node
	Args []Node
}

func (c *CallExpr) nodePos() Pos { return c.Pos }

// IndexExpr represents expr[index].
type IndexExpr struct {
	Pos   Pos
	Expr  Node
	Index Node
}

func (i *IndexExpr) nodePos() Pos { return i.Pos }

// Identifier is a variable or function name.
type Identifier struct {
	Pos  Pos
	Name string
}

func (id *Identifier) nodePos() Pos { return id.Pos }

// NumberLit is a numeric literal.
type NumberLit struct {
	Pos   Pos
	Value string
}

func (n *NumberLit) nodePos() Pos { return n.Pos }

// StringLit is a string literal.
type StringLit struct {
	Pos   Pos
	Value string
}

func (s *StringLit) nodePos() Pos { return s.Pos }

// BoolLit is a boolean literal (True / False).
type BoolLit struct {
	Pos   Pos
	Value bool
}

func (b *BoolLit) nodePos() Pos { return b.Pos }

// ListLit represents [a, b, c].
type ListLit struct {
	Pos      Pos
	Elements []Node
}

func (l *ListLit) nodePos() Pos { return l.Pos }

// HaltStmt represents a bare halt statement.
type HaltStmt struct {
	Pos Pos
}

func (h *HaltStmt) nodePos() Pos { return h.Pos }

// ExprStmt wraps an expression used as a statement (e.g. a function call).
type ExprStmt struct {
	Expr Node
}

func (e *ExprStmt) nodePos() Pos { return e.Expr.nodePos() }

// GlobalStmt represents "global var = expr".
type GlobalStmt struct {
	Pos  Pos
	Name string
	Init Node // may be nil
}

func (g *GlobalStmt) nodePos() Pos { return g.Pos }

// LocalStmt represents "local var = expr".
type LocalStmt struct {
	Pos  Pos
	Name string
	Init Node // may be nil
}

func (l *LocalStmt) nodePos() Pos { return l.Pos }

// NodePos returns the position of any AST node.
func NodePos(n Node) Pos {
	return n.nodePos()
}

// ParseError records a single parse error with position.
type ParseError struct {
	Line    int
	Column  int
	Message string
}

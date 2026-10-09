package parser

import (
	"testing"
)

func TestParseEmptyProgram(t *testing.T) {
	prog, errs := Parse("")
	if len(errs) != 0 {
		t.Errorf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 0 {
		t.Errorf("expected 0 statements, got %d", len(prog.Stmts))
	}
}

func TestParseFunctionDef(t *testing.T) {
	input := `def add(a, b):
  return a + b
end`
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Stmts))
	}

	fn, ok := prog.Stmts[0].(*FunctionDef)
	if !ok {
		t.Fatalf("expected *FunctionDef, got %T", prog.Stmts[0])
	}
	if fn.Name != "add" {
		t.Errorf("expected name 'add', got %q", fn.Name)
	}
	if len(fn.Params) != 2 {
		t.Errorf("expected 2 params, got %d", len(fn.Params))
	}
	if len(fn.Body) != 1 {
		t.Fatalf("expected 1 body statement, got %d", len(fn.Body))
	}
	ret, ok := fn.Body[0].(*ReturnStmt)
	if !ok {
		t.Fatalf("expected *ReturnStmt, got %T", fn.Body[0])
	}
	binOp, ok := ret.Value.(*BinaryExpr)
	if !ok {
		t.Fatalf("expected *BinaryExpr, got %T", ret.Value)
	}
	if binOp.Op != "+" {
		t.Errorf("expected op '+', got %q", binOp.Op)
	}
}

func TestParseIfElifElse(t *testing.T) {
	input := `if x > 0:
  y = 1
elif x == 0:
  y = 0
else:
  y = -1
end`
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Stmts))
	}

	ifStmt, ok := prog.Stmts[0].(*IfStmt)
	if !ok {
		t.Fatalf("expected *IfStmt, got %T", prog.Stmts[0])
	}
	if len(ifStmt.Body) != 1 {
		t.Errorf("expected 1 if body statement, got %d", len(ifStmt.Body))
	}
	if len(ifStmt.ElseIfs) != 1 {
		t.Errorf("expected 1 elif, got %d", len(ifStmt.ElseIfs))
	}
	if len(ifStmt.Else) != 1 {
		t.Errorf("expected 1 else body statement, got %d", len(ifStmt.Else))
	}
}

func TestParseWhile(t *testing.T) {
	input := `while i < 10:
  i = i + 1
end`
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Stmts))
	}

	w, ok := prog.Stmts[0].(*WhileStmt)
	if !ok {
		t.Fatalf("expected *WhileStmt, got %T", prog.Stmts[0])
	}
	if len(w.Body) != 1 {
		t.Errorf("expected 1 body statement, got %d", len(w.Body))
	}
}

func TestParseAssignment(t *testing.T) {
	input := "x = 42"
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Stmts))
	}

	assign, ok := prog.Stmts[0].(*Assignment)
	if !ok {
		t.Fatalf("expected *Assignment, got %T", prog.Stmts[0])
	}
	target, ok := assign.Target.(*Identifier)
	if !ok {
		t.Fatalf("expected *Identifier target, got %T", assign.Target)
	}
	if target.Name != "x" {
		t.Errorf("expected target 'x', got %q", target.Name)
	}
	val, ok := assign.Value.(*NumberLit)
	if !ok {
		t.Fatalf("expected *NumberLit value, got %T", assign.Value)
	}
	if val.Value != "42" {
		t.Errorf("expected value '42', got %q", val.Value)
	}
}

func TestParseFunctionCall(t *testing.T) {
	input := `movej(q, a=1.4, v=1.05)`
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Stmts))
	}

	exprStmt, ok := prog.Stmts[0].(*ExprStmt)
	if !ok {
		t.Fatalf("expected *ExprStmt, got %T", prog.Stmts[0])
	}
	call, ok := exprStmt.Expr.(*CallExpr)
	if !ok {
		t.Fatalf("expected *CallExpr, got %T", exprStmt.Expr)
	}
	fn, ok := call.Func.(*Identifier)
	if !ok {
		t.Fatalf("expected *Identifier func, got %T", call.Func)
	}
	if fn.Name != "movej" {
		t.Errorf("expected func 'movej', got %q", fn.Name)
	}
	if len(call.Args) != 3 {
		t.Errorf("expected 3 args, got %d", len(call.Args))
	}
}

func TestParseGlobalLocal(t *testing.T) {
	input := "global x = 10\nlocal y = 20"
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(prog.Stmts))
	}

	g, ok := prog.Stmts[0].(*GlobalStmt)
	if !ok {
		t.Fatalf("expected *GlobalStmt, got %T", prog.Stmts[0])
	}
	if g.Name != "x" {
		t.Errorf("expected name 'x', got %q", g.Name)
	}

	l, ok := prog.Stmts[1].(*LocalStmt)
	if !ok {
		t.Fatalf("expected *LocalStmt, got %T", prog.Stmts[1])
	}
	if l.Name != "y" {
		t.Errorf("expected name 'y', got %q", l.Name)
	}
}

func TestParseBooleanExpr(t *testing.T) {
	input := "x = True and not False"
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Stmts))
	}

	assign, ok := prog.Stmts[0].(*Assignment)
	if !ok {
		t.Fatalf("expected *Assignment, got %T", prog.Stmts[0])
	}
	andExpr, ok := assign.Value.(*BinaryExpr)
	if !ok {
		t.Fatalf("expected *BinaryExpr, got %T", assign.Value)
	}
	if andExpr.Op != "and" {
		t.Errorf("expected op 'and', got %q", andExpr.Op)
	}
}

func TestParseList(t *testing.T) {
	input := "x = [1, 2, 3]"
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Stmts))
	}

	assign, ok := prog.Stmts[0].(*Assignment)
	if !ok {
		t.Fatalf("expected *Assignment, got %T", prog.Stmts[0])
	}
	list, ok := assign.Value.(*ListLit)
	if !ok {
		t.Fatalf("expected *ListLit, got %T", assign.Value)
	}
	if len(list.Elements) != 3 {
		t.Errorf("expected 3 elements, got %d", len(list.Elements))
	}
}

func TestParseIndexExpr(t *testing.T) {
	input := "x = arr[0]"
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Stmts))
	}

	assign, ok := prog.Stmts[0].(*Assignment)
	if !ok {
		t.Fatalf("expected *Assignment, got %T", prog.Stmts[0])
	}
	idx, ok := assign.Value.(*IndexExpr)
	if !ok {
		t.Fatalf("expected *IndexExpr, got %T", assign.Value)
	}
	_ = idx
}

func TestParseError(t *testing.T) {
	input := `def broken(:
end`
	_, errs := Parse(input)
	if len(errs) == 0 {
		t.Error("expected parse errors, got none")
	}
}

func TestParseThread(t *testing.T) {
	input := `thread myThread():
  sleep(1)
end`
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Stmts))
	}

	th, ok := prog.Stmts[0].(*ThreadDef)
	if !ok {
		t.Fatalf("expected *ThreadDef, got %T", prog.Stmts[0])
	}
	if th.Name != "myThread" {
		t.Errorf("expected name 'myThread', got %q", th.Name)
	}
	if len(th.Body) != 1 {
		t.Errorf("expected 1 body statement, got %d", len(th.Body))
	}
}

func TestParseForLoop(t *testing.T) {
	input := `for i in range(10):
  x = i
end`
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Stmts))
	}

	forStmt, ok := prog.Stmts[0].(*ForStmt)
	if !ok {
		t.Fatalf("expected *ForStmt, got %T", prog.Stmts[0])
	}
	if forStmt.Var != "i" {
		t.Errorf("expected var 'i', got %q", forStmt.Var)
	}
}

func TestParseComplexProgram(t *testing.T) {
	input := `def main():
  global speed = 1.05
  target = p[0.5, 0.3, 0.1, 0, 3.14, 0]
  movej(target, a=1.4, v=speed)
  if get_digital_in(0) == True:
    set_digital_out(1, True)
  else:
    set_digital_out(1, False)
  end
  i = 0
  while i < 5:
    movel(target)
    i = i + 1
  end
end`
	prog, errs := Parse(input)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
	if len(prog.Stmts) != 1 {
		t.Fatalf("expected 1 top-level statement, got %d", len(prog.Stmts))
	}

	fn, ok := prog.Stmts[0].(*FunctionDef)
	if !ok {
		t.Fatalf("expected *FunctionDef, got %T", prog.Stmts[0])
	}
	if fn.Name != "main" {
		t.Errorf("expected name 'main', got %q", fn.Name)
	}
}

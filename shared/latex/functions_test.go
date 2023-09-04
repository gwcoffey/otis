package latex

import "testing"

func TestCommand(t *testing.T) {
	expectCommand(t, `\foo`, Command("foo", nil, nil))
	expectCommand(t, `\foo[a]`, Command("foo", []string{"a"}, nil))
	expectCommand(t, `\foo[aaa,bbb]`, Command("foo", []string{"aaa", "bbb"}, nil))
	expectCommand(t, `\foo{a}`, Command("foo", nil, []string{"a"}))
	expectCommand(t, `\foo{aaa,bbb}`, Command("foo", nil, []string{"aaa", "bbb"}))
	expectCommand(t, `\foo[a]{b}`, Command("foo", []string{"a"}, []string{"b"}))
	expectCommand(t, `\foo[aaa,bbb]{ccc,ddd}`, Command("foo", []string{"aaa", "bbb"}, []string{"ccc", "ddd"}))

	expectCommand(t, `\foo[a{,}b,c]`, Command("foo", []string{"a,b", "c"}, nil))
	expectCommand(t, `\foo[a{]}b]`, Command("foo", []string{"a]b"}, nil))
	expectCommand(t, `\foo[a[b{]}c]`, Command("foo", []string{"a[b]c"}, nil))

	expectCommand(t, `\foo{a{,}b,c}`, Command("foo", nil, []string{"a,b", "c"}))
}

func expectCommand(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Error("expected", actual, "to equal", expected)
	}
}

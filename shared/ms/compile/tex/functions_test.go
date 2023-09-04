package tex

import "testing"

func TestNoEscapes(t *testing.T) {
	expectEscape(t, `clean`, `clean`)
}

func TestBackslashEscapes(t *testing.T) {
	expectEscape(t, `#1`, `\#1`)
	expectEscape(t, `$3.20`, `\$3.20`)
	expectEscape(t, `25%`, `25\%`)
	expectEscape(t, `A&B`, `A\&B`)
	expectEscape(t, `my_word`, `my\_word`)
	expectEscape(t, `{open`, `\{open`)
	expectEscape(t, `close}`, `close\}`)
	expectEscape(t, `mixed # with $ & {some_more}`, `mixed \# with \$ \& \{some\_more\}`)
}

func TestCommandEscapes(t *testing.T) {
	expectEscape(t, `a\b`, `a$\backslash$b`)
	expectEscape(t, `a~b`, `a$\sim$b`)
	expectEscape(t, `a^b`, `a$\textasciicircum$b`)
}

func expectEscape(t *testing.T, unescaped string, expected string) {
	it := escapeText(unescaped)
	if it != expected {
		t.Fatal("expected", it, "to equal", expected)
	}
}

func TestFormatMarkdown(t *testing.T) {
	expectFormat(t, `clean`, `clean`)
	expectFormat(t, `this *is* neat`, `this \emph{is} neat`)
	expectFormat(t, `*this is neat*`, `\emph{this is neat}`)
	expectFormat(t, `> a blockquote`, "\\begin{quotation}\na blockquote\n\\end{quotation}\n")
}

func expectFormat(t *testing.T, unformatted string, expected string) {
	it := formatMarkdown(unformatted)
	if it != expected {
		t.Error("expected", it, "to equal", expected)
	}
}

func TestCommand(t *testing.T) {
	expectCommand(t, `\foo`, command("foo", nil, nil))
	expectCommand(t, `\foo[a]`, command("foo", []string{"a"}, nil))
	expectCommand(t, `\foo[aaa,bbb]`, command("foo", []string{"aaa", "bbb"}, nil))
	expectCommand(t, `\foo{a}`, command("foo", nil, []string{"a"}))
	expectCommand(t, `\foo{aaa,bbb}`, command("foo", nil, []string{"aaa", "bbb"}))
	expectCommand(t, `\foo[a]{b}`, command("foo", []string{"a"}, []string{"b"}))
	expectCommand(t, `\foo[aaa,bbb]{ccc,ddd}`, command("foo", []string{"aaa", "bbb"}, []string{"ccc", "ddd"}))

	expectCommand(t, `\foo[a{,}b,c]`, command("foo", []string{"a,b", "c"}, nil))
	expectCommand(t, `\foo[a{]}b]`, command("foo", []string{"a]b"}, nil))
	expectCommand(t, `\foo[a[b{]}c]`, command("foo", []string{"a[b]c"}, nil))

	expectCommand(t, `\foo{a{,}b,c}`, command("foo", nil, []string{"a,b", "c"}))
}

func expectCommand(t *testing.T, expected string, actual string) {
	if actual != expected {
		t.Error("expected", actual, "to equal", expected)
	}
}

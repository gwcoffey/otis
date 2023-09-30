package text

import "testing"

func TestKebabToSentence(t *testing.T) {
	if expected, actual := "Foo bar", KebabToSentence("foo-bar"); expected != actual {
		t.Errorf("expected %v but got %v", expected, actual)
	}
	if expected, actual := "A longer example", KebabToSentence("a-longer-example"); expected != actual {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}

func TestSentenceToKebab(t *testing.T) {
	if expected, actual := "foo-bar", ToKebab("Foo bar"); expected != actual {
		t.Errorf("expected %v but got %v", expected, actual)
	}
	if expected, actual := "foo-bar", ToKebab(" Foo bar "); expected != actual {
		t.Errorf("expected %v but got %v", expected, actual)
	}
	if expected, actual := "foo-bar", ToKebab("*Foo.^%$bar__"); expected != actual {
		t.Errorf("expected %v but got %v", expected, actual)
	}

}

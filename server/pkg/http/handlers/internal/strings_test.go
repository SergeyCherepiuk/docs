package internal_test

import (
	"testing"

	"github.com/SergeyCherepiuk/docs/pkg/http/handlers/internal"
)

func TestToSentence1(t *testing.T) {
	actual := internal.ToSentence("lorem iPSuM")
	expected := "Lorem iPSuM"

	if actual != expected {
		t.Errorf(`expected: "%s", actual: "%s"`, actual, expected)
	}
}

func TestToSentence2(t *testing.T) {
	actual := internal.ToSentence(" lorem iPSuM ")
	expected := " lorem iPSuM "

	if actual != expected {
		t.Errorf(`expected: "%s", actual: "%s"`, actual, expected)
	}
}

func TestToSentence3(t *testing.T) {
	actual := internal.ToSentence("123lorem iPSuM")
	expected := "123lorem iPSuM"

	if actual != expected {
		t.Errorf(`expected: "%s", actual: "%s"`, actual, expected)
	}
}

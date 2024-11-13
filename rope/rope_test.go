package rope

import (
	"testing"
)

func TestRopeBase(t *testing.T)  {
	input := "abcdefg"
	rope := NewRope(input)
	err, s := rope.String()
	if err != nil {
		t.Fatalf("Error stringifying rope: %s", err)
	}
	err = testString(t, input, s)
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestRopeInsert(t *testing.T)  {
	type InsertTest struct {
		insertedString string
		index int
		expected string
	}
	tests := []InsertTest{{"123", 3, "abc123defg"},
		{"123", 0, "123abcdefg"},
		{"",0, "abcdefg"}}
	for _, tt := range tests {
		rope := NewRope("abcdefg")
		rope.Insert(tt.index, tt.insertedString)
		err, s := rope.String()
		if err != nil {
			t.Fatalf("Error stringifying rope: %s", err)
		}
		err = testString(t, tt.expected, s)
		if err != nil {
			t.Fatalf("%s", err)
		}
	}
}

func TestRopeDelete(t *testing.T)  {
	type DeleteTest struct {
		index int
		length int
		expected string
	}
	tests := []DeleteTest{{3, 0, "abcefg"},
		{ 0, 0, "bcdefg"},
		{0, 3, "efg"},}
	for _, tt := range tests {
		rope := NewRope("abcdefg")
		rope.Delete(tt.index, tt.length)
		err, s := rope.String()
		if err != nil {
			t.Fatalf("Error stringifying rope: %s", err)
		}
		err = testString(t, tt.expected, s)
		if err != nil {
			t.Fatalf("%s", err)
		}
	}
}

func testString(t *testing.T, s1 string, s2 string) (error) {
	t.Helper()
	if s1 != s2 {
		t.Errorf("Strings don't match %s , %s", s1, s2)
	}
	return nil
}

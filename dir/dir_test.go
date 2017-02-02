package dir

import "testing"

type TestDataCleanPath struct {
	in  string
	out string
}

func TestCleanPath(t *testing.T) {
	tests := []TestDataCleanPath{
		{"", ""},
		{"hello", "hello"},
		{"'foo'", "foo"},
		{"'foo", "foo"},
		{"foo'", "foo"},
		{"'f''o''o", "foo"},
	}

	for i, test := range tests {
		actual := CleanPath(test.in)
		if test.out != actual {
			t.Errorf("#%d: CleanPath(%s)=%s; expected %s", i, test.in, actual, test.out)
		}
	}
}

type TestDataSplitPath struct {
	in   string
	org  string
	repo string
	err  string
}

func TestSplitPath(t *testing.T) {
	tests := []TestDataSplitPath{
		{"", "", "", "A path should be: orgname/reponame, got: "},
		{"foo", "", "", "A path should be: orgname/reponame, got: foo"},
		{"foobar/", "foobar", "", ""},
		{"/foobar", "", "foobar", ""},
		{"foo/bar", "foo", "bar", ""},
		{"//foobar", "", "", ""},
		{"///foobar", "", "", ""},
		{"foo/bar/////", "foo", "bar", ""},
	}

	for i, test := range tests {
		org, repo, err := SplitPath(test.in)
		if test.org != org {
			t.Errorf("#%d: org, _ := SplitPath(%s) == %s; expected %s", i, test.in, org, test.org)
		}
		if test.repo != repo {
			t.Errorf("#%d: _, repo := SplitPath(%s) == %s; expected %s", i, test.in, repo, test.repo)
		}
		if (err == nil && test.err != "") || (err != nil && err.Error() != test.err) {
			t.Errorf("#%d: _, _, err := TestSplitPath(%s) == %v; expected %v", i, test.in, err, test.err)
		}
	}
}

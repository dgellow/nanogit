package dir

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dgellow/nanogit/settings"
)

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
			t.Errorf("#%d: org, _, _ := SplitPath(%s) == %s; expected %s", i, test.in, org, test.org)
		}
		if test.repo != repo {
			t.Errorf("#%d: _, repo, _ := SplitPath(%s) == %s; expected %s", i, test.in, repo, test.repo)
		}
		if (err == nil && test.err != "") || (err != nil && err.Error() != test.err) {
			t.Errorf("#%d: _, _, err := SplitPath(%s) == %v; expected %v", i, test.in, err, test.err)
		}
	}
}


type TestDataGetOrgDir struct {
	in string
	out string
	err string
}

func TestGetOrgDir(t *testing.T) {
	// With default data root in configuration file
	testDefault := TestDataGetOrgDir{"foobar", "", "Data root in configuration file is empty"}
	out, err := GetOrgDir(testDefault.in)
	if testDefault.out != out {
		t.Errorf("Default data root: out, _ := GetOrgDir(%s) == %s; expected %s", testDefault.in, out, testDefault.out)
	}
	if (err == nil && testDefault.err != "") || (err != nil && err.Error() != testDefault.err) {
		t.Errorf("Default data root: _, err := GetOrgDir(%s) == %v; expected %v", testDefault.in, err, testDefault.err)
	}

	// Set data root
	settings.ConfInfo.Conf.Server.DataRoot = "./dataroot"
	defer func(){settings.ConfInfo.Conf.Server.DataRoot = ""}()
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Errorf("Error when trying to get absolute path of current directory")
	}

	// Tests after a data root has been set
	tests := []TestDataGetOrgDir{
		{"", "/dataroot", ""},
		{"foo", "/dataroot/foo", ""},
		{" ", "/dataroot/ ", ""},
		{"/foo", "/dataroot/foo", ""},
		{"/foo/bar", "/dataroot/foo/bar", ""},
	}

	for i, test := range tests {
		out, err := GetOrgDir(test.in)
		testOut := filepath.Join(currentDir, test.out)
		if testOut != out {
			t.Errorf("#%d: out, _ := GetOrgDir(%s)=%s; expected %s", i, test.in, out, testOut)
		}
		if (err == nil && test.err != "") || (err != nil && err.Error() != test.err) {
			t.Errorf("#%d: _, err := GetOrgDir(%s) == %v; expected %v", i, test.in, err, test.err)
		}
	}
}

type TestDataGetRepoDir struct {
	org string
	repo string
	out string
	err string
}

func TestGetRepoDir(t *testing.T) {
	// With default data root in configuration file
	testDefault := TestDataGetRepoDir{"foo", "bar", "", "Data root in configuration file is empty"}
	out, err := GetRepoDir(testDefault.org, testDefault.repo)
	if testDefault.out != out {
		t.Errorf("Default data root: out, _ := GetRepoDir(%s, %s) == %s; expected %s", testDefault.org, testDefault.repo, out, testDefault.out)
	}
	if (err == nil && testDefault.err != "") || (err != nil && err.Error() != testDefault.err) {
		t.Errorf("Default data root: _, err := GetRepoDir(%s, %s) == %v; expected %v", testDefault.org, testDefault.repo, err, testDefault.err)
	}

	// Set data root
	settings.ConfInfo.Conf.Server.DataRoot = "./dataroot"
	defer func(){settings.ConfInfo.Conf.Server.DataRoot = ""}()
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Errorf("Error when trying to get absolute path of current directory")
	}

	// Tests after a data root has been set
	tests := []TestDataGetRepoDir{
		{"", "", "/dataroot", ""},
		{"foo", "", "/dataroot/foo", ""},
		{"foo", "bar", "/dataroot/foo/bar", ""},
		{"", "bar", "/dataroot/bar", ""},
		{"/foo/bar", "", "/dataroot/foo/bar", ""},
	}

	for i, test := range tests {
		out, err := GetRepoDir(test.org, test.repo)
		testOut := filepath.Join(currentDir, test.out)
		if testOut != out {
			t.Errorf("#%d: out, _ := GetRepoDir(%s, %s)=%s; expected %s", i, test.org, test.repo, out, testOut)
		}
		if (err == nil && test.err != "") || (err != nil && err.Error() != test.err) {
			t.Errorf("#%d: _, err := GetRepoDir(%s, %s) == %v; expected %v", i, test.org, test.repo, err, test.err)
		}
	}
}

package iniparser

import (
	"os"
	"reflect"
	"testing"
)

const validIni = `[Simple Values]
you can also use=to delimit keys from values
key=value
paces in keys=allowed

[You can use comments]
# like this
; or this
# By default only in an empty line.
# That being said, this can be customized.`

func TestLoadFromString(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected map[string]section
		err      bool
	}{
		{
			name: "test on valid INI",
			data: validIni,
			expected: map[string]section{
				"Simple Values": {
					"you can also use": "to delimit keys from values",
					"key":              "value",
					"paces in keys":    "allowed",
				},
				"You can use comments": {},
			},
			err: false,
		},
		{
			name: "empty section name",
			data: `[   ]
key=value`,
			expected: map[string]section{},
			err:      true,
		},
		{
			name: "empty key name",
			data: `[section]
 =value`,
			expected: map[string]section{},
			err:      true,
		},
		{
			name: "empty val name",
			data: `[section]
 key=  `,
			expected: map[string]section{},
			err:      true,
		},
		{
			name: "duplicate key",
			data: `[section]
 key= val1
 key=val2`,
			expected: map[string]section{},
			err:      true,
		},
	}
	for _, test := range tests {
		p := NewIniParser()
		t.Run(test.name, func(t *testing.T) {
			err := p.LoadFromString(test.data)
			if err != nil && test.err {
				return
			}
			if err == nil && test.err {
				t.Errorf("LoadFromString : error not expected , wanted error : %v , got : %v", test.err, err)
				return
			}
			assertTwoMaps(t, p.GetSections(), test.expected)
		})
	}

}

func TestLoadFromFile(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected map[string]section
		err      bool
	}{
		{
			name:     "test on valid INI",
			filePath: "./testdata/validini.ini",
			expected: map[string]section{
				"Simple Values": {
					"you can also use": "to delimit keys from values",
					"key":              "value",
					"paces in keys":    "allowed",
				},
				"You can use comments": {},
			},
			err: false,
		},
		{
			name:     "empty section name",
			filePath: "./testdata/emptySec.ini",
			expected: map[string]section{},
			err:      true,
		},
		{
			name:     "empty key name",
			filePath: "./testdata/emptyKey.ini",
			expected: map[string]section{},
			err:      true,
		},
		{
			name:     "empty val name",
			filePath: "./testdata/emptyVal.ini",
			expected: map[string]section{},
			err:      true,
		},
		{
			name:     "duplicate key",
			filePath: "./testdata/duplicateKey.ini",
			expected: map[string]section{},
			err:      true,
		},
	}
	for _, test := range tests {
		p := NewIniParser()
		t.Run(test.name, func(t *testing.T) {
			err := p.LoadFromFile(test.filePath)
			if (err != nil) && test.err {
				return

			} else if err == nil && test.err {
				t.Errorf("LoadFromString : error not expected , wanted error : %v , got : %v", test.err, err)

			}
			assertTwoMaps(t, p.GetSections(), test.expected)
		})
	}
}

func TestGetSectionNames(t *testing.T) {

	tests := []struct {
		name     string
		data     string
		expected []string
	}{
		{
			name:     "non-empty sections",
			data:     validIni,
			expected: []string{"Simple Values", "You can use comments"},
		},
		{
			name:     "empty sections",
			data:     "",
			expected: []string{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewIniParser()
			gotSections := p.GetSectionNames()
			if reflect.DeepEqual(gotSections, test.expected) {
				t.Errorf("GetSectionNames : expected: %v , got : %v", test.expected, gotSections)
			}
		},
		)
	}

}

func TestGetSections(t *testing.T) {

	tests := []struct {
		name     string
		data     string
		expected map[string]section
	}{
		{
			name: "non-empty sections",
			data: validIni,
			expected: map[string]section{
				"Simple Values": {
					"you can also use": "to delimit keys from values",
					"key":              "value",
					"paces in keys":    "allowed",
				},
				"You can use comments": {},
			},
		},
		{
			name:     "empty sections",
			data:     "",
			expected: map[string]section{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewIniParser()
			p.LoadFromString(test.data)
			gotSections := p.GetSections()
			assertTwoMaps(t, gotSections, test.expected)
		},
		)
	}

}

func TestGet(t *testing.T) {

	tests := []struct {
		name        string
		sectionName string
		keyName     string
		expected    string
		error       bool
	}{
		{
			name:        "get value from existing section and key",
			sectionName: "Simple Values",
			keyName:     "key",
			expected:    "value",
			error:       false,
		},
		{
			name:        "get value from non existing section",
			sectionName: "any",
			keyName:     "key",
			expected:    "",
			error:       true,
		},
		{
			name:        "get value from non existing key",
			sectionName: "Simple Values",
			keyName:     "any",
			expected:    "",
			error:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewIniParser()
			err := p.LoadFromString(validIni)
			if err != nil {
				t.Errorf("GetSectionNames : can't load from file got : %v", err)
			}
			gotValue, err := p.Get(test.sectionName, test.keyName)
			if err != nil && !test.error {
				t.Errorf("GetSectionNames : expected: %v , got : %v", test.error, err)
			}
			assertStrings(t, gotValue, test.expected)
		},
		)
	}
}

func TestSet(t *testing.T) {

	tests := []struct {
		name        string
		sectionName string
		keyName     string
		value       string
		error       bool
	}{
		{
			name:        "set value in existing section and key",
			sectionName: "Simple Values",
			keyName:     "key",
			value:       "newValueTest",
			error:       false,
		},
		{
			name:        "set value in non existing section",
			sectionName: "any",
			keyName:     "key",
			value:       "",
			error:       true,
		},
		{
			name:        "set value in non existing key",
			sectionName: "Simple Values",
			keyName:     "any",
			value:       "",
			error:       true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := NewIniParser()
			err := p.LoadFromString(validIni)
			if err != nil {
				t.Errorf("String : error not expected , got : %v", err)
			}
			err = p.Set(test.sectionName, test.keyName, test.value)
			if err != nil && !test.error {
				t.Errorf("GetSectionNames : expected: %v , got : %v", test.error, err)
			}
			checkValue, _ := p.Get(test.sectionName, test.keyName)
			assertStrings(t, checkValue, test.value)
		},
		)
	}
}

func TestString(t *testing.T) {

	t.Run("Valid INI", func(t *testing.T) {
		p := NewIniParser()
		err := p.LoadFromString(validIni)
		if err != nil {
			t.Errorf("String : error not expected , got : %v", err)
		}
		got := p.String()
		err = p.LoadFromString(got)
		if err != nil {
			t.Errorf("String : error not expected , got : %v", err)
		}
		validIniExpected := `[Simple Values]
you can also use=to delimit keys from values
key=value
paces in keys=allowed
[You can use comments]
`
		assertStrings(t, got, validIniExpected)
	})
}

func TestSaveToFile(t *testing.T) {

	t.Run("Valid INI", func(t *testing.T) {
		p := NewIniParser()
		err := p.LoadFromString(validIni)
		if err != nil {
			t.Errorf("SaveToFile : error not expected , got : %v", err)
		}
		err = p.SaveToFile("./testdata/simple_example.ini")
		if err != nil {
			t.Errorf("SaveToFile : error not expected , got : %v", err)
		}
		got, err := os.ReadFile("./testdata/simple_example.ini")
		if err != nil {
			t.Errorf("SaveToFile : error not expected , got : %v", err)
		}
		validIniExpected := `[Simple Values]
you can also use=to delimit keys from values
key=value
paces in keys=allowed
[You can use comments]
`
		if reflect.DeepEqual(got, validIniExpected) {
			t.Errorf("SaveToFile : expected %v , got %v", validIniExpected, string(got))
		}
	})
}

func assertTwoMaps(t testing.TB, got, want map[string]section) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got error:\n\t%v \nwant:\n\t%v", got, want)
	}
}

func assertStrings(t testing.TB, got, want string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n\t%q \nwant:\n\t%q", got, want)
	}
}

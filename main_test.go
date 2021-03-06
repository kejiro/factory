package factory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefine(t *testing.T) {
	type User struct {
		FirstName string
		LastName  string
		Dob       time.Time
	}

	err := Define(User{}, Definition{
		"FirstName": "John",
		"LastName":  "Doe",
	})
	require.NoError(t, err)

	user := User{}
	err = Build(&user, Definition{})
	require.NoError(t, err)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.True(t, user.Dob.IsZero())
}

func TestDefineWithFunction(t *testing.T) {
	type User struct {
		FirstName  string
		LastName   string
		Dob        time.Time
		LoginCount int
	}

	counter := 1
	err := Define(&User{}, Definition{
		"FirstName": "John",
		"LastName": func() string {
			return "Doe"
		},
		"Dob": func() time.Time {
			dob, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
			return dob
		},
		"LoginCount": func() int {
			return counter
		},
	})
	require.NoError(t, err)

	user := User{}
	err = Build(&user, Definition{})
	require.NoError(t, err)

	dob, err := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
	require.NoError(t, err)

	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, counter, user.LoginCount)
	assert.Equal(t, dob, user.Dob)
}
func TestDefineWithOverride(t *testing.T) {
	type User struct {
		FirstName  string
		LastName   string
		Dob        time.Time
		LoginCount int
	}

	counter := 1
	err := Define(&User{}, Definition{
		"FirstName": "John",
		"LastName": func() string {
			return "Doe"
		},
		"Dob": func() time.Time {
			dob, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
			return dob
		},
		"LoginCount": func() int {
			return counter
		},
	})
	require.NoError(t, err)

	user := User{}
	err = Build(&user, Definition{"FirstName": "Jack"})
	require.NoError(t, err)

	dob, err := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
	require.NoError(t, err)

	assert.Equal(t, "Jack", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, counter, user.LoginCount)
	assert.Equal(t, dob, user.Dob)
}

func TestDefineWithEmbeddedStruct(t *testing.T) {
	type Model struct {
		ID uint
	}

	type User struct {
		Model
		FirstName  string
		LastName   string
		Dob        time.Time
		LoginCount int
	}

	counter := 1
	err := Define(&User{}, Definition{
		"ID":        uint(42),
		"FirstName": "John",
		"LastName": func() string {
			return "Doe"
		},
		"Dob": func() time.Time {
			dob, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
			return dob
		},
		"LoginCount": func() int {
			return counter
		},
	})
	require.NoError(t, err)

	user := User{}
	err = Build(&user, Definition{})
	require.NoError(t, err)

	assert.Equal(t, uint(42), user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
}

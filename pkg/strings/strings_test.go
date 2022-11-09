package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "not equal",
			value: "not xxx eq 1500",
			want:  false,
		},
		{
			name:  "not equal",
			value: "",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, IsEmpty(tt.value), tt.want)
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue int64
		want         int64
		hasError     bool
	}{
		{
			name:         "1 casting successfully",
			value:        "1",
			defaultValue: 10,
			want:         1,
			hasError:     false,
		},
		{
			name:         "10 casting successfully",
			value:        "10",
			defaultValue: 5,
			want:         10,
			hasError:     false,
		},
		{
			name:         "1000000 casting successfully",
			value:        "1000000",
			defaultValue: 5,
			want:         1000000,
			hasError:     false,
		},
		{
			name:         "with default value",
			value:        "",
			defaultValue: 50,
			want:         50,
			hasError:     false,
		},
		{
			name:         "with incorrect values",
			value:        "2332323wewewewwwewew",
			defaultValue: 50,
			want:         0,
			hasError:     true,
		},
		{
			name:         "with float numbers",
			value:        "25.5",
			defaultValue: 50,
			want:         0,
			hasError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ToInt64(tt.value, tt.defaultValue)

			if tt.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, val, tt.want)
		})
	}
}

func TestContainsAnyString(t *testing.T) {
	type args struct {
		value string
		lists []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "when endpoint is ping",
			args: args{
				value: "/ping",
				lists: []string{"ping", "health"},
			},
			want: true,
		},
		{
			name: "when endpoint is not ping or health",
			args: args{
				value: "/rules",
				lists: []string{"ping", "health"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsAnyString(tt.args.value, tt.args.lists...); got != tt.want {
				t.Errorf("ContainsAnyString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsStringPointerEmpty(t *testing.T) {
	var res bool
	var tests = []struct {
		input *string
		r     bool
	}{
		{StringToStringPointer("test"), false},
		{nil, true},
		{StringToStringPointer(""), true},
		{StringToStringPointer("other test"), false},
	}

	for _, test := range tests {
		res = IsStringPointerEmpty(test.input)
		assert.Equal(t, test.r, res)
	}
}
func TestIsHex(t *testing.T) {

	assert.NoError(t, IsHex("6108753dd8567400011cdc00"))
	assert.Error(t, IsHex("23423"))
	assert.Error(t, IsHex("6108753dd8567400011cdc0z"))
}
func TestStringPointerToString(t *testing.T) {
	var abc = "abc"

	assert.Equal(t, "", StringPointerToString(nil))
	assert.Equal(t, "abc", StringPointerToString(&abc))
}

func TestIsBoolean(t *testing.T) {

	assert.True(t, IsBoolean("true"))
	assert.True(t, IsBoolean("false"))
	assert.False(t, IsBoolean("falsee"))
}
func TestStringArrPointerToStringArr(t *testing.T) {
	var array = []string{"123"}

	assert.Equal(t, []string{"123"}, StringArrPointerToStringArr(&array))
	assert.Equal(t, []string{}, StringArrPointerToStringArr(nil))
	assert.Equal(t, []string{}, StringArrPointerToStringArr(&[]string{}))

}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue float64
		want         float64
		hasError     bool
	}{
		{
			name:         "1 casting successfully",
			value:        "1.0",
			defaultValue: 0,
			want:         1.0,
			hasError:     false,
		},
		{
			name:         "10 casting successfully",
			value:        "10.0",
			defaultValue: 0,
			want:         10.0,
			hasError:     false,
		},
		{
			name:         "1000000 casting successfully",
			value:        "1000000.0",
			defaultValue: 0,
			want:         1000000.0,
			hasError:     false,
		},
		{
			name:         "with default value",
			value:        "",
			defaultValue: 0,
			want:         0,
			hasError:     false,
		},
		{
			name:         "with incorrect values",
			value:        "2332323wewewewwwewew",
			defaultValue: 0,
			want:         0,
			hasError:     true,
		},
		{
			name:         "with float numbers",
			value:        "25.5",
			defaultValue: 0,
			want:         25.5,
			hasError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ToFloat64(tt.value, tt.defaultValue)

			if tt.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, val, tt.want)
		})
	}
}

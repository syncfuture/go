package json

import (
	"testing"
)

type config struct {
	Count    int
	Path     string
	Services []service
}

type service struct {
	Enabled bool
	Name    string
}

const (
	_json string = `{"Count":11,"Path":"/data/es","Services":[{"Enabled":true,"Name":"AAA"},{"Enabled":false,"Name":"BBB"}]}`
)

var _config config

func init() {
	_config = config{
		Count: 11,
		Path:  "/data/es",
		Services: []service{
			service{Enabled: true, Name: "AAA"},
			service{Enabled: false, Name: "BBB"},
		},
	}
}

func TestSerialize(t *testing.T) {
	type args struct {
		objPtr interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			args: args{objPtr: &_config},
			want: _json,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Serialize(tt.args.objPtr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializeToFile(t *testing.T) {
	type args struct {
		objPtr   interface{}
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				objPtr:   &_config,
				filename: "test.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SerializeToFile(tt.args.objPtr, tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("SerializeToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestDeserialize(t *testing.T) {
	type args struct {
		json   string
		objPtr interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				json:   `{"Count":11,"Path":"/data/es","Services":[{"Enabled":true,"Name":"AAA"},{"Enabled":false,"Name":"BBB"}]}`,
				objPtr: new(config),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Deserialize(tt.args.json, tt.args.objPtr); (err != nil) != tt.wantErr {
				t.Errorf("Deserialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeserializeFromFile(t *testing.T) {
	type args struct {
		filename string
		objPtr   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				filename: "test.json",
				objPtr:   new(config),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeserializeFromFile(tt.args.filename, tt.args.objPtr); (err != nil) != tt.wantErr {
				t.Errorf("DeserializeFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

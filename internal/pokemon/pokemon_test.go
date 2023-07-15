package pokemon

import (
	"reflect"
	"testing"
)

func TestGetPokemons(t *testing.T) {
	type args struct {
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		args    args
		want    PokemonList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPokemons(tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPokemons() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPokemons() = %v, want %v", got, tt.want)
			}
		})
	}
}

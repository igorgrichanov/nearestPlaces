package csv

import (
	"nearestPlaces/internal/entity"
	"reflect"
	"testing"
)

var sampleLine = "9\tShKOLA 735\tgorod Moskva, Aviamotornaja ulitsa, dom 51\t(495) 273-21-06\t37.72098869657803\t55.746325696672486\n"

func Test_parseLine(t *testing.T) {
	type args struct {
		line []string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.Restaurant
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				line: []string{
					"9",
					"ShKOLA 735",
					"gorod Moskva, Aviamotornaja ulitsa, dom 51",
					"(495) 273-21-06",
					"37.72098869657803",
					"55.746325696672486",
				},
			},
			want: &entity.Restaurant{
				ID:      "9",
				Name:    "ShKOLA 735",
				Address: "gorod Moskva, Aviamotornaja ulitsa, dom 51",
				Phone:   "(495) 273-21-06",
				Location: struct {
					Lon float64 `json:"lon"`
					Lat float64 `json:"lat"`
				}{Lon: 37.72098869657803, Lat: 55.746325696672486},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCSV(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    []*entity.Restaurant
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				filename: "../../../datasets/test.csv",
			},
			want: []*entity.Restaurant{
				{
					ID:      "0",
					Name:    "SMETANA",
					Address: "gorod Moskva, ulitsa Egora Abakumova, dom 9",
					Phone:   "(499) 183-14-10",
					Location: struct {
						Lon float64 `json:"lon"`
						Lat float64 `json:"lat"`
					}{Lon: 37.71456500043604, Lat: 55.879001531303366},
				},
				{
					ID:      "1",
					Name:    "Rodnik",
					Address: "gorod Moskva, ulitsa Talalihina, dom 2/1, korpus 1",
					Phone:   "(495) 676-55-35",
					Location: struct {
						Lon float64 `json:"lon"`
						Lat float64 `json:"lat"`
					}{Lon: 37.6733061300344, Lat: 55.7382386551547},
				},
				{
					ID:      "2",
					Name:    "Kafe «Akademija»",
					Address: "gorod Moskva, Abel'manovskaja ulitsa, dom 6",
					Phone:   "(495) 662-30-10",
					Location: struct {
						Lon float64 `json:"lon"`
						Lat float64 `json:"lat"`
					}{Lon: 37.6696475969381, Lat: 55.7355114718314},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := New()
			got, err := parser.ParseCSV(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCSV() got = %v, want %v", got, tt.want)
			}
		})
	}
}

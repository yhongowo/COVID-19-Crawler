package service

import (
	"sync"
	"testing"
)

func TestParser_overallParser(t *testing.T) {
	type fields struct {
		WaitGroup sync.WaitGroup
		bodyData  []byte
	}
	type args struct {
		info string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				WaitGroup: tt.fields.WaitGroup,
				bodyData:  tt.fields.bodyData,
			}
			p.overallParser(tt.args.info)
		})
	}
}

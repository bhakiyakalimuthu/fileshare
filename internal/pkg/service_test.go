package pkg

import(
	"go.uber.org/zap"
	"testing"
)

func TestGetFile(t *testing.T){
	type args struct {
		name string
		expected string
		expectedErr bool
	}
	tests := []args{
		{
			name:"Should return expected output if file name is provided",
			expected:"fake.png",
			expectedErr:false,
		},
		{
			name:"Should return error if file name is empty",
			expected:"",
			expectedErr:true,
		},
	}
	for _,tt := range tests {
		s := NewService(zap.NewNop())
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetFile(tt.expected)
			if tt.expectedErr {
				if err == nil {
					t.Errorf("Unexpected error %s",err)
				}
				if got != tt.expected {
					t.Errorf("Got %s Expected %s",got,tt.expected)
				}
				return
			}
			if err!=nil {
				t.Errorf("Unexpected error %s",err)
			}
			if got != tt.expected {
				t.Errorf("Got: %s Expected %s",got,tt.expected)
			}
		} )

	}
}

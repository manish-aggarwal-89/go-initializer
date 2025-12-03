package helper

import (
	"testing"
	"{{MODULE_NAME}}/internal/shared"
	"{{MODULE_NAME}}/internal/shared/entity"
	"{{MODULE_NAME}}/internal/shared/test_var"
)

func TestGetCurrencyFromExtendedData(t *testing.T) {
	type args struct {
		credential entity.CredentialRepository
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test with valid currency",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{
						shared.ExtendedDataKeyCurrency: shared.CURRENCY_IDR,
					},
				},
			},
			want: "IDR",
		},
		{
			name: "Test with empty currency",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCurrencyFromExtendedData(tt.args.credential); got != tt.want {
				t.Errorf("GetCurrencyFromExtendedData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDomainCodeFromExtendedData(t *testing.T) {
	type args struct {
		credential entity.CredentialRepository
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test with valid domain code",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{
						shared.ExtendedDataKeyDomainCode: test_var.CredentialExtendedDataDomainCode,
					},
				},
			},
			want: test_var.CredentialExtendedDataDomainCode,
		},
		{
			name: "Test with empty domain code",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDomainCodeFromExtendedData(tt.args.credential); got != tt.want {
				t.Errorf("GetDomainCodeFromExtendedData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAgentNameFromExtendedData(t *testing.T) {
	type args struct {
		credential entity.CredentialRepository
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test with valid agent name",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{
						shared.ExtendedDataKeyAgentName: test_var.CredentialExtendedDataAgentName,
					},
				},
			},
			want: test_var.CredentialExtendedDataAgentName,
		},
		{
			name: "Test with empty map agent name",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAgentNameFromExtendedData(tt.args.credential); got != tt.want {
				t.Errorf("GetAgentNameFromExtendedData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAgentPasswordFromExtendedData(t *testing.T) {
	type args struct {
		credential entity.CredentialRepository
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test with valid agent password",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{
						shared.ExtendedDataKeyAgentPassword: test_var.CredentialExtendedDataKeyAgentPassword,
					},
				},
			},
			want: test_var.CredentialExtendedDataKeyAgentPassword,
		},
		{
			name: "Test with empty agent password ",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAgenPasswordFromExtendedData(tt.args.credential); got != tt.want {
				t.Errorf("GetAgentPasswordFromExtendedData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOrganizationIdFromExtendedData(t *testing.T) {
	type args struct {
		credential entity.CredentialRepository
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test with valid organization id",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{
						shared.ExtendedDataKeyOrganizationId: test_var.CredentialExtendedDataKeyOrganizationId,
					},
				},
			},
			want: test_var.CredentialExtendedDataKeyOrganizationId,
		},
		{
			name: "Test with empty organization id",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOrganizationIdFromExtendedData(tt.args.credential); got != tt.want {
				t.Errorf("GetOrganizationIdFromExtendedData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPaymentMethodFromExtendedData(t *testing.T) {
	type args struct {
		credential entity.CredentialRepository
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test with valid payment method",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{
						shared.ExtendedDataKeyPaymentMethod: test_var.CredentialExtendedDataKeyPaymentMethod,
					},
				},
			},
			want: test_var.CredentialExtendedDataKeyPaymentMethod,
		},
		{
			name: "Test with empty payment method",
			args: args{
				credential: entity.CredentialRepository{
					ExtendedData: map[string]string{},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPaymentMethodFromExtendedData(tt.args.credential); got != tt.want {
				t.Errorf("GetPaymentMethodFromExtendedData() = %v, want %v", got, tt.want)
			}
		})
	}
}

package service

import (
	"context"
	"github.com/wureny/webook/webook/Internal/domain"
	"github.com/wureny/webook/webook/Internal/repository"
	"reflect"
	"testing"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		repo repository.UserRepository
	}
	tests := []struct {
		name string
		args args
		want UserService
	}{
		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_Edit(t *testing.T) {
	type fields struct {
		repo repository.UserRepository
	}
	type args struct {
		ctx context.Context
		u   domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		//TODO : Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &userService{
				repo: tt.fields.repo,
			}
			if err := svc.Edit(tt.args.ctx, tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("Edit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_FindOrCreate(t *testing.T) {
	type fields struct {
		repo repository.UserRepository
	}
	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &userService{
				repo: tt.fields.repo,
			}
			got, err := svc.FindOrCreate(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOrCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOrCreate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_GetUser(t *testing.T) {
	type fields struct {
		repo repository.UserRepository
	}
	type args struct {
		ctx context.Context
		id  uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &userService{
				repo: tt.fields.repo,
			}
			got, err := svc.GetUser(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_Login(t *testing.T) {
	type fields struct {
		repo repository.UserRepository
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &userService{
				repo: tt.fields.repo,
			}
			got, err := svc.Login(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Login() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_SignUp(t *testing.T) {
	type fields struct {
		repo repository.UserRepository
	}
	type args struct {
		ctx context.Context
		u   domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &userService{
				repo: tt.fields.repo,
			}
			if err := svc.SignUp(tt.args.ctx, tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

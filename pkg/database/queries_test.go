package database

import (
	"reflect"
	"testing"
)

func TestGetConvoID(t *testing.T) {
	type args struct {
		participants  []string
		conversations []*Conversation
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "conversation case 1",
			args: args{
				participants: []string{"975496ca-9bfc-4d71-8736-da4b6383a575", "6d01e668-2642-4e55-af73-46f057b731f9"}, //userIDs for query fake conversation 1 participant 1&2
				conversations: []*Conversation{
					{
						ConvoID: "0675de06-2d2c-444f-9d0a-ffd3303068d8",
						Participants: []User{
							{
								ID: "975496ca-9bfc-4d71-8736-da4b6383a575",
							},
							{
								ID: "6d01e668-2642-4e55-af73-46f057b731f9",
							},
						},
					},
				},
			},
			want:    "0675de06-2d2c-444f-9d0a-ffd3303068d8", //convoID for query fake conversation 1 participant 1
			wantErr: false,
		},
		{
			name: "conversation case 2 - FAIL (no userIDs in conversations)",
			args: args{
				participants:  []string{"975496ca-9bfc-4d71-8736-da4b6383a575", "6d01e668-2642-4e55-af73-46f057b731f9"}, //userIDs for query fake conversation 1 participant 1&2
				conversations: []*Conversation{},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConvoID(tt.args.participants, tt.args.conversations)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConvoID() error = %v, wantErr %v, got:%v", err, tt.wantErr, got)
				return
			}
			if got != tt.want {
				t.Errorf("GetConvoID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_compare(t *testing.T) {
	type fields struct {
		ID        string
		Nickname  string
		Age       string
		Gender    string
		FirstName string
		LastName  string
		Email     string
		Password  string
	}
	type args struct {
		t User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "works",
			fields: fields{
				ID:        "7777777777",
				Nickname:  "nickname",
				Age:       "age",
				Gender:    "gender",
				FirstName: "fname",
				LastName:  "lname",
				Email:     "email",
				Password:  "password",
			},
			args: args{
				t: User{
					ID:        "7777777778",
					Nickname:  "nickname",
					Age:       "age",
					Gender:    "gender",
					FirstName: "fname",
					LastName:  "lname",
					Email:     "email",
					Password:  "password",
				},
			},
			want: []string{
				"Nickname", "Age", "Gender", "FirstName", "LastName", "Email",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				ID:        tt.fields.ID,
				Nickname:  tt.fields.Nickname,
				Age:       tt.fields.Age,
				Gender:    tt.fields.Gender,
				FirstName: tt.fields.FirstName,
				LastName:  tt.fields.LastName,
				Email:     tt.fields.Email,
				Password:  tt.fields.Password,
			}
			if got := u.Compare(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

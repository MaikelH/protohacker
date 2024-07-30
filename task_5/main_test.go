package main

import "testing"

func Test_replaceBoguscoinAddresses(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1",
			args: args{
				input: "Please pay the ticket price of 15 Boguscoins to one of these addresses: 7RX55kE3prCgJP8peeNU1rT2Bf3CwrDQ0Q5 7ocjKK2vuH6M93BYIL1DhcvmjbYK53 76vIekZWN7VKppxonx3XC8vmrgFc1\n",
			},
			want: "Please pay the ticket price of 15 Boguscoins to one of these addresses: 7YWHMfk9JZe0LM0g1ZauHuiSxhI 7YWHMfk9JZe0LM0g1ZauHuiSxhI 7YWHMfk9JZe0LM0g1ZauHuiSxhI\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceBoguscoinAddresses(tt.args.input); got != tt.want {
				t.Errorf("replaceBoguscoinAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

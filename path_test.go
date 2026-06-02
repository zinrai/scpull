package main

import "testing"

func TestLocalPath(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		remote  string
		want    string
		wantErr bool
	}{
		{
			name:   "absolute path preserves structure",
			host:   "192.168.2.5",
			remote: "/var/log/messages",
			want:   "192.168.2.5/var/log/messages",
		},
		{
			name:   "nested path under flash",
			host:   "sw01",
			remote: "/flash/syslog.txt",
			want:   "sw01/flash/syslog.txt",
		},
		{
			name:   "trailing slash cleaned",
			host:   "host",
			remote: "/etc/nginx/",
			want:   "host/etc/nginx",
		},
		{
			name:    "relative path rejected",
			host:    "host",
			remote:  "var/log/messages",
			wantErr: true,
		},
		{
			name:    "parent traversal rejected",
			host:    "host",
			remote:  "/var/../../etc/passwd",
			wantErr: true,
		},
		{
			name:    "empty path rejected",
			host:    "host",
			remote:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := localPath(tt.host, tt.remote)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (result %q)", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("localPath(%q, %q) = %q, want %q", tt.host, tt.remote, got, tt.want)
			}
		})
	}
}

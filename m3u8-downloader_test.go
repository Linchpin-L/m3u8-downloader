package main

import "testing"

func Test_getHost(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		Url  string
		ht   string
		want string
	}{
		{
			name: "1",
			Url:  "https://www.example.com/ab/ac/ad/4477/2025-09-28-07-16-43_2025-09-28-12-56-14.m3u8?Expires=1778840028&key=123",
			ht:   "",
			want: "https://www.example.com/ab/ac/ad/4477/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getHost(tt.Url, tt.ht)
			if err != nil { 
				t.Errorf("getHost() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("getHost() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetURLDir(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "标准文件路径",
			input:    "https://example.com/path/to/file.txt",
			expected: "https://example.com/path/to/",
			wantErr:  false,
		},
		{
			name:     "路径以斜杠结尾",
			input:    "https://example.com/path/to/",
			expected: "https://example.com/path/to/",
			wantErr:  false,
		},
		{
			name:     "根路径",
			input:    "https://example.com/",
			expected: "https://example.com/",
			wantErr:  false,
		},
		{
			name:     "根路径下的文件",
			input:    "https://example.com/file.txt",
			expected: "https://example.com/",
			wantErr:  false,
		},
		{
			name:     "无扩展名的文件路径",
			input:    "https://example.com/path/to/file",
			expected: "https://example.com/path/to/",
			wantErr:  false,
		},
		{
			name:     "带端口、认证信息和查询参数的复杂URL",
			input:    "http://user:pass@example.com:8080/a/b/c?query=1#frag",
			expected: "http://user:pass@example.com:8080/a/b/",
			wantErr:  false,
		},
		{
			name:     "HTTP协议",
			input:    "http://example.com/a/b/c/d",
			expected: "http://example.com/a/b/c/",
			wantErr:  false,
		},
		{
			name:     "本地绝对路径",
			input:    "/local/path/file.go",
			expected: "/local/path/",
			wantErr:  false,
		},
		{
			name:     "本地相对路径",
			input:    "relative/path/file.go",
			expected: "relative/path/",
			wantErr:  false,
		},
		{
			name:     "只有域名没有路径",
			input:    "https://example.com",
			expected: "https://example.com/",
			wantErr:  false,
		},
		{
			name:     "多级嵌套路径",
			input:    "https://example.com/a/b/c/d/e/f/file.txt",
			expected: "https://example.com/a/b/c/d/e/f/",
			wantErr:  false,
		},
		{
			name:     "带查询参数的文件",
			input:    "https://example.com/docs/report.pdf?version=2",
			expected: "https://example.com/docs/",
			wantErr:  false,
		},
		{
			name:     "带锚点的文件",
			input:    "https://example.com/page.html#section1",
			expected: "https://example.com/",
			wantErr:  false,
		},
		{
			name:     "URL编码的路径",
			input:    "https://example.com/path%20with%20spaces/file.txt",
			expected: "https://example.com/path%20with%20spaces/",
			wantErr:  false,
		},
		{
			name:     "子域名",
			input:    "https://sub.domain.example.com/files/doc.pdf",
			expected: "https://sub.domain.example.com/files/",
			wantErr:  false,
		},
		{
			name:     "空路径（仅域名+斜杠）",
			input:    "https://example.com/",
			expected: "https://example.com/",
			wantErr:  false,
		},
		{
			name:     "无效的URL",
			input:    "://invalid-url",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetURLDir(tt.input)
			
			// 检查错误情况
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetURLDir(%q) 期望错误但没有发生", tt.input)
				}
				return
			}
			
			// 检查非错误情况
			if err != nil {
				t.Errorf("GetURLDir(%q) 发生意外错误: %v", tt.input, err)
				return
			}
			
			// 检查结果
			if result != tt.expected {
				t.Errorf("GetURLDir(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGetURLDirEdgeCases 测试边界情况
func TestGetURLDirEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "FTP协议",
			input:    "ftp://files.example.com/pub/file.zip",
			expected: "ftp://files.example.com/pub/",
			wantErr:  false,
		},
		{
			name:     "非标准端口",
			input:    "https://example.com:8443/app/file.js",
			expected: "https://example.com:8443/app/",
			wantErr:  false,
		},
		{
			name:     "路径中有多个点",
			input:    "https://example.com/archive/file.backup.2024.tar.gz",
			expected: "https://example.com/archive/",
			wantErr:  false,
		},
		{
			name:     "IPv6地址",
			input:    "https://[::1]:8080/path/file.html",
			expected: "https://[::1]:8080/path/",
			wantErr:  false,
		},
		{
			name:     "只有文件名的相对路径",
			input:    "file.txt",
			expected: "./",
			wantErr:  false,
		},
		{
			name:     "点号开头的隐藏文件",
			input:    "https://example.com/.hidden/config.yml",
			expected: "https://example.com/.hidden/",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetURLDir(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetURLDir(%q) 期望错误但没有发生", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("GetURLDir(%q) 发生意外错误: %v", tt.input, err)
				return
			}
			
			if result != tt.expected {
				t.Errorf("GetURLDir(%q) = %q, 期望 %q", tt.input, result, tt.expected)
			}
		})
	}
}

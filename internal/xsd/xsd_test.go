package xsd

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const testXSD = `<?xml version="1.0"?>
<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema" targetNamespace="urn:t" elementFormDefault="qualified">
  <xs:element name="r">
    <xs:complexType>
      <xs:simpleContent>
        <xs:extension base="xs:string">
          <xs:attribute name="n" type="xs:positiveInteger" use="required"/>
        </xs:extension>
      </xs:simpleContent>
    </xs:complexType>
  </xs:element>
</xs:schema>`

func parseTestSchema(t *testing.T) *Schema {
	t.Helper()
	s, err := Parse([]byte(testXSD))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Free)
	return s
}

func TestValidate(t *testing.T) {
	s := parseTestSchema(t)

	if err := s.Validate([]byte(`<r xmlns="urn:t" n="1">x</r>`)); err != nil {
		t.Errorf("valid document: %v", err)
	}

	var verr *ValidationError
	err := s.Validate([]byte(`<r xmlns="urn:t" n="0"/>`))
	if !errors.As(err, &verr) || len(verr.Errors) == 0 {
		t.Errorf("schema violation: got %v, want *ValidationError", err)
	}

	err = s.Validate([]byte(`<r xmlns="urn:t" n="1">`))
	if err == nil || errors.As(err, &verr) {
		t.Errorf("malformed document: got %v, want parse error", err)
	}
}

func TestValidateRejectsEntities(t *testing.T) {
	s := parseTestSchema(t)

	secret := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(secret, []byte("s3cr3t"), 0o600); err != nil {
		t.Fatal(err)
	}
	xxe := `<!DOCTYPE r [<!ENTITY x SYSTEM "file://` + secret + `">]><r xmlns="urn:t" n="1">&x;</r>`
	if err := s.Validate([]byte(xxe)); err == nil || strings.Contains(err.Error(), "s3cr3t") {
		t.Errorf("external entity: got %v, want an error that does not contain the file", err)
	}

	laughs := `<!DOCTYPE r [<!ENTITY a "aaaaaaaaaa">`
	prev := "a"
	for _, n := range []string{"b", "c", "d", "e", "f", "g", "h", "i", "j"} {
		laughs += `<!ENTITY ` + n + ` "` + strings.Repeat("&"+prev+";", 10) + `">`
		prev = n
	}
	laughs += `]><r xmlns="urn:t" n="1">&j;</r>`
	if err := s.Validate([]byte(laughs)); err == nil {
		t.Error("entity expansion: got nil error")
	}
}

// libxml2 prints diagnostics straight to file descriptor 2, which Go cannot
// intercept in-process, so this runs Validate in a child process.
func TestValidateDoesNotWriteStderr(t *testing.T) {
	if os.Getenv("XSD_STDERR_CHILD") == "1" {
		s := parseTestSchema(t)
		_ = s.Validate([]byte(`<r xmlns="urn:t" n="1">`))
		_ = s.Validate([]byte(`<r xmlns="urn:t" n="0"/>`))
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestValidateDoesNotWriteStderr$")
	cmd.Env = append(os.Environ(), "XSD_STDERR_CHILD=1")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("child process: %v\n%s", err, stderr.Bytes())
	}
	if bytes.Contains(stderr.Bytes(), []byte("error")) {
		t.Errorf("libxml2 wrote to stderr:\n%s", stderr.Bytes())
	}
}

func TestShippedSchemas(t *testing.T) {
	files, err := filepath.Glob("../../v20*/alpinebits.xsd")
	if err != nil || len(files) == 0 {
		t.Fatalf("no schemas found: %v", err)
	}
	ping := []byte(`<OTA_PingRQ xmlns="http://www.opentravel.org/OTA/2003/05" Version="8.000"><EchoData>x</EchoData></OTA_PingRQ>`)
	for _, f := range files {
		t.Run(filepath.Base(filepath.Dir(f)), func(t *testing.T) {
			b, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			s, err := Parse(b)
			if err != nil {
				t.Fatal(err)
			}
			defer s.Free()
			if err := s.Validate(ping); err != nil {
				t.Errorf("OTA_PingRQ: %v", err)
			}
		})
	}
}

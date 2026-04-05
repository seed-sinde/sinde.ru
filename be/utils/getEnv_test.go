package utils
import "testing"
func TestGetEnvVarReturnsTrimmedConfig(t *testing.T) {
	t.Setenv("POSTGRES_HOST", " 127.0.0.1 ")
	t.Setenv("POSTGRES_PORT", " 5432 ")
	t.Setenv("POSTGRES_USER", " app ")
	t.Setenv("POSTGRES_PASSWORD", " secret ")
	t.Setenv("POSTGRES_DB", " sinde ")
	t.Setenv("POSTGRES_SSL", " disable ")
	cfg := GetEnvVar()
	if cfg == nil {
		t.Fatal("expected config, got nil")
	}
	if cfg.PostgresHost != "127.0.0.1" {
		t.Fatalf("unexpected host: %q", cfg.PostgresHost)
	}
	if cfg.PostgresPort != "5432" {
		t.Fatalf("unexpected port: %q", cfg.PostgresPort)
	}
	if cfg.PostgresUser != "app" {
		t.Fatalf("unexpected user: %q", cfg.PostgresUser)
	}
	if cfg.PostgresPassword != "secret" {
		t.Fatalf("unexpected password: %q", cfg.PostgresPassword)
	}
	if cfg.PostgresDB != "sinde" {
		t.Fatalf("unexpected db: %q", cfg.PostgresDB)
	}
	if cfg.PostgresSSL != "disable" {
		t.Fatalf("unexpected ssl mode: %q", cfg.PostgresSSL)
	}
}
func TestGetEnvVarReturnsNilWhenRequiredValuesAreMissing(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_USER", "")
	t.Setenv("POSTGRES_PASSWORD", "secret")
	t.Setenv("POSTGRES_DB", "sinde")
	t.Setenv("POSTGRES_SSL", "disable")
	cfg := GetEnvVar()
	if cfg != nil {
		t.Fatalf("expected nil config, got %+v", cfg)
	}
}

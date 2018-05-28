package proxy

import "testing"

func TestProxy_Auth(t *testing.T) {
	p := Auth{
		AuthEnable: true,
		Users: map[string]string{
			"foo": "Ys23Ag/5IOWqZCw9QGaVDdHwH00=",
		},
	}

	if !p.AuthLoginPassword("foo", []byte("bar")) {
		t.Error("expected auth true, got false")
	}

	if p.AuthLoginPassword("foo", []byte("bar1")) {
		t.Error("expected auth false, got true")
	}

	if p.AuthLoginPassword("fooq", []byte("bar")) {
		t.Error("expected auth false, got true")
	}
}

func TestProxy_Auth2(t *testing.T) {
	p := Auth{
		AuthEnable: true,
		Users: map[string]string{
			"foo": "Ys23Ag/5IOWqZCw9QGaVDdHwH00=",
			"bar": "C+7Hteo/D9vJXQ3UfzxbwnXaijM=",
		},
	}

	if !p.AuthLoginPassword("foo", []byte("bar")) {
		t.Error("expected auth true, got false")
	}

	if p.AuthLoginPassword("foo", []byte("bar1")) {
		t.Error("expected auth false, got true")
	}

	if p.AuthLoginPassword("fooq", []byte("bar")) {
		t.Error("expected auth false, got true")
	}

	if !p.AuthLoginPassword("bar", []byte("foo")) {
		t.Error("expected auth true, got false")
	}

	if p.AuthLoginPassword("bar", []byte("foo1")) {
		t.Error("expected auth false, got true")
	}

	if p.AuthLoginPassword("bar1", []byte("foo")) {
		t.Error("expected auth false, got true")
	}
}

func BenchmarkAuth_Auth(b *testing.B) {
	p := Auth{
		AuthEnable: true,
		Users: map[string]string{
			"foo": "Ys23Ag/5IOWqZCw9QGaVDdHwH00=",
		},
	}
	pass := []byte("bar")
	for n := 0; n < b.N; n++ {
		if !p.AuthLoginPassword("foo", pass) {
			b.Error("expected auth true, got false")
		}
	}

}

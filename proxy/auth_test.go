package proxy

import "testing"

func TestProxy_Auth(t *testing.T) {
	p := Auth{
		AuthEnable: true,
		Users: map[string]string{
			"foo":"Ys23Ag/5IOWqZCw9QGaVDdHwH00=",
		},
	}

	if !p.Auth("foo", "bar") {
		t.Error("expected auth true, got false")
	}

	if p.Auth("foo", "bar1") {
		t.Error("expected auth false, got true")
	}

	if p.Auth("fooq", "bar") {
		t.Error("expected auth false, got true")
	}
}

func TestProxy_Auth2(t *testing.T) {
	p := Auth{
		AuthEnable: true,
		Users: map[string]string{
			"foo":"Ys23Ag/5IOWqZCw9QGaVDdHwH00=",
			"bar":"C+7Hteo/D9vJXQ3UfzxbwnXaijM=",
		},
	}

	if !p.Auth("foo", "bar") {
		t.Error("expected auth true, got false")
	}

	if p.Auth("foo", "bar1") {
		t.Error("expected auth false, got true")
	}

	if p.Auth("fooq", "bar") {
		t.Error("expected auth false, got true")
	}

	if !p.Auth("bar", "foo") {
		t.Error("expected auth true, got false")
	}

	if p.Auth("bar", "foo1") {
		t.Error("expected auth false, got true")
	}

	if p.Auth("bar1", "foo") {
		t.Error("expected auth false, got true")
	}
}

func BenchmarkAuth_Auth(b *testing.B) {
	p := Auth{
		AuthEnable: true,
		Users: map[string]string{
			"foo":"Ys23Ag/5IOWqZCw9QGaVDdHwH00=",
		},
	}

	for n := 0; n < b.N; n++ {
		if !p.Auth("foo", "bar") {
			b.Error("expected auth true, got false")
		}
	}


}
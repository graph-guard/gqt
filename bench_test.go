package gqt_test

import (
	"testing"

	"github.com/graph-guard/gqt/v4"
)

func Benchmark(b *testing.B) {
	src := []byte(`# variables are in global template scope
	query {
	  # any id
	  user(id: *) {
		# allow a maximum of 1000 nodes to be fetched in one request
		# with a maxium depth of 4
		orders(after: *, limit = $limitOrders: < 1000) {
		  id
		  created
		  ... on DirectOrder {
			status
		  }
		  # Again, a maximum of 1000 nodes must not be exceeded by this request
		  items(after: *, limit = $limitItems: < 1000 / $limitOrders) {
			id
			title
			description
			relatedProducts(
			  after: *,
			  limit: < 1000 / $limitOrders / $limitItems,
			) {
			  id
			  title
			  description
			}
		  }
		}
	  }
	}
	`)
	var errs []gqt.Error
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if _, _, errs = gqt.Parse(src); len(errs) > 0 {
			b.Fatal(errs)
		}
	}
}

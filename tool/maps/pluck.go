package maps

func PluckSliceMap[KT K, VT V](sm []map[KT]VT, keys []KT) []map[KT]VT {
	var nsm []map[KT]VT
	for _, m := range sm {
		nsm = append(nsm, PluckMap(m, keys))
	}

	return nsm
}

func PluckMap[KT K, VT V](m map[KT]VT, keys []KT) map[KT]VT {
	nm := make(map[KT]VT)
	for _, key := range keys {
		nm[key] = m[key]
	}

	return nm
}

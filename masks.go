package dns

const (
	maskQROffset     = 15
	maskOpCodeOffset = 11
	maskAAOffset     = 10
	maskTCOffset     = 9
	maskRDOffset     = 8
	maskRAOffset     = 7
	maskZOffset      = 4

	maskQR     = 1 << maskQROffset
	maskOpCode = 15 << maskOpCodeOffset
	maskAA     = 1 << maskAAOffset
	maskTC     = 1 << maskTCOffset
	maskRD     = 1 << maskRDOffset
	maskRA     = 1 << maskRAOffset
	maskZ      = 7 << maskZOffset
	maskRCode  = 15
)

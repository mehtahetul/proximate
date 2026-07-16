package models

// BucketDistance converts an exact distance in metres into a coarse,
// human-friendly label shown to clients. The exact float is never exposed —
// only used internally for sorting and cursor pagination.
func BucketDistance(metres float64) string {
	switch {
	case metres <= 100:
		return "📍 Very Close (0–100 m)"
	case metres <= 250:
		return "🚶 Short Walk (100–250 m)"
	case metres <= 500:
		return "🚶‍♂️ Nearby (250–500 m)"
	case metres <= 1000:
		return "🏙️ Within 1 km (500 m–1 km)"
	default:
		return "📡 Over 1 km (1 km+)"
	}
}
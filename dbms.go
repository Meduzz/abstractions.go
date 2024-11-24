package abstractions

import (
	"time"

	"github.com/Meduzz/abstractions.go/internal/dbms/caching"
	"github.com/Meduzz/abstractions.go/internal/dbms/csrf"
	"github.com/Meduzz/abstractions.go/internal/dbms/log"
	"github.com/Meduzz/abstractions.go/lib"
	"gorm.io/gorm"
)

func DbmsWorkLogDelegate(db *gorm.DB) (lib.WorkLogDelegate, error) {
	return log.NewDbmsWorkLogDelegate(db)
}

func DbmsCachingDelegate(db *gorm.DB, eviction lib.Eviction, ttl time.Duration, prefix string) (lib.CacheStorageDelegate, error) {
	return caching.NewDbmsCachingDelegate(db, eviction, prefix, ttl)
}

func DbmsCsrfDelegate(db *gorm.DB, ttl time.Duration) (lib.CSRFStorageDelegate, error) {
	return csrf.NewDbmsCSRFStorageDelegate(db, ttl)
}

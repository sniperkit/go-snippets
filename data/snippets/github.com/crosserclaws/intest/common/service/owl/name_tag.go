package owl

import (
	"fmt"
	cache "github.com/crosserclaws/intest/common/ccache"
	owlDb "github.com/crosserclaws/intest/common/db/owl"
	owlModel "github.com/crosserclaws/intest/common/model/owl"
)

type NameTagService struct {
	cache       *cache.CacheCtrl
	cacheConfig *cache.DataCacheConfig
}

func NewNameTagService(cacheConfig cache.DataCacheConfig) *NameTagService {
	return &NameTagService{
		cacheConfig: &cacheConfig,
		cache:       cache.NewCacheCtrl(cache.NewDataCache(cacheConfig)),
	}
}

func (s *NameTagService) GetNameTagById(nameTagId int16) *owlModel.NameTag {
	v := s.cache.MustFetchNativeAndDoNotCacheEmpty(
		nameTagKeyById(nameTagId),
		s.cacheConfig.Duration,
		func() interface{} {
			return owlDb.GetNameTagById(nameTagId)
		},
	)

	if v == nil {
		return nil
	}

	return v.(*owlModel.NameTag)
}

func (s *NameTagService) GetNameTagsByIds(nameTagIds ...int16) []*owlModel.NameTag {
	result := make([]*owlModel.NameTag, 0)

	for _, id := range nameTagIds {
		result = append(result, s.GetNameTagById(id))
	}

	return result
}

func nameTagKeyById(id int16) string {
	return fmt.Sprintf("!nid!%d", id)
}

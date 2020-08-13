package mongo

import(
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"github.com/xxxmicro/base/domain/model"
	breflect "github.com/xxxmicro/base/reflect"
)

func buildQuery(ms *breflect.StructInfo, filters map[string]interface{}) (bson.M, error) {
	bFilters := bson.M{}
	if filters == nil || len(filters) == 0 {
		return bFilters, nil
	}

	for k, v := range filters {
		filterType := model.FilterType(k)

		switch(filterType) {
		case model.FilterType_AND:
			subFilters, ok := v.([]interface{})
			if !ok {
				return nil, errors.New("ERR_MALFORMED_PARAMETERS")
			}
			
			subBFilters := make([]bson.M, len(subFilters))
			for i, sub := range subFilters {
				subFilter, ok := sub.(map[string]interface{})
				if !ok {
					return nil, errors.New("ERR_MALFORMED_PARAMETERS")
				}
				
				subBFilter, err := buildQuery(ms, subFilter)
				if err != nil {
					return nil, err
				}
				subBFilters[i] = subBFilter
			}
			bFilters["$and"] = subBFilters
		case model.FilterType_OR:
			subFilters, ok := v.([]interface{})
			if !ok {
				return nil, errors.New("ERR_MALFORMED_PARAMETERS")
			}
			
			subBFilters := make([]bson.M, len(subFilters))
			for i, sub := range subFilters {
				subFilter, ok := sub.(map[string]interface{})
				if !ok {
					return nil, errors.New("ERR_MALFORMED_PARAMETERS")
				}

				subBFilter, err := buildQuery(ms, subFilter)
				if err != nil {
					return nil, err
				}
				subBFilters[i] = subBFilter
				bFilters["$or"] = subBFilters
			}
		case model.FilterType_NOR:
			subFilters, ok := v.([]interface{})
			if !ok {
				return nil, errors.New("ERR_MALFORMED_PARAMETERS")
			}
			
			subBFilters := make([]bson.M, len(subFilters))
			for i, sub := range subFilters {
				subFilter, ok := sub.(map[string]interface{})
				if !ok {
					return nil, errors.New("ERR_MALFORMED_PARAMETERS")
				}

				subBFilter, err := buildQuery(ms, subFilter)
				if err != nil {
					return nil, err
				}
				subBFilters[i] = subBFilter
				bFilters["$nor"] = subBFilters
			}
		default:
			if _, ok := ms.FieldsMap[k]; !ok {
				err := errors.New(fmt.Sprintf("ERR_DB_UNKNOWN_FIELD %s", k))
				return nil, err
			}

			bFilter, err := buildMongoFilter(ms, v)
			if err != nil {
				return nil, err
			}
			bFilters[k] = bFilter
		}	
	}
	return bFilters, nil
}

func buildMongoFilter(ms *breflect.StructInfo, value interface{}) (bson.M, error) {
	vMap, ok := value.(map[string]interface{})
	if !ok {
		return bson.M{"$eq": value}, nil
	}

	for vKey, vValue := range vMap {
		filterType := model.FilterType(vKey)
		switch filterType {
		case model.FilterType_EQ:
			return bson.M{"$eq": vValue}, nil
		case model.FilterType_NE:
			return bson.M{"$ne": vValue}, nil
		case model.FilterType_GT:
			return bson.M{"$gt": vValue}, nil
		case model.FilterType_GTE:
			return bson.M{"$gte": vValue}, nil	
		case model.FilterType_LT:
			return bson.M{"$lt": vValue}, nil
		case model.FilterType_LTE:
			return bson.M{"$lte": vValue}, nil
		case model.FilterType_LIKE:
			return bson.M{"$regex": vValue}, nil
		case model.FilterType_MATCH:
			return bson.M{"$regex": vValue}, nil	
		case model.FilterType_NOT_LIKE:
			return bson.M{"$not": bson.M{"$regex": vValue}}, nil
		case model.FilterType_IN:
			return bson.M{"$in": vValue}, nil
		case model.FilterType_NOT_IN:
			return bson.M{"$nin": vValue}, nil
		case model.FilterType_IS_NULL:	
			return bson.M{"$exists": false}, nil
		case model.FilterType_NOT_NULL:
			return bson.M{"$exists": true}, nil
		default:
			return nil, errors.New("ERR_MALFORMED_FILTER_TYPE")
		}
	}
	return bson.M{}, nil
}

func buildSort(ms *breflect.StructInfo, sorts []*model.SortSpec) ([]string, error){
	bsorts := []string{}
	if sorts != nil {
		for _, s := range sorts {
			if _, ok := ms.FieldsMap[s.Property]; !ok {
				return nil, errors.New(fmt.Sprintf("ERR_DB_UNKNOWN_FIELD %s", s.Property))
			}
			var s1 string
			switch s.Type {
			case model.SortType_DSC:
				{
					s1 = fmt.Sprintf("-%s", s.Property)
				}
			default: // SortType_ASC
				{
					s1 = s.Property
				}
			}
			bsorts = append(bsorts, s1)
		}
	}
	return bsorts, nil
}

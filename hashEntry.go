package ConcurrentHashMap

type hashEntry[K comparable, V any] struct {
	Datas []*pair[K, V]
}

func (h *hashEntry[K, V]) Get(key K) *V {
	if h.Datas == nil {
		return nil
	}
	for _, v := range h.Datas {
		if v.Key == key {
			return &v.Value
		}
	}
	return nil
}
func (h *hashEntry[K, V]) Put(key K, value V) (int64, error) {
	if h.Datas == nil {
		h.Datas = make([]*pair[K, V], 0)
	}
	for _, element := range h.Datas {
		if element.Key == key {
			element.Value = value
			return 0, nil
		}
	}
	h.Datas = append(h.Datas, &pair[K, V]{
		Key:   key,
		Value: value,
	})
	return 1, nil
}
func (h *hashEntry[K, V]) Remove(key K) (int64, error) {
	for index, element := range h.Datas {
		if element.Key == key {
			h.Datas = append(h.Datas[:index], h.Datas[index+1:]...)
			return 1, nil
		}
	}
	return 0, nil
}

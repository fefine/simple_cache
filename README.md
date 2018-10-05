## cache
* cache sth in memory
* a *get* operation will be refresh the expire time if key not expired.
* only support little count k-v

refer: [https://github.com/patrickmn/go-cache](https://github.com/patrickmn/go-cache)
### usage


    # demo 1
    cache := NewCache(interval)
    cache.Set(key, value, 1 * time.Second)

    v, ok := cache.get(key)
    if ok {
        val := v.(string)
        fmt.Print(val)
    }

    # demo 2
    cache := NewCache(interval)
    cache.Set(key, value, 3 * time.Second)

    time.Sleep(2 * time.Second)
    // refresh
    v, ok := cache.get(key)
    if ok {
        val := v.(string)
        fmt.Print(val)
    }

    time.Sleep(2 * time.Second)
    // also can get value
    v, ok = cache.get(key)
    if ok {
        val := v.(string)
        fmt.Print(val)
    }










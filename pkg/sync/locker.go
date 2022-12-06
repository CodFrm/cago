package sync

// Locker 分布式锁接口,默认使用redis
type Locker interface {
	// TryLock 尝试获取锁,不会阻塞
	TryLock(opts ...LockOption) error
	// Lock 加锁,未获取到锁时会阻塞,可以设置一个超时时长
	Lock(opts ...LockOption) error
	// Unlock 解锁
	Unlock() error
	// TryLockKey 根据某个key去尝试获取锁
	TryLockKey(key string, opts ...LockOption) error
	// LockKey 根据某个key去进行加锁
	LockKey(key string, opts ...LockOption) error
	// UnlockKey 根据某个key去解锁
	UnlockKey(key string) error
}

// NewLocker 新建一个锁对象
func NewLocker(keyPrefix string) Locker {
	return newRedis(keyPrefix)
}

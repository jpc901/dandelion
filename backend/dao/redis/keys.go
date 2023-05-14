package redis

const (
	KeyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time"   // KeyPostTimeZSet ZSet:帖子发帖时间
	KeyPostScoreZSet   = "post:score"  // KeyPostScoreZSet ZSet:帖子投票分数
	KeyPostVotedZSetPF = "post:voted:" // KeyPostVotedZSetPF ZSet:记录用户及投票的类型，参数是postID
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}

package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertToTimestamp(date time.Time) *timestamppb.Timestamp {
	return timestamppb.New(date)
}

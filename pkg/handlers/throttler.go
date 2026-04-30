package handlers

import "journal-migrator/pkg/handlers/throttlers"

type BaseThrottler = throttlers.BaseThrottler
type ReplicaLagThrottler = throttlers.ReplicaLagThrottler
type WaitingTimeThrottler = throttlers.WaitingTimeThrottler

var NewReplicaLagThrottler = throttlers.NewReplicaLagThrottler
var NewWaitingTimeThrottler = throttlers.NewWaitingTimeThrottler

package types

type UpgradeStatus struct {
	AcceptedReleaseInfo *ReleaseInfo
	ReadyToUpgrade      bool
}

type ReleaseInfo struct {
	FuryReleaseTag     string
	UpgradeBlockHeight uint64
}

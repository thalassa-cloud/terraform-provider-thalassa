package dbaas

import "time"

// DbObjectStoreRecoveryOverview is the aggregated recovery overview for a DB object store.
type DbObjectStoreRecoveryOverview struct {
	ObjectStoreIdentity string                                 `json:"objectStoreIdentity"`
	AnalysedAt          time.Time                              `json:"analysedAt"`
	Complete            bool                                   `json:"complete"`
	Truncation          DbObjectStoreRecoveryTruncation        `json:"truncation"`
	Summary             DbObjectStoreRecoveryStoreSummary      `json:"summary"`
	Clusters            []DbObjectStoreRecoveryClusterOverview `json:"clusters"`
}

// DbObjectStoreRecoveryTruncation describes incomplete recovery analysis.
type DbObjectStoreRecoveryTruncation struct {
	Truncated               bool    `json:"truncated"`
	Reason                  *string `json:"reason"`
	MaxWalObjectsConsidered *int    `json:"maxWalObjectsConsidered"`
}

// DbObjectStoreRecoveryStoreSummary aggregates recovery metrics across clusters.
type DbObjectStoreRecoveryStoreSummary struct {
	CompletedBackupCount    int        `json:"completedBackupCount"`
	UnavailableBackupCount  int        `json:"unavailableBackupCount"`
	FailedBackupCount       int        `json:"failedBackupCount"`
	OldestBackupAt          *time.Time `json:"oldestBackupAt"`
	NewestBackupAt          *time.Time `json:"newestBackupAt"`
	EarliestPitrFrom        *time.Time `json:"earliestPitrFrom"`
	LatestPitrThroughApprox *time.Time `json:"latestPitrThroughApprox"`
	LatestArchivedWal       string     `json:"latestArchivedWal"`
	LatestWalArchivedAt     *time.Time `json:"latestWalArchivedAt"`
	Timelines               []string   `json:"timelines"`
	GapCount                int        `json:"gapCount"`
}

// DbObjectStoreRecoveryClusterOverview is per-cluster recovery state within a store.
type DbObjectStoreRecoveryClusterOverview struct {
	ClusterIdentity     string                                `json:"clusterIdentity"`
	ClusterName         string                                `json:"clusterName"`
	ClusterExists       bool                                  `json:"clusterExists"`
	RecoveryWindow      *DbObjectStoreRecoveryWindowDTO       `json:"recoveryWindow"`
	LatestArchivedWal   string                                `json:"latestArchivedWal"`
	LatestWalArchivedAt *time.Time                            `json:"latestWalArchivedAt"`
	Timelines           []string                              `json:"timelines"`
	GapCount            int                                   `json:"gapCount"`
	Backups             []DbObjectStoreRecoveryBackupCoverage `json:"backups"`
}

// DbObjectStoreRecoveryWindowDTO mirrors Barman server recovery window fields in recovery APIs.
type DbObjectStoreRecoveryWindowDTO struct {
	FirstRecoverabilityPoint *time.Time `json:"firstRecoverabilityPoint"`
	LastSuccessfulBackupTime *time.Time `json:"lastSuccessfulBackupTime"`
	LastFailedBackupTime     *time.Time `json:"lastFailedBackupTime"`
}

// DbObjectStoreRecoveryBackupCoverage embeds backup fields plus WAL coverage.
type DbObjectStoreRecoveryBackupCoverage struct {
	BackupIdentity string                                    `json:"backupIdentity"`
	Status         string                                    `json:"status"`
	StatusMessage  *string                                   `json:"statusMessage"`
	CreatedAt      time.Time                                 `json:"createdAt"`
	StartedAt      *time.Time                                `json:"startedAt"`
	StoppedAt      *time.Time                                `json:"stoppedAt"`
	SizeBytes      *int64                                    `json:"sizeBytes"`
	BeginWal       string                                    `json:"beginWal"`
	EndWal         string                                    `json:"endWal"`
	BeginLsn       string                                    `json:"beginLsn"`
	EndLsn         string                                    `json:"endLsn"`
	Coverage       DbObjectStoreRecoveryBackupCoverageDetail `json:"coverage"`
}

// DbObjectStoreRecoveryBackupCoverageDetail is continuous WAL chain analysis for one backup.
type DbObjectStoreRecoveryBackupCoverageDetail struct {
	ChainStatus            DbObjectStoreWalChainStatus       `json:"chainStatus"`
	EndWal                 string                            `json:"endWal"`
	WalAvailableTo         string                            `json:"walAvailableTo"`
	FirstMissingWal        *string                           `json:"firstMissingWal"`
	ContinuousSegmentCount int                               `json:"continuousSegmentCount"`
	WalChainSizeBytes      int64                             `json:"walChainSizeBytes"`
	WalChainObjectCount    int                               `json:"walChainObjectCount"`
	PitrThroughApproxAt    *time.Time                        `json:"pitrThroughApproxAt"`
	PitrThroughSource      DbObjectStoreWalPitrThroughSource `json:"pitrThroughSource"`
}

// DbObjectStoreWalChainStatus describes continuous WAL coverage from a backup's end WAL.
type DbObjectStoreWalChainStatus string

const (
	DbObjectStoreWalChainStatusContinuous    DbObjectStoreWalChainStatus = "continuous"
	DbObjectStoreWalChainStatusGapDetected   DbObjectStoreWalChainStatus = "gap_detected"
	DbObjectStoreWalChainStatusEndWalMissing DbObjectStoreWalChainStatus = "end_wal_missing"
	DbObjectStoreWalChainStatusNoWal         DbObjectStoreWalChainStatus = "no_wal"
	DbObjectStoreWalChainStatusUnknown       DbObjectStoreWalChainStatus = "unknown"
	DbObjectStoreWalChainStatusNotApplicable DbObjectStoreWalChainStatus = "not_applicable"
)

// DbObjectStoreWalPitrThroughSource indicates how an approximate PITR upper bound was derived.
type DbObjectStoreWalPitrThroughSource string

const (
	DbObjectStoreWalPitrSourceWalObjectLastModified DbObjectStoreWalPitrThroughSource = "wal_object_last_modified"
	DbObjectStoreWalPitrSourceBarmanRecoveryWindow  DbObjectStoreWalPitrThroughSource = "barman_recovery_window"
	DbObjectStoreWalPitrSourceBackupStoppedAt       DbObjectStoreWalPitrThroughSource = "backup_stopped_at"
	DbObjectStoreWalPitrSourceUnknown               DbObjectStoreWalPitrThroughSource = "unknown"
)

// DbObjectStoreRecoveryWalHierarchy is the default WAL browser hierarchy view.
type DbObjectStoreRecoveryWalHierarchy struct {
	ClusterIdentity string                             `json:"clusterIdentity"`
	WalDirectory    string                             `json:"walDirectory"`
	AnalysedAt      time.Time                          `json:"analysedAt"`
	Complete        bool                               `json:"complete"`
	Truncation      DbObjectStoreRecoveryTruncation    `json:"truncation"`
	Timelines       []DbObjectStoreRecoveryWalTimeline `json:"timelines"`
}

// DbObjectStoreRecoveryWalTimeline groups WAL log folders under one timeline.
type DbObjectStoreRecoveryWalTimeline struct {
	TimelineHex string                        `json:"timelineHex"`
	Timeline    uint32                        `json:"timeline"`
	Logs        []DbObjectStoreRecoveryWalLog `json:"logs"`
}

// DbObjectStoreRecoveryWalLog summarises one timeline/log folder.
type DbObjectStoreRecoveryWalLog struct {
	LogHex           string     `json:"logHex"`
	Log              uint32     `json:"log"`
	TimelineLogHex   string     `json:"timelineLogHex"`
	ObjectCount      int        `json:"objectCount"`
	SizeBytes        int64      `json:"sizeBytes"`
	OldestArchivedAt *time.Time `json:"oldestArchivedAt"`
	NewestArchivedAt *time.Time `json:"newestArchivedAt"`
	FirstWal         string     `json:"firstWal"`
	LastWal          string     `json:"lastWal"`
	GapCount         int        `json:"gapCount"`
}

// DbObjectStoreRecoveryWalSegments lists objects for one timeline/log folder or a flat view.
type DbObjectStoreRecoveryWalSegments struct {
	ClusterIdentity       string                           `json:"clusterIdentity"`
	TimelineHex           string                           `json:"timelineHex"`
	LogHex                string                           `json:"logHex"`
	Objects               []DbObjectStoreRecoveryWalObject `json:"objects"`
	NextContinuationToken *string                          `json:"nextContinuationToken"`
}

// DbObjectStoreRecoveryWalObject is a single object under the WAL tree.
type DbObjectStoreRecoveryWalObject struct {
	Key                 string               `json:"key"`
	FileName            string               `json:"fileName"`
	Kind                DbObjectStoreWalKind `json:"kind"`
	WalName             string               `json:"walName"`
	TimelineHex         string               `json:"timelineHex"`
	LogHex              string               `json:"logHex"`
	SegmentHex          string               `json:"segmentHex"`
	BackupHistoryOffset *string              `json:"backupHistoryOffset"`
	SizeBytes           int64                `json:"sizeBytes"`
	ArchivedAt          *time.Time           `json:"archivedAt"`
}

// DbObjectStoreWalKind classifies an object under a Barman WAL tree.
type DbObjectStoreWalKind string

const (
	DbObjectStoreWalKindWAL             DbObjectStoreWalKind = "wal"
	DbObjectStoreWalKindBackupHistory   DbObjectStoreWalKind = "backup_history"
	DbObjectStoreWalKindTimelineHistory DbObjectStoreWalKind = "timeline_history"
	DbObjectStoreWalKindUnknown         DbObjectStoreWalKind = "unknown"
)

// DbObjectStoreRecoveryWalSearch is the result of a WAL name search.
type DbObjectStoreRecoveryWalSearch struct {
	ClusterIdentity string                           `json:"clusterIdentity"`
	Query           string                           `json:"query"`
	Objects         []DbObjectStoreRecoveryWalObject `json:"objects"`
}

// GetDbObjectStoreRecoveryOverviewRequest carries optional filters for recovery overview.
type GetDbObjectStoreRecoveryOverviewRequest struct {
	// IncludeUnavailable includes unavailable backups. Nil defaults to the API default (true).
	IncludeUnavailable *bool
	ClusterID          string
	BackupID           string
}

// GetDbObjectStoreRecoveryWalHierarchyRequest carries optional filters for the hierarchy view.
type GetDbObjectStoreRecoveryWalHierarchyRequest struct {
	Timeline string
	Kinds    []string
}

// ListDbObjectStoreRecoveryWalSegmentsRequest carries query params for segment/flat WAL listing.
type ListDbObjectStoreRecoveryWalSegmentsRequest struct {
	// View should be "flat" for a flat listing without timeline/log, or empty when timeline and log are set.
	View              string
	Timeline          string
	Log               string
	ContinuationToken string
	Limit             int
	Kinds             []string
}

// SearchDbObjectStoreRecoveryWalRequest carries query params for WAL search.
type SearchDbObjectStoreRecoveryWalRequest struct {
	Query string
	Limit int
}

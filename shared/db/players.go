package db

// represents common properties of statistics
type BaseStat interface {
	Type() StatType
}

// represents a minimal statistic
type MinStat interface {
	Type() StatType
	Total() uint16
}

// represents a statistic
type Stat interface {
	Type() StatType
	Base() uint16
	Equipment() uint16
	Gift() uint16
	Context() uint16
}

// represents an extended statistic
type ExStat interface {
	Type() StatType
	Total() uint16
	Base() uint16
	Equipment() uint16
	Gift() uint16
	Context() uint16
}

type StatGetter interface {
	GetStat(t StatType) (BaseStat, bool)
}

type StatList interface {
	Stats() map[StatType]BaseStat
}

type Stats interface {
	StatGetter
	StatList
}

// represents a stat type
type StatType int

const (
	InvalidStatType StatType = iota
	Prospection
	Initiative
	ActionPoints
	MovementPoints
	Strength
	Vitality
	Wisdom
	Chance
	Agility
	Intelligence
	RangePoints
	Summons
	Damage
	PhysicalDamage
	WeaponControl
	DamagePer
	HealPoints
	TrapDamage
	TrapDamagePer
	DamageReturn
	CriticalHit
	CriticalFailure
	DodgeActionPoints
	DodgeMovementPoints
	ResistanceNeutral
	ResistancePercentNeutral
	ResistancePvpNeutral
	ResistancePercentPvpNeutral
	ResistanceEarth
	ResistancePercentEarth
	ResistancePvpEarth
	ResistancePercentPvpEarth
	ResistanceWater
	ResistancePercentWater
	ResistancePvpWater
	ResistancePercentPvpWater
	ResistanceWind
	ResistancePercentWind
	ResistancePvpWind
	ResistancePercentPvpWind
	ResistanceFire
	ResistancePercentFire
	ResistancePvpFire
	ResistancePercentPvpFire
)

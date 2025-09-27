package family

const (
	AngelCommand    = "Angel Command"
	Armorloid       = "Armorloid"
	ArmoredDragon   = "Armored Dragon"
	ArmoredWyvern   = "Armored Wyvern"
	BalloonMushroom = "Balloon Mushroom"
	BeastFolk       = "Beast Folk"
	Berserker       = "Berserker"
	BrainJacker     = "Brain Jacker"
	Chimera         = "Chimera"
	ColonyBeetle    = "Colony Beetle"
	CyberCluster    = "Cyber Cluster"
	CyberLord       = "Cyber Lord"
	CyberVirus      = "Cyber Virus"
	DarkLord        = "Dark Lord"
	DeathPuppet     = "Death Puppet"
	DemonCommand    = "Demon Command"
	DevilMask       = "Devil Mask"
	Dragonoid       = "Dragonoid"
	DuneGecko       = "Dune Gecko"
	EarthDragon     = "Earth Dragon"
	EarthEater      = "Earth Eater"
	Family          = "Family"
	FireBird        = "Fire Bird"
	Fish            = "Fish"
	GelFish         = "Gel Fish"
	Giant           = "Giant"
	GiantInsect     = "Giant Insect"
	Gladiator       = "Gladiator"
	Ghost           = "Ghost"
	Guardian        = "Guardian"
	Hedrian         = "Hedrian"
	HornedBeast     = "Horned Beast"
	Human           = "Human"
	Initiate        = "Initiate"
	Leviathan       = "Leviathan"
	LightBringer    = "Light Bringer"
	LiquidPeople    = "Liquid People"
	LivingDead      = "Living Dead"
	MachineEater    = "Machine Eater"
	MechaDelSol     = "Mecha del Sol"
	MechaThunder    = "Mecha Thunder"
	MysteryTotem    = "Mystery Totem"
	ParasiteWorm    = "Parasite Worm"
	RainbowPhantom  = "Rainbow Phantom"
	RockBeast       = "Rock Beast"
	SeaHacker       = "Sea Hacker"
	SnowFaerie      = "Snow Faerie"
	StarlightTree   = "Starlight Tree"
	Survivor        = "Survivor"
	TreeFolk        = "Tree Folk"
	VolcanoDragon   = "Volcano Dragon"
	Xenoparts       = "Xenoparts"
	ZombieDragon    = "Zombie Dragon"
)

var Families = []string{
	AngelCommand,
	Armorloid,
	ArmoredDragon,
	ArmoredWyvern,
	BalloonMushroom,
	BeastFolk,
	Berserker,
	BrainJacker,
	Chimera,
	ColonyBeetle,
	CyberCluster,
	CyberLord,
	CyberVirus,
	DarkLord,
	DeathPuppet,
	DemonCommand,
	DevilMask,
	Dragonoid,
	DuneGecko,
	EarthDragon,
	EarthEater,
	FireBird,
	Fish,
	GelFish,
	Giant,
	GiantInsect,
	Gladiator,
	Ghost,
	Guardian,
	Hedrian,
	HornedBeast,
	Human,
	Initiate,
	Leviathan,
	LightBringer,
	LiquidPeople,
	LivingDead,
	MachineEater,
	MechaDelSol,
	MechaThunder,
	MysteryTotem,
	ParasiteWorm,
	RainbowPhantom,
	RockBeast,
	SeaHacker,
	SnowFaerie,
	StarlightTree,
	Survivor,
	TreeFolk,
	VolcanoDragon,
	Xenoparts,
	ZombieDragon,
}

func GetAllFamilies() []string {
	return Families
}

var Cybers = []string{CyberCluster, CyberLord, CyberVirus}
var Dragons = []string{ArmoredDragon, VolcanoDragon, EarthDragon, ZombieDragon}

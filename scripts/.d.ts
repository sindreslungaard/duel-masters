declare interface Player {
  id: string;
  username: string;
}

declare interface Card {
  id: string;
  name: string;
  power: number;
  civilization: string;
  race: string;
  manaCost: number;
  manaRequirements: string[];
}

declare interface CardConstruct extends Card {
  setName(string);
  setPower(number);
  setCivilization(string);
  setRace(string);
  setManaCost(number);
  setManaReq(string);
  useTrait(string);
}

interface Events {
  card_moved: CardMovedEvent;
  attack_confirmed: {
    attacker: string;
  };
}

type EvType = keyof Events;
type Ev<T extends EvType> = Events[T];

declare interface Match {
  on<T extends EvType>(type: T, callback: (data: Ev<T>) => void);
}

declare global {
  const $self: CardConstruct;
  const $match: Match;
}

// Container
declare enum Container {
  Deck = "deck",
  Hand = "hand",
  Shieldzone = "shieldzone",
  Manazone = "manazone",
  Graveyard = "graveyard",
  Battlezone = "battlezone",
  Spellzone = "spellzone",
  Hiddenzone = "hiddenzone",
}

// Race
declare enum Race {
  MachineEater = "Machine Eater",
  Berserker = "Berserker",
  BeastFolk = "Beast Folk",
  ArmoredDragon = "Armored Dragon",
  ArmoredWyvern = "Armored Wyvern",
  ColonyBeetle = "Colony Beetle",
  HornedBeast = "Horned Beast",
  Leviathan = "Leviathan",
  CyberVirus = "Cyber Virus",
  GelFish = "Gel Fish",
  BalloonMushroom = "Balloon Mushroom",
  LivingDead = "Living Dead",
  AngelCommand = "Angel Command",
  Chimera = "Chimera",
  DemonCommand = "Demon Command",
  CyberLord = "Cyber Lord",
  LiquidPeople = "Liquid People",
  Initiate = "Initiate",
  RockBeast = "Rock Beast",
  GiantInsect = "Giant Insect",
  Ghost = "Ghost",
  Dragonoid = "Dragonoid",
  LightBringer = "Light Bringer",
  StarlightTree = "Starlight Tree",
  Human = "Human",
  Guardian = "Guardian",
  ParasiteWorm = "Parasite Worm",
  Armorloid = "Armorloid",
  TreeFolk = "Tree Folk",
  Fish = "Fish",
  DarkLord = "Dark Lord",
  BrainJacker = "Brain Jacker",
  Giant = "Giant",
  MechaThunder = "Mecha Thunder",
  CyberCluster = "Cyber Cluster",
  Hedrian = "Hedrian",
  FireBird = "Fire Bird",
}

// Conditions
declare enum Condition {
  SummoningSickness = "summoning_sickness",
  DoubleBreaker = "doublebreaker",
  TripleBreaker = "triplebreaker",
  PowerAmplifier = "power_amplifier",
  PowerAttacker = "power_attacker",
  AttackUntapped = "attack_untapped",
  Creature = "creature",
  Spell = "spell",
  Blocker = "blocker",
  ShieldTrigger = "shield_trigger",
  Slayer = "slayer",
  Active = "active",
  CantBeBlocked = "cant_be_blocked",
  CantAttackCreatures = "cant_attack_creatures",
  CantAttackPlayers = "cant_attack_players",
  ReducedCost = "reduced_cost",
  IncreasedCost = "increased_cost",
  Evolution = "evolution",
  Survivor = "survivor",
}

// Civilization
declare enum Civilization {
  Fire = "fire",
  Water = "water",
  Nature = "nature",
  Light = "light",
  Darkness = "darkness",
}

// Events
interface CardMovedEvent {
  cardId: string;
  to: Container;
  from: Container;
}

export {};

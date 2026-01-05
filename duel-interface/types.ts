// Message is the default message struct
export interface Message {
  header: string;
}

// DecksMessage lists the users decks
export interface DecksMessage {
  header: string;
  decks: LegacyDeck[];
}

// ChatMessage stores information about a chat message
export interface ChatMessage {
  header: string;
  message: string;
  sender: string;
  color: string;
}

export const TAPPED_FLAG = 1;
export const PLAYABLE_FLAG = 2;
export const TAP_ABILITY_FLAG = 4;
export const SHIELD_FACE_UP_FLAG = 8;

export const cardHasFlag = (flags: number, flag: number): boolean => {
  return (flags & flag) > 0;
};

// CardState stores information about the state of a card
export interface CardState {
  virtualId: string;
  uid: string;
  name: string;
  civilization: string;
  flags: number;
}

// ShieldState stores information about the state of a shield
export interface ShieldState {
  virtualId: string;
  uid: string;
  flags: number;
}

// PlayerState stores information about the state of the current player
export interface PlayerState {
  username: string;
  color: string;
  deck: number;
  handCount: number;
  hand: CardState[];
  shieldzone: ShieldState[];
  shieldMap: Record<string, number>;
  manazone: CardState[];
  graveyard: CardState[];
  playzone: CardState[];
}

// MatchState stores information about the current state of the match in the eyes of a given player
export interface MatchState {
  myTurn: boolean;
  hasAddedManaThisRound: boolean;
  me: PlayerState;
  opponent: PlayerState;
  spectator: boolean;
}

// MatchStateMessage is the message that should be sent to the client for state updates
export interface MatchStateMessage {
  header: string;
  state: MatchState;
}

// WarningMessage is used to send a warning to a player
export interface WarningMessage {
  header: string;
  message: string;
}

export enum ActionType {
  None = "",
  Count = "count",
  Question = "question",
  Order = "order",
  Searchable = "searchable",
}

// ActionMessage is used to prompt the user to make a selection of the specified cards
export interface ActionMessage {
  header: string;
  actionType: ActionType;
  cards?: CardState[] | Record<string, CardState[]>;
  text: string;
  minSelections: number;
  maxSelections: number;
  cancellable: boolean;
  unselectableCards: CardState[];
  choices: string[];
}

// MultipartActionMessage is used to prompt the user to make a selection of the specified cards
export interface MultipartActionMessage {
  header: string;
  cards: Record<string, CardState[]>;
  text: string;
  minSelections: number;
  maxSelections: number;
  cancellable: boolean;
}

// ActionWarningMessage is used to apply an error
export interface ActionWarningMessage {
  header: string;
  message: string;
}

// WaitMessage is used to send a waiting popup with a message
export interface WaitMessage {
  header: string;
  message: string;
}

// LobbyChatMessage is used to store chat messages
export interface LobbyChatMessage {
  username: string;
  color: string;
  message: string;
  timestamp: number;
  removed?: boolean;
}

// LobbyChatMessages is used to store chat messages
export interface LobbyChatMessages {
  header: string;
  messages: LobbyChatMessage[];
}

// UserMessage holds information about users
export interface UserMessage {
  username: string;
  color: string;
  hub: string;
  permissions: string[];
}

// UserListMessage is used to send a list of online users
export interface UserListMessage {
  header: string;
  users: UserMessage[];
}

// MatchMessage holds information about a match
export interface MatchMessage {
  id: string;
  p1: string;
  p1color: string;
  p2: string;
  p2color: string;
  name: string;
  spectate: boolean;
  matchmaking: boolean;
  format: string;
}

// MatchesListMessage is used to list open matches
export interface MatchesListMessage {
  header: string;
  matches: MatchMessage[];
}

export interface MatchRequestMessage {
  id: string;
  name: string;
  host_id: string;
  host_name: string;
  host_color: string;
  guest_id: string;
  guest_name: string;
  guest_color: string;
  format: string;
  link_code: string;
}

export interface MatchRequestsListMessage {
  header: string;
  requests: MatchRequestMessage[];
}

export interface MatchForwardMessage {
  header: string;
  id: string;
}

// ShowCardsMessage is used to show the user n cards without an action to perform
export interface ShowCardsMessage {
  header: string;
  message: string;
  cards: string[];
}

export interface PinnedMessages {
  header: string;
  messages: string[];
}

export interface PlaySoundMessage {
  header: string;
  sound: string;
}

// Type alias for LegacyDeck (assuming it's imported from db module)
export type LegacyDeck = any; // Replace with actual type definition if available

<template>
  <!-- CardPreviewPopup can be teleported if z-index becomes an issue -->
  <CardPreviewPopup :uid=previewCard?.uid :event=previewCardEvent :xPos=previewCardX :side="'left'" />
  <div ref="decklistref">
    <div class="deck-card-slot" v-for="(card, index) in getCardsForDeck(deck, cards)" :key="index"
      @mouseover="setPreviewCard($event, card)" @mouseleave="previewCard = null">
      <div class="deck-card-counter">
        {{ card.count }}x
      </div>
      <div class="deck-card-oval" :class="'card-' + card.civilization.toLowerCase()" :key="index"
        @click="tryRemoveCard(card)">
        <div class="deck-card-name">
          {{ card.name }}
        </div>
        <div class="deck-card-mana">
          {{ card.manaCost }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import CardPreviewPopup from "../components/CardPreviewPopup.vue";
import { getCardsForDeck, playSound } from "../helpers/utils";
import { ref } from "vue";

const deck = defineModel('deck')
const props = defineProps(['cards'])

const previewCard = ref(null);
const previewCardX = ref(0);
const previewCardEvent = ref(null);
const decklistref = ref(null)

function setPreviewCard(event, card) {
  previewCardEvent.value = event;
  previewCardX.value = decklistref.value.getBoundingClientRect().x;
  previewCard.value = card;
}

function tryRemoveCard(card) {
  let uid = card.uid;
  let toSlice = deck.value.indexOf(uid);
  if (toSlice < 0) return;

  deck.value.splice(toSlice, 1);
  playSound("/assets/sounds/card-removed.wav");
  if (deck.value.indexOf(uid) < 0 && previewCard.value.uid === uid) {
    previewCard.value = null
  }
}
</script>

<style scoped lang="scss">
$civ-nature-base-color: #118141;
$civ-fire-base-color: #D12027;
$civ-water-base-color: #47C6F2;
$civ-light-base-color: #FAD241;
$civ-darkness-base-color: #65696C;

.deck-card-slot {
  display: flex;
  margin-bottom: 4px;
  align-items: center;
  width: 100%;
  cursor: pointer;
}

.deck-card-oval {
  border-radius: 10px / 20px;
  color: #000;
  padding-left: 7px;
  padding-right: 3px;
  border: 3px solid;
  flex-grow: 1;
  display: inline-flex;
  justify-content: space-between;
  align-items: center;
  min-width: 0;
}

.deck-card-name {
  text-overflow: ellipsis;
  overflow: hidden;
  white-space: nowrap;
  font-size: 17px;
  font-family: "Perceval Bold", "Roboto", sans-serif;
  user-select: none;
}

.deck-card-mana {
  display: inline-block;
  background-color: black;
  color: white;
  font-weight: 600;
  border-radius: 50%;
  padding-left: 7px;
  padding-right: 7px;
  margin-top: 2px;
  margin-bottom: 2px;
}

.deck-card-counter {
  border-radius: 10px / 20px;
  color: black;
  padding-left: 6px;
  padding-right: 9px;
  padding-top: 4px;
  padding-bottom: 4px;
  background: radial-gradient(40% 100% at right, transparent 50%, lightgray 51%);
  display: inline-block;
  flex-grow: 0;
  margin-right: -3px;
}

.card-fire {
  border-color: $civ-fire-base-color;
  background: radial-gradient(#F0DDD5 0%, #E39E7F 80%, $civ-fire-base-color 100%);
}

.card-nature {
  border-color: $civ-nature-base-color;
  background: radial-gradient(#E1EBD0 0%, #84C188 80%, $civ-nature-base-color 100%);
}

.card-water {
  border-color: $civ-water-base-color;
  background: radial-gradient(#D4EDF8 0%, #93D4DE 80%, $civ-water-base-color 100%);
}

.card-darkness {
  border-color: $civ-darkness-base-color;
  background: radial-gradient(#CCD1D5 0%, #CACFD2 80%, $civ-darkness-base-color 100%);
}

.card-light {
  border-color: $civ-light-base-color;
  background: radial-gradient(#FCF7D0 0%, #FBF396 80%, $civ-light-base-color 100%);
}
</style>

<template>
  <div class="flex flex-col">
    <div class="catalogue-top-bar">
      <div class="catalogue-filters">
        <input v-model="filterCard" class="catalogue-filter" type="search" placeholder="Type to search" />
        <div v-for="civ in ['light', 'darkness', 'nature', 'fire', 'water']" class="civ-filter"
          :class="['civ-color-' + civ, filterCivilization[civ] ? 'civ-filter-selected' : '']"
          @click="filterCivilization[civ] = !filterCivilization[civ]">
        </div>
        <div style="margin-left: 5px;">
          <select class="catalogue-filter" v-model="filterFamily">
            <option class="family" v-for="(family, index) in families" :key="index" :value="family">{{ family }}
            </option>
          </select>
        </div>
        <div v-for="manaNr in ['1', '2', '3', '4', '5', '6', '7+']" class="mana-filter"
          :class="{ 'mana-filter-selected': filterMana[manaNr] }" @click="filterMana[manaNr] = !filterMana[manaNr]">
          {{ manaNr }}
        </div>
        <select v-model="filterSet" class="catalogue-filter">
          <option class="set" v-for="(set, index) in sets" :key="index" :value="set">{{ set }}</option>
        </select>
        <img class="reset-icon" src="/assets/images/reset-icon.svg" v-tooltip="'Reset all filters'"
          @click="resetFilters" />
      </div>
      <div class="zoom-icons">
        <img class="zoom-icon" @click="modifyCatalogueCardSize(10)" src="/assets/images/zoom-in-icon.svg"
          v-tooltip="'Increase card size'" />
        <img @click="modifyCatalogueCardSize(-10)" class="zoom-icon" src="/assets/images/zoom-out-icon.svg"
          v-tooltip="'Decrease card size'" />
      </div>
    </div>
    <div class="catalogue-cards-wrapper">
      <div class="catalogue-cards">
        <div v-for="card in filteredAndSortedCards" class="catalogue-card" :style="{ 'height': cardSize + 'px' }"
          @click="tryAddCard(card)">
          <v-lazy-image class="max-w-full rounded-2xl" :src="`https://scans.shobu.io/${card.uid}.jpg`"
            src-placeholder="https://scans.shobu.io/backside.jpg" :alt="card.name" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import VLazyImage from "v-lazy-image";
import { compareCards, playSound } from "../helpers/utils";
import { ref, watch, computed } from 'vue';

const ALL_FAMILIES = 'All Races';
const ALL_SETS = 'All Sets'

const permissions = () => {
  let p = localStorage.getItem("permissions");
  if (!p) {
    return [];
  }
  return p;
};

const deck = defineModel('deck')
const props = defineProps({
  unsortedCards: Array,
})

const filterCard = ref("")
const filterFamily = ref(ALL_FAMILIES)
const filterSet = ref(ALL_SETS)
const families = ref([ALL_FAMILIES, "Spell"])
const filterCivilization = ref({
  "light": false,
  "darkness": false,
  "nature": false,
  "fire": false,
  "water": false,
})
const filterMana = ref({
  '1': false,
  '2': false,
  '3': false,
  '4': false,
  '5': false,
  '6': false,
  '7+': false,
})
const sets = ref([])

const allcards = ref([])

const cardSize = ref(360)
function modifyCatalogueCardSize(addition) {
  cardSize.value += addition
}

function resetFilters() {
  Object.keys(filterCivilization.value).forEach(civ => filterCivilization.value[civ] = false);
  Object.keys(filterMana.value).forEach(manaCost => filterMana.value[manaCost] = false);
  filterFamily.value = ALL_FAMILIES;
  filterSet.value = ALL_SETS;
  filterCard.value = '';
}

watch(
  () => props.unsortedCards,
  (newCards) => {
    let newSets = {};
    let cardsCiv = {
      "light": [],
      "darkness": [],
      "nature": [],
      "fire": [],
      "water": [],
    };
    for (let card of newCards) {
      cardsCiv[card.civilization.toLowerCase()].push(card);
      if (!newSets[card.set]) {
        newSets[card.set] = true;
      }
    }

    sets.value = Object.keys(newSets);
    sets.value.push(ALL_SETS);
    sets.value.sort();

    let sortedCards = [];
    Object.values(cardsCiv).forEach(
      civSet => sortedCards.push(...civSet.sort((c1, c2) =>
        compareCards(c1, c2, { by: "manaCost", directionNum: 1 })
      ))
    );

    allcards.value = sortedCards;

    let newFamilies = [];
    for (let c of allcards.value) {
      if (c.family) {
        for (let f of c.family) {
          if (!newFamilies.includes(f)) {
            newFamilies.push(f);
          }
        }
      }
    }
    newFamilies.sort();
    families.value.push(...newFamilies);
  },
  { immediate: true, deep: true },
)

function tryAddCard(card) {
  if (
    props.deck.filter(x => x == card.uid).length >= 4
  ) {
    if (!permissions().includes("admin")) {
      playSound("/assets/sounds/card-limit.wav");
      return;
    }
  }

  deck.value.push(card.uid)
  playSound("/assets/sounds/card-added.mp3");
}

const filteredAndSortedCards = computed(() => {
  let cards = allcards.value.filter(card =>
    card.name.toLowerCase().includes(filterCard.value.toLowerCase()) ||
    card.text.toLowerCase().includes(filterCard.value.toLowerCase())
  );

  if (filterSet.value !== ALL_SETS) {
    cards = cards.filter(card => card.set === filterSet.value);
  }

  let filterCivilizationValues = Object.values(filterCivilization.value);
  if (!(filterCivilizationValues.every(v => v === filterCivilizationValues[0]))) {
    cards = cards.filter(
      card => filterCivilization.value[card.civilization]
    );
  }

  if (filterFamily.value.toLowerCase() !== ALL_FAMILIES.toLowerCase()) {
    cards = cards.filter(
      card =>
        (filterFamily.value.toLowerCase() === "spell" && !card.family) ||
        (card.family && card.family.includes(filterFamily.value))
    );
  }

  let filterManaValues = Object.values(filterMana.value);
  if (!(filterManaValues.every(v => v === filterManaValues[0]))) {
    cards = cards.filter(
      card => {
        if (card.manaCost > 6)
          return filterMana.value['7+'];
        else
          return filterMana.value[card.manaCost.toString()];
      }
    );
  }

  return cards;
})

</script>

<style scoped lang="scss">
$catalogue-main-color: #b8b7b7;
$catalogue-secondary-color: #e0dede;
$civ-nature-base-color: #118141;
$civ-fire-base-color: #D12027;
$civ-water-base-color: #47C6F2;
$civ-light-base-color: #FAD241;
$civ-darkness-base-color: #65696C;

.zoom-icons {
  display: inline-flex;
  align-items: flex-end;
  align-self: end;
  padding: 5px 10px 5px 10px;
  margin-right: 10px;
  margin-left: 5px;
  justify-content: space-between;
  gap: 10px;
  background-color: $catalogue-main-color;
}

.zoom-icon {
  width: 25px;
  cursor: pointer;
}

.reset-icon {
  width: 35px;
  margin-left: 20px;
  background-color: lightsalmon;
  cursor: pointer;
}

.civ-filter {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: inline-block;
  margin-left: 5px;
  cursor: pointer;
  border: 4px solid $catalogue-secondary-color;
}

.civ-filter-selected {
  border: 4px solid black;
}

.civ-color-fire {
  background-color: $civ-fire-base-color;
}

.civ-color-water {
  background-color: $civ-water-base-color;
}

.civ-color-light {
  background-color: $civ-light-base-color;
}

.civ-color-darkness {
  background-color: $civ-darkness-base-color;
}

.civ-color-nature {
  background-color: $civ-nature-base-color;
}

.mana-filter {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-left: 5px;
  cursor: pointer;
  color: white;
  background-color: black;
  font-weight: 600;
  border: 4px solid $catalogue-secondary-color;
}

.mana-filter-selected {
  border: 4px solid orange;
}

.catalogue-cards-wrapper {
  overflow-y: auto;
  width: 100%;
  height: 100%;
  background-color: $catalogue-main-color;
}

.catalogue-cards {
  width: calc(100% - 10px);
  display: inline-flex;
  flex-wrap: wrap;
  gap: 5px;
  justify-content: space-evenly;

  padding: 5px;
}

.catalogue-top-bar {
  display: flex;
  justify-content: space-between;
  background-color: $catalogue-secondary-color;
}

.catalogue-filters {
  margin-left: 10px;
  flex-wrap: wrap;
  display: inline-flex;
  align-items: center;
  padding: 3px 0;
}

.catalogue-filter {
  background-color: black;
  color: white;
}

.catalogue-card {
  aspect-ratio: 264/367;
  margin-bottom: 5px;
  cursor: pointer;
}

input,
select {
  padding: 5px !important;
  width: auto !important;
  margin-left: 5px;
  border-radius: 4px;
  color: #ccc;
  resize: none;
  border: none;
  background: url(/assets/images/overlay_30.png);
}

select option {
  background: #111;
}

input:focus,
select:focus {
  outline: none;
}
</style>

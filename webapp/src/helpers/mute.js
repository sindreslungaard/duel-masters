export const toggleMute = player => {
  const mutedPlayers = getMutedPlayers();
  const isMuted = mutedPlayers.includes(player);

  if (isMuted) {
    setMutedPlayers(mutedPlayers.filter(p => p !== player));
  } else {
    setMutedPlayers([...mutedPlayers, player]);
  }

  return !isMuted;
};

export const isMuted = player => getMutedPlayers().includes(player);
export const getMutedPlayers = () =>
  JSON.parse(localStorage.getItem("muted") ?? "[]");

export const didSeeMuteWarning = () => {
  if (localStorage.getItem("didMute")) {
    return true;
  }

  localStorage.setItem("didMute", "true");
  return false;
};

const setMutedPlayers = players => {
  localStorage.setItem("muted", JSON.stringify(players));

  // "storage" event does not fire in the same window. fire a custom event
  // setTimeout is needed for vue to have time to reconcile state
  setTimeout(() => dispatchEvent(new Event("storageUpdated")));
};
